package clients

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"github.com/grimdork/foreman/api"
	ll "github.com/grimdork/loglines"
)

// Scout configuration.
type Scout struct {
	Client
	// Port of the client. Assumes standard HTTPS if unspecified.
	Port    int16
	alerts  chan api.Message
	cfg     *tls.Config
	dialer  *net.Dialer
	address string
}

// NewScout creates a new scout.
func NewScout(hostname string, alerts chan api.Message, capool *x509.CertPool) *Scout {
	scout := &Scout{
		Client: Client{
			Hostname: hostname,
			Interval: 60,
			Status:   api.StatusWaiting,
			quit:     make(chan any),
		},
		alerts: alerts,
		cfg: &tls.Config{
			RootCAs:            capool,
			InsecureSkipVerify: false,
		},
		dialer: &net.Dialer{Timeout: time.Second * 10},
	}
	return scout
}

// Start the checker.
func (scout *Scout) Start() {
	scout.quit = make(chan any)
	scout.address = net.JoinHostPort(scout.Hostname, fmt.Sprintf("%d", scout.Port))
	go func() {
		scout.Check()
		tick := time.NewTicker(time.Second * time.Duration(scout.Interval))
		for {
			select {
			case <-tick.C:
				scout.Check()

			case <-scout.quit:
				return
			}
		}
	}()
}

// Stop the checker.
func (scout *Scout) Stop() {
	close(scout.quit)
}

// Check checks for connectivity, then validates the certificate.
func (scout *Scout) Check() {
	ll.Msg("Checking %s:%d", scout.Hostname, scout.Port)
	msg := api.Message{
		Name:      scout.Hostname,
		OldStatus: scout.Status,
	}
	conn, err := tls.DialWithDialer(scout.dialer, "tcp", scout.address, scout.cfg)
	if err != nil {
		switch err.(type) {
		case tls.RecordHeaderError:
			msg.NewStatus = api.StatusCritical
			msg.Type = api.FailNoTLS

		case *net.OpError:
			msg.NewStatus = api.StatusCritical
			msg.Type = api.FailUnknownHost
		}
		scout.alerts <- msg
		return
	}

	defer conn.Close()
	ll.Msg("Checking certificate for %s", scout.Hostname)
	err = conn.Handshake()
	if err != nil {
		ll.Msg("Error checking certificate: %s", err)
		return
	}

	err = conn.VerifyHostname(scout.Hostname)
	if err != nil {
		msg.NewStatus = api.StatusCritical
		msg.Type = api.FailNoCert
		scout.alerts <- msg
	}
}
