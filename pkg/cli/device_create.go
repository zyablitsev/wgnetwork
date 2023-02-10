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

	"github.com/google/uuid"
	"github.com/zyablitsev/qrencode-go/qrencode"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"wgnetwork/api/manager"
	"wgnetwork/pkg/pretty"
	"wgnetwork/pkg/rpcapi"
)

// ActionDeviceCreate object.
type ActionDeviceCreate struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	userUUID   *string
	label      *string
	wanForward *bool
	wgPubKey   *string
}

// NewActionDeviceCreate constructor.
func NewActionDeviceCreate(log logger) *ActionDeviceCreate {
	flagset := flag.NewFlagSet(
		"device-create",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	userUUID := flagset.String(
		"user_uuid",
		"",
		"user uuid")
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

	a := &ActionDeviceCreate{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		userUUID:   userUUID,
		label:      label,
		wanForward: wanForward,
		wgPubKey:   wgPubKey,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionDeviceCreate) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionDeviceCreate) Execute(args []string) error {
	logPrefix := "[device-create] Execute"

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

	b := manager.DeviceCreateRequest{
		UserUUID:    *a.userUUID,
		Label:       *a.label,
		WANForward:  *a.wanForward,
		WGPublicKey: *a.wgPubKey,
	}.Marshal()
	b = rpcapi.Request{
		Method: "manager/device/create",
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

	result := manager.DeviceCreateResponse{}
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
	if result.WgDevicePrivKey != "" {
		privKey = result.WgDevicePrivKey
	}
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

	if result.WgDevicePrivKey != "" {
		grid, err := qrencode.Encode(cfg, qrencode.ECLevelL)
		if err != nil {
			return err
		}
		os.Stdout.WriteString("\nscan with wireguard app\n")
		grid.TerminalOutput(os.Stdout)
	}

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionDeviceCreate) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	if a.userUUID == nil {
		return errors.New("user_uuid required")
	}
	_, err := uuid.Parse(*a.userUUID)
	if err != nil {
		return errors.New("bad uuid value")
	}

	if a.label == nil || len(*a.label) == 0 {
		return errors.New("label required")
	}

	if a.wgPubKey != nil && len(*a.wgPubKey) > 0 {
		_, err = wgtypes.ParseKey(*a.wgPubKey)
		if err != nil {
			return errors.New("bad wg_pubkey value")
		}
	}

	return nil
}

func buildWgCfg(
	sk, pk string,
	address, allowedIPs, endpoint string,
) string {
	var buf strings.Builder
	nl := "\n"
	buf.WriteString("[Interface]")
	buf.WriteString(nl)
	buf.WriteString("PrivateKey = " + sk)
	buf.WriteString(nl)
	buf.WriteString("Address = " + address)
	buf.WriteString(nl)
	buf.WriteString(nl)
	buf.WriteString("[Peer]")
	buf.WriteString(nl)
	buf.WriteString("PublicKey = " + pk)
	buf.WriteString(nl)
	buf.WriteString("AllowedIPs = " + allowedIPs)
	buf.WriteString(nl)
	buf.WriteString("Endpoint = " + endpoint)
	buf.WriteString(nl)
	buf.WriteString("PersistentKeepalive = 25")
	buf.WriteString(nl)

	return buf.String()
}
