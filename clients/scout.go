package clients

import (
	"github.com/grimdork/foreman/api"
)

// Scout configuration.
type Scout struct {
	Client
	// Port of the client. Assumes standard HTTPS if unspecified.
	Port int
	// Tries made so far.
	Tries uint8
}

// NewScout creates a new scout.
func NewScout(hostname string, port int, interval int) *Scout {
	return &Scout{
		Client: Client{
			Hostname: hostname,
			Interval: interval,
			Status:   api.StatusWaiting,
			quit:     make(chan any),
		},
		Port: port,
	}
}

// Start the checker.
func (c *Scout) Start() {
	c.quit = make(chan any)
	go func() {
		for {
			select {
			case <-c.quit:
				return
			}
		}
	}()
}

// Stop the checker.
func (c *Scout) Stop() {
	close(c.quit)
}
