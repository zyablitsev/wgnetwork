package wgnetwork

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/miekg/dns"
	bolt "go.etcd.io/bbolt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"wgnetwork/api/auth"
	"wgnetwork/api/manager"
	"wgnetwork/api/system"
	"wgnetwork/femanager"
	"wgnetwork/firewall"
	"wgnetwork/model"
	"wgnetwork/pkg/httpapi"
	"wgnetwork/pkg/iface"
	"wgnetwork/pkg/ipset"
	"wgnetwork/pkg/log"
	"wgnetwork/pkg/rpcapi"
	"wgnetwork/pkg/wgmngr"
	"wgnetwork/resolver"
)

// Service object.
type Service struct {
	ctx context.Context
	cfg config
	log logger

	db  *bolt.DB
	nft *firewall.NFTables

	trustIPSet        ipset.IPSet
	wgManagerIPSet    ipset.IPSet
	wgForwardWanIPSet ipset.IPSet

	wgm     *wgmngr.Manager
	wgpeers wgmngr.PeerSet

	resolver *resolver.Handler
}

// Init service.
func Init(ctx context.Context) (*Service, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("can't load configuration: %v", err)
	}

	log, err := log.New(cfg.LogLevel, os.Stdout, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("can't init logger: %v", err)
	}

	log.Debugf("opening bolt db path %q…", cfg.DBPath)
	var db *bolt.DB
	db, err = bolt.Open(
		cfg.DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	// if initialisation complete on error, gracefully close db connection
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	log.Info("opened bolt db")

	// generate wireguard private key for members network
	var wgsk wgtypes.Key
	wgsk, err = wgPrivateKey(db)
	if err != nil {
		return nil, err
	}

	const wgLinkType = "wireguard"
	err = iface.Create(
		log,
		cfg.WGIface, wgLinkType, cfg.wgIfaceIP, cfg.wgIfaceIPNet)
	if err != nil {
		return nil, err
	}

	var wgm *wgmngr.Manager
	wgm, err = wgmngr.NewManager(
		wgsk, cfg.WGIface, cfg.WGPort)
	if err != nil {
		return nil, err
	}

	var (
		nft *firewall.NFTables

		managerPorts []uint16 = []uint16{
			uint16(cfg.APIHTTPPort), uint16(cfg.FEHTTPPort)}
	)
	nftCfg := firewall.Config{
		Enabled:          cfg.NFTEnabled,
		NetworkNamespace: cfg.NFTNetworkNamespace,
		DefaultPolicy:    cfg.NFTDefaultPolicy,
		WGIface:          cfg.WGIface,
		WGPort:           cfg.WGPort,
		WGIPNet:          cfg.wgIfaceIPNet,
		Ifaces:           cfg.NFTIfaces,
		TrustPorts:       cfg.NFTTrustPorts,
	}
	nft, err = firewall.Init(nftCfg, managerPorts)
	if err != nil {
		return nil, err
	}

	ns := fmt.Sprintf("server.%s", cfg.DNSZone)
	mbox := fmt.Sprintf("hostmaster.server.%s", cfg.DNSZone)
	resolver := resolver.New(
		log, db,
		cfg.DNSResolverAddrs, cfg.DNSZone,
		ns, mbox,
		cfg.wgIfaceIP)

	s := &Service{
		ctx: ctx,
		cfg: cfg,
		log: log,
		db:  db,
		nft: nft,

		trustIPSet:        ipset.IPSet{},
		wgManagerIPSet:    ipset.IPSet{},
		wgForwardWanIPSet: ipset.IPSet{},

		wgm:     wgm,
		wgpeers: wgmngr.PeerSet{},

		resolver: resolver,
	}

	return s, nil
}

