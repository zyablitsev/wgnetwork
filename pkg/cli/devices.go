package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"wgnetwork/api/manager"
	"wgnetwork/pkg/pretty"
	"wgnetwork/pkg/rpcapi"
)

// ActionDevices object.
type ActionDevices struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
}

// NewActionDevices constructor.
func NewActionDevices(log logger) *ActionDevices {
	flagset := flag.NewFlagSet(
		"devices",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")

	a := &ActionDevices{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionDevices) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionDevices) Execute(args []string) error {
	logPrefix := "[devices] Execute"

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

	b := rpcapi.Request{
		Method: "manager/devices",
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

	result := manager.DeviceListResponse{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return err
	}

	table := pretty.NewTable(6)
	table.SetHeader(
		[]string{"ipnetwork", "pubkey", "label", "wan_forward", "allowed ips", "user uuid"})

	for _, d := range result {
		table.AddRow([]string{
			d.IPNetwork,
			d.PubKey,
			d.Label,
			strconv.FormatBool(d.WANForward),
			strings.Join(d.AllowedIPs, "\n"),
			d.UserUUID})
	}

	os.Stdout.WriteString(table.Render())

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionDevices) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	return nil
}
