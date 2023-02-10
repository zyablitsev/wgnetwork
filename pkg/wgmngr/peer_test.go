package wgmngr

import (
	"bytes"
	"net"
	"testing"
)

func TestPeer(t *testing.T) {
	peer := Peer{
		ipnet: net.IPNet{
			IP:   net.ParseIP("10.0.0.1").To4(),
			Mask: net.IPv4Mask(255, 255, 255, 255),
		},
		publicKey: generatePrivateKey(),
		allowedIPs: []net.IPNet{
			{
				IP:   net.ParseIP("172.16.0.10").To4(),
				Mask: net.IPv4Mask(255, 255, 255, 240),
			},
		},
		endpointIP:   net.ParseIP("100.0.0.1").To4(),
		endpointPort: 51820,
	}

	allowedIPs := peer.AllowedIPs()

	expectedAmount := 2
	if len(allowedIPs) != expectedAmount {
		t.Errorf("wrong allowed ips amount %d expected %d",
			len(allowedIPs), expectedAmount)
		return
	}

	expectedNetworkIP := net.ParseIP("10.0.0.1").To4()
	expectedNetworkMask := net.IPv4Mask(255, 255, 255, 255)

	if bytes.Compare(allowedIPs[0].IP, expectedNetworkIP) != 0 {
		t.Errorf("wrong network ip %q, expected %q",
			allowedIPs[0].IP, expectedNetworkIP)
	}
	if bytes.Compare(allowedIPs[0].Mask, expectedNetworkMask) != 0 {
		t.Errorf("wrong network ip %q, expected %q",
			allowedIPs[0].Mask, expectedNetworkIP)
	}

	expectedNetworkIP = net.ParseIP("172.16.0.10").To4()
	expectedNetworkMask = net.IPv4Mask(255, 255, 255, 240)

	if bytes.Compare(allowedIPs[1].IP, expectedNetworkIP) != 0 {
		t.Errorf("wrong network ip %q, expected %q",
			allowedIPs[1].IP, expectedNetworkIP)
	}
	if bytes.Compare(allowedIPs[1].Mask, expectedNetworkMask) != 0 {
		t.Errorf("wrong network ip %q, expected %q",
			allowedIPs[1].Mask, expectedNetworkIP)
	}
}
