package telemetry

import (
	"runtime"

	"github.com/deref/exo/internal/about"
)

type Event struct {
	// Set by event constructor.
	Type            string                 `json:"event_type"`
	EventProperties map[string]interface{} `json:"event_properties,omitempty"`
	UserProperties  map[string]interface{} `json:"user_properties,omitempty"`
	Platform        string                 `json:"platform,omitempty"`
	OSName          string                 `json:"os_name,omitempty"`
	OSVersion       string                 `json:"os_version,omitempty"`
	AppVersion      string                 `json:"app_version,omitempty"`

	// Set by telemetry.
	DeviceID  string `json:"device_id"`
	SessionID int64  `json:"session_id"`
	EventID   int    `json:"event_id"`
	Time      int64  `json:"time"`
}

func SystemInfoIdentifiedEvent() Event {
	return Event{
		Type:       "system-info-identified",
		Platform:   runtime.GOARCH,
		OSName:     runtime.GOOS,
		AppVersion: about.Version,
		EventProperties: map[string]interface{}{
			"cpu_count": runtime.NumCPU(),
		},
	}
}

func OperationsPerformedEvent(operation string, success bool, duration SummaryStatistics) Event {
	return Event{
		Type: "operations-performed",
		EventProperties: map[string]interface{}{
			"operation":       operation,
			"success":         success,
			"occurrences":     duration.Count,
			"duration_max":    duration.Max,
			"duration_min":    duration.Min,
			"duration_mean":   duration.Mean,
			"duration_median": duration.Median,
			"duration_stddev": duration.StdDev,
			"duration_sum":    duration.Sum,
		},
	}
}
