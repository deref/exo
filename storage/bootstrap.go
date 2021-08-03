package storage

var tableSchema = NewSchema(
	ElementDescriptor{
		Name: "table_oid",
		Type: TypeInt32,
	},
	ElementDescriptor{
		Name: "schema_oid",
		Type: TypeInt32,
	},
	ElementDescriptor{
		Name: "table_name",
		Type: TypeUnicode,
	},
)

func tableTable(db *Database) *table {
	return &table{
		db:     db,
		oid:    oidTable,
		name:   "table",
		schema: tableSchema,
		serde:  NewSchematizedRowSerde(tableSchema),
		// Indexes: []IndexOptions{
		// 	SingleColumnIndex("table_name"),
		// },
	}
}

var schemaSchema = NewSchema(
	ElementDescriptor{
		Name: "schema_oid",
		Type: TypeInt32,
	},
	ElementDescriptor{
		Name: "serialized_schema",
		Type: TypeBytes,
	},
)

func schemaTable(db *Database) *table {
	return &table{
		db:     db,
		oid:    oidSchema,
		name:   "schema",
		schema: schemaSchema,
		serde:  NewSchematizedRowSerde(schemaSchema),
	}
}

var indexSchema = NewSchema(
	ElementDescriptor{
		Name: "index_oid",
		Type: TypeInt32,
	},
	ElementDescriptor{
		Name: "column_names",
		Type: TypeBytes,
	},
)

var bootstrapKey = NewTuple(int32(oidBootstrap)).Serialize()

func (db *Database) isBootstrapped() (bool, error) {
	val, err := GetAtomic(db.store, bootstrapKey)
	return val != nil, err
}

func (db *Database) bootstrap() error {
	db.nextOID = 100 // XXX: How to initialize on subsequent runs?

	return Transact(db.store, func(txn WriteTransaction) error {
		// Create table for schemas.
		if err := schemaTable(db).Create(txn); err != nil {
			return err
		}

		// Create table for tables.
		if err := tableTable(db).Create(txn); err != nil {
			return err
		}

		return db.store.Set(txn, bootstrapKey, []byte{0})
	})
}
