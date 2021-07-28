package storage_test

import (
	"fmt"
	"testing"

	"github.com/deref/exo/storage"
	"github.com/stretchr/testify/assert"
)

func TestRawOps(t *testing.T) {
	kv := storage.NewMemoryKVEngine()

	schema := storage.NewSchema(
		storage.ElementDescriptor{
			Name: "partition",
			Type: storage.TypeUint64,
		},
		storage.ElementDescriptor{
			Name: "id",
			Type: storage.TypeUint64,
		},
		storage.ElementDescriptor{
			Name: "name",
			Type: storage.TypeUnicode,
		},
	)
	serde := storage.NewSchematizedRowSerde(schema)

	tup := storage.NewTupleWithSchema(schema)
	tup.SetUint64(0, uint64(3))
	tup.SetUint64(1, uint64(81235))
	tup.SetUnicode(2, "The quick brown fox...")

	key := tup.Without(2).Serialize()
	val, err := serde.Serialize(tup)
	assert.NoError(t, err)

	kv.Set(key, val)

	it := kv.Scan(kv.ReadTransaction(), storage.ScanArgs{
		// Only partition.
		Prefix: tup.Without(2).Without(1).Serialize(),
	})
	entries, err := storage.Collect(it)
	assert.NoError(t, err)

	for _, entry := range entries {
		keyTup, err := storage.Deserialize(entry.Key)
		assert.NoError(t, err)
		obj, err := serde.Deserialize(entry.Value)
		assert.NoError(t, err)
		fmt.Printf("%s\n\t%s\n", keyTup, obj)
	}
	assert.True(t, false)
}

// func TestTable(t *testing.T) {
// 	kv := storage.NewMemoryKVEngine()

// 	schema := storage.NewSchema(
// 		storage.ElementDescriptor{
// 			Name: "id",
// 			Type: storage.TypeUint64,
// 		},
// 		storage.ElementDescriptor{
// 			Name: "name",
// 			Type: storage.TypeUnicode,
// 		},
// 	)
// 	t1 := storage.NewTable(storage.TableOptions{
// 		OID:    1,
// 		Name:   "my-table",
// 		Schema: schema,
// 	})
// }
