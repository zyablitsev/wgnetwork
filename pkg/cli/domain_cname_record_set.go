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
	"wgnetwork/model"
	"wgnetwork/pkg/pretty"
	"wgnetwork/pkg/rpcapi"
)

// ActionDomainCNameRecordSet object.
type ActionDomainCNameRecordSet struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	name       *string
	target     *string
	ttl        *int
}

// NewActionDomainCNameRecordSet constructor.
func NewActionDomainCNameRecordSet(log logger) *ActionDomainCNameRecordSet {
	flagset := flag.NewFlagSet(
		"domain-cname-record-set",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	name := flagset.String(
		"name",
		"",
		"domain name")
	target := flagset.String(
		"target",
		"",
		"record target value")
	ttl := flagset.Int(
		"ttl",
		30,
		"ttl value")

	a := &ActionDomainCNameRecordSet{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		name:       name,
		target:     target,
		ttl:        ttl,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionDomainCNameRecordSet) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionDomainCNameRecordSet) Execute(args []string) error {
	logPrefix := "[domain-cname-record-set] Execute"

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

	params := model.CNAMERecord{
		TTL:    uint32(*a.ttl),
		Target: *a.target,
	}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	b := manager.DomainRecordSetRequest{
		Name: *a.name,
		Type: "cname",
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

func (a *ActionDomainCNameRecordSet) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	if a.name == nil || len(*a.name) == 0 {
		return errors.New("name required")
	}

	if a.target == nil {
		return errors.New("target required")
	}

	return nil
}
