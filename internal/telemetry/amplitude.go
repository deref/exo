package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/deref/exo/internal/util/logging"
)

const (
	amplitudeURL           = "https://api2.amplitude.com/2/httpapi"
	eventChannelBufferSize = 8
	eventBufferSize        = 32
)

var eventSendInterval = time.Second * 15

// AmplitudeClient was inlfuenced by https://github.com/savaki/amplitude-go
// but uses the version 2 HTTP API.
type AmplitudeClient struct {
	apiKey string
	events chan Event
	flush  chan chan struct{}
	client *http.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func NewAmplitudeClient(ctx context.Context, client *http.Client, apiKey string) *AmplitudeClient {
	ctx, cancel := context.WithCancel(ctx)
	ac := &AmplitudeClient{
		apiKey: apiKey,
		events: make(chan Event, eventChannelBufferSize),
		flush:  make(chan chan struct{}),
		ctx:    ctx,
		cancel: cancel,
		client: client,
	}

	go ac.run()

	return ac
}

func (a *AmplitudeClient) Publish(evt Event) error {
	select {
	case a.events <- evt:
		return nil
	default:
		return fmt.Errorf("event channel buffer of size %d is full", eventChannelBufferSize)
	}
}

func (a *AmplitudeClient) Flush() {
	notify := make(chan struct{})
	a.flush <- notify
	// Wait for flush to complete.
	<-notify
}

func (a *AmplitudeClient) Close() {
	a.Flush()
	a.cancel()
}

func (a *AmplitudeClient) run() {
	timer := time.NewTimer(eventSendInterval)

	buffer := make([]Event, eventBufferSize)
	idx := 0
	doPublish := func() {
		if idx > 0 {
			a.publish(buffer[0:idx])
			idx = 0
		}
	}

	for {
		timer.Reset(eventSendInterval)
		select {
		case <-a.ctx.Done():
			return

		case <-timer.C:
			doPublish()

		case evt := <-a.events:
			buffer[idx] = evt
			idx++
			if idx == eventBufferSize {
				doPublish()
			}

		case notify := <-a.flush:
			doPublish()
			notify <- struct{}{}
		}
	}
}

func (a *AmplitudeClient) publish(events []Event) {
	logger := logging.CurrentLogger(a.ctx)
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	if err := e.Encode(UploadRequest{
		APIKey: a.apiKey,
		Events: events,
	}); err != nil {
		logger.Infof("Could not serialize upload request: %w", err)
		return
	}

	res, err := a.client.Post(amplitudeURL, "application/json", &buf)
	if err != nil {
		logger.Infof("Could not perform telemetry request: %w", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logger.Infof("Got unexpected status code: %d", res.StatusCode)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Infof("could not read body of errored request")
			return
		}
		logger.Infof("res: %+v\n", string(body))
	}
}

type UploadRequest struct {
	APIKey string  `json:"api_key"`
	Events []Event `json:"events"`
}
