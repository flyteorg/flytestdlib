package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	key string
	val int
}

func TestDeterministicWeightedRandomStr(t *testing.T) {
	item1 := testData{
		key: "key1",
		val: 1,
	}
	item2 := testData{
		key: "key2",
		val: 2,
	}
	entries := []Entry{
		{
			Item:   item1,
			Weight: 0.4,
		},
		{
			Item:   item2,
			Weight: 0.6,
		},
	}
	rand, err := NewWeightedRandom(entries, "key")
	assert.Nil(t, err)
	retItem, err := rand.GetWithSeed("ab")
	assert.Nil(t, err)
	assert.Equal(t, item1, retItem)

	assert.Nil(t, err)
	for i := 1; i <= 10; i++ {
		retItem, err := rand.GetWithSeed("hi")
		assert.Nil(t, err)
		assert.Equal(t, item2, retItem)
	}
}

func TestDeterministicWeightedRandomInt(t *testing.T) {
	item1 := testData{
		key: "key1",
		val: 4,
	}
	item2 := testData{
		key: "key2",
		val: 3,
	}
	entries := []Entry{
		{
			Item:   item1,
			Weight: 0.4,
		},
		{
			Item:   item2,
			Weight: 0.6,
		},
	}
	rand, err := NewWeightedRandom(entries, "val")
	assert.Nil(t, err)
	retItem, err := rand.GetWithSeed("ab")
	assert.Nil(t, err)
	assert.Equal(t, item2, retItem)

	for i := 1; i <= 10; i++ {
		retItem, err := rand.GetWithSeed("hi")
		assert.Nil(t, err)
		assert.Equal(t, item1, retItem)
	}
}

func TestDeterministicWeightedFewZeroWeight(t *testing.T) {
	item1 := testData{
		key: "key1",
		val: 4,
	}
	item2 := testData{
		key: "key2",
		val: 3,
	}
	entries := []Entry{
		{
			Item:   item1,
			Weight: 0.4,
		},
		{
			Item: item2,
		},
	}
	rand, err := NewWeightedRandom(entries, "val")
	assert.Nil(t, err)
	retItem, err := rand.GetWithSeed("ab")
	assert.Nil(t, err)
	assert.Equal(t, item1, retItem)

	for i := 1; i <= 10; i++ {
		retItem, err := rand.GetWithSeed("hi")
		assert.Nil(t, err)
		assert.Equal(t, item1, retItem)
	}
}

func TestDeterministicWeightedAllZeroWeights(t *testing.T) {
	item1 := testData{
		key: "key1",
		val: 4,
	}
	item2 := testData{
		key: "key2",
		val: 3,
	}
	entries := []Entry{
		{
			Item: item1,
		},
		{
			Item: item2,
		},
	}
	rand, err := NewWeightedRandom(entries, "key")
	assert.Nil(t, err)
	retItem, err := rand.GetWithSeed("hi")
	assert.Nil(t, err)
	assert.Equal(t, item2, retItem)

	for i := 1; i <= 10; i++ {
		retItem, err := rand.GetWithSeed("ab")
		assert.Nil(t, err)
		assert.Equal(t, item1, retItem)
	}
}

func TestDeterministicWeightInvalidSortKey(t *testing.T) {
	item1 := testData{
		key: "key1",
		val: 4,
	}
	item2 := testData{
		key: "key2",
		val: 3,
	}
	entries := []Entry{
		{
			Item: item1,
		},
		{
			Item: item2,
		},
	}
	_, err := NewWeightedRandom(entries, "key1")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid sort key")
}

func TestDeterministicWeightInvalidWeights(t *testing.T) {
	item1 := testData{
		key: "key1",
		val: 4,
	}
	item2 := testData{
		key: "key2",
		val: 3,
	}
	entries := []Entry{
		{
			Item:   item1,
			Weight: -3.0,
		},
		{
			Item: item2,
		},
	}
	_, err := NewWeightedRandom(entries, "key")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid weight -3.000000")
}
