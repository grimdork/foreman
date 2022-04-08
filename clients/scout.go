package clients

import (
	"github.com/grimdork/foreman/api"
)

// Scout configuration.
type Scout struct {
	Client
	// Port of the client. Assumes standard HTTPS if unspecified.
	Port int16
	// Tries made so far.
	Tries uint8
	log   chan string
}

// NewScout creates a new scout.
func NewScout(hostname string, logger chan string) *Scout {
	return &Scout{
		Client: Client{
			Hostname: hostname,
			Status:   api.StatusWaiting,
			quit:     make(chan any),
		},
		log: logger,
	}
}

// Start the checker.
func (scout *Scout) Start() {
	scout.quit = make(chan any)
	go func() {
		scout.log <- "Starting scout for " + scout.Hostname
		for {
			select {
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
