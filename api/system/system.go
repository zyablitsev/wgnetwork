package system

import (
	"context"
	"encoding/json"
	"net/http"

	"wgnetwork/pkg/httpapi"
)

// API object.
type API struct {
	ctx context.Context
	log logger
}

// New constructor for API.
func New(
	ctx context.Context,
	log logger,
) *API {
	api := &API{ctx: ctx, log: log}
	return api
}

// RegisterHandlers on provided api handler.
func (api *API) RegisterHandlers(http *httpapi.API) {
	http.Register("GET", "/heartbeat", api.heartbeat)
	http.Register("POST", "/heartbeat", api.heartbeat)
}

// heartbeat handler
func (api *API) heartbeat(w http.ResponseWriter, r *http.Request) {
	// write successful response
	api.httpOk(w, r.RequestURI)
}

// httpOk sends {200, "OK"} default success response to remote party
func (api *API) httpOk(w http.ResponseWriter, requestURI string) {
	response := &httpResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	}
	api.httpWriteJSON(w, response, http.StatusOK)
}

// httpError is an error handler to call in "if err!=nil" blocks
// writes rest {code, msg} answer if w is not nil
func (api *API) httpError(w http.ResponseWriter, rURI string, httpCode int, err error) {
	api.log.Errorf("%s %s: %v", rURI, http.StatusText(httpCode), err)
	if w == nil {
		return
	}
	response := &httpResponse{
		Code:    httpCode,
		Message: http.StatusText(httpCode),
	}
	api.httpWriteJSON(w, response, httpCode)
}

// writeJSON to response
func (api *API) httpWriteJSON(w http.ResponseWriter, v interface{}, httpCode int) {
	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := json.Marshal(v)
	if err == nil {
		_, err = w.Write(b)
	}
	if err != nil {
		api.log.Error(err)
	}
}

// httpResponse is http operation status response format
type httpResponse struct {
	Code        int                               `json:"code"`
	Message     string                            `json:"message"`
	Description map[string]map[string]interface{} `json:"description,omitempty"`
	Checksum    *string                           `json:"checksum,omitempty"`
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
