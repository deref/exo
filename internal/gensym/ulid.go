package gensym

import (
	"context"
	"math/rand"
	"sync"

	"github.com/deref/exo/internal/chrono"
	"github.com/oklog/ulid/v2"
)

// A concurrency-safe source of monotonic ULIDs.
// See <https://github.com/ulid/spec>.
type ULIDGenerator struct {
	mu      sync.Mutex
	entropy *ulid.MonotonicEntropy
}

func NewULIDGenerator(ctx context.Context) *ULIDGenerator {
	seed := int64(chrono.NowNano(ctx))
	randRead := rand.New(rand.NewSource(seed))
	return &ULIDGenerator{
		entropy: ulid.Monotonic(randRead, 0),
	}
}

func (gen *ULIDGenerator) NextID(ctx context.Context) ([]byte, error) {
	// The math/rand generator is not thread-safe, so we have to guard access with a mutex.
	gen.mu.Lock()
	defer gen.mu.Unlock()
	ts := ulid.Timestamp(chrono.Now(ctx))
	id, err := ulid.New(ts, gen.entropy)
	if err != nil {
		return nil, err
	}
	return id.MarshalBinary()
}
