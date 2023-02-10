package manager

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"wgnetwork/model"
)

// wgCfg handler
func (api *API) wgCfg(
	ctx context.Context, w http.ResponseWriter, _ json.RawMessage,
) (json.RawMessage, error) {
	if api.cfg.AuthRequired {
		tx, err := api.db.Begin(false) // non-writeable tx
		if err != nil {
			return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
		}
		defer tx.Rollback()

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

	response := WgCfgResponse{
		WanIP:    api.cfg.WanIP.String(),
		WgInet:   api.cfg.WgInet.String(),
		WgPort:   api.cfg.WgPort,
		WgPubKey: api.cfg.WgPubKey.String(),
	}.marshal()

	return response, nil
}

// WgCfgResponse model.
type WgCfgResponse struct {
	WanIP    string `json:"wanip"`
	WgInet   string `json:"wg_inet"`
	WgPort   uint16 `json:"wg_port"`
	WgPubKey string `json:"wg_pubkey"`
}

func (s WgCfgResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}
