//go:build !linux
// +build !linux

package firewall

import (
	"net"
)

// NFTables mock.
type NFTables struct {
	wanIP net.IP
}

// Init mock.
func Init(Config, []uint16) (*NFTables, error) {
	nft := &NFTables{
		wanIP: net.IPv4(127, 0, 0, 1).To4(),
	}

	return nft, nil
}

// UpdateTrustIPs mock method.
func (nft *NFTables) UpdateTrustIPs(_, _ []net.IP) error {
	return nil
}

// UpdateWGManagerIPs mock method.
func (nft *NFTables) UpdateWGManagerIPs(_, _ []net.IP) error {
	return nil
}

// UpdateWGForwardWanIPs mock method.
func (nft *NFTables) UpdateWGForwardWanIPs(_, _ []net.IP) error {
	return nil
}

// Cleanup mock method.
func (nft *NFTables) Cleanup() error {
	return nil
}

// WanIP returns ip address of wan interface.
func (nft *NFTables) WanIP() net.IP {
	return nft.wanIP
}

// IfacesIPs returns ip addresses list of additional ifaces.
func (nft *NFTables) IfacesIPs() ([]net.IP, error) {
	return nil, nil
}
