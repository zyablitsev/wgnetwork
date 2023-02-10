// +build !linux

package iface

import "net"

// Create network link for interface mock.
func Create(
	log logger,
	iface, linkType string,
	ip net.IP, ipNet *net.IPNet,
) error {
	return nil
}

// Remove network link for interface mock.
func Remove(log logger, iface string) error {
	return nil
}

// logger desribes interface of log object.
type logger interface {
	Debugf(string, ...interface{})
}
