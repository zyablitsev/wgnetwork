package main

import (
	"flag"
	"fmt"
	"os"

	"wgnetwork/pkg/cli"
	"wgnetwork/pkg/logger"
)

func main() {
	log, err := logger.New(os.Stdout, os.Stderr)
	if err != nil {
		panic(fmt.Errorf("can't init logger: %v", err))
	}

	actionWgCfg := cli.NewActionWgCfg(log)
	actionUserCreate := cli.NewActionUserCreate(log)
	actionUserEdit := cli.NewActionUserEdit(log)
	actionUserRemove := cli.NewActionUserRemove(log)
	actionUser := cli.NewActionUser(log)
	actionUsers := cli.NewActionUsers(log)
	actionDeviceCreate := cli.NewActionDeviceCreate(log)
	actionDeviceEdit := cli.NewActionDeviceEdit(log)
	actionDeviceRemove := cli.NewActionDeviceRemove(log)
	actionDevice := cli.NewActionDevice(log)
	actionDevices := cli.NewActionDevices(log)
	actionTrustIPSetAdd := cli.NewActionTrustIPSetAdd(log)
	actionTrustIPSetRemove := cli.NewActionTrustIPSetRemove(log)
	actionTrustIPSet := cli.NewActionTrustIPSet(log)

	// parse command-line argiments
	flag.Parse()

	if flag.NArg() == 0 {
		actionWgCfg.Usage()
		actionUserCreate.Usage()
		actionUserEdit.Usage()
		actionUserRemove.Usage()
		actionUser.Usage()
		actionUsers.Usage()
		actionDeviceCreate.Usage()
		actionDeviceEdit.Usage()
		actionDeviceRemove.Usage()
		actionDevice.Usage()
		actionDevices.Usage()
		actionTrustIPSetAdd.Usage()
		actionTrustIPSetRemove.Usage()
		actionTrustIPSet.Usage()

		return
	}

	args := flag.Args()
	log.Debugf("found args in the cmd - parsing and running action: %+v",
		args)

	var action executor

	switch args[0] {
	case "wgcfg":
		action = actionWgCfg
	case "user-create":
		action = actionUserCreate
	case "user-edit":
		action = actionUserEdit
	case "user-remove":
		action = actionUserRemove
	case "user":
		action = actionUser
	case "users":
		action = actionUsers
	case "device-create":
		action = actionDeviceCreate
	case "device-edit":
		action = actionDeviceEdit
	case "device-remove":
		action = actionDeviceRemove
	case "device":
		action = actionDevice
	case "devices":
		action = actionDevices
	case "trust-ipset-add":
		action = actionTrustIPSetAdd
	case "trust-ipset-remove":
		action = actionTrustIPSetRemove
	case "trust-ipset":
		action = actionTrustIPSet
	default:
		log.Errorf("unknown action")
		os.Exit(1)
	}

	err = action.Execute(args[1:])
	if err != nil {
		log.Errorf("failed to execute: %v", err)
		os.Exit(1)
	}
}

// executor describes interface of action object.
type executor interface {
	Usage()
	Execute(args []string) error
}
