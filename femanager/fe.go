package femanager

import (
	"fmt"
	"net"
	"net/http"

	bolt "go.etcd.io/bbolt"

	"wgnetwork/model"
)

// Frontend struct.
type Frontend struct {
	log       logger
	db        *bolt.DB
	index     []byte
	css       []byte
	js        []byte
	devAuthIP string
}

// Init constructor.
func Init(
	log logger,
	db *bolt.DB,
	apiUrl, devAuthIP string,
) (*Frontend, error) {
	index, err := loadIndex(apiUrl)
	if err != nil {
		return nil, err
	}

	css, js, err := loadAssets()
	if err != nil {
		return nil, err
	}

	fe := &Frontend{
		log:       log,
		db:        db,
		index:     index,
		css:       css,
		js:        js,
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
