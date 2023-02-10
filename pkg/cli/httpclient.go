package cli

import (
	"context"
	"net"
	"net/http"
	"time"
)

type dialer struct {
	net.Dialer
	socket string
}

// Dial overrides net.Dialer.Dial to force unix socket connection.
func (d *dialer) Dial(network, address string) (net.Conn, error) {
	return d.Dialer.Dial("unix", d.socket)
}

// DialContext overrides net.Dialer.DialContext to force unix socket connection.
func (d *dialer) DialContext(
	ctx context.Context,
	network, address string,
) (net.Conn, error) {
	return d.Dialer.DialContext(ctx, "unix", d.socket)
}

func newHTTPClient(socket string) *http.Client {
	transport := &http.Transport{
		DialContext: (&dialer{Dialer: net.Dialer{
			Timeout: 500 * time.Millisecond,
		}, socket: socket}).DialContext,
		Dial: (&dialer{Dialer: net.Dialer{
			Timeout: 500 * time.Millisecond,
		}, socket: socket}).Dial,
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   1000 * time.Millisecond,
	}

	return client
}
