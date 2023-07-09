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
	"wgnetwork/api/manager"
	"wgnetwork/pkg/pretty"
	"wgnetwork/pkg/rpcapi"
)

// ActionDomain object.
type ActionDomain struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	name       *string
}

// NewActionDomain constructor.
func NewActionDomain(log logger) *ActionDomain {
	flagset := flag.NewFlagSet(
		"domain",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	name := flagset.String(
		"name",
		"",
		"domain name")

	a := &ActionDomain{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		name:       name,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionDomain) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionDomain) Execute(args []string) error {
	logPrefix := "[domain] Execute"

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
		Method: "manager/dns/domain",
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

	result := manager.DomainResponse{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return err
	}

	os.Stdout.WriteString("domain: ")
	os.Stdout.WriteString(result.Name)
	os.Stdout.WriteString("\n")
	table := pretty.NewTable(3)
	table.SetHeader([]string{"type", "value", "ttl"})
	if result.CNAME != nil {
		table.AddRow([]string{"cname", result.CNAME.Target, strconv.FormatUint(uint64(result.CNAME.TTL), 10)})
	}
	for _, r := range result.A {
		table.AddRow([]string{"a", r.A.String(), strconv.FormatUint(uint64(r.TTL), 10)})
	}
	os.Stdout.WriteString(table.Render())

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionDomain) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	if a.name == nil || len(*a.name) == 0 {
		return errors.New("name required")
	}

	return nil
}
