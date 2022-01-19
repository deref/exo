package resolvers

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ProgressResolver struct {
	Current int32
	Total   int32
}

func (r *ProgressResolver) Percent() float64 {
	return float64(r.Current) / float64(r.Total) * 100.0
}

type TaskTracker struct {
	DB *sqlx.DB
}

func NewTaskTracker() *TaskTracker {
	panic("TODO")
}

func (tt *TaskTracker) ReportProgress(ctx context.Context, current, total int32) {
	panic("TODO")
}

func ReportProgress(ctx context.Context, current, total int32) {
	tt := CurrentTaskTracker(ctx)
	if tt == nil {
		return
	}
	tt.ReportProgress(ctx, current, total)
}
