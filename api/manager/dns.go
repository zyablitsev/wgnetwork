package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"wgnetwork/model"
)

func (api *API) domainCreate(
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

	request := new(DomainRequest)
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

	// check if exists already
	_, err = model.LoadDomain(tx, request.Name)
	if err == nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	d := model.NewDomain(request.Name)
	err = d.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store domain: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := DomainResponse{d}

	return response.marshal(), nil
}

// DomainResponse model
type DomainResponse struct {
	model.Domain
}

func (s DomainResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// DomainRequest model.
type DomainRequest struct {
	Name string `json:"name"`
}

func (s *DomainRequest) validate() (string, error) {
	if len(s.Name) == 0 {
		err := errors.New("required")
		return "name", err
	}
	if len(s.Name) > 253 {
		err := errors.New("length should be lower than 253")
		return "name", err
	}

	return "", nil
}

// Marshall returns the json encoding of DeviceRequest.
func (s DomainRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

func (api *API) domainRecordSet(
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

	request := new(DomainRecordSetRequest)
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

	d, err := model.LoadDomain(tx, request.Name)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	rtype := strings.ToLower(request.Type)
	switch rtype {
	case "a":
		r, err := request.GetA()
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}
		d.SetA(r)
	case "cname":
		r, err := request.GetCNAME()
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}
		d.SetCNAME(r)
	}

	err = d.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store domain: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := DomainResponse{d}

	return response.marshal(), nil
}

// DomainRecordSetRequest model.
type DomainRecordSetRequest struct {
	Name string          `json:"name"`
	Type string          `json:"rtype"`
	Data json.RawMessage `json:"rdata"`
}

func (s *DomainRecordSetRequest) validate() (string, error) {
	if len(s.Name) == 0 {
		err := errors.New("required")
		return "name", err
	}
	if len(s.Name) > 253 {
		err := errors.New("length should be lower than 253")
		return "name", err
	}

	if s.Data == nil {
		err := errors.New("required")
		return "rdata", err
	}

	return "", nil
}

// Marshall returns the json encoding of DeviceRecordSetRequest.
func (s DomainRecordSetRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// GetA returns ARecord of data,
func (s DomainRecordSetRequest) GetA() (model.ARecord, error) {
	if strings.ToLower(s.Type) != "a" {
		return model.ARecord{}, errors.New("wrong type")
	}

	r := model.ARecord{}
	err := json.Unmarshal(s.Data, &r)
	if err != nil {
		return model.ARecord{}, err
	}

	return r, nil
}

// GetCNAME returns CNAMERecord of data,
func (s DomainRecordSetRequest) GetCNAME() (model.CNAMERecord, error) {
	if strings.ToLower(s.Type) != "cname" {
		return model.CNAMERecord{}, errors.New("wrong type")
	}

	r := model.CNAMERecord{}
	err := json.Unmarshal(s.Data, &r)
	if err != nil {
		return model.CNAMERecord{}, err
	}

	return r, nil
}

func (api *API) domainRecordRemove(
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

	request := new(DomainRecordRemoveRequest)
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

	d, err := model.LoadDomain(tx, request.Name)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	rtype := strings.ToLower(request.Type)
	switch rtype {
	case "a":
		ip, err := request.GetA()
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}
		d.RemoveA(ip)
	case "cname":
		target, err := request.GetCNAME()
		if err != nil {
			return rpcError{Code: 400, Message: "bad request"}.marshal(), err
		}
		d.RemoveCNAME(target)
	}

	err = d.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store domain: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := DomainResponse{d}

	return response.marshal(), nil
}

// DomainRecordRemoveRequest model.
type DomainRecordRemoveRequest struct {
	Name string          `json:"name"`
	Type string          `json:"rtype"`
	Data json.RawMessage `json:"rdata"`
}

func (s *DomainRecordRemoveRequest) validate() (string, error) {
	if len(s.Name) == 0 {
		err := errors.New("required")
		return "name", err
	}
	if len(s.Name) > 253 {
		err := errors.New("length should be lower than 253")
		return "name", err
	}

	if s.Data == nil {
		err := errors.New("required")
		return "rdata", err
	}

	return "", nil
}

// Marshall returns the json encoding of DeviceRecordRemoveRequest.
func (s DomainRecordRemoveRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// GetA returns ARecord of data,
func (s DomainRecordRemoveRequest) GetA() (net.IP, error) {
	if strings.ToLower(s.Type) != "a" {
		return nil, errors.New("wrong type")
	}

	ip := net.IP{}
	err := json.Unmarshal(s.Data, &ip)
	if err != nil {
		return nil, err
	}
	ip = ip.To4()

	return ip, nil
}

// GetCNAME returns CNAMERecord of data,
func (s DomainRecordRemoveRequest) GetCNAME() (string, error) {
	if strings.ToLower(s.Type) != "cname" {
		return "", errors.New("wrong type")
	}

	target := ""
	err := json.Unmarshal(s.Data, &target)
	if err != nil {
		return "", err
	}

	return target, nil
}

func (api *API) domainRemove(
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

	request := new(DomainRequest)
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

	err = model.RemoveDomain(tx, request.Name)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	return json.RawMessage(`{"msg": "ok"}`), nil
}

func (api *API) domain(
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

	request := new(DomainRequest)
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

	d, err := model.LoadDomain(tx, request.Name)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := DomainResponse{d}

	return response.marshal(), nil
}

func (api *API) domainList(
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

	domains, err := model.LoadDomains(tx)
	if err != nil {
		return rpcError{Code: 400, Message: "bad request"}.marshal(), err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := DomainListResponse(domains)

	return response.marshal(), nil
}

// DomainListResponse model.
type DomainListResponse model.Domains

func (s DomainListResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}
