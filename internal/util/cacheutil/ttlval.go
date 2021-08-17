package cacheutil

import (
	"context"
	"time"
)

type valFunc = func() (interface{}, error)

// TTLVal is a type that holds some value and lazily refreshes it when the
// value is requested if a configurable duration has passed since it was
// last refreshed.
type TTLVal struct {
	val    interface{}
	getVal valFunc
	runAt  time.Time
	ttl    time.Duration
}

func NewTTLVal(getVal valFunc, ttl time.Duration) *TTLVal {
	return &TTLVal{
		getVal: getVal,
		ttl:    ttl,
	}
}

func (v *TTLVal) Get(ctx context.Context) (interface{}, error) {
	if time.Now().Sub(v.runAt) > v.ttl {
		val, err := v.getVal()
		if err != nil {
			return nil, err
		}
		v.runAt = time.Now()
		v.val = val
	}
	return v.val, nil
}
