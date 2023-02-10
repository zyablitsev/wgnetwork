package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"wgnetwork/model"
)

// trustIPSetAdd handler
func (api *API) trustIPSetAdd(
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

	request := new(TrustIPSetRequest)
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

	ipset, err := model.LoadManagerSSHTrustIPSet(tx)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	ip := net.ParseIP(request.IP).To4()
	ipset.Add(ip)

	err = ipset.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store ipset: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	return json.RawMessage(`{"msg": "ok"}`), nil
}

// trustIPSetRemove handler
func (api *API) trustIPSetRemove(
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

	request := new(TrustIPSetRequest)
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

	ipset, err := model.LoadManagerSSHTrustIPSet(tx)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	ip := net.ParseIP(request.IP).To4()
	ipset.Remove(ip)

	err = ipset.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store ipset: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	return json.RawMessage(`{"msg": "ok"}`), nil
}

// TrustIPSetRequest model.
type TrustIPSetRequest struct {
	IP string `json:"ip"`
}

func (s *TrustIPSetRequest) validate() (string, error) {
	if net.ParseIP(s.IP).To4() == nil {
		err := errors.New("required")
		return "ip", err
	}

	return "", nil
}

// Marshall returns the json encoding of TrustIPSetRequest.
func (s TrustIPSetRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// trustIPSet handler
func (api *API) trustIPSet(
	ctx context.Context, w http.ResponseWriter, _ json.RawMessage,
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

	ipset, err := model.LoadManagerSSHTrustIPSet(tx)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	response := TrustIPSetResponse(ipset)

	return response.marshal(), nil
}

// TrustIPSetResponse model.
type TrustIPSetResponse []net.IP

func (s TrustIPSetResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}
