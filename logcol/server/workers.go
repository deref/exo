package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"

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
	AddEvent(ctx context.Context, timestamp uint64, message []byte) error
	// Remove oldest events beyond capacity limit.
	GC(ctx context.Context) error
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
	sink := newBadgerSink(lc.db, lc.idGen, logName)
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
			fmt.Fprintf(os.Stderr, "worker run error: %v\n", wkr.err)
		}
	}()

	go func() {
		for {
			select {
			case <-wkr.done:
				return
			case <-time.After(5 * time.Second):
				if err := wkr.sink.GC(ctx); err != nil {
					// TODO: Panic instead?
					fmt.Fprintf(os.Stderr, "worker evict error: %v\n", wkr.err)
				}
			}
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

	wkr.stop()
}

func (wkr *worker) stop() {
	// XXX Saw a race on shutdown/interrupt where close called on already closed channel.
	// TODO: Make stop safe to call on already stopped worker?
	close(wkr.done)
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

	r := bufio.NewReaderSize(source, api.MaxMessageSize)
	for {
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
		if err := wkr.sink.AddEvent(ctx, timestamp, message); err != nil {
			return fmt.Errorf("adding event: %w", err)
		}
	}
}
