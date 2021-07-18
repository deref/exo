package server

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/deref/exo/logd/store/badger"
)

func (lc *LogCollector) Run(ctx context.Context) error {
	if err := lc.start(ctx); err != nil {
		return err
	}
	if err := lc.loop(ctx); err != nil {
		return err
	}
	lc.stop()
	return nil
}

func (lc *LogCollector) start(ctx context.Context) error {
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

func (lc *LogCollector) loop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Second):
			if err := lc.removeOldEvents(ctx); err != nil {
				return err
			}
		}
	}
}

func (lc *LogCollector) removeOldEvents(ctx context.Context) error {
	lc.debugf("removing old events")
	state, err := lc.derefState()
	if err != nil {
		return err
	}
	for logName := range state.Logs {
		log := lc.store.GetLog(logName)
		if err := log.RemoveOldEvents(ctx); err != nil {
			return fmt.Errorf("removing %q events: %w", logName, err)
		}
	}
	lc.debugf("removed old events")
	return nil
}

func (lc *LogCollector) stop() {
	lc.debugf("Stop")
	lc.mx.Lock()
	defer lc.mx.Unlock()

	for _, worker := range lc.workers {
		worker.stop()
	}

	lc.wg.Wait()
	lc.debugf("stopped")

	lc.debugf("closing store")
	_ = lc.store.Close()
	lc.debugf("store closed")
	lc.store = nil
}
