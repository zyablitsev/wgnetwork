package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	bolt "go.etcd.io/bbolt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"wgnetwork/pkg/rpcapi"
)

// Config object.
type Config struct {
	AuthRequired bool

	WanIP    net.IP
	WgInet   *net.IPNet
	WgIPNet  *net.IPNet
	WgPort   uint16
	WgPubKey wgtypes.Key

	OTPIssuer string

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
	rpc.Register("manager/wg/cfg", api.wgCfg)

	rpc.Register("manager/user/create", api.userCreate)
	rpc.Register("manager/user/edit", api.userEdit)
	rpc.Register("manager/user/remove", api.userRemove)
	rpc.Register("manager/user", api.user)
	rpc.Register("manager/users", api.userList)

	rpc.Register("manager/device/create", api.deviceCreate)
	rpc.Register("manager/device/edit", api.deviceEdit)
	rpc.Register("manager/device/remove", api.deviceRemove)
	rpc.Register("manager/device", api.device)
	rpc.Register("manager/devices", api.deviceList)

	rpc.Register("manager/trust/ipset/add", api.trustIPSetAdd)
	rpc.Register("manager/trust/ipset/remove", api.trustIPSetRemove)
	rpc.Register("manager/trust/ipset", api.trustIPSet)
}

// rpcError object
type rpcError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (s rpcError) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// validateError object
type validateError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (s validateError) marshal() json.RawMessage {
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
