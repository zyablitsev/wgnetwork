package rpcapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"strings"
)

const maxBodyBytes = int64(65536)

// API for http rpc (json-rpc)
type API struct {
	handlers map[string]HandlerFunc
	origin   string
	log      Logger

	devAuthIP string
}

// New return httprouter Handler
func New(log Logger, origin, devAuthIP string) *API {
	api := &API{
		handlers: make(map[string]HandlerFunc),
		origin:   origin,
		log:      log,

		devAuthIP: devAuthIP,
	}
	return api
}

// ServeHTTP implements http.Handler interface
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle preflight cors request
	if r.Method == "OPTIONS" {
		api.commonHeaders(w)
		handleOk(w, nil, json.RawMessage(`{}`))
		return
	}

	// check content-type
	charset := "utf-8"
	mediatype, params, err := mime.ParseMediaType(r.Header.Get("content-type"))
	if err == nil {
		if v, ok := params["charset"]; ok {
			charset = strings.ToLower(v)
		}
	}
	// per net/http doc, means that the length is known and non-null
	if err != nil ||
		(r.ContentLength > 0 &&
			!(strings.Compare(mediatype, "application/json") == 0 &&
				strings.Compare(charset, "utf-8") == 0)) {

		rawResponse(w, "404 page not found", http.StatusNotFound)
		return
	}

	api.commonHeaders(w)

	// discard all not POST requests
	if r.Method != "POST" {
		api.log.Errorf("method not allowed %s", r.Method)
		handleError(w, nil, json.RawMessage(`{"msg": "method not allowed"}`))
		return
	}

	// unpack message
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if r == nil || r.Body == nil {
		api.log.Error("can't unpack message: bad request")
		handleError(w, nil, json.RawMessage(`{"msg": "bad request"}`))
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		api.log.Errorf("can't unpack message: bad request, raw: %v", err)
		handleError(w, nil, json.RawMessage(`{"msg": "bad request"}`))
		return
	}

	request := new(Request)
	err = json.Unmarshal(content, &request)
	if err != nil {
		api.log.Errorf("failed unmarshal request: %v", err)
		handleError(w, nil, json.RawMessage(`{"msg": "bad request"}`))
		return
	}

	handler, ok := api.handlers[request.Method]
	if !ok {
		api.log.Errorf("method %q not found: %v", request.Method, err)
		handleError(w, request, json.RawMessage(`{"msg": "method not found"}`))
		return
	}

	api.log.Infof("got rpc %q request from %q", request.Method, r.RemoteAddr)

	var ip string
	if api.devAuthIP != "" {
		ip = api.devAuthIP
	} else {
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			api.log.Errorf("can't parse remote ip, raw: %v", err)
			handleError(w, nil, json.RawMessage(`{"msg": "bad request"}`))
			return
		}
	}
	remoteIP := []byte(net.ParseIP(ip).To4())
	if remoteIP == nil {
		err = fmt.Errorf("bad remote ip: %s", string(remoteIP))
		api.log.Errorf("can't parse remote ip, raw: %v", err)
		handleError(w, nil, json.RawMessage(`{"msg": "bad request"}`))
		return
	}

	p := Params{
		RemoteAddr: remoteIP,
		Headers:    r.Header,
	}
	ctx := NewContext(r.Context(), p)

	result, err := handler(ctx, w, request.Params)
	if err != nil {
		api.log.Errorf("failed handle request %s, raw: %v", request, err)
		handleError(w, request, result)
		return
	}

	handleOk(w, request, result)
}

func (api *API) commonHeaders(w http.ResponseWriter) {
	w.Header().Set("access-control-allow-origin", api.origin)
	w.Header().Set("access-control-allow-methods", "POST")
	w.Header().Set("access-control-allow-headers", "content-type, x-session")
	w.Header().Set("access-control-allow-credentials", "true")
	w.Header().Set("origin", api.origin)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as RPC handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(
	ctx context.Context,
	w http.ResponseWriter,
	params json.RawMessage,
) (json.RawMessage, error)

// Register route handler
func (api *API) Register(pattern string, handler HandlerFunc) {
	if _, ok := api.handlers[pattern]; !ok {
		api.handlers[pattern] = handler
	} else {
		api.log.Warningf("method %q exists already", pattern)
	}
}

func handleError(w http.ResponseWriter, r *Request, b json.RawMessage) {
	response := &Response{
		Error: b,
	}

	b, _ = json.Marshal(response)
	w.Header().Set("content-type", "application/json")
	w.Write(b)
}

func handleOk(w http.ResponseWriter, r *Request, b json.RawMessage) {
	response := &Response{
		Result: b,
	}

	b, _ = json.Marshal(response)
	w.Header().Set("content-type", "application/json")
	w.Write(b)
}

func rawResponse(
	w http.ResponseWriter,
	msg string, code int,
) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

// Request type
type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

func (s *Request) String() string {
	return fmt.Sprintf("method: %q, params: %v",
		s.Method, string(s.Params))
}

// Marshall returns the json encoding of Request.
func (s Request) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// Response type
type Response struct {
	Result json.RawMessage `json:"result,omitempty"`
	Error  json.RawMessage `json:"error,omitempty"`
}

// Params contains the headers and remote ip address
type Params struct {
	RemoteAddr []byte
	Headers    http.Header
}

type paramsKey struct{}

// NewContext creates a new context with request http.Header value attached.
func NewContext(ctx context.Context, p Params) context.Context {
	return context.WithValue(ctx, paramsKey{}, p)
}

// FromContext returns the http.Header value in ctx if it exists.
func FromContext(ctx context.Context) (p Params, ok bool) {
	p, ok = ctx.Value(paramsKey{}).(Params)
	return
}

// Logger desribes interface of log object
type Logger interface {
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
}
