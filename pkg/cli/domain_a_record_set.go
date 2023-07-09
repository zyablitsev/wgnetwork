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

	"wgnetwork/api/manager"
	"wgnetwork/model"
	"wgnetwork/pkg/pretty"
	"wgnetwork/pkg/rpcapi"
)

// ActionDomainARecordSet object.
type ActionDomainARecordSet struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	name       *string
	ip         *string
	ttl        *int
}

// NewActionDomainARecordSet constructor.
func NewActionDomainARecordSet(log logger) *ActionDomainARecordSet {
	flagset := flag.NewFlagSet(
		"domain-a-record-set",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	name := flagset.String(
		"name",
		"",
		"domain name")
	ip := flagset.String(
		"ip",
		"",
		"record a value")
	ttl := flagset.Int(
		"ttl",
		30,
		"ttl value")

	a := &ActionDomainARecordSet{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		name:       name,
		ip:         ip,
		ttl:        ttl,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionDomainARecordSet) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionDomainARecordSet) Execute(args []string) error {
	logPrefix := "[domain-a-record-set] Execute"

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

	params := model.ARecord{
		TTL: uint32(*a.ttl),
		A:   net.ParseIP(*a.ip).To4(),
	}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	b := manager.DomainRecordSetRequest{
		Name: *a.name,
		Type: "a",
		Data: json.RawMessage(data),
	}.Marshal()
	b = rpcapi.Request{
		Method: "manager/dns/domain/record/set",
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

func (a *ActionDomainARecordSet) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	if a.name == nil || len(*a.name) == 0 {
		return errors.New("name required")
	}

	if a.ip == nil {
		return errors.New("ip required")
	}

	return nil
}
