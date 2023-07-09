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
	"wgnetwork/pkg/rpcapi"
)

// ActionDomainRemove object.
type ActionDomainRemove struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	name       *string
}

// NewActionDomainRemove constructor.
func NewActionDomainRemove(log logger) *ActionDomainRemove {
	flagset := flag.NewFlagSet(
		"domain-remove",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	name := flagset.String(
		"name",
		"",
		"domain name")

	a := &ActionDomainRemove{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		name:       name,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionDomainRemove) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionDomainRemove) Execute(args []string) error {
	logPrefix := "[domain-remove] Execute"

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

	b := manager.DomainRequest{
		Name: *a.name,
	}.Marshal()
	b = rpcapi.Request{
		Method: "manager/dns/domain/remove",
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

func (a *ActionDomainRemove) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	if a.name == nil || len(*a.name) == 0 {
		return errors.New("name required")
	}

	return nil
}
