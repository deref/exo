package resolvers

import (
	"context"
)

type ProgressResolver struct {
	Current int32
	Total   int32
}

type ProgressInput struct {
	Current int32
	Total   int32
}

func (r *ProgressResolver) Percent() float64 {
	return float64(r.Current) / float64(r.Total) * 100.0
}

func (r *MutationResolver) reportProgress(ctx context.Context, progress ProgressInput) {
	t := CurrentTask(ctx)
	if t == nil {
		return
	}
	_, err := r.updateTask(ctx, t.ID, t.WorkerID, &progress)
	r.Logger.Infof("error reporting progress on task %q: %v", err)
}
