package server

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/deref/exo/logd/store/badger"
)

func (lc *LogCollector) Start(ctx context.Context) error {
	lc.debugf("Start")
	lc.mx.Lock()
	defer lc.mx.Unlock()

	if lc.workers != nil {
		return errors.New("already started")
	}

	logsDir := filepath.Join(lc.varDir, "logs")
	var err error
	lc.store, err = badger.Open(ctx, logsDir)
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}

	lc.workers = make(map[string]*collectorWorker)

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

	lc.debugf("closing store")
	_ = lc.store.Close()
	lc.debugf("store closed")
	lc.store = nil
}
