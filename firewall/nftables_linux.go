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

	tFilter  *nftables.Table
	cInput   *nftables.Chain
	cForward *nftables.Chain
	cOutput  *nftables.Chain

	tNAT         *nftables.Table
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

		tFilter:  tFilter,
		cInput:   cInput,
		cForward: cForward,
		cOutput:  cOutput,

		tNAT:         tNAT,
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

	nft.inputLocalIfaceRules(c)
	nft.outputLocalIfaceRules(c)
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
	err = nft.sdnForwardRules(c)
	if err != nil {
		return err
	}
	nft.natRules(c)

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
func (nft *NFTables) inputLocalIfaceRules(c *nftables.Conn) {
	// cmd: nft add rule ip filter input meta iifname "lo" accept
	// --
	// iifname "lo" accept
	exprs := make([]expr.Any, 0, 3)
	exprs = append(exprs, nfutils.SetIIF(loIface)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname != "lo" \
	// ip saddr 127.0.0.0/8 reject
	// --
	// iifname != "lo" ip saddr 127.0.0.0/8 reject with icmp type prot-unreachable
	exprs = make([]expr.Any, 0, 6)
	exprs = append(exprs, nfutils.SetNIIF(loIface)...)
	exprs = append(exprs,
		nfutils.SetSourceNet([]byte{127, 0, 0, 0}, []byte{255, 255, 255, 0})...)
	exprs = append(exprs, nfutils.ExprReject(
		unix.NFT_REJECT_ICMP_UNREACH,
		unix.NFT_REJECT_ICMPX_UNREACH,
	))
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)
}

// outputLocalIfaceRules to apply.
func (nft *NFTables) outputLocalIfaceRules(c *nftables.Conn) {
	// cmd: nft add rule ip filter output meta oifname "lo" accept
	// --
	// oifname "lo" accept
	exprs := make([]expr.Any, 0, 3)
	exprs = append(exprs, nfutils.SetOIF(loIface)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)
}

