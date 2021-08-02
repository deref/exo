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
		val, err := kv.Get([]byte("hello"))
		assert.NoError(t, err)
		assert.Nil(t, val)
	}

	{
		// Set/Get
		assert.NoError(t, kv.Set([]byte("hello"), []byte("world")))
		val, err := kv.Get([]byte("hello"))
		assert.NoError(t, err)
		assert.Equal(t, []byte("world"), val)
	}
}

func TestScan(t *testing.T) {
	kv := storage.NewMemoryKVEngine()

	_ = kv.Set([]byte("hi world"), []byte("1"))
	_ = kv.Set([]byte("hello world"), []byte("2"))
	_ = kv.Set([]byte("hi there"), []byte("3"))

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
