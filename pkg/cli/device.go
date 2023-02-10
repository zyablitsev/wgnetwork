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

// ActionDevice object.
type ActionDevice struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	ip         *string
}

// NewActionDevice constructor.
func NewActionDevice(log logger) *ActionDevice {
	flagset := flag.NewFlagSet(
		"device",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	ip := flagset.String(
		"ip",
		"",
		"device ip")

	a := &ActionDevice{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		ip:         ip,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionDevice) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionDevice) Execute(args []string) error {
	logPrefix := "[device] Execute"

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

	b := manager.DeviceRequest{
		IP: *a.ip,
	}.Marshal()
	b = rpcapi.Request{
		Method: "manager/device",
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
		strings.Join(allowedIPs, ", "),
		addr)
	os.Stdout.WriteString("\ntunnel config:\n")
	os.Stdout.WriteString(cfg)

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionDevice) validate() error {
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

func allowedIPs(wgipnet *net.IPNet, ipnets []*net.IPNet) []*net.IPNet {
	if len(ipnets) != 1 {
		return ipnets
	}

	ones, _ := ipnets[0].Mask.Size()
	if ones > 0 {
		return ipnets
	}

	ipnets = []*net.IPNet{
		wgipnet,
		// 1.0.0.0/8
		{IP: net.IPv4(1, 0, 0, 0).To4(), Mask: net.IPv4Mask(255, 0, 0, 0)},
		// 2.0.0.0/8
		{IP: net.IPv4(2, 0, 0, 0).To4(), Mask: net.IPv4Mask(255, 0, 0, 0)},
		// 3.0.0.0/8
		{IP: net.IPv4(3, 0, 0, 0).To4(), Mask: net.IPv4Mask(255, 0, 0, 0)},
		// 4.0.0.0/6
		{IP: net.IPv4(4, 0, 0, 0).To4(), Mask: net.IPv4Mask(252, 0, 0, 0)},
		// 8.0.0.0/7
		{IP: net.IPv4(8, 0, 0, 0).To4(), Mask: net.IPv4Mask(254, 0, 0, 0)},
		// 11.0.0.0/8
		{IP: net.IPv4(11, 0, 0, 0).To4(), Mask: net.IPv4Mask(255, 0, 0, 0)},
		// 12.0.0.0/6
		{IP: net.IPv4(12, 0, 0, 0).To4(), Mask: net.IPv4Mask(252, 0, 0, 0)},
		// 16.0.0.0/4
		{IP: net.IPv4(16, 0, 0, 0).To4(), Mask: net.IPv4Mask(240, 0, 0, 0)},
		// 32.0.0.0/3
		{IP: net.IPv4(32, 0, 0, 0).To4(), Mask: net.IPv4Mask(224, 0, 0, 0)},
		// 64.0.0.0/2
		{IP: net.IPv4(64, 0, 0, 0).To4(), Mask: net.IPv4Mask(192, 0, 0, 0)},
		// 128.0.0.0/3
		{IP: net.IPv4(128, 0, 0, 0).To4(), Mask: net.IPv4Mask(224, 0, 0, 0)},
		// 160.0.0.0/5
		{IP: net.IPv4(160, 0, 0, 0).To4(), Mask: net.IPv4Mask(248, 0, 0, 0)},
		// 168.0.0.0/6
		{IP: net.IPv4(168, 0, 0, 0).To4(), Mask: net.IPv4Mask(252, 0, 0, 0)},
		// 172.0.0.0/12
		{IP: net.IPv4(172, 0, 0, 0).To4(), Mask: net.IPv4Mask(255, 240, 0, 0)},
		// 172.32.0.0/11
		{IP: net.IPv4(172, 32, 0, 0).To4(), Mask: net.IPv4Mask(255, 224, 0, 0)},
		// 172.64.0.0/10
		{IP: net.IPv4(172, 64, 0, 0).To4(), Mask: net.IPv4Mask(255, 192, 0, 0)},
		// 172.128.0.0/9
		{IP: net.IPv4(172, 128, 0, 0).To4(), Mask: net.IPv4Mask(255, 128, 0, 0)},
		// 173.0.0.0/8
		{IP: net.IPv4(173, 0, 0, 0).To4(), Mask: net.IPv4Mask(255, 0, 0, 0)},
		// 174.0.0.0/7
		{IP: net.IPv4(174, 0, 0, 0).To4(), Mask: net.IPv4Mask(254, 0, 0, 0)},
		// 176.0.0.0/4
		{IP: net.IPv4(176, 0, 0, 0).To4(), Mask: net.IPv4Mask(240, 0, 0, 0)},
		// 192.0.0.0/9
		{IP: net.IPv4(192, 0, 0, 0).To4(), Mask: net.IPv4Mask(255, 128, 0, 0)},
		// 192.128.0.0/11
		{IP: net.IPv4(192, 128, 0, 0).To4(), Mask: net.IPv4Mask(255, 224, 0, 0)},
		// 192.160.0.0/13
		{IP: net.IPv4(192, 160, 0, 0).To4(), Mask: net.IPv4Mask(255, 248, 0, 0)},
		// 192.169.0.0/16
		{IP: net.IPv4(192, 169, 0, 0).To4(), Mask: net.IPv4Mask(255, 255, 0, 0)},
		// 192.170.0.0/15
		{IP: net.IPv4(192, 170, 0, 0).To4(), Mask: net.IPv4Mask(255, 254, 0, 0)},
		// 192.172.0.0/14
		{IP: net.IPv4(192, 172, 0, 0).To4(), Mask: net.IPv4Mask(255, 252, 0, 0)},
		// 192.176.0.0/12
		{IP: net.IPv4(192, 176, 0, 0).To4(), Mask: net.IPv4Mask(255, 240, 0, 0)},
		// 192.192.0.0/10
		{IP: net.IPv4(192, 192, 0, 0).To4(), Mask: net.IPv4Mask(255, 192, 0, 0)},
		// 193.0.0.0/8
		{IP: net.IPv4(193, 0, 0, 0).To4(), Mask: net.IPv4Mask(255, 0, 0, 0)},
		// 194.0.0.0/7
		{IP: net.IPv4(194, 0, 0, 0).To4(), Mask: net.IPv4Mask(254, 0, 0, 0)},
		// 196.0.0.0/6
		{IP: net.IPv4(196, 0, 0, 0).To4(), Mask: net.IPv4Mask(252, 0, 0, 0)},
		// 200.0.0.0/5
		{IP: net.IPv4(200, 0, 0, 0).To4(), Mask: net.IPv4Mask(248, 0, 0, 0)},
		// 208.0.0.0/4
		{IP: net.IPv4(208, 0, 0, 0).To4(), Mask: net.IPv4Mask(240, 0, 0, 0)},
	}

	return ipnets
}
