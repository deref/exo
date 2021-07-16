package server

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v3"
)

func (lc *LogCollector) Start(ctx context.Context) error {
	lc.debugf("Start")
	lc.mx.Lock()
	defer lc.mx.Unlock()

	if lc.workers != nil {
		return errors.New("already started")
	}

	dbDir := filepath.Join(lc.varDir, "logs")
	var err error
	lc.db, err = badger.Open(badger.DefaultOptions(dbDir))
	if err != nil {
		return fmt.Errorf("opening db: %w", err)
	}

	lc.workers = make(map[string]*worker)

	state, err := lc.derefState()
	if err != nil {
		return fmt.Errorf("getting state: %w", err)
	}

	for logName, logState := range state.Logs {
		lc.startWorker(ctx, logName, logState)
	}

	return nil
}

func (lc *LogCollector) Stop(ctx context.Context) {
	lc.debugf("Stop")
	lc.mx.Lock()
	defer lc.mx.Unlock()

	if lc.workers == nil {
		panic("not running")
	}

	for _, worker := range lc.workers {
		worker.stop()
	}

	lc.wg.Wait()
	lc.workers = nil
	lc.debugf("stopped")

	lc.debugf("closing db")
	_ = lc.db.Close()
	lc.debugf("db closed")
	lc.db = nil
}
