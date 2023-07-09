package wgnetwork

import (
	"net"

	bolt "go.etcd.io/bbolt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"wgnetwork/model"
	"wgnetwork/pkg/wgmngr"
)

func wgPrivateKey(db *bolt.DB) (wgtypes.Key, error) {
	tx, err := db.Begin(true) // writeable tx
	if err != nil {
		return wgtypes.Key{}, err
	}
	defer tx.Rollback()

	bname := []byte("wg")
	b, err := tx.CreateBucketIfNotExists(bname)
	if err != nil {
		return wgtypes.Key{}, err
	}

	var sk wgtypes.Key

	key := []byte("cfg")
	v := b.Get(key)
	if v == nil {
		sk, err = wgtypes.GeneratePrivateKey()
		if err == nil {
			err = b.Put(key, sk[:])
		}
	} else {
		sk, err = wgtypes.NewKey(v)
	}

	if err != nil {
		return wgtypes.Key{}, err
	}

	if err := tx.Commit(); err != nil {
		return wgtypes.Key{}, err
	}

	return sk, nil
}

func wgManagerIPs(users model.Users) []net.IP {
	if users == nil {
		return nil
	}

	ips := make([]net.IP, 0, len(users))
	for i := range users {
		if !users[i].IsManager {
			continue
		}

		for j := range users[i].Devices {
			ips = append(ips, users[i].Devices[j])
		}
	}

	return ips
}

func wgForwardWanIPs(devices model.Devices) []net.IP {
	if devices == nil {
		return nil
	}

	ips := make([]net.IP, len(devices))
	i := 0
	for j := range devices {
		if !devices[j].WANForward {
			continue
		}

		ips[i] = devices[j].IPNetwork.IP
		i++
	}

	return ips[:i]
}

func wgPeers(devices model.Devices) (wgmngr.Peers, error) {
	if devices == nil {
		return nil, nil
	}

	peers := make(wgmngr.Peers, len(devices))
	for i := 0; i < len(devices); i++ {
		peer, err := wgmngr.NewPeer(
			*devices[i].CIDR(),
			devices[i].PubKey,
			nil, // allowed ips
			nil, // endpoint ip
			0,   // endpoint port
			0)   // keep alive interval
		if err != nil {
			return nil, err
		}

		peers[i] = peer
	}

	return peers, nil
}