// Run service.
func (s *Service) Run() {
	// init state
	err := s.refresh()
	if err != nil {
		s.log.Error(err)
		return
	}

	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	// prepare mux
	apiTcpMux := s.apiTcpMux(ctx)
	apiUnixMux := s.apiUnixMux(ctx)
	feTcpMux, err := s.feTcpMux(ctx)
	if err != nil {
		s.log.Error(err)
		return
	}

	var (
		lc            *net.ListenConfig
		listenDnsTcp  net.Listener
		listenDnsUdp  net.PacketConn
		listenApiTcp  net.Listener
		listenApiUnix net.Listener
		listenFeTcp   net.Listener
		wg            sync.WaitGroup

		dnsTcp *dns.Server
		dnsUdp *dns.Server

		apiTcp  *http.Server
		apiUnix *http.Server
		feTcp   *http.Server
	)

	// run dns tcp
	lc = &net.ListenConfig{}
	listenDnsTcp, err = lc.Listen(ctx, "tcp", s.cfg.dnsTcpAddr)
	if err != nil {
		s.log.Error(err)
		return
	}
	defer listenDnsTcp.Close()
	dnsTcp = &dns.Server{
		Listener:     listenDnsTcp,
		Handler:      s.resolver,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	wg.Add(1)
	go func() {
		s.log.Info("dns tcp socket serve running…")
		err = dnsTcp.ActivateAndServe()
		if err != nil {
			s.log.Errorf("dns tcp serve has failed %v", err)
		} else {
			s.log.Error("dns tcp serve has stopped")
		}
		cancel()
		wg.Done()
	}()

	// run dns udp
	lc = &net.ListenConfig{}
	listenDnsUdp, err = lc.ListenPacket(ctx, "udp", s.cfg.dnsUdpAddr)
	if err != nil {
		s.log.Error(err)
		return
	}
	defer listenDnsUdp.Close()
	dnsUdp = &dns.Server{
		PacketConn:   listenDnsUdp,
		Handler:      s.resolver,
		UDPSize:      512,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	wg.Add(1)
	go func() {
		s.log.Info("dns udp socket serve running…")
		err = dnsUdp.ActivateAndServe()
		if err != nil {
			s.log.Errorf("dns tcp serve has failed %v", err)
		} else {
			s.log.Error("dns tcp serve has stopped")
		}
		cancel()
		wg.Done()
	}()

	// run api tcp
	lc = &net.ListenConfig{}
	listenApiTcp, err = lc.Listen(ctx, "tcp", s.cfg.apiHTTPAddr)
	if err != nil {
		s.log.Errorf("api tcp socket bind failed: %v", err)
		return
	}
	defer listenApiTcp.Close()
	apiTcp = &http.Server{
		Addr:              s.cfg.apiHTTPAddr,
		Handler:           apiTcpMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	wg.Add(1)
	go func() {
		s.log.Infof("serving tcp socket api on %q…", s.cfg.apiHTTPAddr)
		err = apiTcp.Serve(listenApiTcp)
		if err != http.ErrServerClosed {
			s.log.Errorf("api tcp socket serve failed: %v", err)
		} else {
			s.log.Info("api tcp socket serve stopped")
		}
		cancel()
		wg.Done()
	}()

	// run api unix
	lc = &net.ListenConfig{}
	listenApiUnix, err = lc.Listen(ctx, "unix", s.cfg.APIUnixSocket)
	if err != nil {
		s.log.Errorf("api unix socket bind failed: %v", err)
		return
	}
	defer listenApiUnix.Close()
	apiUnix = &http.Server{
		Handler:           apiUnixMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	wg.Add(1)
	go func() {
		s.log.Infof("serving unix socket api on %q…", s.cfg.APIUnixSocket)
		err = apiUnix.Serve(listenApiUnix)
		if err != http.ErrServerClosed {
			s.log.Errorf("api unix socket serve failed: %v", err)
		} else {
			s.log.Info("api unix socket serve stopped")
		}
		cancel()
		wg.Done()
	}()

	// run fe tcp
	lc = &net.ListenConfig{}
	listenFeTcp, err = lc.Listen(ctx, "tcp", s.cfg.feHTTPAddr)
	if err != nil {
		s.log.Errorf("fe tcp socket bind failed: %v", err)
		return
	}
	defer listenFeTcp.Close()
	feTcp = &http.Server{
		Addr:              s.cfg.feHTTPAddr,
		Handler:           feTcpMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	wg.Add(1)
	go func() {
		s.log.Infof("serving tcp socket fe on %q…", s.cfg.feHTTPAddr)
		err = feTcp.Serve(listenFeTcp)
		if err != http.ErrServerClosed {
			s.log.Errorf("fe tcp socket serve failed: %v", err)
		} else {
			s.log.Info("fe tcp socket serve stopped")
		}
		cancel()
		wg.Done()
	}()

	// run refresh periodically
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		tickerChan := ticker.C
		for {
			select {
			case <-tickerChan:
				err := s.refresh()
				if err != nil {
					s.log.Error(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// handle shutdown
	select {
	case <-ctx.Done(): // on any service failed
	case <-s.ctx.Done():
		s.log.Info("interrupt syscall received")
		cancel()
	}

	s.log.Info("cleaning up…")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = feTcp.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("failed gracefully stop api tcp socket server %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = apiUnix.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("failed gracefully stop api tcp socket server %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = apiTcp.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("failed gracefully stop api tcp socket server %v", err)
	}

	err = dnsUdp.Shutdown()
	if err != nil {
		s.log.Errorf("failed gracefully stop dns udp socket server %v", err)
	}

	err = dnsTcp.Shutdown()
	if err != nil {
		s.log.Errorf("failed gracefully stop dns tcp socket server %v", err)
	}

	s.cleanup()
	s.log.Info("waiting workers to stop…")
	wg.Wait()
	s.log.Info("cleanup done, shutdown")
}

func (s *Service) cleanup() {
	s.db.Close()
	s.wgm.Cleanup()
	err := iface.Remove(s.log, s.cfg.WGIface)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *Service) refresh() error {
	tx, err := s.db.Begin(false) // non-writeable tx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	trustIPSet, err := model.LoadManagerSSHTrustIPSet(tx)
	if err != nil {
		return err
	}
	s.trustIPSet.Replace(trustIPSet)
	removed := s.trustIPSet.Removed()
	added := s.trustIPSet.Added()
	s.trustIPSet = s.trustIPSet.Copy()
	err = s.nft.UpdateTrustIPs(removed, added)
	if err != nil {
		s.trustIPSet = ipset.IPSet{}
		// TODO: flush nft ipset
		return err
	}

	users, err := model.LoadUsers(tx)
	if err != nil {
		return err
	}

	devices, err := model.LoadDevices(tx)
	if err != nil {
		return err
	}

	wgManagerIPs := wgManagerIPs(users)
	s.wgManagerIPSet.Replace(wgManagerIPs)
	removed = s.wgManagerIPSet.Removed()
	added = s.wgManagerIPSet.Added()
	s.wgManagerIPSet = s.wgManagerIPSet.Copy()
	err = s.nft.UpdateWGManagerIPs(removed, added)
	if err != nil {
		s.wgManagerIPSet = ipset.IPSet{}
		// TODO: flush nft ipset
		return err
	}

	wgForwardWanIPs := wgForwardWanIPs(devices)
	s.wgForwardWanIPSet.Replace(wgForwardWanIPs)
	removed = s.wgForwardWanIPSet.Removed()
	added = s.wgForwardWanIPSet.Added()
	s.wgForwardWanIPSet = s.wgForwardWanIPSet.Copy()
	err = s.nft.UpdateWGForwardWanIPs(removed, added)
	if err != nil {
		s.wgForwardWanIPSet = ipset.IPSet{}
		// TODO: flush nft ipset
		return err
	}

	wgpeers, err := wgPeers(devices)
	if err != nil {
		return err
	}

	s.wgpeers.Replace(wgpeers)
	peersRemoved := s.wgpeers.Removed()
	peersAdded := s.wgpeers.Added()
	s.wgpeers = s.wgpeers.Copy()

	err = s.wgm.PeerRemove(peersRemoved)
	if err != nil {
		s.wgpeers = wgmngr.PeerSet{}
		return err
	}
	err = s.wgm.PeerSet(peersAdded)
	if err != nil {
		s.wgpeers = wgmngr.PeerSet{}
		return err
	}

	domains, err := model.LoadDomains(tx)
	if err != nil {
		return err
	}
	m := make(map[string]model.Domain, len(domains))
	for _, d := range domains {
		m[d.Name] = d
	}
	s.resolver.Update(m)

	return nil
}

func (s *Service) apiTcpMux(ctx context.Context) *http.ServeMux {
	// http rpc service api
	httprpc := rpcapi.New(s.log, s.cfg.feHTTPOrigin, s.cfg.DevAuthIP)

	// auth api
	authCfg := auth.Config{
		SessionSecret: s.cfg.SessionSecret,
		SessionTTL:    s.cfg.SessionTTL,
	}
	auth := auth.New(ctx, s.log, authCfg, s.db)
	auth.RegisterHandlers(httprpc)

	// manager api
	managerCfg := manager.Config{
		AuthRequired: true,

		WanIP:    s.nft.WanIP(),
		WgInet:   s.cfg.wgIfaceInet,
		WgIPNet:  s.cfg.wgIfaceIPNet,
		WgPort:   s.cfg.WGPort,
		WgPubKey: s.wgm.PublicKey(),

		OTPIssuer: s.cfg.OTPIssuer,

		SessionSecret: s.cfg.SessionSecret,
		SessionTTL:    s.cfg.SessionTTL,
	}
	manager := manager.New(ctx, s.log, managerCfg, s.db)
	manager.RegisterHandlers(httprpc)

	// http rest service api
	httprest := httpapi.New(s.log)

	system := system.New(ctx, s.log)
	system.RegisterHandlers(httprest)

	// register handlers
	mux := http.NewServeMux()
	mux.Handle("/rpc", httprpc)
	mux.Handle("/rest/", http.StripPrefix("/rest", httprest))

	return mux
}

func (s *Service) apiUnixMux(ctx context.Context) *http.ServeMux {
	// http rpc service api
	httprpc := rpcapi.New(s.log, "", "127.0.0.1")

	// manager api
	cfg := manager.Config{
		AuthRequired: false,

		WanIP:    s.nft.WanIP(),
		WgInet:   s.cfg.wgIfaceInet,
		WgIPNet:  s.cfg.wgIfaceIPNet,
		WgPort:   s.cfg.WGPort,
		WgPubKey: s.wgm.PublicKey(),

		OTPIssuer: s.cfg.OTPIssuer,

		SessionSecret: s.cfg.SessionSecret,
		SessionTTL:    s.cfg.SessionTTL,
	}
	manager := manager.New(ctx, s.log, cfg, s.db)
	manager.RegisterHandlers(httprpc)

	// register rpc handlers
	mux := http.NewServeMux()
	mux.Handle("/rpc", httprpc)

	return mux
}

func (s *Service) feTcpMux(ctx context.Context) (*http.ServeMux, error) {
	// fe service
	apiURL := &url.URL{Scheme: "http", Host: s.cfg.apiHTTPAddr, Path: "rpc"}
	fe, err := femanager.Init(
		s.log,
		s.db,
		apiURL.String(),
		s.cfg.DevAuthIP)
	if err != nil {
		return nil, err
	}

	// register fe handler
	mux := http.NewServeMux()
	mux.Handle("/", fe)

	return mux, nil
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
