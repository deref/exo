package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/deref/exo/chrono"
	"github.com/deref/exo/logd/api"
	"github.com/deref/exo/logd/store"
)

type Worker struct {
	Source string
	Sink   store.Log
	Debug  bool
}

func (wkr *Worker) debugf(format string, v ...interface{}) {
	if wkr.Debug {
		fmt.Fprintln(os.Stderr, "worker", wkr.Source, fmt.Errorf(format, v...))
	}
}

func (wkr *Worker) Run(ctx context.Context) error {
	wkr.debugf("opening fifo")
	fd, err := syscall.Open(wkr.Source, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		wkr.debugf("error opening fifo: %w", err)
		return fmt.Errorf("opening source: %w", err)
	}
	source := os.NewFile(uintptr(fd), wkr.Source)
	wkr.debugf("fifo open")

	go func() {
		<-ctx.Done()
		err := source.Close()
		wkr.debugf("closed source: %v", err)
	}()

	r := bufio.NewReaderSize(source, api.MaxMessageSize)
	for {
		wkr.debugf("reading line")
		message, isPrefix, err := r.ReadLine()
		if isDoneErr(err) {
			return nil
		}
		wkr.debugf("message: %s", message)
		if err != nil {
			return fmt.Errorf("reading message: %w", err)
		}
		// TODO: Do something better with lines that are too long.
		for isPrefix {
			// Skip remainder of message.
			message = append([]byte{}, message...)
			wkr.debugf("skipping line")
			_, isPrefix, err = r.ReadLine()
			if isDoneErr(err) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("truncating message: %w", err)
			}
		}

		timestamp := chrono.NowNano(ctx)
		if err := wkr.Sink.AddEvent(ctx, timestamp, message); err != nil {
			return fmt.Errorf("adding event: %w", err)
		}
	}
}

func isDoneErr(err error) bool {
	return err == io.EOF || (err != nil && strings.HasSuffix(err.Error(), "file already closed"))
}
