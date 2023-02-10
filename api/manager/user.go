package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"

	"wgnetwork/model"
)

// userCreate handler
func (api *API) userCreate(
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

	request := new(UserCreateRequest)
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

	u := model.NewUser(request.Name)

	response := UserCreateResponse{
		UUID: u.UUID,
		Name: u.Name,
	}

	if request.IsManager {
		provisionURI, secret, err := u.SetManager(api.cfg.OTPIssuer)
		if err != nil {
			err = fmt.Errorf("can't set manager: %v", err)
			return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
		}
		response.IsManager = true
		response.Key = secret
		response.ProvisionURI = provisionURI
	}

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

	return response.marshal(), nil
}

// UserCreateRequest model.
type UserCreateRequest struct {
	Name      string `json:"name"`
	IsManager bool   `json:"is_manager"`
}

func (s *UserCreateRequest) validate() (string, error) {
	if len(s.Name) == 0 {
		err := errors.New("required")
		return "name", err
	}
	if len(s.Name) > 255 {
		err := errors.New("length should be lower than 256")
		return "name", err
	}

	return "", nil
}

// Marshall returns the json encoding of UserCreateRequest.
func (s UserCreateRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

type UserCreateResponse struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	IsManager    bool   `json:"is_manager"`
	Key          string `json:"key,omitempty"`
	ProvisionURI string `json:"provision_uri,omitempty"`
}

func (s UserCreateResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// userEdit handler
func (api *API) userEdit(
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

	request := new(UserEditRequest)
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

	u, err := model.LoadUser(tx, request.UUID)
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	if request.Name != nil {
		u.Name = *request.Name
	}
	if request.IsManager != nil {
		if !u.IsManager && *request.IsManager {
			_, _, err := u.SetManager(api.cfg.OTPIssuer)
			if err != nil {
				err = fmt.Errorf("can't set manager: %v", err)
				return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
			}
		} else if u.IsManager && !*request.IsManager {
			u.UnsetManager()
		}
		u.IsManager = *request.IsManager
	}

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

	response := UserResponse{
		UUID:      u.UUID,
		Name:      u.Name,
		IsManager: u.IsManager,
		Devices:   u.Devices,
	}
	if u.IsManager {
		provisionURI, err := u.OTPProvisionURI(api.cfg.OTPIssuer)
		if err != nil {
			return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
		}
		response.Key = u.TFASecret
		response.ProvisionURI = provisionURI
	}

	return response.marshal(), nil
}

// UserEditRequest model.
type UserEditRequest struct {
	UUID      string  `json:"uuid"`
	Name      *string `json:"name"`
	IsManager *bool   `json:"is_manager"`
}

func (s *UserEditRequest) validate() (string, error) {
	_, err := uuid.Parse(s.UUID)
	if err != nil {
		err := errors.New("required")
		return "uuid", err
	}

	if s.Name != nil && (len(*s.Name) == 0 || len(*s.Name) > 255) {
		err := errors.New(
			"length should be lower than 256 and non empty string")
		return "name", err
	}

	return "", nil
}

// Marshall returns the json encoding of UserEditRequest.
func (s UserEditRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// userRemove handler
func (api *API) userRemove(
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

	request := new(UserRemoveRequest)
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
	err = model.RemoveUser(tx, request.UUID)
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

// UserRemoveRequest model.
type UserRemoveRequest struct {
	UUID string `json:"uuid"`
}

func (s *UserRemoveRequest) validate() (string, error) {
	_, err := uuid.Parse(s.UUID)
	if err != nil {
		err := errors.New("required")
		return "uuid", err
	}

	return "", nil
}

// Marshall returns the json encoding of UserRemoveRequest.
func (s UserRemoveRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// user handler
func (api *API) user(
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

	request := new(UserRequest)
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

	u, err := model.LoadUser(tx, request.UUID)
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := UserResponse{
		UUID:      u.UUID,
		Name:      u.Name,
		IsManager: u.IsManager,
		Devices:   u.Devices,
	}
	if u.IsManager {
		provisionURI, err := u.OTPProvisionURI(api.cfg.OTPIssuer)
		if err != nil {
			return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
		}
		response.Key = u.TFASecret
		response.ProvisionURI = provisionURI
	}

	return response.marshal(), nil
}

// UserRequest model.
type UserRequest struct {
	UUID string `json:"uuid"`
}

// Validate fields.
func (s *UserRequest) Validate() (string, error) {
	_, err := uuid.Parse(s.UUID)
	if err != nil {
		err := errors.New("required")
		return "uuid", err
	}

	return "", nil
}

// Marshall returns the json encoding of UserRequest.
func (s UserRequest) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// UserResponse model.
type UserResponse struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`
	IsManager    bool     `json:"is_manager"`
	Key          string   `json:"key,omitempty"`
	ProvisionURI string   `json:"provision_uri,omitempty"`
	Devices      []net.IP `json:"devices"`
}

func (s UserResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// userList handler
func (api *API) userList(
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

	users, err := model.LoadUsers(tx)
	if err != nil {
		return rpcError{Code: 500, Message: "bad gateway"}.marshal(), err
	}

	response := make(UserListResponse, len(users))
	for i, u := range users {
		response[i] = UserListItem{
			UUID:      u.UUID,
			Name:      u.Name,
			IsManager: u.IsManager,
			Devices:   u.Devices,
		}
	}

	return response.marshal(), nil
}

// UserListItem model.
type UserListItem struct {
	UUID      string   `json:"uuid"`
	Name      string   `json:"name"`
	IsManager bool     `json:"is_manager"`
	Devices   []net.IP `json:"devices"`
}

// UserListResponse model.
type UserListResponse []UserListItem

func (s UserListResponse) marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}
