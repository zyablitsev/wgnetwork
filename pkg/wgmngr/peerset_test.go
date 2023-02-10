package wgmngr

import (
	"bytes"
	"net"
	"sort"
	"testing"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func TestPeerSet(t *testing.T) {
	peerset := PeerSet{}

	peers := Peers{
		Peer{
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.1").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			publicKey:    generatePrivateKey(),
			allowedIPs:   nil,
			endpointIP:   net.ParseIP("100.0.0.1").To4(),
			endpointPort: 51820,
		},
		Peer{
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.2").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			publicKey:    generatePrivateKey(),
			allowedIPs:   nil,
			endpointIP:   net.ParseIP("100.0.0.2").To4(),
			endpointPort: 51820,
		},
		Peer{
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.3").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			publicKey:    generatePrivateKey(),
			allowedIPs:   nil,
			endpointIP:   net.ParseIP("100.0.0.3").To4(),
			endpointPort: 51820,
		},
	}
	sort.Sort(peers)
	peerset.Replace(peers)

	removed := peerset.Removed()
	expectedAmount := 0
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	added := peerset.Added()
	expectedAmount = 3
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(added), expectedAmount)
		return
	}

	peerset = peerset.Copy()
	toRemove := peers[2:]

	peers = peers[:2]
	sort.Sort(peers)
	peerset.Replace(peers)

	removed = peerset.Removed()
	expectedAmount = len(toRemove)
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	for i := range removed {
		expected := toRemove[i].PublicKey()
		removed := removed[i].PublicKey()
		if bytes.Compare(removed[:], expected[:]) != 0 {
			t.Errorf("wrong removed %s, expected %s",
				removed, expected)
			return
		}
	}

	added = peerset.Added()
	expectedAmount = 0
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	peerset = peerset.Copy()
	toRemove = peers[:]

	peers = peers[:0]
	sort.Sort(peers)
	peerset.Replace(peers)

	removed = peerset.Removed()
	expectedAmount = len(toRemove)
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	for i := range removed {
		expected := toRemove[i].PublicKey()
		removed := removed[i].PublicKey()
		if bytes.Compare(removed[:], expected[:]) != 0 {
			t.Errorf("wrong removed %s, expected %s",
				removed, expected)
			return
		}
	}

	added = peerset.Added()
	expectedAmount = 0
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	peerset = peerset.Copy()

	peers = Peers{
		Peer{
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.1").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			publicKey:    generatePrivateKey(),
			allowedIPs:   nil,
			endpointIP:   net.ParseIP("100.0.0.1").To4(),
			endpointPort: 51820,
		},
		Peer{
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.2").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			publicKey:    generatePrivateKey(),
			allowedIPs:   nil,
			endpointIP:   net.ParseIP("100.0.0.2").To4(),
			endpointPort: 51820,
		},
		Peer{
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.3").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			publicKey:    generatePrivateKey(),
			allowedIPs:   nil,
			endpointIP:   net.ParseIP("100.0.0.3").To4(),
			endpointPort: 51820,
		},
	}
	sort.Sort(peers)
	peerset.Replace(peers)

	removed = peerset.Removed()
	expectedAmount = 0
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	added = peerset.Added()
	expectedAmount = len(peers)
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	for i := range added {
		expected := peers[i].PublicKey()
		added := added[i].PublicKey()
		if bytes.Compare(added[:], expected[:]) != 0 {
			t.Errorf("wrong added %s, expected %s",
				added, expected)
			return
		}
	}

	peerset = peerset.Copy()

	toAdd := Peers{
		Peer{
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.4").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			publicKey:    generatePrivateKey(),
			allowedIPs:   nil,
			endpointIP:   net.ParseIP("100.0.0.4").To4(),
			endpointPort: 51820,
		},
		Peer{
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.5").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			publicKey:    generatePrivateKey(),
			allowedIPs:   nil,
			endpointIP:   net.ParseIP("100.0.0.5").To4(),
			endpointPort: 51820,
		},
	}
	sort.Sort(toAdd)
	peers = append(peers, toAdd...)
	sort.Sort(peers)
	peerset.Replace(peers)

	removed = peerset.Removed()
	expectedAmount = 0
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	added = peerset.Added()
	expectedAmount = len(toAdd)
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	for i := range added {
		expected := toAdd[i].PublicKey()
		added := added[i].PublicKey()
		if bytes.Compare(added[:], expected[:]) != 0 {
			t.Errorf("wrong added %s, expected %s",
				added, expected)
			return
		}
	}

	peerset = peerset.Copy()

	toRemove = Peers{peers[0], peers[2], peers[4]}
	sort.Sort(toRemove)
	toAdd = Peers{
		Peer{
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.7").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			publicKey:    generatePrivateKey(),
			allowedIPs:   nil,
			endpointIP:   net.ParseIP("100.0.0.7").To4(),
			endpointPort: 51820,
		},
	}
	sort.Sort(toAdd)
	peers = Peers{peers[1], peers[3]}
	peers = append(peers, toAdd...)
	sort.Sort(peers)

	peerset.Replace(peers)

	removed = peerset.Removed()
	expectedAmount = len(toRemove)
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	for i := range removed {
		expected := toRemove[i].PublicKey()
		removed := removed[i].PublicKey()
		if bytes.Compare(removed[:], expected[:]) != 0 {
			t.Errorf("wrong removed %s, expected %s",
				removed, expected)
			return
		}
	}

	added = peerset.Added()
	expectedAmount = len(toAdd)
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	for i := range added {
		expected := toAdd[i].PublicKey()
		added := added[i].PublicKey()
		if bytes.Compare(added[:], expected[:]) != 0 {
			t.Errorf("wrong added %s, expected %s",
				added, expected)
			return
		}
	}
}

func generatePrivateKey() wgtypes.Key {
	k, _ := wgtypes.GeneratePrivateKey()
	return k
}
