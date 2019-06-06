package fairqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedupingQueue(t *testing.T) {
	t.Run("get-empty", func(t *testing.T) {

		q := newDedupingQueue()
		assert.Equal(t, nil, q.Get())
	})

	t.Run("add-normal", func(t *testing.T) {
		q := newDedupingQueue()
		assert.False(t, q.dirty.has("x"))
		q.Add("x")
		assert.True(t, q.dirty.has("x"))
		assert.False(t, q.processing.has("x"))
		assert.Equal(t, 1, len(q.queue))
	})

	t.Run("add-when-processing", func(t *testing.T) {
		q := newDedupingQueue()
		q.Add("x")
		assert.True(t, q.dirty.has("x"))
		assert.False(t, q.processing.has("x"))
		assert.Equal(t, 1, len(q.queue))
		assert.Equal(t, 1, q.Len())

		assert.Equal(t, "x", q.Get())

		assert.False(t, q.dirty.has("x"))
		assert.True(t, q.processing.has("x"))
		assert.Equal(t, 0, len(q.queue))
		assert.Equal(t, 0, q.Len())

		q.Add("x")
		assert.True(t, q.dirty.has("x"))
		assert.True(t, q.processing.has("x"))
		assert.Equal(t, 0, len(q.queue))
		assert.Equal(t, 0, q.Len())

		q.Done("x")
		assert.True(t, q.dirty.has("x"))
		assert.False(t, q.processing.has("x"))
		assert.Equal(t, 1, len(q.queue))
		assert.Equal(t, 1, q.Len())

	})

	t.Run("add-when-dirty", func(t *testing.T) {
		q := newDedupingQueue()
		q.Add("x")
		assert.True(t, q.dirty.has("x"))
		assert.False(t, q.processing.has("x"))
		assert.Equal(t, 1, len(q.queue))
		assert.Equal(t, 1, q.Len())

		q.Add("x")
		assert.True(t, q.dirty.has("x"))
		assert.False(t, q.processing.has("x"))
		assert.Equal(t, 1, len(q.queue))
		assert.Equal(t, 1, q.Len())

		q.Done("x")
		assert.True(t, q.dirty.has("x"))
		assert.False(t, q.processing.has("x"))
		assert.Equal(t, 1, len(q.queue))
		assert.Equal(t, 1, q.Len())

	})

	t.Run("done-processing-one", func(t *testing.T) {
		q := newDedupingQueue()
		q.Add("x")
		assert.True(t, q.dirty.has("x"))
		assert.False(t, q.processing.has("x"))
		assert.Equal(t, 1, len(q.queue))
		assert.Equal(t, 1, q.Len())

		assert.Equal(t, "x", q.Get())

		assert.False(t, q.dirty.has("x"))
		assert.True(t, q.processing.has("x"))
		assert.Equal(t, 0, len(q.queue))
		assert.Equal(t, 0, q.Len())

		q.Done("x")
		assert.False(t, q.dirty.has("x"))
		assert.False(t, q.processing.has("x"))
		assert.Equal(t, 0, len(q.queue))
		assert.Equal(t, 0, q.Len())

	})
}
