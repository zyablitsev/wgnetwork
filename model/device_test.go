package model

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
	"wgnetwork/pkg/ipcalc"

	bolt "go.etcd.io/bbolt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func TestDevice(t *testing.T) {
	dbpath := "test.db"
	db, err := bolt.Open(
		dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Errorf("can't open db: %v", err)
		return
	}
	defer db.Close()

	bname := []byte("devices")
	err = deleteBucket(db, bname)
	if err != nil {
		t.Error(err)
		return
	}

	ip, ipnet, err := ipcalc.ParseCIDR("172.16.0.1/29")
	if err != nil {
		t.Error(err)
		return
	}
	ipnet = &net.IPNet{IP: ip, Mask: ipnet.Mask}

	for i := 0; i < 5; i++ {
		sk, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			t.Error(err)
			return
		}
		pk := sk.PublicKey()

		d, err := createDevice(db, ipnet, pk, "server1", "")
		if err != nil {
			t.Error(err)
		}
		_ = d
	}

	tx, err := db.Begin(false) // non-writeable tx
	if err != nil {
		t.Error(err)
		return
	}
	bname = []byte("devices")
	bucket := tx.Bucket(bname)
	c := bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		d := Device{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			t.Error(err)
			return
		}
	}
	tx.Rollback()

	ips := ipcalc.NetworkHostIPs(ipnet, 0)
	err = removeDevice(db, ips[2])
	if err != nil {
		t.Error(err)
		return
	}

	tx, err = db.Begin(false) // non-writeable tx
	if err != nil {
		t.Error(err)
		return
	}
	bname = []byte("devices")
	bucket = tx.Bucket(bname)
	c = bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		d := Device{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			t.Error(err)
			return
		}
	}
	tx.Rollback()

	for i := 0; i < 1; i++ {
		sk, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			t.Error(err)
			return
		}
		pk := sk.PublicKey()

		d, err := createDevice(db, ipnet, pk, "server1", "")
		if err != nil {
			t.Error(err)
		}
		_ = d
	}

	tx, err = db.Begin(false) // non-writeable tx
	if err != nil {
		t.Error(err)
		return
	}
	bname = []byte("devices")
	bucket = tx.Bucket(bname)
	c = bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		d := Device{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			t.Error(err)
			return
		}
	}
	tx.Rollback()

	ips = ipcalc.NetworkHostIPs(ipnet, 0)
	err = removeDevice(db, ips[2])
	if err != nil {
		t.Error(err)
		return
	}
	err = removeDevice(db, ips[4])
	if err != nil {
		t.Error(err)
		return
	}

	tx, err = db.Begin(false) // non-writeable tx
	if err != nil {
		t.Error(err)
		return
	}
	bname = []byte("devices")
	bucket = tx.Bucket(bname)
	c = bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		d := Device{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			t.Error(err)
			return
		}
	}
	tx.Rollback()

	for i := 0; i < 2; i++ {
		sk, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			t.Error(err)
			return
		}
		pk := sk.PublicKey()

		d, err := createDevice(db, ipnet, pk, "server1", "")
		if err != nil {
			t.Error(err)
		}
		_ = d
	}

	tx, err = db.Begin(false) // non-writeable tx
	if err != nil {
		t.Error(err)
		return
	}
	bname = []byte("devices")
	bucket = tx.Bucket(bname)
	c = bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		d := Device{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			t.Error(err)
			return
		}
	}
	tx.Rollback()

	ips = ipcalc.NetworkHostIPs(ipnet, 0)
	err = removeDevice(db, ips[1])
	if err != nil {
		t.Error(err)
		return
	}
	err = removeDevice(db, ips[2])
	if err != nil {
		t.Error(err)
		return
	}
	err = removeDevice(db, ips[5])
	if err != nil {
		t.Error(err)
		return
	}

	tx, err = db.Begin(false) // non-writeable tx
	if err != nil {
		t.Error(err)
		return
	}
	bname = []byte("devices")
	bucket = tx.Bucket(bname)
	c = bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		d := Device{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			t.Error(err)
			return
		}
	}
	tx.Rollback()

	for i := 0; i < 2; i++ {
		sk, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			t.Error(err)
			return
		}
		pk := sk.PublicKey()

		d, err := createDevice(db, ipnet, pk, "server1", "")
		if err != nil {
			t.Error(err)
		}
		_ = d
	}

	tx, err = db.Begin(false) // non-writeable tx
	if err != nil {
		t.Error(err)
		return
	}
	bname = []byte("devices")
	bucket = tx.Bucket(bname)
	c = bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		d := Device{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			t.Error(err)
			return
		}
	}
	tx.Rollback()
}

func createDevice(
	db *bolt.DB,
	ipnet *net.IPNet,
	pk wgtypes.Key,
	label string,
	userUUID string,
) (Device, error) {
	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return Device{}, err
	}
	defer tx.Rollback()

	ipNetwork, err := AllocateIP(tx, ipnet, pk)
	if err != nil {
		return Device{}, err
	}

	d := Device{
		IPNetwork:  ipNetwork,
		PubKey:     pk,
		Label:      label,
		WANForward: false,

		UserUUID: userUUID}

	err = d.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store device: %v", err)
		return Device{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return Device{}, err
	}

	return d, nil
}

func storeDevice(db *bolt.DB, d Device) (Device, error) {
	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return Device{}, err
	}
	defer tx.Rollback()

	err = d.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store device: %v", err)
		return Device{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return Device{}, err
	}

	return d, nil
}

func removeDevice(db *bolt.DB, ip net.IP) error {
	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = RemoveDevice(tx, ip)
	if err != nil {
		err = fmt.Errorf("can't store device: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return err
	}

	return nil
}
