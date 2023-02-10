package ipset

import (
	"bytes"
	"net"
	"sort"
)

// IPSet object.
type IPSet struct {
	state map[string]net.IP

	removed map[string]net.IP
	added   map[string]net.IP

	frozen bool
}

// Copy of IPSet object.
func (s *IPSet) Copy() IPSet {
	cp := IPSet{
		state:  make(map[string]net.IP, len(s.state)),
		frozen: false,
	}

	for k, v := range s.state {
		cp.state[k] = make(net.IP, len(v))
		copy(cp.state[k], v)
	}

	return cp
}

// Replace object state with new set of ip addresses.
func (s *IPSet) Replace(p []net.IP) {
	if s.frozen {
		return
	}
	s.frozen = true

	prevState := make(map[string]net.IP, len(s.state))
	for k, v := range s.state {
		prevState[k] = make(net.IP, len(v))
		copy(prevState[k], v)
	}

	s.state = make(map[string]net.IP, len(p))

	s.added = make(map[string]net.IP)
	for _, v := range p {
		// set current state
		s.state[v.String()] = make(net.IP, len(v))
		copy(s.state[v.String()], v)

		_, ok := prevState[v.String()]
		if ok {
			continue
		}

		s.added[v.String()] = make(net.IP, len(v))
		copy(s.added[v.String()], v)
	}

	s.removed = make(map[string]net.IP)
	for k, v := range prevState {
		_, ok := s.state[k]
		if ok {
			continue
		}

		s.removed[k] = make(net.IP, len(v))
		copy(s.removed[k], v)
	}
}

// Added ips to the state after Replace applied.
func (s *IPSet) Added() []net.IP {
	if !s.frozen {
		return nil
	}

	ips := make([]net.IP, len(s.added))
	idx := 0
	for _, v := range s.added {
		ips[idx] = make(net.IP, len(v))
		copy(ips[idx], v)
		idx++
	}
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})

	return ips
}

// Removed ips from the state after Replace applied.
func (s *IPSet) Removed() []net.IP {
	if !s.frozen {
		return nil
	}

	ips := make([]net.IP, len(s.removed))
	idx := 0
	for _, v := range s.removed {
		ips[idx] = make(net.IP, len(v))
		copy(ips[idx], v)
		idx++
	}
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})

	return ips
}
