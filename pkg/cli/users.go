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

// ActionUsers object.
type ActionUsers struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
}

// NewActionUsers constructor.
func NewActionUsers(log logger) *ActionUsers {
	flagset := flag.NewFlagSet(
		"users",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")

	a := &ActionUsers{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionUsers) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionUsers) Execute(args []string) error {
	logPrefix := "[users] Execute"

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
		Method: "manager/users",
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

	result := manager.UserListResponse{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return err
	}

	table := pretty.NewTable(4)
	table.SetHeader([]string{"uuid", "name", "is manager", "devices"})

	for _, u := range result {
		devices := make([]string, len(u.Devices))
		for i := 0; i < len(u.Devices); i++ {
			devices[i] = u.Devices[i].String()
		}

		table.AddRow([]string{
			u.UUID,
			u.Name,
			strconv.FormatBool(u.IsManager),
			strings.Join(devices, "\n")})
	}

	os.Stdout.WriteString(table.Render())

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionUsers) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	return nil
}
