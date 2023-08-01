//go:build linux
// +build linux

package firewall

import (
	"net"
	"runtime"

	"golang.org/x/sys/unix"

	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
	"github.com/google/nftables/expr"

	"github.com/vishvananda/netns"

	"wgnetwork/model"
	"wgnetwork/pkg/ifconfig"
	"wgnetwork/pkg/nfutils"
)

const loIface = "lo"

// NFTables struct.
type NFTables struct {
	cfg Config

	originNetNS netns.NsHandle
	targetNetNS netns.NsHandle

	wanIface string
	wanIP    net.IP
	wgIface  string
	wgPort   uint16

	tFilter         *nftables.Table
	cInput          *nftables.Chain
	cForwardManaged *nftables.Chain
	cForward        *nftables.Chain
	cOutput         *nftables.Chain

	tNAT         *nftables.Table
	cPrerouting  *nftables.Chain
	cPostrouting *nftables.Chain

	filterSetTrustIP     *nftables.Set
	filterSetWGManagerIP *nftables.Set
	filterSetWGForwardIP *nftables.Set

	managerPorts []uint16

	applied bool
}

// Init nftables firewall.
func Init(
	cfg Config,
	managerPorts []uint16,
) (*NFTables, error) {
	// obtain default interface name, ip address and gateway ip address
	wanIface, _, wanIP, err := ifconfig.IPAddr()
	if err != nil {
		return nil, err
	}

	defaultPolicy := nftables.ChainPolicyDrop
	if cfg.DefaultPolicy == "accept" {
		defaultPolicy = nftables.ChainPolicyAccept
	}

	tFilter := &nftables.Table{Family: nftables.TableFamilyIPv4, Name: "filter"}
	cInput := &nftables.Chain{
		Name:     "input",
		Table:    tFilter,
		Type:     nftables.ChainTypeFilter,
		Priority: nftables.ChainPriorityFilter,
		Hooknum:  nftables.ChainHookInput,
		Policy:   &defaultPolicy,
	}

	// var (
	// 	chainPriorityFilterManaged = nftables.ChainPriorityRef(
	// 		*nftables.ChainPriorityFilter - 1)
	// 	defaultPolicyManaged = nftables.ChainPolicyAccept
	// )
	cForwardManaged := &nftables.Chain{
		Name:  "forward_managed",
		Table: tFilter,
		Type:  nftables.ChainTypeFilter,
		// Priority: chainPriorityFilterManaged,
		// Hooknum:  nftables.ChainHookForward,
		// Policy:   &defaultPolicyManaged,
	}
	cForward := &nftables.Chain{
		Name:     "forward",
		Table:    tFilter,
		Type:     nftables.ChainTypeFilter,
		Priority: nftables.ChainPriorityFilter,
		Hooknum:  nftables.ChainHookForward,
		Policy:   &defaultPolicy,
	}
	cOutput := &nftables.Chain{
		Name:     "output",
		Table:    tFilter,
		Type:     nftables.ChainTypeFilter,
		Priority: nftables.ChainPriorityFilter,
		Hooknum:  nftables.ChainHookOutput,
		Policy:   &defaultPolicy,
	}

	tNAT := &nftables.Table{Family: nftables.TableFamilyIPv4, Name: "nat"}
	cPrerouting := &nftables.Chain{
		Name:     "prerouting",
		Table:    tNAT,
		Type:     nftables.ChainTypeNAT,
		Priority: nftables.ChainPriorityNATDest,
		// Priority: nftables.ChainPriorityFilter,
		Hooknum: nftables.ChainHookPrerouting,
	}
	cPostrouting := &nftables.Chain{
		Name:     "postrouting",
		Table:    tNAT,
		Type:     nftables.ChainTypeNAT,
		Priority: nftables.ChainPriorityNATSource,
		Hooknum:  nftables.ChainHookPostrouting,
	}

	filterSetTrustIP := &nftables.Set{
		Name:    "trust_ipset",
		Table:   tFilter,
		KeyType: nftables.TypeIPAddr,
	}
	filterSetWGManagerIP := &nftables.Set{
		Name:    "wgmanager_ipset",
		Table:   tFilter,
		KeyType: nftables.TypeIPAddr,
	}
	filterSetWGForwardIP := &nftables.Set{
		Name:    "wgforward_ipset",
		Table:   tFilter,
		KeyType: nftables.TypeIPAddr,
	}

	nft := &NFTables{
		cfg: cfg,

		wanIface: wanIface,
		wanIP:    wanIP,
		wgIface:  cfg.WGIface,
		wgPort:   cfg.WGPort,

		tFilter:         tFilter,
		cInput:          cInput,
		cForwardManaged: cForwardManaged,
		cForward:        cForward,
		cOutput:         cOutput,

		tNAT:         tNAT,
		cPrerouting:  cPrerouting,
		cPostrouting: cPostrouting,

		filterSetTrustIP:     filterSetTrustIP,
		filterSetWGManagerIP: filterSetWGManagerIP,
		filterSetWGForwardIP: filterSetWGForwardIP,

		managerPorts: managerPorts,
	}

	err = nft.apply()
	if err != nil {
		return nil, err
	}

	return nft, nil
}

