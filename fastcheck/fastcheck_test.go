package fastcheck

import (
	"context"
	"github.com/flyteorg/flytestdlib/contextutils"
	"github.com/flyteorg/flytestdlib/promutils"
	"github.com/flyteorg/flytestdlib/promutils/labeled"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T)  {
	ctx := context.TODO()

	lru, err := NewLRUCacheFilter(2, promutils.NewTestScope())
	assert.NoError(t, err)
	assert.NotNil(t, lru)
	oppo, err := NewOppoBloomFilter(2, promutils.NewTestScope())
	assert.NoError(t, err)
	assert.NotNil(t, oppo)

	twentyNineId := []byte{27, 28, 29}
	thirtyId := []byte{27, 28, 30}
	thirtyThreeId := []byte{27, 28, 33}
	// the last 2 byte arrays have the same hash value

	assert.False(t, lru.Contains(ctx, twentyNineId))
	assert.False(t, oppo.Contains(ctx, twentyNineId))
	lru.Add(ctx, twentyNineId)
	oppo.Add(ctx, twentyNineId)
	assert.True(t, lru.Contains(ctx, twentyNineId))
	assert.True(t, oppo.Contains(ctx, twentyNineId))

	assert.False(t, lru.Contains(ctx, thirtyId))
	assert.False(t, oppo.Contains(ctx, thirtyId))

	assert.False(t, lru.Contains(ctx, thirtyThreeId))
	assert.False(t, oppo.Contains(ctx, thirtyThreeId))

	// Now that they have the same hash value
	lru.Add(ctx, thirtyId)
	oppo.Add(ctx, thirtyId)
	assert.True(t, lru.Contains(ctx, thirtyId))
	assert.True(t, oppo.Contains(ctx, thirtyId))
	// LRU should not contain it and oppo should also return false
	assert.False(t, lru.Contains(ctx, thirtyThreeId))
	assert.False(t, oppo.Contains(ctx, thirtyThreeId))

	lru.Add(ctx, thirtyThreeId)
	oppo.Add(ctx, thirtyThreeId)

	// LRU will evict first entered, while oppo will evict matching hash
	assert.True(t, lru.Contains(ctx, thirtyId))
	assert.False(t, oppo.Contains(ctx, thirtyId))

}

func TestSizeRounding(t *testing.T) {
	f, _ := NewOppoBloomFilter(3, promutils.NewTestScope())
	if len(f.(*oppoBloomFilter).array) != 4 {
		t.Errorf("3 should round to 4")
	}
	f, _ = NewOppoBloomFilter(4, promutils.NewTestScope())
	if len(f.(*oppoBloomFilter).array) != 4 {
		t.Errorf("4 should round to 4")
	}
	f, _ = NewOppoBloomFilter(129, promutils.NewTestScope())
	if len(f.(*oppoBloomFilter).array) != 256 {
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

func init()  {
	labeled.SetMetricKeys(contextutils.MetricKeysFromStrings([]string{"test"})...)
}