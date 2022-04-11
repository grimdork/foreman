package clients

import (
	"github.com/grimdork/foreman/api"
)

// Canary sends a heartbeat pulse at intervals.
type Canary struct {
	Client
	// Key is the client's secret key.
	Key    string
	alerts chan api.Message
}

// NewCanary creates a new canary.
func NewCanary(name string, alerts chan api.Message) *Canary {
	return &Canary{
		Client: Client{
			Hostname: name,
			Interval: 60,
			Status:   api.StatusWaiting,
			quit:     make(chan any),
		},
		alerts: alerts,
	}
}
