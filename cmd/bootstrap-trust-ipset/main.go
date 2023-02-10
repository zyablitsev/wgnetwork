package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"

	"wgnetwork/model"
	"wgnetwork/pkg/logger"
)

func main() {
	log, err := logger.New(os.Stdout, os.Stderr)
	if err != nil {
		panic(fmt.Errorf("can't init logger: %v", err))
	}

	// define args
	dbpathFlag := flag.String("dbpath", "", "dbpath")
	ipsFlag := stringsFlag{}
	flag.Var(&ipsFlag, "trustip", "device ip")

	// parse args
	flag.Parse()

	// validate args
	if dbpathFlag == nil {
		flag.Usage()

		return
	}

	dbpath := *dbpathFlag
	if len(dbpath) == 0 {
		err := errors.New("dbpath required")
		log.Errorf("failed to parse params: %v", err)
		flag.Usage()
		return
	}

	ips := make([]net.IP, len(ipsFlag.v))
	for i := 0; i < len(ipsFlag.v); i++ {
		ip := net.ParseIP(ipsFlag.v[i]).To4()
		if ip == nil {
			err = fmt.Errorf("bad ip address value %q", ipsFlag.v[i])
			log.Errorf("failed to parse params: %v", err)
			flag.Usage()
			return
		}

		ips[i] = ip
	}
	if len(ips) == 0 {
		err := errors.New("trustip required")
		log.Errorf("failed to parse params: %v", err)
		flag.Usage()
		return
	}

	// execute
	err = execute(log, dbpath, ips)
	if err != nil {
		log.Errorf("failed to execute: %v", err)
		os.Exit(1)
	}
}

func execute(log *logger.Logger, dbpath string, ips []net.IP) error {
	db, err := bolt.Open(
		dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	// if initialisation complete on error, gracefully close db connection
	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ipset, err := model.LoadManagerSSHTrustIPSet(tx)
	if err != nil {
		return err
	}

	for i := 0; i < len(ips); i++ {
		ipset.Add(ips[i])
		os.Stdout.WriteString("added: " + ips[i].String() + "\n")
	}

	err = ipset.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store ipset: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return err
	}

	log.Debug("done")

	return nil
}

type stringsFlag struct {
	v []string
}

func (f *stringsFlag) String() string {
	return strings.Join(f.v, ", ")
}

func (f *stringsFlag) Set(s string) error {
	f.v = append(f.v, s)
	return nil
}