// networkNamespaceBind target by name.
func (nft *NFTables) networkNamespaceBind() (*nftables.Conn, error) {
	if nft.cfg.NetworkNamespace == "" {
		return &nftables.Conn{NetNS: int(nft.originNetNS)}, nil
	}

	// Lock the OS Thread so we don't accidentally switch namespaces
	runtime.LockOSThread()

	origin, err := netns.Get()
	if err != nil {
		nft.networkNamespaceRelease()
		return nil, err
	}

	nft.originNetNS = origin

	target, err := netns.GetFromName(nft.cfg.NetworkNamespace)
	if err != nil {
		nft.networkNamespaceRelease()
		return nil, err
	}

	// switch to target network namespace
	err = netns.Set(target)
	if err != nil {
		nft.networkNamespaceRelease()
		return nil, err
	}
	nft.targetNetNS = target

	return &nftables.Conn{NetNS: int(nft.targetNetNS)}, nil
}

// networkNamespaceRelease to origin.
func (nft *NFTables) networkNamespaceRelease() error {
	if nft.cfg.NetworkNamespace == "" {
		return nil
	}

	// finally unlock os thread
	defer runtime.UnlockOSThread()

	// switch back to the original namespace
	err := netns.Set(nft.originNetNS)
	if err != nil {
		return err
	}

	// close fd to origin and dev ns
	nft.originNetNS.Close()
	nft.targetNetNS.Close()

	nft.targetNetNS = 0

	return nil
}

// apply rules
func (nft *NFTables) apply() error {
	if !nft.cfg.Enabled {
		return nil
	}

	// bind network namespace if it was set in config
	c, err := nft.networkNamespaceBind()
	if err != nil {
		return err
	}
	// release network namespace finally
	defer nft.networkNamespaceRelease()

	c.FlushRuleset()
	//
	// Init Tables and Chains.
	//

	// add filter table
	// cmd: nft add table ip filter
	c.AddTable(nft.tFilter)
	// add input chain of filter table
	// cmd: nft add chain ip filter input \
	// { type filter hook input priority 0 \; policy drop\; }
	c.AddChain(nft.cInput)
	// add forward_managed chain
	// cmd: nft add chain ip filter forward_managed \
	// { type filter hook forward priority -1 \; policy drop\; }
	c.AddChain(nft.cForwardManaged)
	// add forward chain
	// cmd: nft add chain ip filter forward \
	// { type filter hook forward priority 0 \; policy drop\; }
	c.AddChain(nft.cForward)
	// add output chain
	// cmd: nft add chain ip filter output \
	// { type filter hook output priority 0 \; policy drop\; }
	c.AddChain(nft.cOutput)

	// add nat table
	// cmd: nft add table ip nat
	c.AddTable(nft.tNAT)
	// add prerouting chain
	// cmd: nft add chain ip nat prerouting \
	// { type nat hook prerouting priority -100 \; }
	c.AddChain(nft.cPrerouting)
	// add postrouting chain
	// cmd: nft add chain ip nat postrouting \
	// { type nat hook postrouting priority 100 \; }
	c.AddChain(nft.cPostrouting)

	//
	// Init sets.
	//

	// add trust_ipset
	// cmd: nft add set ip filter trust_ipset { type ipv4_addr\; }
	// --
	// set trust_ipset {
	//         type ipv4_addr
	// }
	err = c.AddSet(nft.filterSetTrustIP, nil)
	if err != nil {
		return err
	}

	// add wgmanager_ipset
	// cmd: nft add set ip filter wgmanager_ipset { type ipv4_addr\; }
	// --
	// set wgmanager_ipset {
	//         type ipv4_addr
	// }
	err = c.AddSet(nft.filterSetWGManagerIP, nil)
	if err != nil {
		return err
	}

	// add wgforward_ipset
	// cmd: nft add set ip filter wgforward_ipset { type ipv4_addr\; }
	// --
	// set wgforward_ipset {
	//         type ipv4_addr
	// }
	err = c.AddSet(nft.filterSetWGForwardIP, nil)
	if err != nil {
		return err
	}

	//
	// Init filter rules.
	//

	err = nft.inputLocalIfaceRules(c)
	if err != nil {
		return err
	}
	err = nft.outputLocalIfaceRules(c)
	if err != nil {
		return err
	}
	err = nft.inputHostBaseRules(c, nft.wanIface)
	if err != nil {
		return err
	}
	err = nft.outputHostBaseRules(c, nft.wanIface)
	if err != nil {
		return err
	}
	err = nft.inputTrustIPSetRules(c, nft.wanIface)
	if err != nil {
		return err
	}
	err = nft.outputTrustIPSetRules(c, nft.wanIface)
	if err != nil {
		return err
	}
	err = nft.inputPublicRules(c, nft.wanIface)
	if err != nil {
		return err
	}
	err = nft.outputPublicRules(c, nft.wanIface)
	if err != nil {
		return err
	}
	err = nft.sdnRules(c)
	if err != nil {
		return err
	}
	err = nft.sdnForwardManagedRules(c)
	if err != nil {
		return err
	}
	err = nft.sdnForwardRules(c)
	if err != nil {
		return err
	}
	err = nft.dnatRules(c)
	if err != nil {
		return err
	}
	err = nft.snatRules(c)
	if err != nil {
		return err
	}

	for _, iface := range nft.cfg.Ifaces {
		if iface == nft.wanIface {
			continue
		}

		err = nft.inputHostBaseRules(c, iface)
		if err != nil {
			return err
		}
		err = nft.outputHostBaseRules(c, iface)
		if err != nil {
			return err
		}
		err = nft.inputTrustIPSetRules(c, iface)
		if err != nil {
			return err
		}
		err = nft.outputTrustIPSetRules(c, iface)
		if err != nil {
			return err
		}
		err = nft.inputPublicRules(c, iface)
		if err != nil {
			return err
		}
		err = nft.outputPublicRules(c, iface)
		if err != nil {
			return err
		}
	}

	// apply configuration
	err = c.Flush()
	if err != nil {
		return err
	}
	nft.applied = true

	return nil
}

