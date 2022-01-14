package resolvers

type TaskResolver struct {
	Q RootResolver
	TaskRow
}

type TaskRow struct {
	ID       string  `db:"id"`
	ParentID *string `db:"parent_id"`
}
