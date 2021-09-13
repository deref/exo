package contextutil

import "context"

type cloneFunc = func(src, dest context.Context) context.Context

var cloneFuncs []cloneFunc

// DeriveBackgroundContext creates a new context from `ctx` that contains all of the global
// dependencies present in `ctx` but does not have any cancellation or timeout associated
// with it. This is useful when we have a context tied to some request, and we need to spawn
// an asynchronous job that may still need access to the context global values.
//
// Note that to clone the values from the old context to the new context, a module needs to
// call RegisterContextCloner with a function that is able to copy the context global that
// it manages from the old context to the new one.
func DeriveBackgroundContext(ctx context.Context) context.Context {
	newCtx := context.Background()
	for _, cf := range cloneFuncs {
		newCtx = cf(ctx, newCtx)
	}
	return newCtx
}

func RegisterContextCloner(fn cloneFunc) {
	cloneFuncs = append(cloneFuncs, fn)
}
