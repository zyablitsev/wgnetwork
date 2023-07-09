package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"wgnetwork/api/manager"
	"wgnetwork/pkg/pretty"
	"wgnetwork/pkg/rpcapi"
)

// ActionDeviceEdit object.
type ActionDeviceEdit struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	ip         *string
	label      *string
	wanForward *bool
	wgPubKey   *string
}

// NewActionDeviceEdit constructor.
func NewActionDeviceEdit(log logger) *ActionDeviceEdit {
	flagset := flag.NewFlagSet(
		"device-edit",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	ip := flagset.String(
		"ip",
		"",
		"device ip")
	label := flagset.String(
		"label",
		"",
		"label")
	wanForward := flagset.Bool(
		"wan_forward",
		false,
		"wan_forward")
	wgPubKey := flagset.String(
		"wg_pubkey",
		"",
		"wireguard public key")

	a := &ActionDeviceEdit{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		ip:         ip,
		label:      label,
		wanForward: wanForward,
		wgPubKey:   wgPubKey,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionDeviceEdit) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionDeviceEdit) Execute(args []string) error {
	logPrefix := "[device-edit] Execute"

	a.log.Debugf("%s: trying to parse args: %v…", logPrefix, args)
	err := a.flagset.Parse(args)
	if err != nil {
		return errors.New("can't parse args")
	}

	// validate arguments
	err = a.validate()
	if err != nil {
		return err
	}

	client := newHTTPClient(*a.unixSocket)

	request := manager.DeviceEditRequest{IP: *a.ip}
	if a.label != nil && len(*a.label) > 0 {
		request.Label = a.label
	}
	if a.wanForward != nil {
		request.WANForward = a.wanForward
	}
	if a.wgPubKey != nil && len(*a.wgPubKey) > 0 {
		request.WGPublicKey = a.wgPubKey
	}
	b := request.Marshal()
	b = rpcapi.Request{
		Method: "manager/device/edit",
		Params: b,
	}.Marshal()
	br := bytes.NewBuffer(b)

	resp, err := client.Post(
		"http://localhost/rpc",
		"application/json; charset=utf-8",
		br)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("bad response status code: %d", resp.StatusCode)
		return err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	response := &rpcapi.Response{}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return err
	}

	result := manager.DeviceResponse{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return err
	}

	table := pretty.NewTable(5)
	table.SetHeader([]string{"label", "wan forward", "allowed ips", "user name", "user uuid"})

	table.AddRow([]string{
		result.Label,
		strconv.FormatBool(result.WANForward),
		strings.Join(result.WgDeviceAllowedIPs, "\n"),
		result.UserName,
		result.UserUUID})

	os.Stdout.WriteString(table.Render())

	privKey := "<PLACEHOLDER>"
	addr := fmt.Sprintf("%s:%d", result.WanIP, result.WgPort)

	ipnets := make([]*net.IPNet, len(result.WgDeviceAllowedIPs))
	for i := 0; i < len(result.WgDeviceAllowedIPs); i++ {
		_, ipnet, err := net.ParseCIDR(result.WgDeviceAllowedIPs[i])
		if err != nil {
			return err
		}
		ipnets[i] = ipnet
	}
	_, wgipnet, err := net.ParseCIDR(result.WgIPNet)
	if err != nil {
		return err
	}
	ipnets = allowedIPs(wgipnet, ipnets)

	allowedIPs := make([]string, len(ipnets))
	for i := range ipnets {
		allowedIPs[i] = ipnets[i].String()
	}

	cfg := buildWgCfg(
		privKey,
		result.WgPubKey,
		result.WgDeviceInet,
		result.WgIP,
		strings.Join(allowedIPs, ", "),
		addr)
	os.Stdout.WriteString("\ntunnel config:\n")
	os.Stdout.WriteString(cfg)

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionDeviceEdit) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	if a.ip == nil {
		return errors.New("ip required")
	}

	if net.ParseIP(*a.ip).To4() == nil {
		return errors.New("bad ip value")
	}

	return nil
}
