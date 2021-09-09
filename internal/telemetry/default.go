package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/deref/exo"
	"github.com/deref/exo/internal/util/cacheutil"
	"github.com/deref/exo/internal/util/logging"
)

type defaultTelemetry struct {
	ctx            context.Context
	client         *http.Client
	installationID string
	sessionID      string
	operationGauge *SummaryGauge

	latestVersion *cacheutil.TTLVal
	getSession    sync.Once
}

func (t *defaultTelemetry) IsEnabled() bool {
	// defaultTelemetry can only be created via the NewTelemetry() factory
	// function, which will only return a defaultTelemetry if enabled.
	return true
}

func (t *defaultTelemetry) StartSession(ctx context.Context) {
	go t.ensureSession(ctx)
}

func (t *defaultTelemetry) LatestVersion(ctx context.Context) (string, error) {
	v, err := t.latestVersion.Get(ctx)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (t *defaultTelemetry) SendEvent(ctx context.Context, evt event) {
	// TODO: Limit concurrent sends/throttle.
	go func() {
		t.ensureSession(ctx)
		var buf bytes.Buffer
		req := telemetryRequest{
			Method: "record-event",
			Data: map[string]interface{}{
				"id":        evt.ID(),
				"payload":   evt.Payload(),
				"sessionId": t.sessionID,
			},
		}
		e := json.NewEncoder(&buf)
		e.Encode(req)
		// Ignore response.
		_, _ = t.client.Post(exo.TelemetryEndpoint, "application/json", &buf)
	}()
}

func (t *defaultTelemetry) RecordOperation(op OperationInvocation) {
	var success string
	if op.Success {
		success = "y"
	} else {
		success = "n"
	}

	t.operationGauge.Observe(Tags{
		"operation": op.Operation,
		"success":   success,
	}, float64(op.DurationMicros))
}

type startSessionResponse struct {
	SessionID string `json:"sessionId"`
}

func (t *defaultTelemetry) ensureSession(ctx context.Context) {
	if t.sessionID != "" {
		return
	}

	t.getSession.Do(func() {
		logger := logging.CurrentLogger(ctx)
		errSetSession := func(res *http.Response, err error) {
			var msg strings.Builder
			msg.WriteString("Error creating session. Setting null UUID.")
			if res != nil && res.StatusCode > 0 {
				var responseText string
				if responseBody, resErr := ioutil.ReadAll(res.Request.Body); resErr != nil {
					responseText = "cannot read response"
				} else {
					responseText = string(responseBody)
				}
				msg.WriteString(fmt.Sprintf(" - HTTP %d: %q", res.StatusCode, responseText))
			}

			if err != nil {
				msg.WriteString(fmt.Sprintf(": %v", err))
			}

			logger.Infof(msg.String())
			t.sessionID = "00000000-0000-0000-0000-000000000000"
		}

		var buf bytes.Buffer
		req := telemetryRequest{
			Method: "start-session",
		}
		e := json.NewEncoder(&buf)
		e.Encode(req)
		res, err := t.client.Post(exo.TelemetryEndpoint, "application/json", &buf)
		if err != nil {
			errSetSession(res, err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			errSetSession(res, nil)
			return
		}

		typedRes := startSessionResponse{}
		if err := json.NewDecoder(res.Body).Decode(&typedRes); err != nil {
			errSetSession(res, err)
			return
		}

		t.sessionID = typedRes.SessionID
		logger.Infof("Started session: %q", t.sessionID)
	})
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

func (t *defaultTelemetry) sendRecordedTelemetry() {
	if len(t.operationGauge.buckets) == 0 {
		return
	}

	newGauge := newOperationGauge()
	oldPtr := (*unsafe.Pointer)(unsafe.Pointer(&t.operationGauge))
	swappedPtr := atomic.SwapPointer(oldPtr, unsafe.Pointer(newGauge))
	g := (*SummaryGauge)(swappedPtr)

	logging.CurrentLogger(t.ctx).Infof("Sending telemetry")

	for _, bucket := range g.buckets {
		tags := bucket.Tags()
		op := tags["operation"]
		success := tags["success"] == "y"
		summary := bucket.Summarize()

		t.SendEvent(t.ctx, &OperationsPerformed{
			Operation:       op,
			Success:         success,
			DurationSummary: summary,
		})
	}
}

func newOperationGauge() *SummaryGauge {
	return NewSummaryGauge([]string{"operation", "success"})
}
