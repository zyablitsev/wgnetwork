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

	"github.com/google/uuid"

	"wgnetwork/api/manager"
	"wgnetwork/pkg/rpcapi"
)

// ActionUserRemove object.
type ActionUserRemove struct {
	flagset *flag.FlagSet
	log     logger

	unixSocket *string
	uuid       *string
}

// NewActionUserRemove constructor.
func NewActionUserRemove(log logger) *ActionUserRemove {
	flagset := flag.NewFlagSet(
		"user-remove",
		flag.ExitOnError)

	unixSocket := flagset.String(
		"unix-socket",
		"/tmp/wgmanager.sock",
		"unix-socket")
	uuid := flagset.String(
		"uuid",
		"",
		"user uuid")

	a := &ActionUserRemove{
		flagset: flagset,
		log:     log,

		unixSocket: unixSocket,
		uuid:       uuid,
	}

	return a
}

// Usage prints out flagset usage.
func (a *ActionUserRemove) Usage() {
	a.flagset.Usage()
}

// Execute action.
func (a *ActionUserRemove) Execute(args []string) error {
	logPrefix := "[user-remove] Execute"

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

	b := manager.UserRemoveRequest{
		UUID: *a.uuid,
	}.Marshal()
	b = rpcapi.Request{
		Method: "manager/user/remove",
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

func (a *ActionUserRemove) validate() error {
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
