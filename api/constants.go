package api

// Status of a monitored host.
const (
	// StatusWaiting means we're in the process of starting monitoring.
	StatusWaiting uint8 = iota
	// StatusOK if last check was fine.
	StatusOK
	// StatusWarning if it missed one check after retries.
	StatusWarning
	// StatusCritical if it missed more than one check after retries.
	StatusCritical
)

// Check failure types.
const (
	// FailOK means the check was successful.
	FailOK uint8 = iota
	// FailUnknownHost means the hostname was not found.
	FailUnknownHost
	// FailNoTLS means the host port is not secure.
	FailNoTLS
	// FailNoCert means the host certificate was not found or invalid.
	FailNoCert
)

// Message is sent when status changes from or to OK.
type Message struct {
	// OldStatus is the previous status.
	OldStatus uint8
	// NewStatus is the new status. May be the same as old if it's just an update.
	NewStatus uint8
	// Type of error.
	Type uint8
	// Name of the monitored host.
	Name string
}

// StatusString returns a string representation of the status.
func StatusString(status uint8) string {
	switch status {
	case StatusWaiting:
		return "waiting"
	case StatusOK:
		return "ok"
	case StatusWarning:
		return "warning"
	case StatusCritical:
		return "critical"
	default:
		return "unknown"
	}
}
