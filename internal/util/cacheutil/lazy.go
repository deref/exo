package cacheutil

import "sync"

type Lazy[T any] struct {
	once  sync.Once
	value T
	thunk func() T
}

func NewLazy[T any](thunk func() T) *Lazy[T] {
	return &Lazy[T]{
		thunk: thunk,
	}
}

func (z *Lazy[T]) Force() T {
	z.once.Do(func() {
		z.value = z.thunk()
	})
	return z.value
}
