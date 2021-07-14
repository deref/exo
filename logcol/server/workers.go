package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/deref/exo/chrono"
	"github.com/deref/exo/logcol/api"
)

type worker struct {
	sourcePath string
	sink       sink
	done       chan struct{}

	err    error
	source *os.File
}

type sink interface {
	AddEvent(ctx context.Context, sid uint64, timestamp uint64, message []byte) error
}

func (lc *LogCollector) ensureWorker(ctx context.Context, logName string, state LogState) {
	lc.mx.Lock()
	defer lc.mx.Unlock()
	if lc.workers == nil {
		// No worker support in peer mode.
		return
	}
	lc.startWorker(ctx, logName, state)
}

func (lc *LogCollector) startWorker(ctx context.Context, logName string, state LogState) {
	wkr, exists := lc.workers[logName]
	if exists {
		return
	}
	sink := newBadgerSink(lc.db, logName)
	wkr = &worker{
		sourcePath: state.Source,
		sink:       sink,
		done:       make(chan struct{}, 0),
	}
	lc.workers[logName] = wkr
	go func() {
		wkr.err = wkr.run(ctx)
		if wkr.err != nil {
			// TODO: Panic instead.
			fmt.Fprintf(os.Stderr, "worker error: %v\n", wkr.err)
		}
	}()
}

func (lc *LogCollector) stopWorker(logName string) {
	lc.mx.Lock()
	defer lc.mx.Unlock()

	wkr := lc.workers[logName]
	if wkr == nil {
		return
	}

	_ = wkr.source.Close()
}

func (wkr *worker) run(ctx context.Context) error {
	source, err := os.Open(wkr.sourcePath + ".fifo")
	if err != nil {
		return fmt.Errorf("opening source: %w", err)
	}

	go func() {
		<-wkr.done
		_ = source.Close()
	}()

	sid := uint64(0) // XXX get sequence id from last line in file.

	r := bufio.NewReaderSize(source, api.MaxMessageSize)
	for {
		sid++
		message, isPrefix, err := r.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading: %w", err)
		}
		// TODO: Do something better with lines that are too long.
		for isPrefix {
			// Skip remainder of message.
			message = append([]byte{}, message...)
			_, isPrefix, err = r.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("reading: %w", err)
			}
		}

		timestamp := chrono.NowNano(ctx)
		if err := wkr.sink.AddEvent(ctx, sid, timestamp, message); err != nil {
			return fmt.Errorf("adding event: %w", err)
		}
	}
}
