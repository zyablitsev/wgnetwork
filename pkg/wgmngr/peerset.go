package wgmngr

import "sort"

// PeerSet object
type PeerSet struct {
	state map[string]Peer

	removed map[string]Peer
	added   map[string]Peer

	frozen bool
}

// Copy of PeerSet object.
func (s *PeerSet) Copy() PeerSet {
	cp := PeerSet{
		state:  make(map[string]Peer, len(s.state)),
		frozen: false,
	}

	for k, v := range s.state {
		cp.state[k] = v.Copy()
	}

	return cp
}

// Replace object state with new set of peers.
func (s *PeerSet) Replace(p Peers) {
	if s.frozen {
		return
	}
	s.frozen = true

	prevState := make(map[string]Peer, len(s.state))
	for k, v := range s.state {
		prevState[k] = v.Copy()
	}

	s.state = make(map[string]Peer, len(p))

	s.added = make(map[string]Peer)
	for _, v := range p {
		// set current state
		publicKey := v.PublicKey().String()
		s.state[publicKey] = v.Copy()

		_, ok := prevState[publicKey]
		if ok {
			continue
		}

		s.added[publicKey] = v.Copy()
	}

	s.removed = make(map[string]Peer)
	for k, v := range prevState {
		_, ok := s.state[k]
		if ok {
			continue
		}

		s.removed[k] = v.Copy()
	}
}

// Added peers to the state after Replace applied.
func (s *PeerSet) Added() Peers {
	if !s.frozen {
		return nil
	}

	peers := make(Peers, len(s.added))
	idx := 0
	for _, v := range s.added {
		peers[idx] = v.Copy()
		idx++
	}
	sort.Sort(peers)

	return peers
}

// Removed peers from the state after Replace applied.
func (s *PeerSet) Removed() Peers {
	if !s.frozen {
		return nil
	}

	peers := make(Peers, len(s.removed))
	idx := 0
	for _, v := range s.removed {
		peers[idx] = v.Copy()
		idx++
	}
	sort.Sort(peers)

	return peers
}