// inputLocalIfaceRules to apply.
func (nft *NFTables) inputLocalIfaceRules(c *nftables.Conn) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter input meta iifname "lo" accept
	// --
	// iifname "lo" accept
	fr = &model.FilterRule{
		IIF:    []string{loIface},
		Action: "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter input meta iifname != "lo" \
	// ip saddr 127.0.0.0/8 reject
	// --
	// iifname != "lo" ip saddr 127.0.0.0/8 reject with icmp type prot-unreachable
	fr = &model.FilterRule{
		IIF:    []string{loIface},
		IIFNeq: true,
		SAddr: [][2]*net.IPNet{
			{
				&net.IPNet{
					IP:   net.IPv4(127, 0, 0, 0).To4(),
					Mask: net.IPv4Mask(255, 0, 0, 0),
				},
				nil,
			},
		},
		RejectActionType: unix.NFT_REJECT_ICMP_UNREACH,
		RejectActionCode: unix.NFT_REJECT_ICMPX_UNREACH,
		Action:           "reject",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// outputLocalIfaceRules to apply.
func (nft *NFTables) outputLocalIfaceRules(c *nftables.Conn) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter output meta oifname "lo" accept
	// --
	// oifname "lo" accept
	fr = &model.FilterRule{
		OIF:    []string{loIface},
		Action: "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// inputHostBaseRules to apply.
func (nft *NFTables) inputHostBaseRules(c *nftables.Conn, iface string) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter input meta iifname "eth0" ip protocol icmp \
	// ct state { established, related } accept
	// --
	// iifname "eth0" ip protocol icmp ct state { established, related } accept
	fr = &model.FilterRule{
		IIF:       []string{iface},
		Proto:     []string{"icmp"},
		Conntrack: []string{"established", "related"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol udp udp sport 53 \
	// ct state established accept
	// --
	// iifname "eth0" { udp, tcp } sport domain ct state established accept
	fr = &model.FilterRule{
		IIF:       []string{iface},
		Proto:     []string{"tcp", "udp"},
		SPort:     [][2]uint16{{53, 0}},
		Conntrack: []string{"established"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol udp udp sport 123 \
	// ct state established accept
	// --
	// iifname "eth0" { udp } sport ntp ct state established accept
	fr = &model.FilterRule{
		IIF:       []string{iface},
		Proto:     []string{"udp"},
		SPort:     [][2]uint16{{123, 0}},
		Conntrack: []string{"established"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol tcp tcp sport { 80, 443 } \
	// ct state established accept
	// --
	// iifname "eth0" tcp sport { http, https } ct state established accept
	fr = &model.FilterRule{
		IIF:       []string{iface},
		Proto:     []string{"tcp"},
		SPort:     [][2]uint16{{80, 0}, {443, 0}},
		Conntrack: []string{"established"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// outputHostBaseRules to apply.
func (nft *NFTables) outputHostBaseRules(c *nftables.Conn, iface string) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter output meta oifname "eth0" ip protocol icmp \
	// ct state { new, established } accept
	// --
	// oifname "eth0" ip protocol icmp ct state { established, new } accept
	fr = &model.FilterRule{
		OIF:       []string{iface},
		Proto:     []string{"icmp"},
		Conntrack: []string{"established", "new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol udp udp dport 53 \
	// ct state { new, established } accept
	// --
	// oifname "eth0" { udp, tcp } dport domain ct state { established, new } accept
	fr = &model.FilterRule{
		OIF:       []string{iface},
		Proto:     []string{"tcp", "udp"},
		DPort:     [][2]uint16{{53, 0}},
		Conntrack: []string{"established", "new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol udp udp dport 123 \
	// ct state { new, established } accept
	// --
	// oifname "eth0" { udp } dport ntp ct state { established, new } accept
	fr = &model.FilterRule{
		OIF:       []string{iface},
		Proto:     []string{"udp"},
		DPort:     [][2]uint16{{123, 0}},
		Conntrack: []string{"established", "new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol tcp tcp dport { 80, 443 } \
	// ct state { new, established } accept
	// --
	// oifname "eth0" tcp dport { http, https } ct state { established, new } accept
	fr = &model.FilterRule{
		OIF:       []string{iface},
		Proto:     []string{"tcp"},
		DPort:     [][2]uint16{{80, 0}, {443, 0}},
		Conntrack: []string{"established", "new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// inputTrustIPSetRules to apply.
func (nft *NFTables) inputTrustIPSetRules(c *nftables.Conn, iface string) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter input meta iifname "eth0" ip protocol icmp \
	// icmp type echo-request ip saddr @trust_ipset ct state new accept
	// --
	// iifname "eth0" icmp type echo-request ip saddr @trust_ipset ct state new accept
	fr = &model.FilterRule{
		IIF:       []string{iface},
		Proto:     []string{"icmp"},
		ICMPType:  []byte{0x08},
		SAddrSet:  nft.filterSetTrustIP,
		Conntrack: []string{"new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol tcp tcp dport { 5522 } ip saddr @trust_ipset \
	// ct state { new, established } accept
	// --
	// iifname "eth0" tcp dport { 5522 } ip saddr @trust_ipset ct state { established, new } accept
	trustPorts := make([][2]uint16, len(nft.cfg.TrustPorts))
	for i, p := range nft.cfg.TrustPorts {
		trustPorts[i] = [2]uint16{p, 0}
	}

	fr = &model.FilterRule{
		IIF:       []string{iface},
		Proto:     []string{"tcp"},
		SAddrSet:  nft.filterSetTrustIP,
		DPort:     trustPorts,
		Conntrack: []string{"established", "new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// outputTrustIPSetRules to apply.
func (nft *NFTables) outputTrustIPSetRules(c *nftables.Conn, iface string) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol tcp tcp sport { 5522 } ip daddr @trust_ipset \
	// ct state established accept
	// --
	// oifname "eth0" tcp sport { 5522 } ip daddr @trust_ipset ct state established accept
	trustPorts := make([][2]uint16, len(nft.cfg.TrustPorts))
	for i, p := range nft.cfg.TrustPorts {
		trustPorts[i] = [2]uint16{p, 0}
	}

	fr = &model.FilterRule{
		OIF:       []string{iface},
		Proto:     []string{"tcp"},
		DAddrSet:  nft.filterSetTrustIP,
		SPort:     trustPorts,
		Conntrack: []string{"established"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// inputPublicRules to apply.
func (nft *NFTables) inputPublicRules(c *nftables.Conn, iface string) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol udp udp dport 51820 accept
	// --
	// iifname "eth0" udp dport 51820 accept
	fr = &model.FilterRule{
		IIF:    []string{iface},
		Proto:  []string{"udp"},
		DPort:  [][2]uint16{{nft.wgPort, 0}},
		Action: "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// outputPublicRules to apply.
func (nft *NFTables) outputPublicRules(c *nftables.Conn, iface string) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol udp udp sport 51820 accept
	// --
	// oifname "eth0" udp sport 51820 accept
	fr = &model.FilterRule{
		OIF:    []string{iface},
		Proto:  []string{"udp"},
		SPort:  [][2]uint16{{nft.wgPort, 0}},
		Action: "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// sdnRules to apply.
func (nft *NFTables) sdnRules(c *nftables.Conn) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter input meta iifname "wg0" ip protocol icmp \
	// icmp type echo-request ct state new accept
	// --
	// iifname "wg0" icmp type echo-request ct state new accept
	fr = &model.FilterRule{
		IIF:       []string{nft.wgIface},
		Proto:     []string{"icmp"},
		ICMPType:  []byte{0x08},
		Conntrack: []string{"new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter input meta iifname "wg0" ip protocol icmp \
	// ct state { established, related } accept
	// --
	// iifname "wg0" ip protocol icmp ct state { established, related } accept
	fr = &model.FilterRule{
		IIF:       []string{nft.wgIface},
		Proto:     []string{"icmp"},
		Conntrack: []string{"established", "related"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter input meta iifname "wg0" \
	// ip protocol tcp tcp dport { 80, 8080 } ip saddr @wgmanager_ipset \
	// ct state { new, established } accept
	// --
	// iifname "wg0" tcp dport { https, 8443 } ip saddr @wgmanager_ipset ct state { established, new } accept
	managerPorts := make([][2]uint16, len(nft.managerPorts)+len(nft.cfg.TrustPorts))
	i := 0
	for _, p := range nft.managerPorts {
		managerPorts[i] = [2]uint16{p, 0}
		i++
	}
	for _, p := range nft.cfg.TrustPorts {
		managerPorts[i] = [2]uint16{p, 0}
		i++
	}

	fr = &model.FilterRule{
		IIF:      []string{nft.wgIface},
		Proto:    []string{"tcp"},
		SAddrSet: nft.filterSetWGManagerIP,
		DPort:    managerPorts,

		Conntrack: []string{"established", "new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter input meta iifname "wg0" \
	// ip protocol udp udp dport 53 \
	// ct state established accept
	// --
	// iifname "eth0" { udp, tcp } sport domain ct state established accept
	fr = &model.FilterRule{
		IIF:       []string{nft.wgIface},
		Proto:     []string{"tcp", "udp"},
		SAddr:     [][2]*net.IPNet{{nft.cfg.WGIPNet, nil}},
		DPort:     [][2]uint16{{53, 0}},
		Conntrack: []string{"established", "new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter output meta oifname "wg0" ip protocol icmp \
	// ct state { new, established } accept
	// --
	// oifname "wg0" ip protocol icmp ct state { established, new } accept
	fr = &model.FilterRule{
		OIF:   []string{nft.wgIface},
		Proto: []string{"icmp"},

		Conntrack: []string{"established", "new"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter output meta oifname "wg0" \
	// ip protocol tcp tcp sport { 80, 8080 } ip daddr @wgmanager_ipset \
	// ct state established accept
	// --
	// oifname "wg0" tcp sport { https, 8443 } ct state established accept
	managerPorts = make([][2]uint16, len(nft.managerPorts)+len(nft.cfg.TrustPorts))
	i = 0
	for _, p := range nft.managerPorts {
		managerPorts[i] = [2]uint16{p, 0}
		i++
	}
	for _, p := range nft.cfg.TrustPorts {
		managerPorts[i] = [2]uint16{p, 0}
		i++
	}

	fr = &model.FilterRule{
		OIF:      []string{nft.wgIface},
		Proto:    []string{"tcp"},
		SPort:    managerPorts,
		DAddrSet: nft.filterSetWGManagerIP,

		Conntrack: []string{"established"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter input meta iifname "wg0" \
	// ip protocol udp udp dport 53 \
	// ct state established accept
	// --
	// iifname "eth0" { udp, tcp } sport domain ct state established accept
	fr = &model.FilterRule{
		OIF:       []string{nft.wgIface},
		Proto:     []string{"tcp", "udp"},
		SPort:     [][2]uint16{{53, 0}},
		DAddr:     [][2]*net.IPNet{{nft.cfg.WGIPNet, nil}},
		Conntrack: []string{"established"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// sdnForwardManagedRules to apply.
func (nft *NFTables) sdnForwardManagedRules(c *nftables.Conn) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter forward_managed \
	// meta iifname "wg0" \
	// meta oifname "wg0" \
	// accept
	// --
	// iifname "wg0" oifname "wg0" accept;
	fr = &model.FilterRule{
		IIF:    []string{nft.wgIface},
		OIF:    []string{nft.wgIface},
		Action: "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForwardManaged,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// sdnForwardRulen to apply.
func (nft *NFTables) sdnForwardRules(c *nftables.Conn) error {
	var (
		fr    *model.FilterRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip filter forward \
	// ip protocol tcp tcp sport 25 drop
	// --
	// tcp sport smtp drop;
	fr = &model.FilterRule{
		Proto: []string{"tcp"},
		SPort: [][2]uint16{{25, 0}},

		Action: "drop",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter forward \
	// meta iifname "wg0" \
	// ip saddr @wgforward_ipset \
	// meta oifname "eth0" \
	// accept
	// --
	// iifname "wg0" oifname "eth0" accept;
	fr = &model.FilterRule{
		IIF:      []string{nft.wgIface},
		OIF:      []string{nft.wanIface},
		SAddrSet: nft.filterSetWGForwardIP,
		Action:   "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter forward \
	// ct state { established, related } accept
	// --
	// ct state { established, related } accept;
	fr = &model.FilterRule{
		IIF:       []string{nft.wanIface},
		OIF:       []string{nft.wgIface},
		DAddrSet:  nft.filterSetWGForwardIP,
		Conntrack: []string{"established", "related"},
		Action:    "accept",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter forward \
	// meta iifname "eth0" \
	// meta oifname "wg0" \
	// jump forward_managed
	// --
	// iifname "wg0" oifname "wg0" accept;
	fr = &model.FilterRule{
		IIF:             []string{nft.wanIface},
		OIF:             []string{nft.wgIface},
		Action:          "jump",
		JumpActionChain: "forward_managed",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	// cmd: nft add rule ip filter forward \
	// meta iifname "wg0" \
	// meta oifname "wg0" \
	// jump forward_managed
	// --
	// iifname "wg0" oifname "wg0" accept;
	fr = &model.FilterRule{
		IIF:             []string{nft.wgIface},
		OIF:             []string{nft.wgIface},
		Action:          "jump",
		JumpActionChain: "forward_managed",
	}
	exprs, err = filterExpressions(c, nft.tFilter, fr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// snatRules to apply.
func (nft *NFTables) snatRules(c *nftables.Conn) error {
	var (
		sr    *model.SNATRule
		rule  *nftables.Rule
		exprs []expr.Any
		err   error
	)

	// cmd: nft add rule ip nat postrouting meta oifname "eth0" \
	// snat 192.168.0.1
	// --
	// oifname "eth0" snat to 192.168.15.11
	sr = &model.SNATRule{
		OIF: []string{nft.wanIface},
		IIF: []string{nft.wgIface},
		// SAddrSet: nft.filterSetWGManagerIP,
		SNATToIP: nft.wanIP,
	}
	exprs, err = snatExpressions(c, nft.tNAT, sr)
	if err != nil {
		return err
	}
	rule = &nftables.Rule{
		Table: nft.tNAT,
		Chain: nft.cPostrouting,
		Exprs: exprs}
	rule = c.AddRule(rule)
	_ = rule

	return nil
}

// dnatRules to apply.
func (nft *NFTables) dnatRules(c *nftables.Conn) error {

	return nil
}

// UpdateTrustIPs updates filterSetTrustIP.
func (nft *NFTables) UpdateTrustIPs(del, add []net.IP) error {
	if !nft.applied {
		return nil
	}

	return nft.updateIPSet(nft.filterSetTrustIP, del, add)
}

// UpdateWGManagerIPs updates filterSetWGManagerIP.
func (nft *NFTables) UpdateWGManagerIPs(del, add []net.IP) error {
	if !nft.applied {
		return nil
	}

	return nft.updateIPSet(nft.filterSetWGManagerIP, del, add)
}

// UpdateWGForwardWanIPs updates filterSetWGForwardIP.
func (nft *NFTables) UpdateWGForwardWanIPs(del, add []net.IP) error {
	if !nft.applied {
		return nil
	}

	return nft.updateIPSet(nft.filterSetWGForwardIP, del, add)
}

func (nft *NFTables) updateIPSet(set *nftables.Set, del, add []net.IP) error {
	// bind network namespace if it was set in config
	c, err := nft.networkNamespaceBind()
	if err != nil {
		return err
	}
	// release network namespace finally
	defer nft.networkNamespaceRelease()

	if len(del) > 0 {
		elements := make([]nftables.SetElement, len(del))
		for i, v := range del {
			elements[i] = nftables.SetElement{Key: v}
		}
		err = c.SetDeleteElements(set, elements)
		if err != nil {
			return err
		}
	}

	if len(add) > 0 {
		elements := make([]nftables.SetElement, len(add))
		for i, v := range add {
			elements[i] = nftables.SetElement{Key: v}
		}
		err = c.SetAddElements(set, elements)
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

// Cleanup rules to default policy filtering.
func (nft *NFTables) Cleanup() error {
	if !nft.cfg.Enabled {
		return nil
	}
	// bind network namespace if it was set in config
	c, err := nft.networkNamespaceBind()
	if err != nil {
		return err
	}
	// release network namespace finally
	defer nft.networkNamespaceRelease()

	filterSetTrustElements, _ := c.GetSetElements(nft.filterSetTrustIP) // omit error

	c.FlushRuleset()

	// add filter table
	// cmd: nft add table ip filter
	c.AddTable(nft.tFilter)
	// add input chain of filter table
	// cmd: nft add chain ip filter input \
	// { type filter hook input priority 0 \; policy drop\; }
	c.AddChain(nft.cInput)
	// add forward chain
	// cmd: nft add chain ip filter forward \
	// { type filter hook forward priority 0 \; policy drop\; }
	c.AddChain(nft.cForward)
	// add output chain
	// cmd: nft add chain ip filter output \
	// { type filter hook output priority 0 \; policy drop\; }
	c.AddChain(nft.cOutput)

	// add trust_ipset
	// cmd: nft add set ip filter trust_ipset { type ipv4_addr\; }
	err = c.AddSet(nft.filterSetTrustIP, nil)
	if err != nil {
		return err
	}

	if filterSetTrustElements != nil {
		_ = c.SetAddElements(nft.filterSetTrustIP, filterSetTrustElements) // omit error
	}

	_ = nft.inputLocalIfaceRules(c)                // omit error
	_ = nft.outputLocalIfaceRules(c)               // omit error
	_ = nft.inputHostBaseRules(c, nft.wanIface)    // omit error
	_ = nft.outputHostBaseRules(c, nft.wanIface)   // omit error
	_ = nft.inputTrustIPSetRules(c, nft.wanIface)  // omit error
	_ = nft.outputTrustIPSetRules(c, nft.wanIface) // omit error
	for _, iface := range nft.cfg.Ifaces {
		if iface == nft.wanIface {
			continue
		}

		_ = nft.inputHostBaseRules(c, iface)    // omit error
		_ = nft.outputHostBaseRules(c, iface)   // omit error
		_ = nft.inputTrustIPSetRules(c, iface)  // omit error
		_ = nft.outputTrustIPSetRules(c, iface) // omit error
	}

	// apply configuration
	err = c.Flush()
	if err != nil {
		return err
	}
	nft.applied = false

	return nil
}

// WanIP returns ip address of wan interface.
func (nft *NFTables) WanIP() net.IP {
	return nft.wanIP
}

// IfacesIPs returns ip addresses list of additional ifaces.
func (nft *NFTables) IfacesIPs() ([]net.IP, error) {
	ips := make([]net.IP, 0, len(nft.cfg.Ifaces))

	for _, v := range nft.cfg.Ifaces {
		if v == nft.wanIface || v == nft.wgIface {
			continue
		}

		iface, err := net.InterfaceByName(v)
		if err != nil {
			return nil, err
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipnet.IP.To4()
			if ip != nil {
				ips = append(ips, ip)
			}
		}
	}

	return ips, nil
}

func filterExpressions(
	c *nftables.Conn, t *nftables.Table, r *model.FilterRule,
) ([]expr.Any, error) {
	iif, err := nfutils.IFaceExpressions(c, t, expr.MetaKeyIIFNAME, r.IIF, r.IIFNeq)
	if err != nil {
		return nil, err
	}

	oif, err := nfutils.IFaceExpressions(c, t, expr.MetaKeyOIFNAME, r.OIF, r.OIFNeq)
	if err != nil {
		return nil, err
	}

	proto, err := nfutils.L4ProtoExpressions(c, t, r.Proto, r.ProtoNeq)
	if err != nil {
		return nil, err
	}

	icmpType, err := nfutils.ICMPTypeExpressions(c, t, r.ICMPType, r.ICMPTypeNeq)
	if err != nil {
		return nil, err
	}

	isDest := false
	sport, err := nfutils.TransportExpressions(c, t, r.SPort, r.SPortNeq, isDest)
	if err != nil {
		return nil, err
	}
	saddr, err := nfutils.NetworkExpressions(c, t, r.SAddrSet, r.SAddr, r.SAddrNeq, isDest)
	if err != nil {
		return nil, err
	}

	isDest = true
	dport, err := nfutils.TransportExpressions(c, t, r.DPort, r.DPortNeq, isDest)
	if err != nil {
		return nil, err
	}
	daddr, err := nfutils.NetworkExpressions(c, t, r.DAddrSet, r.DAddr, r.DAddrNeq, isDest)
	if err != nil {
		return nil, err
	}

	conntrack, err := nfutils.ConntrackExpressions(c, t, r.Conntrack)
	if err != nil {
		return nil, err
	}

	action := filterActionExpressions(c, t, r)

	cnt := len(iif) + len(oif) +
		len(proto) + len(icmpType) + len(sport) + len(dport) +
		len(saddr) + len(daddr) +
		len(conntrack)
	if action != nil {
		cnt += 1
	}

	exprs := make([]expr.Any, cnt)

	i := 0
	for j := range iif {
		exprs[i] = iif[j]
		i++
	}

	for j := range oif {
		exprs[i] = oif[j]
		i++
	}

	for j := range proto {
		exprs[i] = proto[j]
		i++
	}

	for j := range icmpType {
		exprs[i] = icmpType[j]
		i++
	}

	for j := range sport {
		exprs[i] = sport[j]
		i++
	}

	for j := range saddr {
		exprs[i] = saddr[j]
		i++
	}

	for j := range dport {
		exprs[i] = dport[j]
		i++
	}

	for j := range daddr {
		exprs[i] = daddr[j]
		i++
	}

	for j := range conntrack {
		exprs[i] = conntrack[j]
		i++
	}

	if action != nil {
		exprs[cnt-1] = action
	}

	return exprs, nil
}

func filterActionExpressions(
	c *nftables.Conn, t *nftables.Table, r *model.FilterRule,
) expr.Any {
	switch r.Action {
	case "accept":
		return nfutils.AcceptExpression()
	case "drop":
		return nfutils.DropExpression()
	case "jump":
		if len(r.JumpActionChain) == 0 {
			return nil
		}

		return nfutils.JumpExpression(r.JumpActionChain)
	case "reject":
		return nfutils.RejectExpression(
			r.RejectActionType,
			r.RejectActionCode,
		)
	}

	return nil
}

func snatExpressions(
	c *nftables.Conn, t *nftables.Table, r *model.SNATRule,
) ([]expr.Any, error) {
	iif, err := nfutils.IFaceExpressions(c, t, expr.MetaKeyIIFNAME, r.IIF, false)
	if err != nil {
		return nil, err
	}

	oif, err := nfutils.IFaceExpressions(c, t, expr.MetaKeyOIFNAME, r.OIF, false)
	if err != nil {
		return nil, err
	}

	saddr, err := nfutils.NetworkExpressions(c, t, r.SAddrSet, nil, false, false)
	if err != nil {
		return nil, err
	}

	action := snatActionExpressions(c, t, r.SNATToIP)

	cnt := len(iif) + len(oif) + len(saddr) + len(action)

	exprs := make([]expr.Any, cnt)

	i := 0
	for j := range iif {
		exprs[i] = iif[j]
		i++
	}

	for j := range oif {
		exprs[i] = oif[j]
		i++
	}

	for j := range saddr {
		exprs[i] = saddr[j]
		i++
	}

	for j := range action {
		exprs[i] = action[j]
		i++
	}

	return exprs, nil
}

func snatActionExpressions(
	c *nftables.Conn, t *nftables.Table, ip net.IP,
) []expr.Any {
	ipReg := uint32(1)
	exprs := []expr.Any{
		&expr.Immediate{
			Register: ipReg,
			Data:     ip.To4(),
		},
		&expr.NAT{
			Type:       expr.NATTypeSourceNAT,
			Family:     unix.NFPROTO_IPV4,
			RegAddrMin: ipReg,
		},
	}

	return exprs
}

func dnatExpressions(
	c *nftables.Conn, t *nftables.Table, r *model.DNATRule,
) ([]expr.Any, error) {
	iif, err := nfutils.IFaceExpressions(c, t, expr.MetaKeyIIFNAME, r.IIF, false)
	if err != nil {
		return nil, err
	}

	oif, err := nfutils.IFaceExpressions(c, t, expr.MetaKeyOIFNAME, r.OIF, false)
	if err != nil {
		return nil, err
	}

	proto, err := nfutils.L4ProtoExpressions(c, t, r.Proto, false)
	if err != nil {
		return nil, err
	}

	saddr, err := nfutils.NetworkExpressions(c, t, r.SAddrSet, r.SAddr, false, false)
	if err != nil {
		return nil, err
	}

	isDest := false
	sport, err := nfutils.TransportExpressions(c, t, r.SPort, false, isDest)
	if err != nil {
		return nil, err
	}

	isDest = true
	dport, err := nfutils.TransportExpressions(c, t, r.DPort, false, isDest)
	if err != nil {
		return nil, err
	}

	action := dnatActionExpressions(c, t, r.DNATToIP, r.DNATToPort)

	cnt := len(iif) + len(oif) +
		len(proto) + len(saddr) + len(sport) + len(dport) + len(action)

	exprs := make([]expr.Any, cnt)

	i := 0
	for j := range iif {
		exprs[i] = iif[j]
		i++
	}

	for j := range oif {
		exprs[i] = oif[j]
		i++
	}

	for j := range proto {
		exprs[i] = proto[j]
		i++
	}

	for j := range saddr {
		exprs[i] = saddr[j]
		i++
	}

	for j := range sport {
		exprs[i] = sport[j]
		i++
	}

	for j := range dport {
		exprs[i] = dport[j]
		i++
	}

	for j := range action {
		exprs[i] = action[j]
		i++
	}

	return exprs, nil
}

func dnatActionExpressions(
	c *nftables.Conn, t *nftables.Table, ip net.IP, port uint16,
) []expr.Any {
	ipReg := uint32(1)
	portReg := uint32(2)
	exprs := []expr.Any{
		&expr.Immediate{
			Register: ipReg,
			Data:     ip.To4(),
		},
		&expr.Immediate{
			Register: portReg,
			Data:     binaryutil.BigEndian.PutUint16(port),
		},
		&expr.NAT{
			Type:        expr.NATTypeDestNAT,
			Family:      unix.NFPROTO_IPV4,
			RegAddrMin:  ipReg,
			RegProtoMin: portReg,
		},
	}

	return exprs
}
