package wgnetwork

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"wgnetwork/api/auth"
	"wgnetwork/api/manager"
	"wgnetwork/api/system"
	"wgnetwork/firewall"
	"wgnetwork/model"
	"wgnetwork/pkg/envconfig"
	"wgnetwork/pkg/httpapi"
	"wgnetwork/pkg/iface"
	"wgnetwork/pkg/ipcalc"
	"wgnetwork/pkg/ipset"
	"wgnetwork/pkg/rpcapi"
	"wgnetwork/pkg/wgmngr"
)

// config for service.
type config struct {
	DBPath string `env:"DB_PATH" default:"wgnetwork.db"`

	WGBinary string `env:"WG_BINARY" default:"/usr/bin/wg"`
	WGIface  string `env:"WG_IFACE" default:"wg0"`
	WGPort   uint16 `env:"WG_PORT" default:"51820"`
	WGCIDR   string `env:"WG_CIDR" default:"172.16.0.1/24"`

	NFTEnabled          bool   `env:"NFT_ENABLED" default:"false"`
	NFTNetworkNamespace string `env:"NFT_NETWORK_NAMESPACE"`
	NFTDefaultPolicy    string `env:"NFT_DEFAULT_POLICY" default:"drop"`

	NFTIfaces []string `env:"NFT_IFACES"`

	NFTTrustPorts []uint16 `env:"NFT_TRUST_PORTS" default:"22"`

	FEHTTPPort    int    `env:"FE_HTTP_PORT" default:"80"`
	APIHTTPPort   int    `env:"API_HTTP_PORT" default:"8080"`
	APIUnixSocket string `env:"API_UNIX_SOCKET" default:"/tmp/wgmanager.sock"`

	OTPIssuer     string        `env:"OTP_ISSUER" default:"wgnetwork"`
	SessionSecret string        `env:"SESSION_SECRET" default:"secret"`
	SessionTTL    time.Duration `env:"SESSION_TTL" default:"5m"`

	wgIfaceIP    net.IP
	wgIfaceIPNet *net.IPNet
	wgIfaceInet  *net.IPNet

	apiHTTPAddr  string
	feHTTPAddr   string
	feHTTPOrigin string

	// for development environment only
	DevHostname   string `env:"DEV_HOSTNAME"`
	DevHTTPOrigin string `env:"DEV_HTTPORIGIN"`
	DevAuthIP     string `env:"DEV_AUTHIP"`
}

// loadConfig reads configuration from environment variables
func loadConfig() (config, error) {
	cfg := config{}
	err := envconfig.ReadEnv(&cfg)
	if err != nil {
		return config{}, err
	}

	wgIfaceIP, wgIfaceIPNet, err := ipcalc.ParseCIDR(cfg.WGCIDR)
	if err != nil {
		return config{}, err
	}
	cfg.wgIfaceIP = wgIfaceIP
	cfg.wgIfaceIPNet = wgIfaceIPNet
	cfg.wgIfaceInet = &net.IPNet{IP: wgIfaceIP, Mask: wgIfaceIPNet.Mask}

	hostname := cfg.wgIfaceIP.String()
	if cfg.DevHostname != "" {
		hostname = cfg.DevHostname
	}
	cfg.apiHTTPAddr = fmt.Sprintf("%s:%d", hostname, cfg.APIHTTPPort)

	cfg.feHTTPAddr = fmt.Sprintf("%s:%d", hostname, cfg.FEHTTPPort)
	if cfg.DevHTTPOrigin != "" {
		cfg.feHTTPOrigin = cfg.DevHTTPOrigin
	} else if cfg.FEHTTPPort != 80 {
		cfg.feHTTPOrigin = fmt.Sprintf(
			"http://%s:%d", hostname, cfg.FEHTTPPort)
	} else {
		cfg.feHTTPOrigin = fmt.Sprintf(
			"http://%s", hostname)
	}

	return cfg, err
}

// Service object.
type Service struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    config
	log    logger

	db  *bolt.DB
	nft *firewall.NFTables

	trustIPSet        ipset.IPSet
	wgManagerIPSet    ipset.IPSet
	wgForwardWanIPSet ipset.IPSet

	wgm     *wgmngr.Manager
	wgpeers wgmngr.PeerSet

	fehttp        http.Server
	apihttp       http.Server
	apihttpsocket http.Server
}

// Init service.
func Init(
	ctx context.Context,
	cancel context.CancelFunc,
	log logger,
) (*Service, error) {
	var (
		err error
		cfg config
	)

	cfg, err = loadConfig()
	if err != nil {
		return nil, err
	}

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
		Ifaces:           cfg.NFTIfaces,
		TrustPorts:       cfg.NFTTrustPorts,
	}
	nft, err = firewall.Init(nftCfg, managerPorts)
	if err != nil {
		return nil, err
	}

	s := &Service{
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
		log:    log,

		db:  db,
		nft: nft,

		trustIPSet:        ipset.IPSet{},
		wgManagerIPSet:    ipset.IPSet{},
		wgForwardWanIPSet: ipset.IPSet{},

		wgm:     wgm,
		wgpeers: wgmngr.PeerSet{},
	}

	return s, nil
}

