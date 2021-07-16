package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/deref/exo/chrono"
	"github.com/deref/exo/logcol/api"
)

type worker struct {
	sourcePath string
	sink       sink
	debug      bool

	err      error
	source   *os.File
	shutdown chan struct{}
}

func (wkr *worker) debugf(format string, v ...interface{}) {
	if wkr.debug {
		fmt.Fprintln(os.Stderr, "worker", wkr.sourcePath, fmt.Errorf(format, v...))
	}
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
		debug:      lc.debug,
		shutdown:   make(chan struct{}),
	}
	lc.workers[logName] = wkr

	done := make(chan struct{})
	lc.wg.Add(2)

	go func() {
		defer wkr.debugf("run done")
		defer lc.wg.Done()
		wkr.err = wkr.run(ctx)
		if wkr.err != nil {
			// TODO: Panic instead.
			fmt.Fprintf(os.Stderr, "worker run error: %v\n", wkr.err)
		}
		close(done)
	}()

	go func() {
		defer wkr.debugf("gc loop done")
		defer lc.wg.Done()
		for {
			select {
			case <-done:
				return
			case <-time.After(5 * time.Second):
				wkr.debugf("gc start")
				if err := wkr.sink.GC(ctx); err != nil {
					// TODO: Panic instead?
					fmt.Fprintf(os.Stderr, "worker evict error: %v\n", wkr.err)
				}
				wkr.debugf("gc done")
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
	wkr.debugf("stop")
	close(wkr.shutdown)
}

func (wkr *worker) run(ctx context.Context) error {
	// This will block until there is at least one writer on the fifo.
	// Shutdown codepath will open-for-write, then close this path to unblock.
	wkr.debugf("opening fifo")
	fd, err := syscall.Open(wkr.sourcePath, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		wkr.debugf("error opening fifo: %w", err)
		return fmt.Errorf("opening source: %w", err)
	}
	wkr.source = os.NewFile(uintptr(fd), wkr.sourcePath)
	wkr.debugf("fifo open")

	go func() {
		<-wkr.shutdown
		err := wkr.source.Close()
		wkr.debugf("closed source: %v", err)
	}()

	r := bufio.NewReaderSize(wkr.source, api.MaxMessageSize)
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
		if err := wkr.sink.AddEvent(ctx, timestamp, message); err != nil {
			return fmt.Errorf("adding event: %w", err)
		}
	}
}

func isDoneErr(err error) bool {
	return err == io.EOF || (err != nil && strings.HasSuffix(err.Error(), "file already closed"))
}
