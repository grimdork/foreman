package clients

import "time"

// Client contains fields common to all watchers.
type Client struct {
	// Hostname is used in alerts.
	Hostname string
	// Interval to expect check-ins in seconds (default 60, give or take 2).
	Interval int
	// LastCheck time.
	LastCheck time.Time
	// Failed is the time failures started if status is not OK.
	Failed time.Time
	// Status of the client.
	Status uint8
	// Acknowledgement of error status.
	Acknowledgement bool
	// Assignee of the acknowledged error.
	Assignee string

	quit chan any
}
