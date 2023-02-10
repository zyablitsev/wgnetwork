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
	"wgnetwork/api/manager"
	"wgnetwork/pkg/pretty"
	"wgnetwork/pkg/rpcapi"
)

// ActionTrustIPSet object.
type ActionTrustIPSet struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
}

// NewActionTrustIPSet constructor.
func NewActionTrustIPSet(log logger) *ActionTrustIPSet {
	flagset := flag.NewFlagSet(
		"trust-ipset",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")

	a := &ActionTrustIPSet{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionTrustIPSet) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionTrustIPSet) Execute(args []string) error {
	logPrefix := "[trust-ipset] Execute"

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
		Method: "manager/trust/ipset",
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

	result := manager.TrustIPSetResponse{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return err
	}

	ips := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		ips[i] = result[i].String()
	}

	table := pretty.NewTable(1)
	table.SetHeader([]string{"trusted ips"})

	for _, ip := range ips {
		table.AddRow([]string{ip})
	}

	os.Stdout.WriteString(table.Render())

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionTrustIPSet) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	return nil
}
