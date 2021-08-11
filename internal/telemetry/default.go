package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/deref/exo"
	"github.com/deref/exo/internal/config"
	"github.com/deref/exo/internal/util/cacheutil"
)

type defaultTelemetry struct {
	cfg       *config.TelemetryConfig
	client    *http.Client
	sessionID string

	latestVersion *cacheutil.TTLVal
}

func (t *defaultTelemetry) IsEnabled() bool {
	// defaultTelemetry can only be created via the NewTelemetry() factory
	// function, which will only return a defaultTelemetry if enabled.
	return true
}

func (t *defaultTelemetry) StartSession() {
	t.ensureSession()
}

func (t *defaultTelemetry) LatestVersion() (string, error) {
	v, err := t.latestVersion.Get()
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (t *defaultTelemetry) SendEvent(evt event) {
	t.ensureSession()
	var buf bytes.Buffer
	req := telemetryRequest{
		Method: "record-event",
		Data: map[string]interface{}{
			"id":      evt.ID(),
			"payload": evt.Payload(),
		},
	}
	e := json.NewEncoder(&buf)
	e.Encode(req)
	// Ignore response.
	_, _ = t.client.Post(exo.TelemetryEndpoint, "application/json", &buf)
}

func (t *defaultTelemetry) ensureSession() {
	if t.sessionID == "" {
		t.startSession()
	}
}

func (t *defaultTelemetry) startSession() {
	// TODO: Start a session at the telemetry endpoint, and update the internal sessionID.
	// This needs to be thread-safe.
}

func (t *defaultTelemetry) getLatestVersion() (interface{}, error) {
	resp, err := t.client.Get(exo.CheckVersionEndpoint)
	if err != nil {
		return "", fmt.Errorf("fetching latest version: %w", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading latest version: %w", err)
	}

	return string(body), nil
}
