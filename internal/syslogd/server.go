package syslogd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/eventd/api"
	"github.com/deref/exo/internal/util/logging"
	"github.com/influxdata/go-syslog/v3"
	"github.com/influxdata/go-syslog/v3/rfc5424"
)

// Server implements a UDP-based Syslog server backed by a log collector store.
type Server struct {
	Logger     logging.Logger
	SyslogPort uint
	api.Store
}

func (svr *Server) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", svr.SyslogPort)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("listening: %w", err)
	}
	svr.Logger.Infof("listening for syslog at udp %s", addr)

	errC := make(chan error, 1)
	go func() {
		maxPacketSize := 8192 // RFC5425#section-4.3.1
		buffer := make([]byte, maxPacketSize)
		syslogMachine := rfc5424.NewMachine()
		for {
			packetSize, _, err := conn.ReadFrom(buffer)
			if err != nil {
				errC <- err
				return
			}
			syslogMessage, err := syslogMachine.Parse(buffer[:packetSize])
			if err != nil {
				svr.Logger.Infof("parsing syslog message: %v", err)
				continue
			}
			event, err := syslogToEvent(syslogMessage)
			if err != nil {
				svr.Logger.Infof("interpreting syslog message: %v", err)
				continue
			}
			if _, err := svr.AddEvent(ctx, event); err != nil {
				errC <- fmt.Errorf("adding event: %w", err)
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errC:
		return err
	}
}

// See supervise implementation for details on Syslog field usage.
func syslogToEvent(syslogMessage syslog.Message) (*api.AddEventInput, error) {
	rfc5425Message, ok := syslogMessage.(*rfc5424.SyslogMessage)
	if !ok {
		panic("unexpected syslog message type")
	}
	if rfc5425Message.Appname == nil {
		return nil, errors.New("expected APP-NAME")
	}
	if rfc5425Message.MsgID == nil {
		return nil, errors.New("expected MSGID")
	}
	if rfc5425Message.Timestamp == nil {
		return nil, errors.New("expected TIMESTAMP")
	}

	appname := *rfc5425Message.Appname
	msgID := *rfc5425Message.MsgID

	var streamName string
	switch msgID {
	case "out", "err":
		// provider=unix
		streamName = fmt.Sprintf("%s:%s", appname, msgID) // See note: [LOG_COMPONENTS].
	default:
		if msgID == appname {
			// provider=docker
			streamName = appname
		} else {
			return nil, fmt.Errorf("unexpected MSGID: %q", msgID)
		}
	}

	message := ""
	if rfc5425Message.Message != nil {
		message = strings.TrimSuffix(*rfc5425Message.Message, "\n")
	}

	return &api.AddEventInput{
		Stream:    streamName,
		Timestamp: rfc5425Message.Timestamp.Format(chrono.RFC3339MicroUTC),
		Message:   message,
	}, nil
}
