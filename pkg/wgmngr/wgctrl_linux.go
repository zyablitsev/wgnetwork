//go:build linux
// +build linux

package wgmngr

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func configureDevice(iface string, cfg wgtypes.Config) error {
	ctrl, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer ctrl.Close()

	return ctrl.ConfigureDevice(iface, cfg)
}
