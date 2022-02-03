package resolvers

import (
	"context"

	"github.com/deref/exo/internal/api"
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
	ctxVars := api.CurrentContextVariables(ctx)
	if ctxVars == nil || ctxVars.TaskID == "" {
		return
	}
	_, err := r.updateTask(ctx, ctxVars.TaskID, ctxVars.WorkerID, &progress)
	r.SystemLog.Infof("error reporting progress on task %q: %v", err)
}
