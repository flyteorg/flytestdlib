package cache

import (
	"context"
	"time"

	"github.com/lyft/flytestdlib/logger"

	"github.com/lyft/flytestdlib/utils"

	lru "github.com/hashicorp/golang-lru"
	"github.com/lyft/flytestdlib/promutils"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/util/wait"
)

type ItemID = string

// AutoRefresh with regular GetOrCreate and Delete along with background asynchronous refresh. Caller provides
// callbacks for create, refresh and delete item.
// The cache doesn't provide apis to update items.
type AutoRefresh interface {
	// starts background refresh of items
	Start(ctx context.Context)

	// Get item by id if exists else null
	Get(id ItemID) Item

	// Get object if exists else create it
	GetOrCreate(item Item) (Item, error)
}

type metrics struct {
	SystemErrors prometheus.Counter
	Evictions    prometheus.Counter
	scope        promutils.Scope
}

type Item interface {
	ID() ItemID
}

type ItemSyncResponse struct {
	Item   Item
	Action SyncAction
}

// Possible actions for the cache to take as a result of running the sync function on any given cache item
type SyncAction int

const (
	Unchanged SyncAction = iota

	// The item returned has been updated and should be updated in the cache
	Update
)

// Your implementation of this function for your cache instance is responsible for returning
//   1. The new Item, and
//   2. What action should be taken.  The sync function has no insight into your object, and needs to be
//      told explicitly if the new item is different from the old one.
type SyncFunc func(ctx context.Context, batch []Item) (
	updatedBatch []ItemSyncResponse, err error)

// Your implementation of this function for your cache instance is responsible for subdividing
// the list of cache items into batches.
type CreateBatchesFunc func(ctx context.Context, snapshot []Item) (batches [][]Item, err error)

// Thread-safe general purpose auto-refresh cache that watches for updates asynchronously for the keys after they are added to
// the cache. An item can be inserted only once.
//
// Get reads from sync.map while refresh is invoked on a snapshot of keys. Cache eventually catches up on deleted items.
//
// Sync is run as a fixed-interval-scheduled-task, and is skipped if sync from previous cycle is still running.
type autoRefresh struct {
	metrics         metrics
	syncCb          SyncFunc
	createBatchesCb CreateBatchesFunc
	lruMap          *lru.Cache
	syncRateLimiter utils.RateLimiter
	syncPeriod      time.Duration
}

func getEvictionFunction(counter prometheus.Counter) func(key interface{}, value interface{}) {
	return func(_ interface{}, _ interface{}) {
		counter.Inc()
	}
}

func SingleItemBatches(ctx context.Context, snapshot []Item) (batches [][]Item, err error) {
	res := make([][]Item, 0, len(snapshot))
	for _, item := range snapshot {
		res = append(res, []Item{item})
	}

	return res, nil
}

func newMetrics(scope promutils.Scope) metrics {
	return metrics{
		Evictions:    scope.MustNewCounter("lru_evictions", "Counter for evictions from LRU."),
		SystemErrors: scope.MustNewCounter("sync_errors", "Counter for sync errors."),
		scope:        scope,
	}
}

func NewAutoRefreshBatchedCache(createBatches CreateBatchesFunc, syncCb SyncFunc, syncRateLimiter utils.RateLimiter, resyncPeriod time.Duration,
	size int, scope promutils.Scope) (AutoRefresh, error) {

	metrics := newMetrics(scope)
	lruCache, err := lru.NewWithEvict(size, getEvictionFunction(metrics.Evictions))
	if err != nil {
		return nil, err
	}

	cache := &autoRefresh{
		metrics:         metrics,
		createBatchesCb: createBatches,
		syncCb:          syncCb,
		lruMap:          lruCache,
		syncRateLimiter: syncRateLimiter,
		syncPeriod:      resyncPeriod,
	}

	return cache, nil
}

func NewAutoRefreshCache(syncCb SyncFunc, syncRateLimiter utils.RateLimiter, resyncPeriod time.Duration,
	size int, scope promutils.Scope) (AutoRefresh, error) {
	return NewAutoRefreshBatchedCache(SingleItemBatches, syncCb, syncRateLimiter, resyncPeriod, size, scope)
}

func (w *autoRefresh) Start(ctx context.Context) {
	go wait.Until(func() {
		err := w.syncRateLimiter.Wait(ctx)
		if err != nil {
			logger.Errorf(ctx, "Failed to get sync token. Error: %v", err)
		}

		err = w.sync(ctx)
		if err != nil {
			logger.Errorf(ctx, "Failed to sync. Error: %v", err)
		}
	}, w.syncPeriod, ctx.Done())
}

func (w *autoRefresh) Get(id string) Item {
	if val, ok := w.lruMap.Get(id); ok {
		return val.(Item)
	}
	return nil
}

// Return the item if exists else create it.
// Create should be invoked only once. recreating the object is not supported.
func (w *autoRefresh) GetOrCreate(item Item) (Item, error) {
	if val, ok := w.lruMap.Get(item.ID()); ok {
		return val.(Item), nil
	}

	w.lruMap.Add(item.ID(), item)
	return item, nil
}

// This function is called internally by its own timer. Roughly, it will,
//  - List keys
//  - Create batches of keys based on createBatchesCb
//  - For each batch of the keys, call syncCb, which tells us if the items have been updated
//    - If any has, then overwrite the item in the cache.
//
// What happens when the number of things that a user is trying to keep track of exceeds the size
// of the cache?  Trivial case where the cache is size 1 and we're trying to keep track of two things.
//  * Plugin asks for update on item 1 - cache evicts item 2, stores 1 and returns it unchanged
//  * Plugin asks for update on item 2 - cache evicts item 1, stores 2 and returns it unchanged
//  * Sync loop updates item 2, repeat
func (w *autoRefresh) sync(ctx context.Context) error {
	keys := w.lruMap.Keys()
	snapshot := make([]Item, 0, len(keys))
	for _, k := range keys {
		// If not ok, it means evicted between the item was evicted between getting the keys and this update loop
		// which is fine, we can just ignore.
		if value, ok := w.lruMap.Peek(k); ok {
			snapshot = append(snapshot, value.(Item))
		}
	}

	batches, err := w.createBatchesCb(ctx, snapshot)
	if err != nil {
		return err
	}

	for _, batch := range batches {
		updatedBatch, err := w.syncCb(ctx, batch)
		if err != nil {
			logger.Error(ctx, "failed to get latest copy of a batch. Error: %v", err)
			continue
		}

		for _, item := range updatedBatch {
			if item.Action == Update {
				// Add adds the item if it has been evicted or updates an existing one.
				w.lruMap.Add(item.Item.ID(), item.Item)
			}
		}
	}

	return nil
}
