package wgnetwork

import (
	"bytes"
	"embed"
	"fmt"
	"net"
	"net/http"
	"text/template"

	bolt "go.etcd.io/bbolt"

	"wgnetwork/model"
)

// Frontend struct.
type Frontend struct {
	log logger
	db  *bolt.DB

	index []byte
	js    []byte
	css   []byte

	devAuthIP string
}

// InitFrontend constructor.
func InitFrontend(
	log logger,
	db *bolt.DB,
	url, devAuthIP string,
) (*Frontend, error) {
	index, err := getFEIndex(url)
	if err != nil {
		return nil, err
	}

	js, css, err := getFEAssets()
	if err != nil {
		return nil, err
	}

	fe := &Frontend{
		log: log,
		db:  db,

		index: index,
		js:    js,
		css:   css,

		devAuthIP: devAuthIP,
	}

	return fe, nil
}

// ServeHTTP implements http.Handler interface
func (fe *Frontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check client ip
	var (
		ip  string
		err error
	)
	if fe.devAuthIP != "" {
		ip = fe.devAuthIP
	} else {
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fe.log.Errorf("can't parse remote ip, raw: %v", err)
			w.Header().Set("content-type", "text/html; charset=utf-8")
			response(w, []byte("not found"), http.StatusNotFound)
			return
		}
	}
	remoteIP := []byte(net.ParseIP(ip).To4())
	if remoteIP == nil {
		err = fmt.Errorf("bad remote ip: %s", string(remoteIP))
		fe.log.Errorf("can't parse remote ip, raw: %v", err)
		w.Header().Set("content-type", "text/html; charset=utf-8")
		response(w, []byte("not found"), http.StatusNotFound)
		return
	}

	tx, err := fe.db.Begin(false) // non-writeable tx
	if err != nil {
		fe.log.Errorf("can't begin tx: %v", err)
		w.Header().Set("content-type", "text/html; charset=utf-8")
		response(w, []byte("not found"), http.StatusNotFound)
		return
	}
	defer tx.Rollback()

	d, err := model.LoadDevice(tx, remoteIP)
	if err != nil {
		fe.log.Errorf("can't load device: %v", err)
		w.Header().Set("content-type", "text/html; charset=utf-8")
		response(w, []byte("not found"), http.StatusNotFound)
		return
	}

	u, err := model.LoadUser(tx, d.UserUUID)
	if err != nil {
		fe.log.Errorf("can't load user: %v", err)
		w.Header().Set("content-type", "text/html; charset=utf-8")
		response(w, []byte("not found"), http.StatusNotFound)
		return
	}

	if !u.IsManager {
		fe.log.Error("access forbidden for user: %s", u.Name)
		w.Header().Set("content-type", "text/html; charset=utf-8")
		response(w, []byte("not found"), http.StatusNotFound)
		return
	}

	// discard all not GET requests
	if r.Method != "GET" {
		fe.log.Errorf("method not allowed %s", r.Method)
		w.Header().Set("content-type", "text/html; charset=utf-8")
		response(w, []byte("not found"), http.StatusNotFound)
		return
	}

	if r.URL.Path == "/assets/index.js" {
		w.Header().Set("content-type", "text/javascript")
		response(w, fe.js, http.StatusOK)
		return
	} else if r.URL.Path == "/assets/index.css" {
		w.Header().Set("content-type", "text/css")
		response(w, fe.css, http.StatusOK)
		return
	}

	w.Header().Set("content-type", "text/html; charset=utf-8")
	response(w, fe.index, http.StatusOK)
}

func response(
	w http.ResponseWriter,
	b []byte, code int,
) {
	w.WriteHeader(code)
	w.Write(b)
}

//go:embed fe/dist/assets/index.js
//go:embed fe/dist/assets/index.css
var assets embed.FS

//go:embed fe/index.gohtml
var FEIndexTemplate string

func getFEIndex(url string) ([]byte, error) {
	tmpl, err := template.New("index.template").Parse(FEIndexTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	data := struct {
		APIUrl string
	}{
		APIUrl: url,
	}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getFEAssets() ([]byte, []byte, error) {
	var (
		js, css []byte
		err     error
	)

	js, err = assets.ReadFile("fe/dist/assets/index.js")
	if err != nil {
		return nil, nil, err
	}

	css, err = assets.ReadFile("fe/dist/assets/index.css")
	if err != nil {
		return nil, nil, err
	}

	return js, css, nil
}
