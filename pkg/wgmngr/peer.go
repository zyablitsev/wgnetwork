package wgmngr

import (
	"bytes"
	"net"
	"sort"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Peer object.
type Peer struct {
	ipnet net.IPNet

	publicKey    wgtypes.Key
	allowedIPs   []net.IPNet
	endpointIP   net.IP
	endpointPort uint16
	keepAlive    time.Duration
}

// NewPeer constructor.
func NewPeer(
	ipnet net.IPNet,
	publicKey wgtypes.Key,
	allowedIPs []net.IPNet,
	endpointIP net.IP,
	endpointPort uint16,
	keepAlive time.Duration,
) (Peer, error) {
	// TODO: check if allowdIPs contains ipnet already and exclude then

	p := Peer{
		ipnet: ipnet,

		publicKey:    publicKey,
		allowedIPs:   allowedIPs,
		endpointIP:   endpointIP,
		endpointPort: endpointPort,
		keepAlive:    keepAlive,
	}

	return p, nil
}

// Copy object
func (p *Peer) Copy() Peer {
	ipnet := net.IPNet{
		IP:   make(net.IP, len(p.ipnet.IP)),
		Mask: make(net.IPMask, len(p.ipnet.Mask)),
	}
	copy(ipnet.IP, p.ipnet.IP)
	copy(ipnet.Mask, p.ipnet.Mask)

	publicKey := p.publicKey

	allowedIPs := make([]net.IPNet, len(p.allowedIPs))
	for i := 0; i < len(p.allowedIPs); i++ {
		allowedIPs[i] = net.IPNet{
			IP:   make(net.IP, len(p.allowedIPs[i].IP)),
			Mask: make(net.IPMask, len(p.allowedIPs[i].Mask)),
		}
		copy(allowedIPs[i].IP, p.allowedIPs[i].IP)
		copy(allowedIPs[i].Mask, p.allowedIPs[i].Mask)
	}

	endpointIP := make(net.IP, len(p.endpointIP))
	copy(endpointIP, p.endpointIP)

	cp := Peer{
		ipnet: ipnet,

		publicKey:    publicKey,
		allowedIPs:   allowedIPs,
		endpointIP:   endpointIP,
		endpointPort: p.endpointPort,
		keepAlive:    p.keepAlive,
	}

	return cp
}

// AllowedIPs returns array of CIDR and Routes
func (p *Peer) AllowedIPs() []net.IPNet {
	allowedIPs := make([]net.IPNet, len(p.allowedIPs)+1)

	allowedIPs[0] = p.ipnet

	for i := 0; i < len(p.allowedIPs); i++ {
		j := i + 1
		allowedIPs[j] = net.IPNet{
			IP:   make(net.IP, len(p.allowedIPs[i].IP)),
			Mask: make(net.IPMask, len(p.allowedIPs[i].Mask)),
		}
		copy(allowedIPs[j].IP, p.allowedIPs[i].IP)
		copy(allowedIPs[j].Mask, p.allowedIPs[i].Mask)
	}

	sort.Slice(allowedIPs, func(i, j int) bool {
		return (bytes.Compare(allowedIPs[i].IP, allowedIPs[j].IP) < 0 &&
			bytes.Compare(allowedIPs[i].Mask, allowedIPs[j].Mask) < 0)
	})

	return allowedIPs
}

// PublicKey value
func (p *Peer) PublicKey() wgtypes.Key {
	return p.publicKey
}

// Endpoint address.
func (p *Peer) Endpoint() *net.UDPAddr {
	if p.endpointIP == nil || p.endpointPort == 0 {
		return nil
	}

	return &net.UDPAddr{IP: p.endpointIP, Port: int(p.endpointPort)}
}

// KeepAlive interval value
func (p *Peer) KeepAlive() time.Duration {
	return p.keepAlive
}

// Peers object
type Peers []Peer

// Copy slice of peers
func (p Peers) Copy() Peers {
	cp := make(Peers, len(p))
	for i := 0; i < len(p); i++ {
		cp[i] = p[i].Copy()
	}

	return cp
}

// Len to satisfy Sort interface type
func (p Peers) Len() int {
	return len(p)
}

// Less to satisfy Sort interface type
func (p Peers) Less(i, j int) bool {
	return bytes.Compare(p[i].publicKey[:], p[j].publicKey[:]) < 0
}

// Swap to satisfy Sort interface type
func (p Peers) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// IPNets of peers.
func (p Peers) IPNets() []net.IPNet {
	ipnets := make([]net.IPNet, len(p))

	for i := 0; i < len(p); i++ {
		ipnets[i] = net.IPNet{
			IP:   make(net.IP, len(p[i].ipnet.IP)),
			Mask: make(net.IPMask, len(p[i].ipnet.Mask)),
		}
		copy(ipnets[i].IP, p[i].ipnet.IP)
		copy(ipnets[i].Mask, p[i].ipnet.Mask)
	}

	sort.Slice(ipnets, func(i, j int) bool {
		return (bytes.Compare(ipnets[i].IP, ipnets[j].IP) < 0 &&
			bytes.Compare(ipnets[i].Mask, ipnets[j].Mask) < 0)
	})

	return ipnets
}

// EndpointIPs of peers.
func (p Peers) EndpointIPs() []net.IP {
	endpointips := make([]net.IP, len(p))

	for i := 0; i < len(p); i++ {
		endpointips[i] = make(net.IP, len(p[i].ipnet.IP))
		copy(endpointips[i], p[i].endpointIP)
	}

	sort.Slice(endpointips, func(i, j int) bool {
		return bytes.Compare(endpointips[i], endpointips[j]) < 0
	})

	return endpointips
}
