package storage

type IndexOptions struct {
	ColumnNames []string
}

func MultiColumnIndex(columns ...string) IndexOptions {
	return IndexOptions{
		ColumnNames: columns,
	}
}

func SingleColumnIndex(column string) IndexOptions {
	return MultiColumnIndex(column)
}
