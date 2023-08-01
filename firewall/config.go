package firewall

import "net"

// Config for nftables.
type Config struct {
	Enabled          bool
	NetworkNamespace string
	DefaultPolicy    string
	WGIface          string
	WGPort           uint16
	WGIPNet          *net.IPNet
	Ifaces           []string
	TrustPorts       []uint16
}
