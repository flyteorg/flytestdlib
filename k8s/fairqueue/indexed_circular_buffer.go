package fairqueue

import "container/ring"

// This is an implementation of a circular buffer using the container/ring package in golang
// This buffer is specialized to allow accessing a specific entry by its key.
// Iterating through the buffer using the Next method is stateful and always returns the sequentially next item in the buffer
type IndexedCircularBuffer struct {
	head    *ring.Ring
	current *ring.Ring
	index   map[string]*ring.Ring
}

// Returns sequentially next item in the Circular buffer. This method is stateful
func (c *IndexedCircularBuffer) Next() interface{} {
	if c.current.Next() == c.head {
		c.current = c.head
	}
	v := c.current.Next()
	c.current = v
	return v.Value
}

// Returns the value in the circular buffer that matches the given Key. If the key is not found, returns a false
func (c *IndexedCircularBuffer) Get(key string) (interface{}, bool) {
	v, ok := c.index[key]
	if !ok {
		return nil, false
	}
	return v.Value, true
}

// Returns the value in the circular buffer that matches the given Key and true indicating that it existed.
// If the key is not found, uses the provided function to get a default item, adds it to the buffer and returns it with a false
// indicating that the value was created
func (c *IndexedCircularBuffer) GetOrDefault(key string, defaultItemGetter func() interface{}) (interface{}, bool) {
	v, ok := c.index[key]
	if !ok {
		i := defaultItemGetter()
		c.Add(key, i)
		return i, false
	}
	return v.Value, true
}

// Adds a new Key,item at the last position in the buffer in an order preserving way
func (c *IndexedCircularBuffer) Add(key string, item interface{}) bool {
	if _, ok := c.index[key]; ok {
		return false
	}
	r := ring.New(1)
	r.Value = item
	c.index[key] = r
	last := c.head.Prev()
	last.Link(r)
	r.Link(c.head)
	return true
}

// Returns a length of the circular buffer
func (c *IndexedCircularBuffer) Len() int {
	return len(c.index)
}

// Iterates over all the elements in the circular buffer in the order of insertion
func (c *IndexedCircularBuffer) Range(do func(v interface{}) bool) {
	for ptr := c.head.Next(); ptr != c.head; ptr = ptr.Next() {
		if do(ptr.Value) == false {
			return
		}
	}
}

// Checks if the current is pointing to the first element in the buffer
func (c *IndexedCircularBuffer) IsCurrentAtHead() bool {
	return c.current == c.head
}

// Iterates through the buffer order-preserving. Stops iterating when the provided function returns false or
func (c *IndexedCircularBuffer) RangeNext(do func(v interface{}) bool) {
	ptr := c.current.Next()
	for ; ptr != c.current; ptr = ptr.Next() {
		if ptr != c.head {
			if do(ptr.Value) == false {
				// Move the current to the iterated location
				c.current = ptr
				return
			}
		}
	}
	if ptr != c.head {
		_ = do(ptr.Value)
	}
}

func NewIndexedCircularBuffer() *IndexedCircularBuffer {
	head := ring.New(1)
	return &IndexedCircularBuffer{
		head:    head,
		current: head,
		index:   map[string]*ring.Ring{},
	}
}
