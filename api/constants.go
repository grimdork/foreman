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
