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

	"github.com/zyablitsev/qrencode-go/qrencode"

	"wgnetwork/api/manager"
	"wgnetwork/pkg/pretty"
	"wgnetwork/pkg/rpcapi"
)

// ActionUserCreate object.
type ActionUserCreate struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	name       *string
	isManager  *bool
}

// NewActionUserCreate constructor.
func NewActionUserCreate(log logger) *ActionUserCreate {
	flagset := flag.NewFlagSet(
		"user-create",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	name := flagset.String(
		"name",
		"",
		"name")
	isManager := flagset.Bool(
		"is_manager",
		false,
		"is manager flag")

	a := &ActionUserCreate{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		name:       name,
		isManager:  isManager,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionUserCreate) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionUserCreate) Execute(args []string) error {
	logPrefix := "[user-create] Execute"

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

	b := manager.UserCreateRequest{
		Name:      *a.name,
		IsManager: *a.isManager,
	}.Marshal()
	b = rpcapi.Request{
		Method: "manager/user/create",
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

	result := manager.UserCreateResponse{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return err
	}

	var table *pretty.Table
	if result.IsManager {
		table = pretty.NewTable(5)
		table.SetHeader(
			[]string{"uuid", "name", "is manager", "key", "provision uri"})

		table.AddRow([]string{
			result.UUID,
			result.Name,
			strconv.FormatBool(result.IsManager),
			result.Key,
			result.ProvisionURI})
	} else {
		table = pretty.NewTable(3)
		table.SetHeader([]string{"uuid", "name", "is manager"})

		table.AddRow([]string{
			result.UUID,
			result.Name,
			strconv.FormatBool(result.IsManager)})
	}

	os.Stdout.WriteString(table.Render())

	if result.IsManager {
		grid, err := qrencode.Encode(result.ProvisionURI, qrencode.ECLevelL)
		if err != nil {
			return err
		}
		os.Stdout.WriteString("\nscan with authenticator app\n")
		grid.TerminalOutput(os.Stdout)
	}

	a.log.Debugf("%s: done", logPrefix)

	return nil
}

func (a *ActionUserCreate) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	if a.name == nil || len(*a.name) == 0 {
		return errors.New("name required")
	}

	return nil
}
