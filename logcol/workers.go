package logcol

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type worker struct {
	sourcePath string

	err    error
	source *os.File
	sink   *os.File
}

func (svc *service) startWorker(logName string, state LogState) {
	svc.mx.Lock()
	defer svc.mx.Unlock()

	wkr, exists := svc.workers[logName]
	if !exists {
		wkr = &worker{
			sourcePath: state.SourcePath,
		}
		svc.workers[logName] = wkr
	}
	go func() {
		wkr.err = wkr.run()
		// TODO: Don't log and do figure out how to handle.
		fmt.Fprintf(os.Stderr, "worker error: %v\n", wkr.err)
	}()
}

func (svc *service) stopWorker(logName string) {
	svc.mx.Lock()
	defer svc.mx.Unlock()

	wkr := svc.workers[logName]
	if wkr == nil {
		return
	}

	_ = wkr.source.Close()
	_ = wkr.sink.Close()
}

func (wkr *worker) run() error {
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

		timestamp := time.Now().Format("2006-01-02T15:04:05.999999999Z")

		if _, err := fmt.Fprintf(sink, "%d %s %s\n", sid, timestamp, message); err != nil {
			return fmt.Errorf("writing: %w", err)
		}
	}
}

func makeChunkPath(sourcePath string, chunkIndex int) string {
	return fmt.Sprintf("%s.%d", sourcePath, chunkIndex)
}
