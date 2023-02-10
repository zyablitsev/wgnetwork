package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"sort"

	bolt "go.etcd.io/bbolt"
)

// ManagerSSHTrustIPSet model.
type ManagerSSHTrustIPSet []net.IP

func (ips *ManagerSSHTrustIPSet) Add(ip net.IP) {
	if ips == nil {
		return
	}

	ip = ip.To4()
	if ip == nil {
		return
	}

	_, found := ips.isExists(ip)
	if found {
		return
	}

	*ips = append(*ips, ip)
	ips.sort()
}

func (ips *ManagerSSHTrustIPSet) Remove(ip net.IP) {
	if ips == nil {
		return
	}

	ip = ip.To4()
	if ip == nil {
		return
	}

	i, found := ips.isExists(ip)
	if !found {
		return
	}

	(*ips)[i] = (*ips)[len(*ips)-1]
	(*ips)[len(*ips)-1] = nil
	(*ips) = (*ips)[:len(*ips)-1]
	ips.sort()
}

func (ips ManagerSSHTrustIPSet) sort() {
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})
}

func (ips ManagerSSHTrustIPSet) isExists(ip net.IP) (int, bool) {
	i, found := sort.Find(len(ips), func(i int) int {
		return bytes.Compare(ip, ips[i])
	})

	return i, found
}

func (ips ManagerSSHTrustIPSet) Store(tx *bolt.Tx) error {
	if !tx.Writable() {
		return errors.New("tx not writable")
	}

	bname := []byte("managers")
	bucket, err := tx.CreateBucketIfNotExists(bname)
	if err != nil {
		return err
	}

	key := []byte("ssh_trust_ipset")
	value, err := json.Marshal(ips)
	if err != nil {
		return err
	}

	return bucket.Put(key, value)
}

// LoadManagerSSHTrustIPSet returns all ips from database.
func LoadManagerSSHTrustIPSet(tx *bolt.Tx) (ManagerSSHTrustIPSet, error) {
	bname := []byte("managers")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return ManagerSSHTrustIPSet{}, nil
	}

	key := []byte("ssh_trust_ipset")
	v := bucket.Get(key)
	if v == nil {
		return ManagerSSHTrustIPSet{}, nil
	}

	ips := []net.IP{}
	err := json.Unmarshal(v, &ips)
	if err != nil {
		return nil, err
	}

	ipset := make(ManagerSSHTrustIPSet, len(ips))
	for i := range ips {
		ipset[i] = ips[i].To4()
	}

	return ipset, nil
}
