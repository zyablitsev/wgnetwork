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

// ActionWgCfg object.
type ActionWgCfg struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
}

// NewActionWgCfg constructor.
func NewActionWgCfg(log logger) *ActionWgCfg {
	flagset := flag.NewFlagSet(
		"wgcfg",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")

	a := &ActionWgCfg{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionWgCfg) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionWgCfg) Execute(args []string) error {
	logPrefix := "[wgcfg] Execute"

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
		Method: "manager/wg/cfg",
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

	result := manager.WgCfgResponse{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return err
	}

	table := pretty.NewTable(4)
	table.SetHeader([]string{"wan ip", "wg inet", "wg port", "wg pubkey"})

	table.AddRow([]string{
		result.WanIP,
		result.WgInet,
		strconv.FormatUint(uint64(result.WgPort), 10),
		result.WgPubKey})

	os.Stdout.WriteString(table.Render())

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionWgCfg) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	return nil
}
