package firewall

// Config for nftables.
type Config struct {
	Enabled          bool
	NetworkNamespace string
	DefaultPolicy    string
	WGIface          string
	WGPort           uint16
	Ifaces           []string
	TrustPorts       []uint16
}
