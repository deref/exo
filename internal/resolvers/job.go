package resolvers

type JobResolver struct {
	Q RootResolver
	JobRow
}

type JobRow struct {
	ID string `db:"id"`
}
