package fastcheck

import (
	"context"

	"github.com/flyteorg/flytestdlib/promutils"
	cache "github.com/hashicorp/golang-lru"
)

// Implements the fastcheck.Filter interface using an underlying LRUCache from cache.Cache
type lruCacheFilter struct {
	lru     *cache.Cache
	metrics Metrics
}

// Simply uses Contains from the lruCacheFilter
func (l lruCacheFilter) Contains(ctx context.Context, id []byte) bool {
	v := l.lru.Contains(string(id))
	if v {
		l.metrics.Hit.Inc(ctx)
		return true
	}
	l.metrics.Miss.Inc(ctx)
	return false
}

func (l lruCacheFilter) Add(_ context.Context, id []byte) bool {
	return l.lru.Add(string(id), nil)
}

// Create a new fastcheck.Filter using an LRU cache of a fixed size
func NewLRUCacheFilter(size int, scope promutils.Scope) (Filter, error) {
	c, err := cache.New(size)
	if err != nil {
		return nil, err
	}
	return lruCacheFilter{
		lru:     c,
		metrics: NewMetrics(scope),
	}, nil
}
