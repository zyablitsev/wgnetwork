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

	"wgnetwork/api/manager"
	"wgnetwork/pkg/rpcapi"
)

// ActionTrustIPSetAdd object.
type ActionTrustIPSetAdd struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	ip         *string
}

// NewActionTrustIPSetAdd constructor.
func NewActionTrustIPSetAdd(log logger) *ActionTrustIPSetAdd {
	flagset := flag.NewFlagSet(
		"trust-ipset-add",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	ip := flagset.String(
		"ip",
		"",
		"device ip")

	a := &ActionTrustIPSetAdd{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		ip:         ip,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionTrustIPSetAdd) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionTrustIPSetAdd) Execute(args []string) error {
	logPrefix := "[trust-ipset-add] Execute"

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

	b := manager.TrustIPSetRequest{
		IP: *a.ip,
	}.Marshal()
	b = rpcapi.Request{
		Method: "manager/trust/ipset/add",
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

	os.Stdout.WriteString(string(response.Result) + "\n")

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionTrustIPSetAdd) validate() error {
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
