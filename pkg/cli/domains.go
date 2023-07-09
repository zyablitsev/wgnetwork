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

// ActionDomains object.
type ActionDomains struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
}

// NewActionDomains constructor.
func NewActionDomains(log logger) *ActionDomains {
	flagset := flag.NewFlagSet(
		"domains",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")

	a := &ActionDomains{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionDomains) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionDomains) Execute(args []string) error {
	logPrefix := "[domains] Execute"

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
		Method: "manager/dns/domains",
		Params: nil,
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

	result := manager.DomainListResponse{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return err
	}

	if len(result) == 0 {
		os.Stdout.WriteString("no domains\n")
		a.log.Debugf("%s: done", logPrefix)

		return nil
	}

	for _, d := range result {
		os.Stdout.WriteString("\ndomain: ")
		os.Stdout.WriteString(d.Name)
		os.Stdout.WriteString("\n")
		table := pretty.NewTable(3)
		table.SetHeader([]string{"type", "value", "ttl"})
		if d.CNAME != nil {
			table.AddRow([]string{"cname", d.CNAME.Target, strconv.FormatUint(uint64(d.CNAME.TTL), 10)})
		}
		for _, r := range d.A {
			table.AddRow([]string{"a", r.A.String(), strconv.FormatUint(uint64(r.TTL), 10)})
		}
		os.Stdout.WriteString(table.Render())
	}

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionDomains) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	return nil
}
