package api

import "time"

// Endpoints
const (
	// EPPulse is the endpoint for heartbeats from clients.
	EPPulse = "/pulse"
	// EPCanaries is the endpoint for canaries.
	EPCanaries = "/canaries"
	// EPCanary is the endpoint for a single canary.
	EPCanary = "/canary"
	// EPScouts is the endpoint for listing active checks.
	EPScouts = "/scouts"
	// EPScout is for watcher management.
	EPScout = "/scout"
	// EPKeys lists all keys.
	EPKeys = "/keys"
	// EPKey is for authentication management.
	EPKey = "/key"
	// EPHealth is the preferred health check endpoint for an active check.
	EPHealth = "/health"
)

// Request header variables.
type Request map[string]string

// ErrorResponse is a JSON-encoded error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Endpoints list
const Endpoints = `{
	"endpoints": [
		{
			"name": "/api/keys",
			"parameters": ["id", "key"]
			"methods":
			{
				"GET": []
			}
		},
		{
			"name": "/api/key",
			"parameters": ["id", "key"]
			"methods":
			{
				"POST": ["keyid", "value"],
				"DELETE": ["keyid"]
			}
		},
		{
			"name": "/pulse",
			"parameters": ["id", "key"],
			"methods":
			{
				"POST": []
			}
		},
		{
			"name": "/api/canaries",
			"parameters": ["id", "key"]
			"methods":
			{
				"GET": []
			}
		},
		{
			"name": "/api/canary",
			"parameters": ["id", "key"]
			"methods":
			{
				"GET": [""],
				"POST": [],
				"DELETE": ["hostname"]
			}
		},
		{
			"name": "/api/scouts",
			"parameters": ["id", "key"]
			"methods":
			{
				"GET": ["hostname"]
			}
		},
		{
			"name": "/api/scout",
			"parameters": ["id", "key"]
			"methods":
			{
				"GET": ["hostname"],
				"POST": ["hostname", "port", "interval", ],
				"DELETE": ["hostname"]
			}
		},
		{
			"name": "/health",
			"methods":
			{
				"GET": []
			}
		}
	]
}
`

// CanaryList structure.
type CanaryList struct {
	Canaries []CanaryListEntry `json:"canaries"`
}

// CanaryListEntry is a condensed version of the Scout struct.
type CanaryListEntry struct {
	Hostname        string    `json:"hostname"`
	Assignee        string    `json:"assignee"`
	Key             string    `json:"key"`
	LastCheck       time.Time `json:"lastcheck"`
	Failed          time.Time `json:"failed"`
	Interval        int       `json:"interval"`
	Status          uint8     `json:"status"`
	Acknowledgement bool      `json:"acknowledgement"`
}

// ScoutList structure.
type ScoutList struct {
	Scouts []ScoutListEntry `json:"scouts"`
}

// ScoutListEntry is a condensed version of the Scout struct.
type ScoutListEntry struct {
	Hostname        string    `json:"hostname"`
	Assignee        string    `json:"assignee"`
	LastCheck       time.Time `json:"lastcheck"`
	Failed          time.Time `json:"failed"`
	Interval        int       `json:"interval"`
	Port            int16     `json:"port"`
	Status          uint8     `json:"status"`
	Acknowledgement bool      `json:"acknowledgement"`
}

// KeyList structure.
type KeyList struct {
	Keys []Key `json:"keys"`
}

// Key structure.
type Key struct {
	// ID is the name of the key user.
	ID string `json:"id"`
	// Key is the secret key.
	Key string `json:"key"`
	// Admin is true if this key can modify and view records.
	Admin bool `json:"admin"`
	// Expiry is the time at which the key expires.
	Expiry time.Time `json:"expiry"`
}
