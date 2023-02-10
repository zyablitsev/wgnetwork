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

	"github.com/google/uuid"
	"github.com/zyablitsev/qrencode-go/qrencode"

	"wgnetwork/api/manager"
	"wgnetwork/pkg/pretty"
	"wgnetwork/pkg/rpcapi"
)

// ActionUserEdit object.
type ActionUserEdit struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	uuid       *string
	name       *string
	isManager  *bool
}

// NewActionUserEdit constructor.
func NewActionUserEdit(log logger) *ActionUserEdit {
	flagset := flag.NewFlagSet(
		"user-edit",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	uuid := flagset.String(
		"uuid",
		"",
		"user uuid")
	name := flagset.String(
		"name",
		"",
		"name")
	isManager := flagset.Bool(
		"is_manager",
		false, // TODO: make optional
		"is manager flag")

	a := &ActionUserEdit{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		uuid:       uuid,
		name:       name,
		isManager:  isManager,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionUserEdit) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionUserEdit) Execute(args []string) error {
	logPrefix := "[user-edit] Execute"

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

	request := manager.UserEditRequest{UUID: *a.uuid}
	if a.name != nil && len(*a.name) > 0 {
		request.Name = a.name
	}
	if a.isManager != nil {
		request.IsManager = a.isManager
	}
	b := request.Marshal()
	b = rpcapi.Request{
		Method: "manager/user/edit",
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

	result := manager.UserResponse{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return err
	}

	var table *pretty.Table
	if result.IsManager {
		table = pretty.NewTable(6)
		table.SetHeader(
			[]string{
				"uuid", "name", "is manager",
				"key", "provision uri",
				"devices"})

		devices := make([]string, len(result.Devices))
		for i := 0; i < len(result.Devices); i++ {
			devices[i] = result.Devices[i].String()
		}

		table.AddRow([]string{
			result.UUID,
			result.Name,
			strconv.FormatBool(result.IsManager),
			result.Key,
			result.ProvisionURI,
			strings.Join(devices, "\n")})
	} else {
		table = pretty.NewTable(4)
		table.SetHeader([]string{"uuid", "name", "is manager", "devices"})

		devices := make([]string, len(result.Devices))
		for i := 0; i < len(result.Devices); i++ {
			devices[i] = result.Devices[i].String()
		}

		table.AddRow([]string{
			result.UUID,
			result.Name,
			strconv.FormatBool(result.IsManager),
			strings.Join(devices, "\n")})
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

func (a *ActionUserEdit) validate() error {
	if a.unixSocket == nil || len(*a.unixSocket) == 0 {
		return errors.New("unix-socket required")
	}

	if a.uuid == nil {
		return errors.New("uuid required")
	}
	_, err := uuid.Parse(*a.uuid)
	if err != nil {
		return errors.New("bad uuid value")
	}

	return nil
}
