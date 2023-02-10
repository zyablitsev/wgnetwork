package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	bolt "go.etcd.io/bbolt"

	"wgnetwork/model"
	"wgnetwork/pkg/rpcapi"
)

// Config object.
type Config struct {
	SessionSecret string
	SessionTTL    time.Duration
}

// API object.
type API struct {
	ctx context.Context
	log logger
	cfg Config
	db  *bolt.DB
}

// New constructor for API.
func New(
	ctx context.Context,
	log logger,
	cfg Config,
	db *bolt.DB,
) *API {
	api := &API{
		ctx: ctx,
		log: log,
		cfg: cfg,
		db:  db,
	}

	return api
}

// RegisterHandlers on provided api handler.
func (api *API) RegisterHandlers(rpc *rpcapi.API) {
	rpc.Register("auth/signin", api.signIn)
	rpc.Register("auth/signout", api.signOut)
	rpc.Register("auth/session/check", api.sessionCheck)
}

// signIn handler
func (api *API) signIn(
	ctx context.Context, w http.ResponseWriter, r json.RawMessage,
) (json.RawMessage, error) {
	ip, _, err := sessionCtx(ctx)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	tx, err := api.db.Begin(true) // writeable tx
	if err != nil {
		return rpcError{500, "bad gateway"}.marshal(), err
	}
	defer tx.Rollback()

	d, err := model.LoadDevice(tx, ip)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	u, err := model.LoadUser(tx, d.UserUUID)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	if !u.IsManager {
		return rpcError{400, "bad request"}.marshal(), err
	}

	request := new(signInRequest)
	err = json.Unmarshal(r, &request)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	err = u.OTPCheck(request.Code)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	s, err := u.CreateSession(api.cfg.SessionSecret, api.cfg.SessionTTL)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	err = u.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store user: %v", err)
		return rpcError{500, "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{500, "bad gateway"}.marshal(), err
	}

	w.Header().Set("x-session", s)
	response := signInResponse{Session: s}.marshal()

	return response, nil
}

type signInRequest struct {
	Code string `json:"code"`
}

func (s signInRequest) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

type signInResponse struct {
	Session string `json:"session"`
}

func (s signInResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// signOut handler
func (api *API) signOut(
	ctx context.Context, w http.ResponseWriter, _ json.RawMessage,
) (json.RawMessage, error) {
	ip, s, err := sessionCtx(ctx)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}
	if s == "" {
		err := errors.New("session not found")
		return rpcError{400, "bad request"}.marshal(), err
	}

	tx, err := api.db.Begin(true) // writeable tx
	if err != nil {
		return rpcError{500, "bad gateway"}.marshal(), err
	}
	defer tx.Rollback()

	d, err := model.LoadDevice(tx, ip)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	u, err := model.SessionUser(tx, api.cfg.SessionSecret, s)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	if u.UUID != d.UserUUID {
		return rpcError{400, "bad request"}.marshal(), err
	}

	u.DestroySession()

	err = u.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store user: %v", err)
		return rpcError{500, "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{500, "bad gateway"}.marshal(), err
	}

	return json.RawMessage(`{"msg": "ok"}`), nil
}

// sessionCheck handler
func (api *API) sessionCheck(
	ctx context.Context, w http.ResponseWriter, _ json.RawMessage,
) (json.RawMessage, error) {
	ip, s, err := sessionCtx(ctx)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}
	if s == "" {
		err := errors.New("session not found")
		return rpcError{400, "bad request"}.marshal(), err
	}

	tx, err := api.db.Begin(true) // writeable tx
	if err != nil {
		return rpcError{500, "bad gateway"}.marshal(), err
	}
	defer tx.Rollback()

	d, err := model.LoadDevice(tx, ip)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	u, err := model.SessionUser(tx, api.cfg.SessionSecret, s)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	if u.UUID != d.UserUUID {
		return rpcError{400, "bad request"}.marshal(), err
	}

	s, err = u.ProlongSession(api.cfg.SessionSecret, api.cfg.SessionTTL)
	if err != nil {
		return rpcError{400, "bad request"}.marshal(), err
	}

	err = u.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store user: %v", err)
		return rpcError{500, "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{500, "bad gateway"}.marshal(), err
	}

	w.Header().Set("x-session", s)
	response := sessionCheckResponse{Session: s}.marshal()

	return response, nil
}

type sessionCheckResponse struct {
	Session string `json:"session"`
}

func (s sessionCheckResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// rpcError object
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s rpcError) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

func sessionCtx(ctx context.Context) (net.IP, string, error) {
	params, ok := rpcapi.FromContext(ctx)
	if !ok {
		err := fmt.Errorf("context params value doesn't exists")
		return nil, "", err
	}

	ip := net.IP(params.RemoteAddr).To4()
	if ip == nil {
		err := fmt.Errorf("bad remote addr %q", string(params.RemoteAddr))
		return nil, "", err
	}

	s := params.Headers.Get("x-session")
	return ip, s, nil
}

// logger desribes interface of log object.
type logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warning(...interface{})
	Warningf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}