// Run service.
func (s *Service) Run() {
	go s.apihttpserve()
	go s.apihttpsocketserve()
	go s.fehttpserve()

	err := s.refresh()
	if err != nil {
		s.log.Error(err)
	}

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
		}
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

	return nil
}

// Stop service and cleanup
func (s *Service) Stop() {
	s.fehttpshutdown()
	s.apihttpsocketshutdown()
	s.apihttpshutdown()
	s.db.Close()
	s.wgm.Cleanup()
	err := iface.Remove(s.log, s.cfg.WGIface)
	if err != nil {
		s.log.Error(err)
	}
	s.cancel()
}

func (s *Service) apihttpserve() {
	mux := http.NewServeMux()

	// http rpc service api
	httprpc := rpcapi.New(s.log, s.cfg.feHTTPOrigin, s.cfg.DevAuthIP)

	// auth api
	authCfg := auth.Config{
		SessionSecret: s.cfg.SessionSecret,
		SessionTTL:    s.cfg.SessionTTL,
	}
	auth := auth.New(s.ctx, s.log, authCfg, s.db)
	auth.RegisterHandlers(httprpc)

	// common api
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
	manager := manager.New(s.ctx, s.log, managerCfg, s.db)
	manager.RegisterHandlers(httprpc)

	// register rpc handlers
	mux.Handle("/rpc", httprpc)

	// http rest service api
	httprest := httpapi.New(s.log)

	system := system.New(s.ctx, s.log)
	system.RegisterHandlers(httprest)

	// register rest handlers
	mux.Handle("/rest/", http.StripPrefix("/rest", httprest))

	listen, err := net.Listen("tcp", s.cfg.apiHTTPAddr)
	if err != nil {
		s.log.Errorf("failed to listen tcp addr %q: %v",
			s.cfg.apiHTTPAddr, err)
		s.cancel()
		return
	}
	defer listen.Close()

	s.log.Infof("serving api on %q…", s.cfg.apiHTTPAddr)
	defer func() {
		s.log.Infof("serving api on %q has been stopped", s.cfg.apiHTTPAddr)
	}()

	s.apihttp = http.Server{
		Addr:              s.cfg.apiHTTPAddr,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return s.ctx
		},
	}

	err = s.apihttp.Serve(listen)
	if err != http.ErrServerClosed {
		s.log.Errorf("http serve has failed %v", err)
	}
}

func (s *Service) apihttpshutdown() {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	err := s.apihttp.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("failed gracefully stop http server %v", err)
	}
}

func (s *Service) apihttpsocketserve() {
	err := os.RemoveAll(s.cfg.APIUnixSocket)
	if err != nil {
		s.log.Errorf("failed to remove unix socket %q: %v",
			s.cfg.APIUnixSocket, err)
		s.cancel()
		return
	}

	mux := http.NewServeMux()

	// http rpc service api
	httprpc := rpcapi.New(s.log, "", "127.0.0.1")

	// common api
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
	manager := manager.New(s.ctx, s.log, cfg, s.db)
	manager.RegisterHandlers(httprpc)

	// register rpc handlers
	mux.Handle("/rpc", httprpc)

	listen, err := net.Listen("unix", s.cfg.APIUnixSocket)
	if err != nil {
		s.log.Errorf("failed to listen unix socket %q: %v",
			s.cfg.APIUnixSocket, err)
		s.cancel()
		return
	}
	defer listen.Close()

	s.log.Infof("serving api on %q…", s.cfg.APIUnixSocket)
	defer func() {
		s.log.Infof("serving api on %q has been stopped", s.cfg.APIUnixSocket)
	}()

	s.apihttpsocket = http.Server{
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return s.ctx
		},
	}

	err = s.apihttpsocket.Serve(listen)
	if err != http.ErrServerClosed {
		s.log.Errorf("http socket serve has failed %v", err)
	}
}

func (s *Service) apihttpsocketshutdown() {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	err := s.apihttpsocket.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("failed gracefully stop http socket server %v", err)
	}
}

func (s *Service) fehttpserve() {
	apiURL := &url.URL{
		Scheme: "http",
		Host:   s.cfg.apiHTTPAddr,
		Path:   "rpc",
	}
	fe, err := InitFrontend(s.log, s.db, apiURL.String(), s.cfg.DevAuthIP)
	if err != nil {
		s.log.Errorf("can't init frontend: %v", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/", fe)

	listen, err := net.Listen("tcp", s.cfg.feHTTPAddr)
	if err != nil {
		s.log.Errorf("failed to listen tcp addr %q: %v",
			s.cfg.feHTTPAddr, err)
		s.cancel()
		return
	}
	defer listen.Close()

	s.log.Infof("serving fe on %q…", s.cfg.feHTTPAddr)
	defer func() {
		s.log.Infof("serving fe on %q has been stopped", s.cfg.feHTTPAddr)
	}()

	s.fehttp = http.Server{
		Addr:              s.cfg.feHTTPAddr,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return s.ctx
		},
	}

	err = s.fehttp.Serve(listen)
	if err != http.ErrServerClosed {
		s.log.Errorf("http serve has failed %v", err)
	}
}

func (s *Service) fehttpshutdown() {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	err := s.fehttp.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("failed gracefully stop http server %v", err)
	}
}

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
