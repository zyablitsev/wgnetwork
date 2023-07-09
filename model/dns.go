package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"sort"

	bolt "go.etcd.io/bbolt"
)

// Domain model.
type Domain struct {
	Name  string       `json:"name"`
	A     []ARecord    `json:"a,omitempty"`
	CNAME *CNAMERecord `json:"cname,omitempty"`
}

// NewDomain constructor.
func NewDomain(name string) Domain {
	d := Domain{Name: name}
	return d
}

// LoadDomain constructor
func LoadDomain(tx *bolt.Tx, name string) (Domain, error) {
	bname := []byte("dns")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return Domain{}, errors.New("not found")
	}

	key := []byte(name)
	v := bucket.Get(key)
	if v == nil {
		return Domain{}, errors.New("not found")
	}

	d := Domain{}
	err := json.Unmarshal(v, &d)
	if err != nil {
		return Domain{}, err
	}

	return d, nil
}

// SetA record.
func (d *Domain) SetA(r ARecord) {
	ip := r.A.To4()
	if ip == nil {
		return
	}

	if d.A == nil {
		d.A = make([]ARecord, 0, 1)
	}

	i, found := d.isAExists(ip)
	if found {
		d.A[i] = r
		return
	}

	d.A = append(d.A, r)
	d.sortA()

	d.CNAME = nil
}

// RemoveA record.
func (d *Domain) RemoveA(ip net.IP) {
	ip = ip.To4()
	if ip == nil {
		return
	}

	i, found := d.isAExists(ip)
	if !found {
		return
	}

	if len(d.A) == 1 {
		d.A = nil
		return
	}

	d.A[i] = d.A[len(d.A)-1]
	d.A[len(d.A)-1] = ARecord{}
	d.A = d.A[:len(d.A)-1]
	d.sortA()
}

// SetCNAME record.
func (d *Domain) SetCNAME(r CNAMERecord) {
	d.CNAME = &r
	d.A = nil
}

// RemoveCNAME record.
func (d *Domain) RemoveCNAME(target string) {
	d.CNAME = nil
}

// Store to database.
func (d *Domain) Store(tx *bolt.Tx) error {
	if !tx.Writable() {
		return errors.New("tx not writable")
	}

	bname := []byte("dns")
	bucket, err := tx.CreateBucketIfNotExists(bname)
	if err != nil {
		return err
	}

	key := []byte(d.Name)
	value, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return bucket.Put(key, value)
}

func (d *Domain) sortA() {
	sort.Slice(d.A, func(i, j int) bool {
		return bytes.Compare(d.A[i].A, d.A[j].A) < 0
	})
}

func (d *Domain) isAExists(ip net.IP) (int, bool) {
	i, found := sort.Find(len(d.A), func(i int) int {
		return bytes.Compare(ip, d.A[i].A)
	})

	return i, found
}

// RemoveDomain from database
func RemoveDomain(tx *bolt.Tx, name string) error {
	if !tx.Writable() {
		return errors.New("tx not writable")
	}

	bname := []byte("dns")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return errors.New("not found")
	}

	key := []byte(name)
	return bucket.Delete(key)
}

// Domains type
type Domains []Domain

// LoadDomains returls all domains fom database.
func LoadDomains(tx *bolt.Tx) (Domains, error) {
	bname := []byte("dns")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return nil, nil
	}

	cnt := bucket.Stats().KeyN
	domains := make(Domains, cnt)

	i := 0
	c := bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		d := Domain{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			return nil, err
		}

		domains[i] = d
		i++
	}

	return domains, nil
}

// SOARecord model.
type SOARecord struct {
	TTL     uint32 `json:"ttl"`
	Ns      string `json:"ns"`
	Mbox    string `json:"mbox"`
	Serial  uint32 `json:"serial"`
	Refresh uint32 `json:"refresh"`
	Retry   uint32 `json:"retry"`
	Expire  uint32 `json:"expire"`
	Minttl  uint32 `json:"minttl"`
}

// NSRecord model.
type NSRecord struct {
	TTL uint32 `json:"ttl"`
	Ns  string `json:"ns"`
}

// ARecord model.
type ARecord struct {
	TTL uint32 `json:"ttl"`
	A   net.IP `json:"a"`
}

// CNAMERecord model.
type CNAMERecord struct {
	TTL    uint32 `json:"ttl"`
	Target string `json:"target"`
}
