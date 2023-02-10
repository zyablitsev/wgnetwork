package model

import (
	"fmt"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
)

func TestUser(t *testing.T) {
	dbpath := "test.db"
	db, err := bolt.Open(
		dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Errorf("can't open db: %v", err)
		return
	}
	defer db.Close()

	bname := []byte("users")
	err = deleteBucket(db, bname)
	if err != nil {
		t.Error(err)
		return
	}

	u, err := createUser(db, "john")
	if err != nil {
		t.Error(err)
		return
	}

	_, _, err = u.SetManager("test")
	if err != nil {
		t.Error(err)
		return
	}

	u, err = storeUser(db, u)
	if err != nil {
		t.Error(err)
		return
	}
}

func deleteBucket(db *bolt.DB, name []byte) error {
	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(name)
	if bucket == nil {
		return nil
	}

	err = tx.DeleteBucket(name)
	if err != nil {
		err = fmt.Errorf("can't delete bucket: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return err
	}

	return nil
}

func createUser(db *bolt.DB, name string) (User, error) {
	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	u := NewUser(name)
	err = u.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store user: %v", err)
		return User{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return User{}, err
	}

	return u, nil
}

func storeUser(db *bolt.DB, u User) (User, error) {
	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	err = u.Store(tx)
	if err != nil {
		err = fmt.Errorf("can't store user: %v", err)
		return User{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("can't commit tx: %v", err)
		return User{}, err
	}

	return u, nil
}
