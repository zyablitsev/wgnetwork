package resolver

import (
	"net"
	"sync"

	"github.com/miekg/dns"
	bolt "go.etcd.io/bbolt"

	"wgnetwork/model"
)

// Handler domain names service.
type Handler struct {
	log logger
	db  *bolt.DB
	rr  *roundrobin
	c   *dns.Client

	zone      string
	ns        string
	mbox      string
	wgIfaceIP net.IP

	m map[string]model.Domain

	sync.RWMutex
}

func New(
	log logger,
	db *bolt.DB,
	servers []string,
	zone string,
	ns string,
	mbox string,
	wgIfaceIP net.IP,
) *Handler {
	rr := &roundrobin{addrs: servers}
	c := new(dns.Client)
	m := map[string]model.Domain{}
	s := &Handler{
		log: log,
		db:  db,
		rr:  rr,
		c:   c,

		zone:      zone,
		ns:        ns,
		mbox:      mbox,
		wgIfaceIP: wgIfaceIP,

		m: m}

	return s
}

// Update set new domains data map.
func (s *Handler) Update(m map[string]model.Domain) {
	s.Lock()
	s.m = m
	s.Unlock()
}

// ServeDNS implements resolver interface.
func (s *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	var unknown = make([]dns.Question, 0, len(r.Question))
	var resolved = make([]dns.RR, 0, len(r.Question))
	var extra = make([]dns.RR, 0)
	for _, q := range r.Question {
		rr, ok := s.getBase(q)
		if ok {
			resolved = append(resolved, rr...)
			continue
		}

		domain, ok := s.getDomain(q.Name)
		if !ok {
			unknown = append(unknown, q)
			continue
		}

		if q.Qtype == dns.TypeSOA {
			resolved = append(resolved, s.rrSoa(q.Name))
			continue
		}
		if q.Qtype == dns.TypeNS {
			resolved = append(resolved, s.rrNs(q.Name))
			continue
		}
		if q.Qtype == dns.TypeCNAME {
			if domain.CNAME == nil {
				continue
			}
			cname := &dns.CNAME{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeCNAME,
					Class:  dns.ClassINET,
					Ttl:    domain.CNAME.TTL,
				},
				Target: domain.CNAME.Target,
			}
			resolved = append(resolved, cname)
			continue
		}
		if q.Qtype == dns.TypeA {
			if len(domain.A) > 0 {
				for _, v := range domain.A {
					a := &dns.A{
						Hdr: dns.RR_Header{
							Name:   q.Name,
							Rrtype: dns.TypeA,
							Class:  dns.ClassINET,
							Ttl:    v.TTL,
						},
						A: v.A,
					}
					resolved = append(resolved, a)
				}
				continue
			}

			if domain.CNAME != nil {
				cname := &dns.CNAME{
					Hdr: dns.RR_Header{
						Name:   q.Name,
						Rrtype: dns.TypeCNAME,
						Class:  dns.ClassINET,
						Ttl:    domain.CNAME.TTL,
					},
					Target: domain.CNAME.Target,
				}
				resolved = append(resolved, cname)

				name := domain.CNAME.Target
				for x := 0; x < 5; x++ {
					domain, ok := s.getDomain(name)
					if !ok {
						break
					}

					if domain.CNAME != nil {
						cname := &dns.CNAME{
							Hdr: dns.RR_Header{
								Name:   name,
								Rrtype: dns.TypeCNAME,
								Class:  dns.ClassINET,
								Ttl:    domain.CNAME.TTL,
							},
							Target: domain.CNAME.Target,
						}
						resolved = append(resolved, cname)
						extra = append(extra, cname)
						name = domain.CNAME.Target
						continue
					}
					for _, v := range domain.A {
						a := &dns.A{
							Hdr: dns.RR_Header{
								Name:   name,
								Rrtype: dns.TypeA,
								Class:  dns.ClassINET,
								Ttl:    v.TTL,
							},
							A: v.A,
						}
						resolved = append(resolved, a)
						extra = append(extra, a)
					}
					break
				}
			}
		}
	}

	if len(unknown) == 0 {
		result := &dns.Msg{}
		result.SetReply(r)
		result.MsgHdr.RecursionAvailable = true
		if len(resolved) > 0 {
			result.MsgHdr.Authoritative = true
		}
		result.Answer = resolved
		result.Ns = []dns.RR{s.rrNs(s.zone)}
		a := &dns.A{
			Hdr: dns.RR_Header{
				Name:   "server.wgn.",
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			A: s.wgIfaceIP,
		}
		result.Extra = []dns.RR{a}
		w.WriteMsg(result)
		return
	}

	r.Question = unknown
	result, _, err := s.c.Exchange(r, s.rr.get())
	if err != nil {
		s.log.Errorf("failed to resolve: %v", err)
		result = &dns.Msg{}
		result.SetReply(r)
	}

	w.WriteMsg(result)
	return
}

func (s *Handler) getBase(q dns.Question) ([]dns.RR, bool) {
	if q.Name != s.ns {
		return nil, false
	}

	if q.Qtype == dns.TypeSOA {
		return []dns.RR{s.rrSoa(q.Name)}, true
	} else if q.Qtype == dns.TypeA {
		a := &dns.A{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			A: s.wgIfaceIP,
		}
		return []dns.RR{a}, true
	}
	return []dns.RR{}, true
}

func (s *Handler) getDomain(name string) (model.Domain, bool) {
	s.RLock()
	domain, ok := s.m[name]
	s.RUnlock()
	return domain, ok
}

func (s *Handler) rrSoa(name string) dns.RR {
	soa := &dns.SOA{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeSOA,
			Class:  dns.ClassINET,
			Ttl:    3600,
		},
		Ns:      s.ns,
		Mbox:    s.mbox,
		Serial:  1,
		Refresh: 86400,
		Retry:   7200,
		Expire:  4000000,
		Minttl:  11200,
	}
	return soa
}

func (s *Handler) rrNs(name string) dns.RR {
	ns := &dns.NS{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeNS,
			Class:  dns.ClassINET,
			Ttl:    300,
		},
		Ns: s.ns,
	}
	return ns
}

type roundrobin struct {
	addrs []string
	idx   int

	sync.Mutex
}

func (rr *roundrobin) get() string {
	var addr string
	rr.Lock()
	i := rr.idx
	addr = rr.addrs[i]
	i += 1
	if i > len(rr.addrs)-1 {
		i = 0
	}
	rr.idx = i
	rr.Unlock()
	return addr
}

// logger desribes interface of log object.
type logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warning(...interface{})
	Warningf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}
