package httpapi

import (
	"bytes"
	"mime"
	"net/http"
	"strings"
	"time"
)

const maxBodyBytes = int64(65536)

const (
	get  = "GET"
	post = "POST"
)

// API for http
type API struct {
	get  map[string]http.HandlerFunc
	post map[string]http.HandlerFunc
	log  Logger
}

// New return httprouter Handler
func New(log Logger) *API {
	api := &API{
		get:  make(map[string]http.HandlerFunc),
		post: make(map[string]http.HandlerFunc),
		log:  log,
	}
	return api
}

// ServeHTTP implements http.Handler interface
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

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

		notFound(w, r)
		return
	}

	var (
		handler http.HandlerFunc
		ok      bool
	)
	switch r.Method {
	case get:
		handler, ok = api.get[r.URL.Path]
	case post:
		handler, ok = api.post[r.URL.Path]
	}
	if !ok {
		notFound(w, r)
		return
	}

	// write request event log message
	api.log.Infof("got rest %q request from %s", r.RequestURI, r.RemoteAddr)

	// wrap ResponseWriter
	ww := &rw{
		status:         http.StatusOK,
		ResponseWriter: w,
		header:         make(http.Header),
	}
	// call the wrapped handler
	start := time.Now()
	handler.ServeHTTP(ww, r)
	elapsed := time.Since(start)
	w.Header().Set("X-Response-Time", elapsed.String())
	ww.WriteResponse()
}

// Register route handler
func (api *API) Register(method, pattern string, handler http.HandlerFunc) {
	routes := (map[string]http.HandlerFunc)(nil)
	switch method {
	case get:
		routes = api.get
	case post:
		routes = api.post
	default:
		return
	}
	if _, ok := routes[pattern]; !ok {
		routes[pattern] = handler
	} else {
		api.log.Infof("route %q exists already", pattern)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	handleError(w, "404 page not found", http.StatusNotFound)
}

func handleError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

type rw struct {
	status int
	body   bytes.Buffer
	header http.Header
	http.ResponseWriter
}

// Header returns & satisfies the http.ResponseWriter interface
func (w *rw) Header() http.Header {
	return w.header
}

// Write satisfies the http.ResponseWriter interface and
// captures data written, in bytes
func (w *rw) Write(data []byte) (int, error) {
	written, err := w.body.Write(data)
	return written, err
}

// WriteHeader satisfies the http.ResponseWriter interface and
// allows us to cach the status code
func (w *rw) WriteHeader(statusCode int) {
	w.status = statusCode
}

func (w *rw) WriteResponse() {
	header := w.ResponseWriter.Header()
	for k, v := range w.header {
		header[k] = v
	}
	w.ResponseWriter.WriteHeader(w.status)
	w.body.WriteTo(w.ResponseWriter)
}

// Logger desribes interface of log object
type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}
