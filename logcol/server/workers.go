package server

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/deref/exo/chrono"
)

type worker struct {
	sourcePath string

	err    error
	source *os.File
	sink   *os.File
}

func (lc *logCollector) startWorker(ctx context.Context, logName string, state LogState) {
	lc.mx.Lock()
	defer lc.mx.Unlock()

	wkr, exists := lc.workers[logName]
	if !exists {
		wkr = &worker{
			sourcePath: state.Source,
		}
		lc.workers[logName] = wkr
	}
	go func() {
		wkr.err = wkr.run(ctx)
		// TODO: Don't log and do figure out how to handle.
		fmt.Fprintf(os.Stderr, "worker error: %v\n", wkr.err)
	}()
}

func (lc *logCollector) stopWorker(logName string) {
	lc.mx.Lock()
	defer lc.mx.Unlock()

	wkr := lc.workers[logName]
	if wkr == nil {
		return
	}

	_ = wkr.source.Close()
	_ = wkr.sink.Close()
}

func (wkr *worker) run(ctx context.Context) error {
	source, err := os.Open(wkr.sourcePath)
	if err != nil {
		return fmt.Errorf("opening source: %w", err)
	}

	chunkIndex := 0 // TODO: Log rotation.
	chunkPath := fmt.Sprintf("%s.%d", wkr.sourcePath, chunkIndex)
	sink, err := os.OpenFile(chunkPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("opening sink: %w", err)
	}

	sid := 0 // XXX get sequence id from last line in file.

	r := bufio.NewReader(source)
	for {
		sid++
		message, isPrefix, err := r.ReadLine()
		if err != nil {
			return fmt.Errorf("reading: %w", err)
		}
		// TODO: Do something better with lines that are too long.
		for isPrefix {
			// Skip remainder of message.
			message = append([]byte{}, message...)
			_, isPrefix, err = r.ReadLine()
			if err != nil {
				return fmt.Errorf("reading: %w", err)
			}
		}

		timestamp := chrono.NowString(ctx)
		if _, err := fmt.Fprintf(sink, "%d %s %s\n", sid, timestamp, message); err != nil {
			return fmt.Errorf("writing: %w", err)
		}
	}
}

func makeChunkPath(source string, chunk int) string {
	return fmt.Sprintf("%s.%d", source, chunk)
}
