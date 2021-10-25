package telemetry

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/deref/exo/internal/about"
	"github.com/deref/exo/internal/util/cacheutil"
)

type Telemetry interface {
	IsEnabled() bool
	LatestVersion(context.Context) (string, error)
	StartSession(context.Context)
	SendEvent(context.Context, Event)
	RecordOperation(OperationInvocation)
}

type Config struct {
	Disable  bool
	DeviceID string
}

func New(ctx context.Context, cfg Config) Telemetry {
	if cfg.Disable {
		return &Nop{}
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	t := &defaultTelemetry{
		ctx:            ctx,
		deviceID:       cfg.DeviceID,
		client:         httpClient,
		ampClient:      NewAmplitudeClient(ctx, httpClient, about.AmplitudeAPIKey),
		operationGauge: newOperationGauge(),
	}
	t.latestVersion = cacheutil.NewTTLVal(t.getLatestVersion, 5*time.Minute)

	go func() {
		timerRand := rand.New(rand.NewSource(time.Now().UnixNano()))
		for {
			// Wait between 120 and 600 seconds.
			nextWait := timerRand.Intn(600-120) + 120
			timer := time.NewTimer(time.Second * time.Duration(nextWait))
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				t.sendRecordedTelemetry()
			}
		}
	}()

	return t
}

type OperationInvocation struct {
	Operation      string
	DurationMicros int
	Success        bool
}

type telemetryRequest struct {
	Method string                 `json:"method"`
	Data   map[string]interface{} `json:"data"`
}
