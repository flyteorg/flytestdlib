package synccache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lyft/flytestdlib/utils"

	"github.com/stretchr/testify/assert"
)

const fakeCacheItemValueLimit = 10

type fakeCacheItem struct {
	id  string
	val int
}

func (f fakeCacheItem) ID() string {
	return f.id
}

func syncFakeItem(_ context.Context, obj CacheItem) (CacheItem, CacheSyncAction, error) {
	item := obj.(fakeCacheItem)
	if item.val == fakeCacheItemValueLimit {
		// After the item has gone through ten update cycles, leave it unchanged
		return obj, Unchanged, nil
	}

	return fakeCacheItem{id: item.ID(), val: item.val + 1}, Update, nil
}

func syncFakeItemAlwaysDelete(_ context.Context, obj CacheItem) (CacheItem, CacheSyncAction, error) {
	return obj, Delete, nil
}

func TestCacheTwo(t *testing.T) {
	testResyncPeriod := time.Millisecond
	rateLimiter := utils.NewRateLimiter("mockLimiter", 100, 1)

	t.Run("normal operation", func(t *testing.T) {
		// the size of the cache is at least as large as the number of items we're storing
		cache, err := NewAutoRefreshCache(syncFakeItem, rateLimiter, testResyncPeriod, 10, 1, nil)
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cache.Start(ctx)

		// Create ten items in the cache
		for i := 1; i <= 10; i++ {
			_, err := cache.GetOrCreate(fakeCacheItem{
				id:  fmt.Sprintf("%d", i),
				val: 0,
			})
			assert.NoError(t, err)
		}

		// Wait half a second for all resync periods to complete
		time.Sleep(500 * time.Millisecond)
		for i := 1; i <= 10; i++ {
			item, err := cache.Get(fmt.Sprintf("%d", i))
			assert.NoError(t, err)
			assert.Equal(t, 10, item.(fakeCacheItem).val)
		}
		cancel()
	})

	t.Run("deleting objects from cache", func(t *testing.T) {
		// the size of the cache is at least as large as the number of items we're storing
		cache, err := NewAutoRefreshCache(syncFakeItemAlwaysDelete, rateLimiter, testResyncPeriod, 10, 1, nil)
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cache.Start(ctx)

		// Create ten items in the cache
		for i := 1; i <= 10; i++ {
			_, err = cache.GetOrCreate(fakeCacheItem{
				id:  fmt.Sprintf("%d", i),
				val: 0,
			})
			assert.NoError(t, err)
		}

		// Wait for all resync periods to complete
		time.Sleep(50 * time.Millisecond)
		for i := 1; i <= 10; i++ {
			obj, err := cache.Get(fmt.Sprintf("%d", i))
			assert.NoError(t, err)
			assert.Nil(t, obj)
		}
		cancel()
	})
}
