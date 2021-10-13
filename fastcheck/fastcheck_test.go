package fastcheck

import (
	"context"
	"math/rand"
	"sync"
	"testing"

	"github.com/flyteorg/flytestdlib/contextutils"
	"github.com/flyteorg/flytestdlib/promutils"
	"github.com/flyteorg/flytestdlib/promutils/labeled"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	ctx := context.TODO()

	lru, err := NewLRUCacheFilter(2, promutils.NewTestScope())
	assert.NoError(t, err)
	assert.NotNil(t, lru)
	oppo, err := NewOppoBloomFilter(2, promutils.NewTestScope())
	assert.NoError(t, err)
	assert.NotNil(t, oppo)

	twentyNineID := []byte{27, 28, 29}
	thirtyID := []byte{27, 28, 30}
	thirtyThreeID := []byte{27, 28, 33}
	// the last 2 byte arrays have the same hash value

	assert.False(t, lru.Contains(ctx, twentyNineID))
	assert.False(t, oppo.Contains(ctx, twentyNineID))
	lru.Add(ctx, twentyNineID)
	oppo.Add(ctx, twentyNineID)
	assert.True(t, lru.Contains(ctx, twentyNineID))
	assert.True(t, oppo.Contains(ctx, twentyNineID))

	assert.False(t, lru.Contains(ctx, thirtyID))
	assert.False(t, oppo.Contains(ctx, thirtyID))

	assert.False(t, lru.Contains(ctx, thirtyThreeID))
	assert.False(t, oppo.Contains(ctx, thirtyThreeID))

	// Now that they have the same hash value
	lru.Add(ctx, thirtyID)
	oppo.Add(ctx, thirtyID)
	assert.True(t, lru.Contains(ctx, thirtyID))
	assert.True(t, oppo.Contains(ctx, thirtyID))
	// LRU should not contain it and oppo should also return false
	assert.False(t, lru.Contains(ctx, thirtyThreeID))
	assert.False(t, oppo.Contains(ctx, thirtyThreeID))

	lru.Add(ctx, thirtyThreeID)
	oppo.Add(ctx, thirtyThreeID)

	// LRU will evict first entered, while oppo will evict matching hash
	assert.True(t, lru.Contains(ctx, thirtyID))
	assert.False(t, oppo.Contains(ctx, thirtyID))

}

func TestSizeRounding(t *testing.T) {
	f, _ := NewOppoBloomFilter(3, promutils.NewTestScope())
	if len(f.array) != 4 {
		t.Errorf("3 should round to 4")
	}
	f, _ = NewOppoBloomFilter(4, promutils.NewTestScope())
	if len(f.array) != 4 {
		t.Errorf("4 should round to 4")
	}
	f, _ = NewOppoBloomFilter(129, promutils.NewTestScope())
	if len(f.array) != 256 {
		t.Errorf("129 should round to 256")
	}
}

func TestTooLargeSize(t *testing.T) {
	size := (1 << 30) + 1
	f, err := NewOppoBloomFilter(size, promutils.NewTestScope())
	if err != ErrSizeTooLarge {
		t.Errorf("did not error out on a too-large filter size")
	}
	if f != nil {
		t.Errorf("did not return nil on a too-large filter size")
	}
}

func TestTooSmallSize(t *testing.T) {
	f, err := NewOppoBloomFilter(0, promutils.NewTestScope())
	if err != ErrSizeTooSmall {
		t.Errorf("did not error out on a too small filter size")
	}
	if f != nil {
		t.Errorf("did not return nil on a too small filter size")
	}
}

func benchmarkFilter(b *testing.B, ids [][]byte, filter Filter, nGoroutines int) {
	ctx := context.TODO()
	wg := sync.WaitGroup{}
	wg.Add(nGoroutines)
	f := func() {
		for _, id := range ids {
			filter.Contains(ctx, id)
			filter.Add(ctx, id)
		}
		wg.Done()
	}

	for i := 0; i < nGoroutines; i++ {
		go f()
	}
	wg.Wait()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandBytesGenerator(n int) []byte {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

// RandBytesSliceGenerator generates a slice of size m with random bytes of size n
func RandBytesSliceGenerator(n int, m int) [][]byte {
	byteSlice := make([][]byte, 0, m)
	for i := 0; i < m; i++ {
		byteSlice = append(byteSlice, RandBytesGenerator(n))
	}
	return byteSlice
}

func BenchmarkFilter(b *testing.B) {
	bf, err := NewOppoBloomFilter(5000, promutils.NewTestScope())
	assert.NoError(b, err)
	ids := RandBytesSliceGenerator(100, 10000)
	nGoroutines := 20

	b.Run("oppoFilterBenchmark", func(b *testing.B) {
		benchmarkFilter(b, ids, bf, nGoroutines)
	})

	b.Run("lruFilterBenchmark", func(b *testing.B) {
		benchmarkFilter(b, ids, bf, nGoroutines)
	})
}

func init() {
	labeled.SetMetricKeys(contextutils.MetricKeysFromStrings([]string{"test"})...)
}
