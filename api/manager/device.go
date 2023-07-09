package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"wgnetwork/model"
)

// deviceCreate handler
func (api *API) deviceCreate(
	ctx context.Context, w http.ResponseWriter, r json.RawMessage,
) (json.RawMessage, error) {
	tx, err := api.db.Begin(true) // writeable tx
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}
	defer tx.Rollback()

	if api.cfg.AuthRequired {
		ip, s, err := sessionCtx(ctx)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}
		if s == "" {
			err := errors.New("session not found")
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		d, err := model.LoadDevice(tx, ip)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		u, err := model.SessionUser(tx, api.cfg.SessionSecret, s)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		if u.UUID != d.UserUUID {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		w.Header().Set("x-session", s)
	}

	request := new(DeviceCreateRequest)
	err = json.Unmarshal(r, &request)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	field, err := request.validate()
	if err != nil {
		msg := err.Error()
		err = errors.New("validation error")
		b := validateError{field, msg}.marshal()
		b = rpcError{Code: 401, Message: "bad request", Data: b}.marshal()
		return b, err
	}

	u, err := model.LoadUser(tx, request.UserUUID)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	var (
		sk wgtypes.Key
		pk wgtypes.Key
	)
	if len(request.WGPublicKey) > 0 {
		// omit error check, we already validate this field
		pk, _ = wgtypes.ParseKey(request.WGPublicKey)
	} else {
		sk, err = wgtypes.GeneratePrivateKey()
		if err != nil {
			return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
		}
		pk = sk.PublicKey()
	}

	ipNetwork, err := model.AllocateIP(tx, api.cfg.WgInet, pk)
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	d := model.Device{
		IPNetwork:  ipNetwork,
		PubKey:     pk,
		Label:      request.Label,
		WANForward: request.WANForward,

		UserUUID: u.UUID}

	err = d.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store device: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	u.AddDevice(ipNetwork.IP)
	err = u.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store user: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := DeviceCreateResponse{
		UserUUID:   u.UUID,
		UserName:   u.Name,
		Label:      d.Label,
		WANForward: d.WANForward,

		WgDeviceInet:       d.CIDR().String(),
		WgDevicePort:       api.cfg.WgPort,
		WgDevicePubKey:     pk.String(),
		WgDeviceAllowedIPs: make([]string, len(d.AllowedIPs())),

		WgInet:   api.cfg.WgInet.String(),
		WgIPNet:  api.cfg.WgIPNet.String(),
		WgIP:     api.cfg.WgInet.IP.String(),
		WgPort:   api.cfg.WgPort,
		WgPubKey: api.cfg.WgPubKey.String(),

		WanIP: api.cfg.WanIP.String(),
	}
	if sk != [wgtypes.KeyLen]byte{} {
		response.WgDevicePrivKey = sk.String()
	}
	for idx, item := range d.AllowedIPs() {
		response.WgDeviceAllowedIPs[idx] = item.String()
	}

	return response.marshal(), nil
}

// DeviceCreateRequest model.
type DeviceCreateRequest struct {
	UserUUID    string `json:"user_uuid"`
	Label       string `json:"label"`
	WANForward  bool   `json:"wan_forward"`
	WGPublicKey string `json:"wg_public_key"`
}

func (s *DeviceCreateRequest) validate() (string, error) {
	_, err := uuid.Parse(s.UserUUID)
	if err != nil {
		err := errors.New("required")
		return "user_uuid", err
	}

	if len(s.Label) == 0 {
		err := errors.New("required")
		return "label", err
	}
	if len(s.Label) > 255 {
		err := errors.New("length should be lower than 256")
		return "label", err
	}

	if len(s.WGPublicKey) > 0 {
		_, err := wgtypes.ParseKey(s.WGPublicKey)
		if err != nil {
			err = errors.New("bad value")
			return "wg_public_key", err
		}
	}

	return "", nil
}

