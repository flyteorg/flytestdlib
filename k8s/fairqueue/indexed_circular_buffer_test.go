package fairqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexedCircularBuffer_RegularOperations(t *testing.T) {
	b := NewIndexedCircularBuffer()
	assert.Nil(t, b.Next())
	v, ok := b.Get("x")
	assert.False(t, ok)
	assert.Nil(t, v)
	assert.Equal(t, 0, b.Len())

	v, ok = b.GetOrDefault("x", func() interface{} {
		return "y"
	})
	assert.Equal(t, "y", v)
	assert.False(t, ok)

	v, ok = b.Get("x")
	assert.True(t, ok)
	assert.Equal(t, "y", v)

	assert.Equal(t, 1, b.Len())

	v, ok = b.GetOrDefault("m", func() interface{} {
		return "n"
	})
	assert.Equal(t, "n", v)
	assert.False(t, ok)

	v, ok = b.Get("m")
	assert.True(t, ok)
	assert.Equal(t, "n", v)

	assert.Equal(t, 2, b.Len())

	v, ok = b.GetOrDefault("m", func() interface{} {
		assert.FailNow(t, "should not be called")
		return nil
	})
	assert.Equal(t, "n", v)
	assert.True(t, ok)

	assert.Equal(t, 2, b.Len())
}

func TestIndexedCircularBuffer_TestCircularBuffer(t *testing.T) {
	b := NewIndexedCircularBuffer()
	assert.Nil(t, b.Next())

	_, ok := b.GetOrDefault("x", func() interface{} {
		return "x"
	})
	assert.False(t, ok)

	_, ok = b.GetOrDefault("y", func() interface{} {
		return "y"
	})
	assert.False(t, ok)

	_, ok = b.GetOrDefault("z", func() interface{} {
		return "z"
	})
	assert.False(t, ok)

	assert.Equal(t, "x", b.Next())
	assert.Equal(t, "y", b.Next())
	assert.Equal(t, "z", b.Next())
	_, ok = b.GetOrDefault("m", func() interface{} {
		return "m"
	})
	assert.False(t, ok)
	assert.Equal(t, "m", b.Next())
	_, ok = b.GetOrDefault("n", func() interface{} {
		return "n"
	})
	assert.False(t, ok)
	assert.Equal(t, "n", b.Next())
	assert.Equal(t, "x", b.Next())
	assert.Equal(t, "y", b.Next())
	assert.Equal(t, "z", b.Next())
}

func TestIndexedCircularBuffer_Range(t *testing.T) {
	b := NewIndexedCircularBuffer()
	assert.Nil(t, b.Next())
	b.Range(func(v interface{}) bool {
		assert.FailNow(t, "should not be invoked for an empty buffer")
		return false
	})

	assert.True(t, b.Add("x", "x"))
	assert.False(t, b.Add("x", "m"))
	assert.True(t, b.Add("y", "y"))
	assert.True(t, b.Add("z", "z"))

	t.Run("complete-iteration", func(t *testing.T) {
		arr := make([]string, 0, 3)
		b.Range(func(v interface{}) bool {
			s, _ := v.(string)
			arr = append(arr, s)
			return true
		})
		assert.Equal(t, []string{"x", "y", "z"}, arr)
	})

	t.Run("partial-iteration", func(t *testing.T) {
		arr := make([]string, 0, 1)
		b.Range(func(v interface{}) bool {
			s, _ := v.(string)
			if s == "y" {
				return false
			}
			arr = append(arr, s)
			return true
		})
		assert.Equal(t, []string{"x"}, arr)
	})
}

func TestIndexedCircularBuffer_RangeNext(t *testing.T) {
	b := NewIndexedCircularBuffer()
	assert.Nil(t, b.Next())
	b.Range(func(v interface{}) bool {
		assert.FailNow(t, "should not be invoked for an empty buffer")
		return false
	})

	assert.True(t, b.Add("x", "x"))
	assert.False(t, b.Add("x", "m"))
	assert.True(t, b.Add("y", "y"))
	assert.True(t, b.Add("z", "z"))

	t.Run("complete-iteration", func(t *testing.T) {
		arr := make([]string, 0, 3)
		b.RangeNext(func(v interface{}) bool {
			s, _ := v.(string)
			arr = append(arr, s)
			return true
		})
		assert.Equal(t, []string{"x", "y", "z"}, arr)
	})

	t.Run("partial-iteration", func(t *testing.T) {
		arr := make([]string, 0, 1)
		b.RangeNext(func(v interface{}) bool {
			s, _ := v.(string)
			if s == "y" {
				return false
			}
			arr = append(arr, s)
			return true
		})
		assert.Equal(t, []string{"x"}, arr)
	})
}

func TestIndexedCircularBuffer_IsCurrentAtHead(t *testing.T) {
	b := NewIndexedCircularBuffer()
	assert.Nil(t, b.Next())
	assert.True(t, b.IsCurrentAtHead())

	assert.True(t, b.Add("x", "x"))
	assert.Equal(t, "x", b.Next())
	assert.False(t, b.IsCurrentAtHead())

	assert.True(t, b.Add("y", "y"))
	assert.Equal(t, "y", b.Next())
	assert.False(t, b.IsCurrentAtHead())

	assert.True(t, b.Add("z", "z"))
	assert.Equal(t, "z", b.Next())
	assert.False(t, b.IsCurrentAtHead())

	{
		arr := make([]string, 0, 3)
		b.RangeNext(func(v interface{}) bool {
			s, _ := v.(string)
			arr = append(arr, s)
			return true
		})
		assert.Equal(t, []string{"x", "y", "z"}, arr)
	}

	{
		b.Next()
		arr := make([]string, 0, 3)
		b.RangeNext(func(v interface{}) bool {
			s, _ := v.(string)
			arr = append(arr, s)
			return true
		})
		assert.Equal(t, []string{"y", "z", "x"}, arr)
	}

	{
		b.Next()
		arr := make([]string, 0, 3)
		b.RangeNext(func(v interface{}) bool {
			s, _ := v.(string)
			arr = append(arr, s)
			return true
		})
		assert.Equal(t, []string{"z", "x", "y"}, arr)
	}
}
