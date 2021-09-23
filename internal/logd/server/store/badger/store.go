package badger

import (
	"context"

	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/logd/server/store"
	"github.com/deref/exo/internal/util/logging"
	"github.com/dgraph-io/badger/v3"
)

type Store struct {
	db    *badger.DB
	idGen *gensym.ULIDGenerator
}

type Log struct {
	db    *badger.DB
	idGen *gensym.ULIDGenerator
	name  string
}

// Ambiguous notions of "logs" here.
// logger is for Badger logging.
// logsDir is where the logs we're capturing are being stored.
func Open(ctx context.Context, logger logging.Logger, logsDir string) (*Store, error) {
	db, err := badger.Open(
		badger.DefaultOptions(logsDir).
			WithLogger(newLogger(logger, defaultLogLevel)),
	)
	if err != nil {
		return nil, err
	}
	return &Store{
		db:    db,
		idGen: gensym.NewULIDGenerator(ctx),
	}, nil
}

func (sto *Store) Close() error {
	return sto.db.Close()
}

func (sto *Store) GetLog(name string) store.Log {
	return &Log{
		db:    sto.db,
		idGen: sto.idGen,
		name:  name,
	}
}

func (sto *Store) NextLog(after store.Log) (store.Log, error) {
	prefix := []byte{}
	if after != nil {
		prefix = append([]byte(after.Name()), 255)
	}
	var nextName string
	if err := sto.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		// TODO: Something better than a linear key scan.
		it.Seek(prefix)
		for it.Valid() {
			it.Next()
			if !it.ValidForPrefix(prefix) {
				if it.Valid() {
					nextName = mustLogFromKey(it.Item().Key())
				}
				break
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if nextName == "" {
		return nil, nil
	}
	return sto.GetLog(nextName), nil
}

func (log *Log) Name() string {
	return log.name
}