// Marshall returns the json encoding of DeviceCreateRequest.
func (s DeviceCreateRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// DeviceCreateResponse model.
type DeviceCreateResponse struct {
	UserUUID   string `json:"user_uuid"`
	UserName   string `json:"user_name"`
	Label      string `json:"label"`
	WANForward bool   `json:"wan_forward"`

	WgDeviceInet       string   `json:"wg_device_inet"`
	WgDevicePort       uint16   `json:"wg_device_port"`
	WgDevicePubKey     string   `json:"wg_device_pubkey"`
	WgDevicePrivKey    string   `json:"wg_device_privkey,omitempty"`
	WgDeviceAllowedIPs []string `json:"wg_device_allowed_ips"`

	WgInet   string `json:"wg_server_inet"`
	WgIPNet  string `json:"wg_server_ipnet"`
	WgIP     string `json:"wg_server_ip"`
	WgPort   uint16 `json:"wg_server_port"`
	WgPubKey string `json:"wg_server_pubkey"`

	WanIP string `json:"server_wanip"`
}

func (s DeviceCreateResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// deviceEdit handler
func (api *API) deviceEdit(
	ctx context.Context, w http.ResponseWriter, r json.RawMessage,
) (json.RawMessage, error) {
	tx, err := api.db.Begin(true) // writeable tx
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}
	defer tx.Rollback()

	if api.cfg.AuthRequired {
		ip, s, err := sessionCtx(ctx)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}
		if s == "" {
			err := errors.New("session not found")
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		d, err := model.LoadDevice(tx, ip)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		u, err := model.SessionUser(tx, api.cfg.SessionSecret, s)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		if u.UUID != d.UserUUID {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		w.Header().Set("x-session", s)
	}

	request := new(DeviceEditRequest)
	err = json.Unmarshal(r, &request)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	field, err := request.validate()
	if err != nil {
		msg := err.Error()
		err = errors.New("validation error")
		b := validateError{field, msg}.marshal()
		b = rpcError{Code: 401, Message: "bad request", Data: b}.marshal()
		return b, err
	}

	ip := net.ParseIP(request.IP).To4()
	d, err := model.LoadDevice(tx, ip)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	u, err := model.LoadUser(tx, d.UserUUID)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	if request.Label != nil {
		d.Label = *request.Label
	}
	if request.WANForward != nil {
		d.WANForward = *request.WANForward
	}
	if request.WGPublicKey != nil {
		// omit error check, we already validate this field
		pk, _ := wgtypes.ParseKey(*request.WGPublicKey)
		d.PubKey = pk
	}

	err = d.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store device: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := DeviceResponse{
		UserUUID:           u.UUID,
		UserName:           u.Name,
		Label:              d.Label,
		WANForward:         d.WANForward,
		WgDeviceAllowedIPs: make([]string, len(d.AllowedIPs())),

		WgDeviceInet:   d.CIDR().String(),
		WgDevicePort:   api.cfg.WgPort,
		WgDevicePubKey: d.PubKey.String(),

		WgInet:   api.cfg.WgInet.String(),
		WgIPNet:  api.cfg.WgIPNet.String(),
		WgIP:     api.cfg.WgInet.IP.String(),
		WgPort:   api.cfg.WgPort,
		WgPubKey: api.cfg.WgPubKey.String(),

		WanIP: api.cfg.WanIP.String(),
	}

	for idx, item := range d.AllowedIPs() {
		response.WgDeviceAllowedIPs[idx] = item.String()
	}

	return response.marshal(), nil
}

// DeviceEditRequest model.
type DeviceEditRequest struct {
	IP          string  `json:"ip"`
	Label       *string `json:"label"`
	WANForward  *bool   `json:"wan_forward"`
	WGPublicKey *string `json:"wg_public_key"`
}

func (s *DeviceEditRequest) validate() (string, error) {
	if net.ParseIP(s.IP).To4() == nil {
		err := errors.New("required")
		return "ip", err
	}

	if s.Label != nil && (len(*s.Label) == 0 || len(*s.Label) > 255) {
		err := errors.New(
			"length should be lower than 256 and non empty string")
		return "label", err
	}

	if s.WGPublicKey != nil {
		_, err := wgtypes.ParseKey(*s.WGPublicKey)
		if err != nil {
			err = errors.New("bad value")
			return "wg_public_key", err
		}
	}

	return "", nil
}

// Marshall returns the json encoding of DeviceEditRequest.
func (s DeviceEditRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// deviceRemove handler
func (api *API) deviceRemove(
	ctx context.Context, w http.ResponseWriter, r json.RawMessage,
) (json.RawMessage, error) {
	tx, err := api.db.Begin(true) // writeable tx
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}
	defer tx.Rollback()

	if api.cfg.AuthRequired {
		ip, s, err := sessionCtx(ctx)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}
		if s == "" {
			err := errors.New("session not found")
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		d, err := model.LoadDevice(tx, ip)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		u, err := model.SessionUser(tx, api.cfg.SessionSecret, s)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		if u.UUID != d.UserUUID {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		w.Header().Set("x-session", s)
	}

	request := new(DeviceRemoveRequest)
	err = json.Unmarshal(r, &request)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	field, err := request.validate()
	if err != nil {
		msg := err.Error()
		err = errors.New("validation error")
		b := validateError{field, msg}.marshal()
		b = rpcError{Code: 401, Message: "bad request", Data: b}.marshal()
		return b, err
	}

	ip := net.ParseIP(request.IP).To4()
	d, err := model.LoadDevice(tx, ip)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	u, err := model.LoadUser(tx, d.UserUUID)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	err = model.RemoveDevice(tx, ip)
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	u.RemoveDevice(ip)
	err = u.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store user: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	return json.RawMessage(`{"msg": "ok"}`), nil
}

// DeviceRemoveRequest model.
type DeviceRemoveRequest struct {
	IP string `json:"ip"`
}

func (s *DeviceRemoveRequest) validate() (string, error) {
	if net.ParseIP(s.IP).To4() == nil {
		err := errors.New("required")
		return "ip", err
	}

	return "", nil
}

// Marshall returns the json encoding of DeviceRemoveRequest.
func (s DeviceRemoveRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// device handler
func (api *API) device(
	ctx context.Context, w http.ResponseWriter, r json.RawMessage,
) (json.RawMessage, error) {
	tx, err := api.db.Begin(false) // non-writeable tx
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}
	defer tx.Rollback()

	if api.cfg.AuthRequired {
		ip, s, err := sessionCtx(ctx)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}
		if s == "" {
			err := errors.New("session not found")
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		d, err := model.LoadDevice(tx, ip)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		u, err := model.SessionUser(tx, api.cfg.SessionSecret, s)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		if u.UUID != d.UserUUID {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		w.Header().Set("x-session", s)
	}

	request := new(DeviceRequest)
	err = json.Unmarshal(r, &request)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	field, err := request.Validate()
	if err != nil {
		msg := err.Error()
		err = errors.New("validation error")
		b := validateError{field, msg}.marshal()
		b = rpcError{Code: 401, Message: "bad request", Data: b}.marshal()
		return b, err
	}

	ip := net.ParseIP(request.IP).To4()
	d, err := model.LoadDevice(tx, ip)
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	u, err := model.LoadUser(tx, d.UserUUID)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	response := DeviceResponse{
		UserUUID:   u.UUID,
		UserName:   u.Name,
		Label:      d.Label,
		WANForward: d.WANForward,

		WgDeviceInet:       d.CIDR().String(),
		WgDevicePort:       api.cfg.WgPort,
		WgDevicePubKey:     d.PubKey.String(),
		WgDeviceAllowedIPs: make([]string, len(d.AllowedIPs())),

		WgInet:   api.cfg.WgInet.String(),
		WgIPNet:  api.cfg.WgIPNet.String(),
		WgIP:     api.cfg.WgInet.IP.String(),
		WgPort:   api.cfg.WgPort,
		WgPubKey: api.cfg.WgPubKey.String(),

		WanIP: api.cfg.WanIP.String(),
	}
	for idx, item := range d.AllowedIPs() {
		response.WgDeviceAllowedIPs[idx] = item.String()
	}

	return response.marshal(), nil
}

// DeviceRequest model.
type DeviceRequest struct {
	IP string `json:"ip"`
}

// Validate fields.
func (s *DeviceRequest) Validate() (string, error) {
	if net.ParseIP(s.IP).To4() == nil {
		err := errors.New("required")
		return "ip", err
	}

	return "", nil
}

// Marshall returns the json encoding of DeviceRequest.
func (s DeviceRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// DeviceResponse model
type DeviceResponse struct {
	UserUUID   string `json:"user_uuid"`
	UserName   string `json:"user_name"`
	Label      string `json:"label"`
	WANForward bool   `json:"wan_forward"`

	WgDeviceInet       string   `json:"wg_device_inet"`
	WgDevicePort       uint16   `json:"wg_device_port"`
	WgDevicePubKey     string   `json:"wg_device_pubkey"`
	WgDeviceAllowedIPs []string `json:"wg_device_allowed_ips"`

	WgInet   string `json:"wg_server_inet"`
	WgIPNet  string `json:"wg_server_ipnet"`
	WgIP     string `json:"wg_server_ip"`
	WgPort   uint16 `json:"wg_server_port"`
	WgPubKey string `json:"wg_server_pubkey"`

	WanIP string `json:"server_wanip"`
}

func (s DeviceResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// deviceList handler
func (api *API) deviceList(
	ctx context.Context, w http.ResponseWriter, _ json.RawMessage,
) (json.RawMessage, error) {
	tx, err := api.db.Begin(false)
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}
	defer tx.Rollback()

	if api.cfg.AuthRequired {
		ip, s, err := sessionCtx(ctx)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}
		if s == "" {
			err := errors.New("session not found")
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		d, err := model.LoadDevice(tx, ip)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		u, err := model.SessionUser(tx, api.cfg.SessionSecret, s)
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		if u.UUID != d.UserUUID {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}

		w.Header().Set("x-session", s)
	}

	devices, err := model.LoadDevices(tx)
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := make(DeviceListResponse, len(devices))
	for i, d := range devices {
		response[i] = DeviceListItem{
			IPNetwork:  d.IPNetwork.CIDR().String(),
			PubKey:     d.PubKey.String(),
			Label:      d.Label,
			WANForward: d.WANForward,
			AllowedIPs: make([]string, len(d.AllowedIPs())),

			UserUUID: d.UserUUID,
		}

		for j, item := range d.AllowedIPs() {
			response[i].AllowedIPs[j] = item.String()
		}
	}

	return response.marshal(), nil
}

// DeviceListItem model.
type DeviceListItem struct {
	IPNetwork  string   `json:"ipnetwork"`
	PubKey     string   `json:"pub_key"`
	Label      string   `json:"label"`
	WANForward bool     `json:"wan_forward"`
	AllowedIPs []string `json:"allowed_ips"`

	UserUUID string `json:"user_uuid"`
}

// DeviceListResponse model.
type DeviceListResponse []DeviceListItem

func (s DeviceListResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}
