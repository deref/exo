package storage

import "fmt"

var tableSchema = NewSchema(
	ElementDescriptor{
		Name: "oid",
		Type: TypeUint32,
	},
	ElementDescriptor{
		Name: "table_oid",
		Type: TypeUint32,
	},
	ElementDescriptor{
		Name: "schema_oid",
		Type: TypeUint32,
	},
	ElementDescriptor{
		Name: "table_name",
		Type: TypeUnicode,
	},
)

var bootstrapSchema = NewSchema(
	ElementDescriptor{
		Name: "oid",
		Type: TypeUint32,
	},
)

func isBootstrapped(store KVEngine) (bool, error) {
	t := NewTupleWithSchema(bootstrapSchema)
	t.SetUint32(0, uint32(oidBootstrap))

	val, err := store.Get(t.Serialize())

	return val != nil, err
}

func bootstrapStore(store KVEngine) error {
	fmt.Println("Bootstrapping store")

	// TODO: install other objects.

	t := NewTupleWithSchema(bootstrapSchema)
	t.SetUint32(0, uint32(oidBootstrap))
	return store.Set(t.Serialize(), []byte{0})
}
