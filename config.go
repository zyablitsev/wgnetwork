package wgnetwork

import (
	"fmt"
	"net"
	"time"

	"wgnetwork/pkg/envconfig"
	"wgnetwork/pkg/ipcalc"
)

// config for service.
type config struct {
	LogLevel string `env:"LOG_LEVEL" default:"debug"`

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

	DNSTcpPort       int      `env:"DNS_TCP_PORT" default:"53"`
	DNSUdpPort       int      `env:"DNS_UDP_PORT" default:"53"`
	DNSResolverAddrs []string `env:"DNS_RESOLVER_ADDRS" default:"8.8.8.8:53,8.8.4.4:53,1.1.1.1:53"`
	DNSZone          string   `env:"DNS_ZONE" default:"wgn."`
	FEHTTPPort       int      `env:"FE_HTTP_PORT" default:"80"`
	APIHTTPPort      int      `env:"API_HTTP_PORT" default:"8080"`
	APIUnixSocket    string   `env:"API_UNIX_SOCKET" default:"/tmp/wgmanager.sock"`

	OTPIssuer     string        `env:"OTP_ISSUER" default:"wgnetwork"`
	HTTPOrigin    string        `env:"HTTPORIGIN"`
	SessionSecret string        `env:"SESSION_SECRET" default:"secret"`
	SessionTTL    time.Duration `env:"SESSION_TTL" default:"5m"`

	wgIfaceIP    net.IP
	wgIfaceIPNet *net.IPNet
	wgIfaceInet  *net.IPNet

	dnsTcpAddr   string
	dnsUdpAddr   string
	apiHTTPAddr  string
	feHTTPAddr   string
	feHTTPOrigin string

	// for development environment only
	DevHostname string `env:"DEV_HOSTNAME"`
	DevAuthIP   string `env:"DEV_AUTHIP"`
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
	cfg.dnsTcpAddr = fmt.Sprintf("%s:%d", hostname, cfg.DNSTcpPort)
	cfg.dnsUdpAddr = fmt.Sprintf("%s:%d", hostname, cfg.DNSUdpPort)
	cfg.apiHTTPAddr = fmt.Sprintf("%s:%d", hostname, cfg.APIHTTPPort)

	cfg.feHTTPAddr = fmt.Sprintf("%s:%d", hostname, cfg.FEHTTPPort)
	if cfg.HTTPOrigin != "" {
		cfg.feHTTPOrigin = cfg.HTTPOrigin
	} else if cfg.FEHTTPPort != 80 {
		cfg.feHTTPOrigin = fmt.Sprintf(
			"http://%s:%d", hostname, cfg.FEHTTPPort)
	} else {
		cfg.feHTTPOrigin = fmt.Sprintf(
			"http://%s", hostname)
	}

	return cfg, err
}
