package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lyft/flytestdlib/promutils"

	"github.com/lyft/flytestdlib/utils"

	"github.com/stretchr/testify/assert"
)

const fakeCacheItemValueLimit = 10

type fakeCacheItem struct {
	val int
}

func syncFakeItem(_ context.Context, batch Batch) ([]ItemSyncResponse, error) {
	items := make([]ItemSyncResponse, 0, len(batch))
	for _, obj := range batch {
		item := obj.GetItem().(fakeCacheItem)
		if item.val == fakeCacheItemValueLimit {
			// After the item has gone through ten update cycles, leave it unchanged
			continue
		}

		items = append(items, ItemSyncResponse{
			ID: obj.GetID(),
			Item: fakeCacheItem{
				val: item.val + 1,
			},
			Action: Update,
		})
	}

	return items, nil
}

func TestCacheTwo(t *testing.T) {
	testResyncPeriod := time.Millisecond
	rateLimiter := utils.NewRateLimiter("mockLimiter", 100, 100)

	t.Run("normal operation", func(t *testing.T) {
		// the size of the cache is at least as large as the number of items we're storing
		cache, err := NewAutoRefreshCache(syncFakeItem, rateLimiter, testResyncPeriod, 10, promutils.NewTestScope())
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cache.Start(ctx)

		// Create ten items in the cache
		for i := 1; i <= 10; i++ {
			_, err := cache.GetOrCreate(fmt.Sprintf("%d", i), fakeCacheItem{
				val: 0,
			})
			assert.NoError(t, err)
		}

		// Wait half a second for all resync periods to complete
		time.Sleep(500 * time.Millisecond)
		for i := 1; i <= 10; i++ {
			item := cache.Get(fmt.Sprintf("%d", i))
			assert.Equal(t, 10, item.(fakeCacheItem).val)
		}
		cancel()
	})
}
