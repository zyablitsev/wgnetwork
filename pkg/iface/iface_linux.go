//go:build linux
// +build linux

package iface

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// Create network link for interface.
func Create(
	log logger,
	iface, linkType string,
	ip net.IP, ipNet *net.IPNet,
) error {
	log.Debugf("%q creating…", iface)

	_, err := net.InterfaceByName(iface)
	if err == nil {
		log.Debugf("%q already exists", iface)
		// we should remove it first
		err = Remove(log, iface)
		if err != nil {
			return err
		}
	}

	la := netlink.NewLinkAttrs()
	la.Name = iface
	link := &netlink.GenericLink{LinkAttrs: la, LinkType: linkType}
	err = netlink.LinkAdd(link)
	if err != nil {
		return fmt.Errorf("%q can't add link: %s", iface, err)
	}
	log.Debugf("%q link added", iface)

	addr := &netlink.Addr{IPNet: &net.IPNet{IP: ip, Mask: ipNet.Mask}}
	err = netlink.AddrAdd(link, addr)
	if err != nil {
		return fmt.Errorf("%q can't add addr: %v", iface, err)
	}
	log.Debugf("%q ip %q, net %q was set", iface, ip, ipNet)

	err = netlink.LinkSetUp(link)
	if err != nil {
		return fmt.Errorf("%s can't link set up: %s", iface, err)
	}
	log.Debugf("%q link is up", iface)

	return nil
}

// Remove network link for interface.
func Remove(log logger, iface string) error {
	log.Debugf("%q removing…", iface)

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("%q can't find: %v", iface, err)
	}

	err = netlink.LinkSetDown(link)
	if err != nil {
		return fmt.Errorf("%s can't link set down: %s", iface, err)
	}
	log.Debugf("%q link is down", iface)

	err = netlink.LinkDel(link)
	if err != nil {
		return fmt.Errorf("%q can't del link: %s", iface, err)
	}
	log.Debugf("%q link removed", iface)

	return nil
}

// logger desribes interface of log object.
type logger interface {
	Debugf(string, ...interface{})
}
