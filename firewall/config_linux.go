//go:build linux
// +build linux

package firewall

import (
	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
)

func (c *Config) trustPorts() []nftables.SetElement {
	elems := make([]nftables.SetElement, len(c.TrustPorts))
	for i, p := range c.TrustPorts {
		elems[i] = nftables.SetElement{Key: binaryutil.BigEndian.PutUint16(p)}
	}

	return elems
}
