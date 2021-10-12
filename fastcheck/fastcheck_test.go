package fastcheck

import (
	"context"
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

func init() {
	labeled.SetMetricKeys(contextutils.MetricKeysFromStrings([]string{"test"})...)
}
