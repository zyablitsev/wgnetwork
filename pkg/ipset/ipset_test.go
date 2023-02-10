package ipset

import (
	"bytes"
	"net"
	"testing"
)

func TestIPSet(t *testing.T) {
	ipset := IPSet{}

	ips := []net.IP{
		net.ParseIP("10.0.0.1").To4(),
		net.ParseIP("10.0.0.2").To4(),
		net.ParseIP("10.0.0.3").To4(),
	}
	ipset.Replace(ips)

	removed := ipset.Removed()
	expectedAmount := 0
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	added := ipset.Added()
	expectedAmount = 3
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(added), expectedAmount)
		return
	}

	ipset = ipset.Copy()

	ips = []net.IP{
		net.ParseIP("10.0.0.1").To4(),
		net.ParseIP("10.0.0.2").To4(),
	}
	ipset.Replace(ips)

	removed = ipset.Removed()
	expectedAmount = 1
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	expectedIP := net.ParseIP("10.0.0.3").To4()
	if bytes.Compare(removed[0], expectedIP) != 0 {
		t.Errorf("wrong removed ip %s, expected %s",
			removed[0], expectedIP)
		return
	}

	added = ipset.Added()
	expectedAmount = 0
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	ipset = ipset.Copy()

	ips = []net.IP{}
	ipset.Replace(ips)

	removed = ipset.Removed()
	expectedAmount = 2
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	expectedIP = net.ParseIP("10.0.0.1").To4()
	if bytes.Compare(removed[0], expectedIP) != 0 {
		t.Errorf("wrong removed ip %s, expected %s",
			removed[0], expectedIP)
		return
	}

	expectedIP = net.ParseIP("10.0.0.2").To4()
	if bytes.Compare(removed[1], expectedIP) != 0 {
		t.Errorf("wrong removed ip %s, expected %s",
			removed[1], expectedIP)
		return
	}

	added = ipset.Added()
	expectedAmount = 0
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(added), expectedAmount)
		return
	}

	ipset = ipset.Copy()

	ips = []net.IP{
		net.ParseIP("10.0.0.1").To4(),
		net.ParseIP("10.0.0.2").To4(),
		net.ParseIP("10.0.0.3").To4(),
	}
	ipset.Replace(ips)

	removed = ipset.Removed()
	expectedAmount = 0
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	added = ipset.Added()
	expectedAmount = 3
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(added), expectedAmount)
		return
	}

	expectedIP = net.ParseIP("10.0.0.1").To4()
	if bytes.Compare(added[0], expectedIP) != 0 {
		t.Errorf("wrong added ip %s, expected %s",
			added[0], expectedIP)
		return
	}

	expectedIP = net.ParseIP("10.0.0.2").To4()
	if bytes.Compare(added[1], expectedIP) != 0 {
		t.Errorf("wrong added ip %s, expected %s",
			added[1], expectedIP)
		return
	}

	expectedIP = net.ParseIP("10.0.0.3").To4()
	if bytes.Compare(added[2], expectedIP) != 0 {
		t.Errorf("wrong added ip %s, expected %s",
			added[2], expectedIP)
		return
	}

	ipset = ipset.Copy()

	ips = []net.IP{
		net.ParseIP("10.0.0.1").To4(),
		net.ParseIP("10.0.0.2").To4(),
		net.ParseIP("10.0.0.3").To4(),
		net.ParseIP("10.0.0.4").To4(),
		net.ParseIP("10.0.0.5").To4(),
	}
	ipset.Replace(ips)

	removed = ipset.Removed()
	expectedAmount = 0
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	added = ipset.Added()
	expectedAmount = 2
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	expectedIP = net.ParseIP("10.0.0.4").To4()
	if bytes.Compare(added[0], expectedIP) != 0 {
		t.Errorf("wrong added ip %s, expected %s",
			added[0], expectedIP)
		return
	}

	expectedIP = net.ParseIP("10.0.0.5").To4()
	if bytes.Compare(added[1], expectedIP) != 0 {
		t.Errorf("wrong added ip %s, expected %s",
			added[1], expectedIP)
		return
	}

	ipset = ipset.Copy()

	ips = []net.IP{
		net.ParseIP("10.0.0.2").To4(),
		net.ParseIP("10.0.0.4").To4(),
		net.ParseIP("10.0.0.7").To4(),
	}
	ipset.Replace(ips)

	removed = ipset.Removed()
	expectedAmount = 3
	if len(removed) != expectedAmount {
		t.Errorf("wrong removed amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	expectedIP = net.ParseIP("10.0.0.1").To4()
	if bytes.Compare(removed[0], expectedIP) != 0 {
		t.Errorf("wrong removed ip %s, expected %s",
			removed[0], expectedIP)
		return
	}

	expectedIP = net.ParseIP("10.0.0.3").To4()
	if bytes.Compare(removed[1], expectedIP) != 0 {
		t.Errorf("wrong removed ip %s, expected %s",
			removed[1], expectedIP)
		return
	}

	expectedIP = net.ParseIP("10.0.0.5").To4()
	if bytes.Compare(removed[2], expectedIP) != 0 {
		t.Errorf("wrong removed ip %s, expected %s",
			removed[2], expectedIP)
		return
	}

	added = ipset.Added()
	expectedAmount = 1
	if len(added) != expectedAmount {
		t.Errorf("wrong added amount %d, expected %d",
			len(removed), expectedAmount)
		return
	}

	expectedIP = net.ParseIP("10.0.0.7").To4()
	if bytes.Compare(added[0], expectedIP) != 0 {
		t.Errorf("wrong added ip %s, expected %s",
			added[0], expectedIP)
		return
	}
}
