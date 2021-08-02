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
	entry := entries[0]

	keyTup, err := storage.Deserialize(entry.Key)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NewTuple(uint64(3), uint64(81235)).Serialize(),
		keyTup.Serialize(),
	)

	obj, err := serde.Deserialize(entry.Value)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NewTuple(uint64(3), uint64(81235), "The quick brown fox...").Serialize(),
		obj.Serialize(),
	)

}

func TestTable(t *testing.T) {
	kv := storage.NewMemoryKVEngine()

	db, err := storage.NewDatabase(kv)
	assert.NoError(t, err)

	tbl := db.NewTable(storage.TableOptions{
		Name: "my_table",
		Schema: storage.NewSchema(
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
		),
		Indexes: []storage.IndexOptions{
			storage.MultiColumnIndex("partition", "id"), // TODO: Mark as primary key.
			storage.SingleColumnIndex("name"),
		},
	})
	err = tbl.Create()
	if !assert.NoError(t, err) {
		return
	}

	err = tbl.Insert(map[string]interface{}{
		"partition": uint64(0),
		"id":        uint64(80831),
		"name":      "Javier",
	})
	if !assert.NoError(t, err) {
		return
	}

	tbl, err = db.Table("my_table")
	if !assert.NoError(t, err) {
		return
	}

	tbl.Scan(func(t *storage.Tuple) bool {
		fmt.Println("Scanning:", t)
		return false
	})
	// err = tbl.InsertAll([]map[string]interface{}{
	// 	{
	// 		"partition": uint64(1),
	// 		"id":        uint64(123),
	// 		"name":      "Andrew",
	// 	},
	// 	{
	// 		"partition": uint64(1),
	// 		"id":        uint64(456),
	// 		"name":      "Diana",
	// 	},
	// 	{
	// 		"partition": uint64(2),
	// 		"id":        uint64(111),
	// 		"name":      "Andrew",
	// 	},
	// })
	// assert.NoError(t, err)

	// rows, err := db.Table("my-table").Where(storage.EQ("name", "Andrew")).Find()
	// assert.NoError(t, err)
	// assert.Equal(t, []map[string]interface{}{
	// 	{
	// 		"partition": uint64(1),
	// 		"id":        uint64(123),
	// 		"name":      "Andrew",
	// 	},
	// 	{
	// 		"partition": uint64(2),
	// 		"id":        uint64(111),
	// 		"name":      "Andrew",
	// 	},
	// }, rows)
}
