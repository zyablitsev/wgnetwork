// +build !linux

package wgmngr

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

func configureDevice(iface string, cfg wgtypes.Config) error {
	return nil
}
