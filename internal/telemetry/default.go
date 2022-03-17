package telemetry

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/deref/exo/internal/about"
	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/util/cacheutil"
	"github.com/deref/exo/internal/util/logging"
)

type defaultTelemetry struct {
	ctx               context.Context
	client            *http.Client
	deviceID          string
	sessionID         int64
	operationGauge    *SummaryGauge
	ampClient         *AmplitudeClient
	derefInternalUser bool

	latestVersion *cacheutil.TTLVal
	getSession    sync.Once

	idMu    sync.Mutex
	eventID int
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

func (t *defaultTelemetry) SendEvent(ctx context.Context, evt Event) {
	t.ensureSession(ctx)

	evt.DeviceID = t.deviceID
	evt.SessionID = t.sessionID
	evt.EventID = t.nextEventID()
	evt.Time = chrono.NowMillisecond(ctx)

	if evt.UserProperties == nil {
		evt.UserProperties = make(map[string]any)
	}
	evt.UserProperties["isDerefInternalUser"] = t.derefInternalUser
	evt.UserProperties["osArch"] = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	if err := t.ampClient.Publish(evt); err != nil {
		logging.CurrentLogger(ctx).Infof("Could not publish telemetry event %q: %v", evt.Type, err)
	}
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

func (t *defaultTelemetry) ensureSession(ctx context.Context) {
	if t.sessionID != 0 {
		return
	}

	t.getSession.Do(func() {
		logger := logging.CurrentLogger(ctx)
		t.sessionID = chrono.NowMillisecond(ctx)
		logger.Infof("Started session: %d", t.sessionID)
	})
}

func (t *defaultTelemetry) getLatestVersion() (any, error) {
	resp, err := t.client.Get(about.CheckVersionEndpoint)
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

		t.SendEvent(t.ctx, OperationsPerformedEvent(op, success, summary))
	}
}

func (t *defaultTelemetry) nextEventID() int {
	t.idMu.Lock()
	defer t.idMu.Unlock()

	nextID := t.eventID
	t.eventID++
	return nextID
}

func newOperationGauge() *SummaryGauge {
	return NewSummaryGauge([]string{"operation", "success"})
}
