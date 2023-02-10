package wgmngr

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Manager object.
type Manager struct {
	privateKey wgtypes.Key
	publicKey  wgtypes.Key
	iface      string
	port       int
}

// NewManager constructor.
func NewManager(
	privateKey wgtypes.Key,
	iface string,
	port uint16,
) (*Manager, error) {
	wgm := &Manager{
		privateKey: privateKey,
		publicKey:  privateKey.PublicKey(),
		iface:      iface,
		port:       int(port),
	}

	err := wgm.init()
	if err != nil {
		return nil, err
	}

	return wgm, nil
}

func (wgm *Manager) init() error {
	cfg := wgtypes.Config{
		PrivateKey:   &wgm.privateKey,
		ListenPort:   &wgm.port,
		ReplacePeers: true,
	}

	return configureDevice(wgm.iface, cfg)
}

// PeerSet configuration.
func (wgm *Manager) PeerSet(peers Peers) error {
	p := make([]wgtypes.PeerConfig, len(peers))
	for i := 0; i < len(p); i++ {
		p[i] = wgtypes.PeerConfig{
			PublicKey:         peers[i].publicKey,
			AllowedIPs:        peers[i].AllowedIPs(),
			ReplaceAllowedIPs: true,
		}
		endpoint := peers[i].Endpoint()
		if endpoint != nil {
			p[i].Endpoint = endpoint
		}
		keepAlive := peers[i].KeepAlive()
		if keepAlive > 0 {
			p[i].PersistentKeepaliveInterval = &keepAlive
		}
	}

	cfg := wgtypes.Config{Peers: p}

	return configureDevice(wgm.iface, cfg)
}

// PeerRemove configuration.
func (wgm *Manager) PeerRemove(peers Peers) error {
	p := make([]wgtypes.PeerConfig, len(peers))
	for i := 0; i < len(p); i++ {
		p[i] = wgtypes.PeerConfig{
			PublicKey: peers[i].publicKey,
			Remove:    true,
		}
	}

	cfg := wgtypes.Config{Peers: p}

	return configureDevice(wgm.iface, cfg)
}

// PeerReplace configuration.
func (wgm *Manager) PeerReplace(peers Peers) error {
	p := make([]wgtypes.PeerConfig, len(peers))
	for i := 0; i < len(p); i++ {
		p[i] = wgtypes.PeerConfig{
			PublicKey:  peers[i].publicKey,
			AllowedIPs: peers[i].AllowedIPs(),
		}
		endpoint := peers[i].Endpoint()
		if endpoint != nil {
			p[i].Endpoint = endpoint
		}
		keepAlive := peers[i].KeepAlive()
		if keepAlive > 0 {
			p[i].PersistentKeepaliveInterval = &keepAlive
		}
	}

	cfg := wgtypes.Config{
		Peers:        p,
		ReplacePeers: true,
	}

	return configureDevice(wgm.iface, cfg)
}

// PublicKey value.
func (wgm *Manager) PublicKey() wgtypes.Key {
	return wgm.publicKey
}

// Cleanup configuration.
func (wgm *Manager) Cleanup() error {
	zerofm := 0
	zerosk := wgtypes.Key{}
	cfg := wgtypes.Config{
		PrivateKey:   &zerosk,
		ReplacePeers: true,
		FirewallMark: &zerofm,
		Peers:        []wgtypes.PeerConfig{},
	}

	return configureDevice(wgm.iface, cfg)
}
