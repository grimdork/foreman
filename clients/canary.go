package clients

import (
	"github.com/grimdork/foreman/api"
)

// Canary sends a heartbeat pulse at intervals.
type Canary struct {
	Client
	// Key is the client's secret key.
	Key string
	log chan string
}

// NewCanary creates a new canary.
func NewCanary(name string, log chan string) *Canary {
	return &Canary{
		Client: Client{
			Hostname: name,
			Status:   api.StatusWaiting,
			quit:     make(chan any),
		},
		log: log,
	}
}
