package model

import (
	"fmt"
	"net"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
)

func TestDomain(t *testing.T) {
	dbpath := "test.db"
	db, err := bolt.Open(
		dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Errorf("can't open db: %v", err)
		return
	}
	defer db.Close()

	bname := []byte("dns")
	err = deleteBucket(db, bname)
	if err != nil {
		t.Error(err)
		return
	}

	d, err := createDomain(db, "chat.wgnetwork.")
	if err != nil {
		t.Error(err)
		return
	}

	d.SetA(
		ARecord{TTL: 30, A: net.IPv4(172, 16, 0, 10)},
	)

	d, err = storeDomain(db, d)
	if err != nil {
		t.Error(err)
		return
	}
}

func createDomain(db *bolt.DB, name string) (Domain, error) {
	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return Domain{}, err
	}
	defer tx.Rollback()

	d := NewDomain(name)
	err = d.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store domain: %v", err)
		return Domain{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return Domain{}, err
	}

	return d, nil
}

func storeDomain(db *bolt.DB, d Domain) (Domain, error) {
	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return Domain{}, err
	}
	defer tx.Rollback()

	err = d.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store domain: %v", err)
		return Domain{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return Domain{}, err
	}

	return d, nil
}
