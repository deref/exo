package storage

type TableOptions struct {
	OID     int64
	Name    string
	Schema  *Schema
	Indexes []IndexOptions
}

func NewTable(opts TableOptions) *table {
	return &table{
		opts: opts,
	}
}

type table struct {
	opts TableOptions
}
