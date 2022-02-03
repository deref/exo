package api

import "context"

type contextKey int

const (
	variablesContextKey contextKey = iota + 1
)

type ContextVariables struct {
	JobID    string `json:"jobId"`
	TaskID   string `json:"id"`
	WorkerID string `json:"workerId"`
}

func ContextWithVariables(ctx context.Context, vars ContextVariables) context.Context {
	return context.WithValue(ctx, variablesContextKey, &vars)
}

func CurrentContextVariables(ctx context.Context) *ContextVariables {
	vars, _ := ctx.Value(variablesContextKey).(*ContextVariables)
	return vars
}
