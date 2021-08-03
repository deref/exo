package storage_test

import (
	"testing"

	"github.com/deref/exo/storage"
	"github.com/stretchr/testify/assert"
)

func TestKVEngine(t *testing.T) {
	kv := storage.NewMemoryKVEngine()

	{
		// Key not yet set.
		val, err := storage.GetAtomic(kv, []byte("hello"))
		assert.NoError(t, err)
		assert.Nil(t, val)
	}

	{
		// Set/Get
		assert.NoError(t, storage.SetAtomic(kv, []byte("hello"), []byte("world")))
		val, err := storage.GetAtomic(kv, []byte("hello"))
		assert.NoError(t, err)
		assert.Equal(t, []byte("world"), val)
	}
}

func TestScan(t *testing.T) {
	kv := storage.NewMemoryKVEngine()

	storage.Transact(kv, func(txn storage.WriteTransaction) error {
		_ = kv.Set(txn, []byte("hi world"), []byte("1"))
		_ = kv.Set(txn, []byte("hello world"), []byte("2"))
		_ = kv.Set(txn, []byte("hi there"), []byte("3"))
		return nil
	})

	it := kv.Scan(kv.ReadTransaction(), storage.ScanArgs{
		Prefix: []byte("hi"),
	})
	entries, err := storage.Collect(it)
	assert.NoError(t, err)

	assert.Equal(t, []storage.ScanEntry{
		{
			Key:   []byte("hi there"),
			Value: []byte("3"),
		},
		{
			Key:   []byte("hi world"),
			Value: []byte("1"),
		},
	}, entries)
}
