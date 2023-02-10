package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"

	bolt "go.etcd.io/bbolt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"wgnetwork/pkg/ipcalc"
)

// Device model.
type Device struct {
	IPNetwork  IPNetwork    `json:"ipnetwork"`
	PubKey     wgtypes.Key  `json:"pub_key"`
	Label      string       `json:"label"`
	WANForward bool         `json:"wan_forward"`
	allowedIPs []*net.IPNet `json:"allowed_ips"`

	UserUUID string `json:"user_uuid"`
}

// NewDevice constructor
func NewDevice(
	ipnet IPNetwork,
	pk wgtypes.Key,
	label string,
	userUUID string,
	wanForward bool,
) Device {
	device := Device{
		IPNetwork: ipnet,
		PubKey:    pk,
		Label:     label,

		UserUUID:   userUUID,
		WANForward: wanForward,
	}

	device.IPNetwork.IP = device.IPNetwork.IP.To4()
	device.IPNetwork.Net.IP = device.IPNetwork.Net.IP.To4()

	return device
}

// LoadDevice constructor
func LoadDevice(tx *bolt.Tx, ip net.IP) (Device, error) {
	bname := []byte("devices")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return Device{}, errors.New("not found")
	}

	key := []byte(ip)
	v := bucket.Get(key)
	if v == nil {
		return Device{}, errors.New("not found")
	}

	d := Device{}
	err := json.Unmarshal(v, &d)
	if err != nil {
		return Device{}, err
	}

	d.IPNetwork.IP = d.IPNetwork.IP.To4()
	d.IPNetwork.Net.IP = d.IPNetwork.Net.IP.To4()

	for i := range d.allowedIPs {
		d.allowedIPs[i].IP = d.allowedIPs[i].IP.To4()
	}

	return d, nil
}

// CIDR for device.
func (d *Device) CIDR() *net.IPNet {
	return &net.IPNet{IP: d.IPNetwork.IP, Mask: net.IPv4Mask(255, 255, 255, 255)}
}

// Store to database.
func (d *Device) Store(tx *bolt.Tx) error {
	if !tx.Writable() {
		return errors.New("tx not writable")
	}

	bname := []byte("devices")
	bucket, err := tx.CreateBucketIfNotExists(bname)
	if err != nil {
		return err
	}

	key := []byte(d.IPNetwork.IP)
	value, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return bucket.Put(key, value)
}

// AllowedIPs for device
func (d *Device) AllowedIPs() []*net.IPNet {
	var ipnets []*net.IPNet

	if len(d.allowedIPs) > 0 {
		for i := range d.allowedIPs {
			ipnets = append(ipnets, d.allowedIPs[i])
		}
	} else if d.WANForward {
		ipnets = []*net.IPNet{
			&net.IPNet{
				IP:   net.IPv4(0, 0, 0, 0).To4(),
				Mask: net.IPv4Mask(0, 0, 0, 0)},
		}
	} else {
		ipnets = []*net.IPNet{
			&net.IPNet{
				IP:   d.IPNetwork.Net.IP,
				Mask: d.IPNetwork.Net.Mask},
		}
	}

	return ipnets
}

// RemoveDevice from database
func RemoveDevice(tx *bolt.Tx, ip net.IP) error {
	if !tx.Writable() {
		return errors.New("tx not writable")
	}

	bname := []byte("devices")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return nil
	}

	key := []byte(ip)
	return bucket.Delete(key)
}

// Devices type
type Devices []Device

// LoadDevices returns all devices from database.
func LoadDevices(tx *bolt.Tx) (Devices, error) {
	bname := []byte("devices")
	bucket := tx.Bucket(bname)
	if bucket == nil {
		return nil, nil
	}

	cnt := bucket.Stats().KeyN
	devices := make(Devices, cnt)

	i := 0
	c := bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		d := Device{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			return nil, err
		}

		d.IPNetwork.IP = d.IPNetwork.IP.To4()
		d.IPNetwork.Net.IP = d.IPNetwork.Net.IP.To4()

		for i := range d.allowedIPs {
			d.allowedIPs[i].IP = d.allowedIPs[i].IP.To4()
		}

		devices[i] = d
		i++
	}

	return devices, nil
}

// NFRule model.
type NFRule struct {
	IPProto   string   `json:"ipproto"`
	Direction string   `json:"direction"`
	SPorts    []uint16 `json:"sports"`
	DPorts    []uint16 `json:"dports"`
	SAddrs    []net.IP `json:"saddrs"`
	DAddrs    []net.IP `json:"daddrs"`
}

// IPNetwork model.
type IPNetwork struct {
	IP  net.IP     `json:"ip"`
	Net *net.IPNet `json:"net"`
}

// CIDR for ipnetwork.
func (ipn *IPNetwork) CIDR() *net.IPNet {
	return &net.IPNet{IP: ipn.IP, Mask: ipn.Net.Mask}
}

// AllocateIP for specified network and the public key.
func AllocateIP(
	tx *bolt.Tx,
	ipnet *net.IPNet,
	pk wgtypes.Key,
) (IPNetwork, error) {
	ips := ipcalc.NetworkHostIPs(ipnet, 0)
	if len(ips) < 2 {
		err := errors.New("network size to small")
		return IPNetwork{}, err
	}

	if !tx.Writable() {
		return IPNetwork{}, errors.New("tx not writable")
	}

	bname := []byte("devices")
	bucket, err := tx.CreateBucketIfNotExists(bname)
	if err != nil {
		return IPNetwork{}, err
	}

	var ip net.IP

	cnt := bucket.Stats().KeyN
	_, size := ipcalc.NetworkSize(ipnet)
	if cnt == int(size)-1 {
		err = errors.New("can't allocate ip, address space is full")
		return IPNetwork{}, err
	} else if cnt == 0 {
		if bytes.Compare(ips[0][:], ipnet.IP[:]) == 0 {
			ip = ips[1]
		} else {
			ip = ips[0]
		}
		return IPNetwork{IP: ip, Net: ipnet}, nil
	}

	devices, err := LoadDevices(tx)
	if err != nil {
		return IPNetwork{}, err
	}

	firstip := ips[0]
	if bytes.Compare(firstip[:], ipnet.IP[:]) == 0 {
		firstip = ips[1]
	}
	if bytes.Compare(firstip, devices[0].IPNetwork.IP[:]) < 0 {
		return IPNetwork{IP: firstip, Net: ipnet}, nil
	}

	for i := 0; i < len(devices); i++ {
		if bytes.Compare(pk[:], devices[i].PubKey[:]) == 0 {
			err = errors.New("public key already exists")
			return IPNetwork{}, err
		}

		if ip != nil {
			continue
		}

		next := ipcalc.NextIP(devices[i].IPNetwork.IP)
		if bytes.Compare(next[:], ipnet.IP[:]) == 0 {
			next = ipcalc.NextIP(next)
		}

		if i == len(devices)-1 {
			ip = next
		} else if bytes.Compare(next[:], devices[i+1].IPNetwork.IP[:]) < 0 {
			ip = next
		}
	}

	netaddr := ipcalc.FirstIP(ipnet)
	result := IPNetwork{IP: ip, Net: &net.IPNet{IP: netaddr, Mask: ipnet.Mask}}

	return result, nil
}
