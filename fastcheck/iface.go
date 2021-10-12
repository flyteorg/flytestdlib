package fastcheck

import (
	"context"
	"github.com/flyteorg/flytestdlib/promutils"
	"github.com/flyteorg/flytestdlib/promutils/labeled"
)

// Filter provides an interface to check if a Key of type []byte was ever seen.
// The size of the cache is dependent on the id size and the initialization. It may also vary based on the implementation
// For example an LRU cache, may have an overhead because of the use of a HashMap with loading factor and collision
// resolution
type Filter interface {
	// Contains returns a True if the id was previously seen or false otherwise
	// It may return a false, even if a item may have previously occurred.
	Contains(ctx context.Context, id []byte) bool

	// Adds the element id to the Filter
	Add(ctx context.Context, id []byte) (evicted bool)
}

// Every implementation of the Filter Interface provides these metrics
type Metrics struct {
	Hit labeled.Counter
	Miss labeled.Counter
}

func NewMetrics(scope promutils.Scope) Metrics {
	return Metrics{
		Hit: labeled.NewCounter("cache_hit", "Indicates that the item was found in the cache", scope),
		Miss: labeled.NewCounter("cache_miss", "Indicates that the item was found in the cache", scope),
	}
}