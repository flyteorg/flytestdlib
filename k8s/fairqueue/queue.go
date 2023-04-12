package fairqueue

import (
	"context"
	"sync"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/lyft/flytestdlib/logger"
)

// Implementation of a FairQueue per Namespace.
// The idea behind this queue is that in the case of asymmetrical load across namespaces it penalizes the
// namespace that has larger load, but tries to be fair across namespaces.
// In case a system is modeled where every tenant has a separate namespace, each tenant is promised an execution
// guarantee which is proportional to the number of workers available and does not rely on the length of the queue
// contributed by other tenants.
// The queue also implements workqueue.Interface
type PerNamespaceFairQueue struct {
	perNamespaceQueue *IndexedCircularBuffer
	cond              *sync.Cond
	shuttingDown      bool
}

// Adds an item to the queue for the namespace where the item belongs. The method uses ache.SplitMetaNamespaceKey
// to identify the namespace of the item. Default or no qualified namespace gets added to a common queue signified by
// an empty namespace
func (m *PerNamespaceFairQueue) Add(item interface{}) {

	key, ok := item.(string)
	if !ok {
		logger.Error(context.TODO(), "failed to add item to the workQueue, item type is not string.")
		return
	}
	// Convert the namespace/name string into a distinct namespace and name
	namespace, _, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		logger.Error(context.TODO(), "failed to add item to the workQueue, item incorrectly formatted. Error: %s", err.Error())
		return
	}
	tenantID := namespace

	m.cond.L.Lock()
	defer m.cond.L.Unlock()
	qObj, _ := m.perNamespaceQueue.GetOrDefault(tenantID, func() interface{} {
		return newDedupingQueue()
	})

	if q, ok := qObj.(*dedupingQueue); !ok {
		logger.Fatal(context.TODO(), "incorrect queue type for tenant ID [%s]", tenantID)
	} else {
		q.Add(item)
	}
	m.cond.Signal()
}

// Provides a length of all the available items across all queues
func (m *PerNamespaceFairQueue) Len() int {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()
	length := 0
	m.perNamespaceQueue.Range(func(v interface{}) bool {
		q, ok := v.(*dedupingQueue)
		if !ok {
			return false
		}
		length += q.Len()
		return true
	})
	return length
}

// Retrieves an item from the next logical namespace in a circular buffer if available. It iterates through all the
// namespaces in this order, until it reaches the current position. Once back at the current position, this call `Get`
// blocks for a new item to be inserted in any one of the namespaces.
func (m *PerNamespaceFairQueue) Get() (item interface{}, shutdown bool) {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()
	for {
		if m.shuttingDown {
			return nil, m.shuttingDown
		}
		var item interface{}
		m.perNamespaceQueue.RangeNext(func(v interface{}) bool {
			q, ok := v.(*dedupingQueue)
			if !ok {
				logger.Errorf(context.TODO(), "found a non conforming type in queue, skipping")
			}
			if i := q.Get(); i != nil {
				item = i
				return false
			}
			return true
		})
		if item != nil {
			return item, m.shuttingDown
		}
		m.cond.Wait()
	}

}

// Removes the item from the queue permanently
func (m *PerNamespaceFairQueue) Done(item interface{}) {

	//m.metrics.done(item)

	key, ok := item.(string)
	if !ok {
		logger.Error(context.TODO(), "failed to add item to the workQueue, item type is not string.")
		return
	}
	// Convert the namespace/name string into a distinct namespace and name
	namespace, _, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		logger.Error(context.TODO(), "failed to add item to the workQueue, item incorrectly formatted. Error: %s", err.Error())
		return
	}
	tenantID := namespace

	m.cond.L.Lock()
	defer m.cond.L.Unlock()
	v, ok := m.perNamespaceQueue.Get(tenantID)
	if !ok {
		return
	}
	q, ok := v.(*dedupingQueue)
	if !ok {
		return
	}
	if q.Done(item) == true {
		m.cond.Signal()
	}
}

// ShutDown will cause q to ignore all new items added to it. As soon as the
// worker goroutines have drained the existing items in the queue, they will be
// instructed to exit.
func (m *PerNamespaceFairQueue) ShutDown() {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()
	m.shuttingDown = true
	m.cond.Broadcast()
}

// Returns a boolean that indicates if a shutdown has been initiated for the queue.
func (m *PerNamespaceFairQueue) ShuttingDown() bool {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()

	return m.shuttingDown
}

func New() workqueue.Interface {
	return &PerNamespaceFairQueue{
		perNamespaceQueue: NewIndexedCircularBuffer(),
		cond:              sync.NewCond(&sync.Mutex{}),
	}
}
