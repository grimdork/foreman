package api

// Message is sent when status changes from or to OK.
type Message struct {
	// EventID for keeping track.
	EventID int64
	// OldStatus is the previous status.
	OldStatus uint8
	// NewStatus is the new status. May be the same as old if it's just an update.
	NewStatus uint8
	// Name of the monitored host.
	Name string
	// Text to go in the message.
	Text string
}
