package telemetry

import (
	"net/http"
	"time"

	"github.com/deref/exo/internal/config"
	"github.com/deref/exo/internal/util/cacheutil"
)

type Telemetry interface {
	IsEnabled() bool
	LatestVersion() (string, error)
	StartSession()
	SendEvent(event)
}

func New(cfg *config.TelemetryConfig) Telemetry {
	if cfg.Disable {
		return &noOpTelemetry{}
	}

	t := &defaultTelemetry{
		cfg: cfg,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
	t.latestVersion = cacheutil.NewTTLVal(t.getLatestVersion, 5*time.Minute)

	return t
}

type telemetryRequest struct {
	Method string
	Data   map[string]interface{}
}
