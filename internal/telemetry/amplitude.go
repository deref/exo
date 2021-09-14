package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	defer close(notify)

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

func (a *AmplitudeClient) publish(events []Event) error {
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	if err := e.Encode(UploadRequest{
		APIKey: a.apiKey,
		Events: events,
	}); err != nil {
		return fmt.Errorf("serializing upload request: %w", err)
	}

	res, err := a.client.Post(amplitudeURL, "application/json", &buf)
	if err != nil {
		return fmt.Errorf("performing telemetry request: %w", err)
	}
	defer res.Body.Close()
	// TODO: Do something with response, potentially log errors, etc.

	return nil
}

type UploadRequest struct {
	APIKey string  `json:"api_key"`
	Events []Event `json:"events"`
}
