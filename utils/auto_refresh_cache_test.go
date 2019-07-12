package utils

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

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

func BenchmarkCache(b *testing.B) {
	testResyncPeriod := time.Millisecond
	rateLimiter := NewRateLimiter("mockLimiter", 100, 1)
	// the size of the cache is at least as large as the number of items we're storing
	itemCount := b.N
	cache, err := NewAutoRefreshCache(syncFakeItem, rateLimiter, testResyncPeriod, itemCount*2, nil)
	assert.NoError(b, err)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	//cache.Start(ctx)

	startIdx := 1

	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(itemCount)
		// Create items in the cache
		for i := 1; i <= itemCount; i++ {
			go func(itemId int) {
				defer wg.Done()
				_, err := cache.GetOrCreate(fakeCacheItem{
					id:  fmt.Sprintf("%d", itemId),
					val: itemId,
				})

				assert.NoError(b, err)
			}(i + startIdx)
		}

		wg.Wait()

		// Wait half a second for all resync periods to complete
		wg.Add(itemCount)
		for i := 1; i <= itemCount; i++ {
			go func(itemId int) {
				defer wg.Done()
				item := cache.Get(fmt.Sprintf("%d", itemId))
				assert.NotNil(b, item, "item #%v", itemId)
				if item != nil {
					assert.Equal(b, strconv.Itoa(itemId), item.(fakeCacheItem).ID())
				}
			}(i + startIdx)
		}
		wg.Wait()
		startIdx += itemCount
	}
}

func TestCacheTwo(t *testing.T) {
	testResyncPeriod := time.Millisecond
	rateLimiter := NewRateLimiter("mockLimiter", 100, 1)

	t.Run("normal operation", func(t *testing.T) {
		// the size of the cache is at least as large as the number of items we're storing
		cache, err := NewAutoRefreshCache(syncFakeItem, rateLimiter, testResyncPeriod, 10, nil)
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
			item := cache.Get(fmt.Sprintf("%d", i))
			assert.Equal(t, 10, item.(fakeCacheItem).val)
		}
		cancel()
	})

	t.Run("deleting objects from cache", func(t *testing.T) {
		// the size of the cache is at least as large as the number of items we're storing
		cache, err := NewAutoRefreshCache(syncFakeItemAlwaysDelete, rateLimiter, testResyncPeriod, 10, nil)
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
			obj := cache.Get(fmt.Sprintf("%d", i))
			assert.Nil(t, obj)
		}
		cancel()
	})
}
