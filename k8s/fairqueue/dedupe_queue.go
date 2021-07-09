package fairqueue

// This is a non thread-safe version of of the workqueue.Type. It has been liberally copied
// with some functions eliminated. The thread-safe version should be built on top of this implementation.
type dedupingQueue struct {
	// queue defines the order in which we will work on items. Every
	// element of queue should be in the dirty set and not in the
	// processing set.
	queue []t

	// dirty defines all of the items that need to be processed.
	dirty set

	// Things that are currently being processed are in the processing set.
	// These things may be simultaneously in the dirty set. When we finish
	// processing something and remove it from this set, we'll check if
	// it's in the dirty set, and if so, add it to the queue.
	processing set
}

type empty struct{}
type t interface{}
type set map[t]empty

func (s set) has(item t) bool {
	_, exists := s[item]
	return exists
}

func (s set) insert(item t) {
	s[item] = empty{}
}

func (s set) delete(item t) {
	delete(s, item)
}

// Add marks item as needing processing.
func (q *dedupingQueue) Add(item interface{}) {
	if q.dirty.has(item) {
		return
	}

	q.dirty.insert(item)
	if q.processing.has(item) {
		return
	}

	q.queue = append(q.queue, item)
}

// Len returns the current queue length, for informational purposes only. You
// shouldn't e.g. gate a call to Add() or Get() on Len() being a particular
// value, that can't be synchronized properly.
func (q *dedupingQueue) Len() int {
	return len(q.queue)
}

// Get blocks until it can return an item to be processed. If shutdown = true,
// the caller should end their goroutine. You must call Done with item when you
// have finished processing it.
func (q *dedupingQueue) Get() (item interface{}) {
	if len(q.queue) == 0 {
		return nil
	}

	item, q.queue = q.queue[0], q.queue[1:]

	q.processing.insert(item)
	q.dirty.delete(item)

	return item
}

// Done marks item as done processing, and if it has been marked as dirty again
// while it was being processed, it will be re-added to the queue for
// re-processing. Returns if the item was re-added for processing.
func (q *dedupingQueue) Done(item interface{}) bool {
	if q.processing.has(item) {
		q.processing.delete(item)
		if q.dirty.has(item) {
			q.queue = append(q.queue, item)
			return true
		}
	}
	return false
}

func newDedupingQueue() *dedupingQueue {
	return &dedupingQueue{
		dirty:      set{},
		processing: set{},
	}
}
