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

// ActionUser object.
type ActionUser struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	uuid       *string
}

// NewActionUser constructor.
func NewActionUser(log logger) *ActionUser {
	flagset := flag.NewFlagSet(
		"user",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	uuid := flagset.String(
		"uuid",
		"",
		"user uuid")

	a := &ActionUser{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		uuid:       uuid,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionUser) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionUser) Execute(args []string) error {
	logPrefix := "[user] Execute"

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

	b := manager.UserRequest{
		UUID: *a.uuid,
	}.Marshal()
	b = rpcapi.Request{
		Method: "manager/user",
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

func (a *ActionUser) validate() error {
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