// inputHostBaseRules to apply.
func (nft *NFTables) inputHostBaseRules(c *nftables.Conn, iface string) error {
	// cmd: nft add rule ip filter input meta iifname "eth0" ip protocol icmp \
	// ct state { established, related } accept
	// --
	// iifname "eth0" ip protocol icmp ct state { established, related } accept
	ctStateSet := nfutils.GetConntrackStateSet(nft.tFilter)
	elems := nfutils.GetConntrackStateSetElems(
		[]string{"established", "related"})
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs := make([]expr.Any, 0, 7)
	exprs = append(exprs, nfutils.SetIIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoICMP()...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol udp udp sport 53 \
	// ct state established accept
	// --
	// iifname "eth0" udp sport domain ct state established accept
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetIIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoUDP()...)
	exprs = append(exprs, nfutils.SetSPort(53)...)
	exprs = append(exprs, nfutils.SetConntrackStateEstablished()...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol tcp tcp sport 53 \
	// ct state established accept
	// --
	// iifname "eth0" tcp sport domain ct state established accept
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetIIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetSPort(53)...)
	exprs = append(exprs, nfutils.SetConntrackStateEstablished()...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol tcp tcp sport { 80, 443 } \
	// ct state established accept
	// --
	// iifname "eth0" tcp sport { http, https } ct state established accept
	portSet := nfutils.GetPortSet(nft.tFilter)
	// portSet := &nftables.Set{Anonymous: true, Constant: true,
	// 	Table: nft.tFilter, KeyType: nftables.TypeInetService}
	elems = nfutils.GetPortElems([]uint16{80, 443})
	err = c.AddSet(portSet, elems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetIIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetSPortSet(portSet)...)
	exprs = append(exprs, nfutils.SetConntrackStateEstablished()...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	return nil
}

// outputHostBaseRules to apply.
func (nft *NFTables) outputHostBaseRules(c *nftables.Conn, iface string) error {
	// cmd: nft add rule ip filter output meta oifname "eth0" ip protocol icmp \
	// ct state { new, established } accept
	// --
	// oifname "eth0" ip protocol icmp ct state { established, new } accept
	ctStateSet := nfutils.GetConntrackStateSet(nft.tFilter)
	elems := nfutils.GetConntrackStateSetElems(
		[]string{"new", "established"})
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs := make([]expr.Any, 0, 7)
	exprs = append(exprs, nfutils.SetOIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoICMP()...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol udp udp dport 53 \
	// ct state { new, established } accept
	// --
	// oifname "eth0" udp dport domain ct state { established, new } accept
	ctStateSet = nfutils.GetConntrackStateSet(nft.tFilter)
	elems = nfutils.GetConntrackStateSetElems(
		[]string{"new", "established"})
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetOIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoUDP()...)
	exprs = append(exprs, nfutils.SetDPort(53)...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol tcp tcp dport 53 \
	// ct state { new, established } accept
	// --
	// oifname "eth0" tcp dport domain ct state { established, new } accept
	ctStateSet = nfutils.GetConntrackStateSet(nft.tFilter)
	elems = nfutils.GetConntrackStateSetElems(
		[]string{"new", "established"})
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetOIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetDPort(53)...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol tcp tcp dport { 80, 443 } \
	// ct state { new, established } accept
	// --
	// oifname "eth0" tcp dport { http, https } ct state { established, new } accept
	portSet := nfutils.GetPortSet(nft.tFilter)
	elems = nfutils.GetPortElems([]uint16{80, 443})
	err = c.AddSet(portSet, elems)
	if err != nil {
		return err
	}

	ctStateSet = nfutils.GetConntrackStateSet(nft.tFilter)
	elems = nfutils.GetConntrackStateSetElems(
		[]string{"new", "established"})
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetOIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetDPortSet(portSet)...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	return nil
}

// inputTrustIPSetRules to apply.
func (nft *NFTables) inputTrustIPSetRules(c *nftables.Conn, iface string) error {
	// cmd: nft add rule ip filter input meta iifname "eth0" ip protocol icmp \
	// icmp type echo-request ip saddr @trust_ipset ct state new accept
	// --
	// iifname "eth0" icmp type echo-request ip saddr @trust_ipset ct state new accept
	exprs := make([]expr.Any, 0, 12)
	exprs = append(exprs, nfutils.SetIIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoICMP()...)
	exprs = append(exprs, nfutils.SetICMPTypeEchoRequest()...)
	exprs = append(exprs, nfutils.SetSAddrSet(nft.filterSetTrustIP)...)
	exprs = append(exprs, nfutils.SetConntrackStateNew()...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol tcp tcp dport { 5522 } ip saddr @trust_ipset \
	// ct state { new, established } accept
	// --
	// iifname "eth0" tcp dport { 5522 } ip saddr @trust_ipset ct state { established, new } accept
	ctStateSet := nfutils.GetConntrackStateSet(nft.tFilter)
	elems := nfutils.GetConntrackStateSetElems(
		[]string{"new", "established"})
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	portSet := nfutils.GetPortSet(nft.tFilter)
	err = c.AddSet(portSet, nft.cfg.trustPorts())
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 11)
	exprs = append(exprs, nfutils.SetIIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetDPortSet(portSet)...)
	exprs = append(exprs, nfutils.SetSAddrSet(nft.filterSetTrustIP)...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	return nil
}

// outputTrustIPSetRules to apply.
func (nft *NFTables) outputTrustIPSetRules(c *nftables.Conn, iface string) error {
	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol tcp tcp sport { 5522 } ip daddr @trust_ipset \
	// ct state established accept
	// --
	// oifname "eth0" tcp sport { 5522 } ip daddr @trust_ipset ct state established accept
	portSet := nfutils.GetPortSet(nft.tFilter)
	err := c.AddSet(portSet, nft.cfg.trustPorts())
	if err != nil {
		return err
	}

	exprs := make([]expr.Any, 0, 12)
	exprs = append(exprs, nfutils.SetOIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetSPortSet(portSet)...)
	exprs = append(exprs, nfutils.SetDAddrSet(nft.filterSetTrustIP)...)
	exprs = append(exprs, nfutils.SetConntrackStateEstablished()...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	return nil
}

// inputPublicRules to apply.
func (nft *NFTables) inputPublicRules(c *nftables.Conn, iface string) error {
	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol udp udp dport 51820 accept
	// --
	// iifname "eth0" udp dport 51820 accept

	exprs := make([]expr.Any, 0, 9)
	exprs = append(exprs, nfutils.SetIIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoUDP()...)
	exprs = append(exprs, nfutils.SetDPort(nft.wgPort)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	return nil
}

// outputPublicRules to apply.
func (nft *NFTables) outputPublicRules(c *nftables.Conn, iface string) error {
	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol udp udp sport 51820 accept
	// --
	// oifname "eth0" udp sport 51820 accept

	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetOIF(iface)...)
	exprs = append(exprs, nfutils.SetProtoUDP()...)
	exprs = append(exprs, nfutils.SetSPort(nft.wgPort)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	return nil
}

// sdnRules to apply.
func (nft *NFTables) sdnRules(c *nftables.Conn) error {
	// cmd: nft add rule ip filter input meta iifname "wg0" ip protocol icmp \
	// icmp type echo-request ct state new accept
	// --
	// iifname "wg0" icmp type echo-request ct state new accept
	exprs := make([]expr.Any, 0, 12)
	exprs = append(exprs, nfutils.SetIIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetProtoICMP()...)
	exprs = append(exprs, nfutils.SetICMPTypeEchoRequest()...)
	exprs = append(exprs, nfutils.SetConntrackStateNew()...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "wg0" ip protocol icmp \
	// ct state { established, related } accept
	// --
	// iifname "wg0" ip protocol icmp ct state { established, related } accept
	ctStateSet := nfutils.GetConntrackStateSet(nft.tFilter)
	elems := nfutils.GetConntrackStateSetElems(
		[]string{"established", "related"})
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 7)
	exprs = append(exprs, nfutils.SetIIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetProtoICMP()...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "wg0" \
	// ip protocol tcp tcp dport { 80, 8080 } ip saddr @wgmanager_ipset \
	// ct state { new, established } accept
	// --
	// iifname "wg0" tcp dport { https, 8443 } ip saddr @wgmanager_ipset ct state { established, new } accept
	ctStateSet = nfutils.GetConntrackStateSet(nft.tFilter)
	elems = nfutils.GetConntrackStateSetElems(
		[]string{"new", "established"})
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	portSet := nfutils.GetPortSet(nft.tFilter)
	portSetElems := make([]nftables.SetElement, len(nft.managerPorts))
	for i, p := range nft.managerPorts {
		portSetElems[i] = nftables.SetElement{
			Key: binaryutil.BigEndian.PutUint16(p)}
	}
	err = c.AddSet(portSet, portSetElems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 9)
	exprs = append(exprs, nfutils.SetIIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetDPortSet(portSet)...)
	exprs = append(exprs, nfutils.SetSAddrSet(nft.filterSetWGManagerIP)...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	ctStateSet = nfutils.GetConntrackStateSet(nft.tFilter)
	elems = nfutils.GetConntrackStateSetElems(
		[]string{"new", "established"})
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetIIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetProtoUDP()...)
	exprs = append(exprs, nfutils.SetDPort(53)...)
	// exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	ctStateSet = nfutils.GetConntrackStateSet(nft.tFilter)
	elems = nfutils.GetConntrackStateSetElems(
		[]string{"new", "established"})
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetIIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetDPort(53)...)
	// exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter output meta oifname "wg0" ip protocol icmp \
	// ct state { new, established } accept
	// --
	// oifname "wg0" ip protocol icmp ct state { established, new } accept
	ctStateSet = nfutils.GetConntrackStateSet(nft.tFilter)
	elems = nfutils.GetConntrackStateSetElems(
		[]string{"new", "established"})
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 7)
	exprs = append(exprs, nfutils.SetOIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetProtoICMP()...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter output meta oifname "wg0" \
	// ip protocol tcp tcp sport { 80, 8080 } ip daddr @wgmanager_ipset \
	// ct state established accept
	// --
	// oifname "wg0" tcp sport { https, 8443 } ct state established accept
	portSet = nfutils.GetPortSet(nft.tFilter)
	portSetElems = make([]nftables.SetElement, len(nft.managerPorts))
	for i, p := range nft.managerPorts {
		portSetElems[i] = nftables.SetElement{
			Key: binaryutil.BigEndian.PutUint16(p)}
	}
	err = c.AddSet(portSet, portSetElems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetOIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetSPortSet(portSet)...)
	exprs = append(exprs, nfutils.SetDAddrSet(nft.filterSetWGManagerIP)...)
	exprs = append(exprs, nfutils.SetConntrackStateEstablished()...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetOIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetProtoUDP()...)
	exprs = append(exprs, nfutils.SetSPort(53)...)
	// exprs = append(exprs, nfutils.SetConntrackStateEstablished()...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetOIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetSPort(53)...)
	// exprs = append(exprs, nfutils.SetConntrackStateEstablished()...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs}
	c.AddRule(rule)

	return nil
}

// sdnForwardRules to apply.
func (nft *NFTables) sdnForwardRules(c *nftables.Conn) error {
	// cmd: nft add rule ip filter forward \
	// ip protocol tcp tcp sport 25 drop
	// --
	// tcp sport smtp drop;
	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetProtoTCP()...)
	exprs = append(exprs, nfutils.SetSPort(25)...)
	exprs = append(exprs, nfutils.ExprDrop())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter forward \
	// meta iifname "wg0" \
	// ip saddr @wgforward_ipset \
	// meta oifname "eth0" \
	// accept
	// --
	// iifname "wg0" oifname "eth0" accept;
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetIIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetSAddrSet(nft.filterSetWGForwardIP)...)
	exprs = append(exprs, nfutils.SetOIF(nft.wanIface)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter forward \
	// ct state { established, related } accept
	// --
	// ct state { established, related } accept;
	ctStateSet := nfutils.GetConntrackStateSet(nft.tFilter)
	elems := nfutils.GetConntrackStateSetElems(
		[]string{"established", "related"})
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetIIF(nft.wanIface)...)
	exprs = append(exprs, nfutils.SetDAddrSet(nft.filterSetWGForwardIP)...)
	exprs = append(exprs, nfutils.SetOIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs}
	c.AddRule(rule)

	// cmd: nft add rule ip filter forward \
	// meta iifname "wg0" \
	// meta oifname "wg0" \
	// accept
	// --
	// iifname "wg0" oifname "wg0" accept;
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetIIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.SetOIF(nft.wgIface)...)
	exprs = append(exprs, nfutils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs}
	c.AddRule(rule)

	return nil
}

// natRules to apply.
func (nft *NFTables) natRules(c *nftables.Conn) {
	// cmd: nft add rule ip nat postrouting meta oifname "eth0" \
	// snat 192.168.0.1
	// --
	// oifname "eth0" snat to 192.168.15.11
	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, nfutils.SetOIF(nft.wanIface)...)
	exprs = append(exprs, nfutils.ExprImmediate(nft.wanIP))
	exprs = append(exprs, nfutils.ExprSNAT(1, 0))
	rule := &nftables.Rule{
		Table: nft.tNAT,
		Chain: nft.cPostrouting,
		Exprs: exprs}
	c.AddRule(rule)
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

	nft.inputLocalIfaceRules(c)
	nft.outputLocalIfaceRules(c)
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
