package random

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"sort"
	"time"
)

//go:generate mockery -all -case=underscore

// Interface to use the Weighted Random
type WeightedRandom interface {
	Get() interface{}
	GetWithSeed(seed string) (interface{}, error)
}

// Interface for items that can be used along with Weighted Random
type Comparable interface {
	Compare(to Comparable) bool
}

// Structure of each entry to select from
type Entry struct {
	Item   Comparable
	Weight float32
}

type internalEntry struct {
	entry        Entry
	currentTotal float32
}

// WeightedRandom selects elements randomly from the list taking into account individual weights.
// Weight has to be assigned between 0 and 1.
// Support deterministic results given a particular seed and sortKey
type weightedRandomImpl struct {
	entries     []internalEntry
	totalWeight float32
}

func validateEntries(entries []Entry) error {
	if len(entries) == 0 {
		return fmt.Errorf("entries is empty")
	}
	for _, entry := range entries {
		if entry.Weight < 0 || entry.Weight > float32(1) {
			return fmt.Errorf("invalid weight %f", entry.Weight)
		}
	}
	return nil
}

// Given a list of entries and sortKey, return WeightedRandom
// The sortKey indicates the field in the object to be used for sorting.
// This enables deterministic results for same seed and sortKey
func NewWeightedRandom(entries []Entry) (WeightedRandom, error) {
	err := validateEntries(entries)
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Item.Compare(entries[j].Item)
	})
	var internalEntries []internalEntry
	numberOfEntries := len(entries)
	totalWeight := float32(0)
	for _, e := range entries {
		totalWeight += e.Weight
	}

	currentTotal := float32(0)
	for _, e := range entries {
		if totalWeight == 0 {
			// This indicates that none of the entries have weight assigned.
			// We will assign equal weights to everyone
			currentTotal += 1.0 / float32(numberOfEntries)
		} else if e.Weight == 0 {
			// Entries which have zero weight are ignored
			continue
		}

		currentTotal += e.Weight
		internalEntries = append(internalEntries, internalEntry{
			entry:        e,
			currentTotal: currentTotal,
		})
	}

	return &weightedRandomImpl{
		entries:     internalEntries,
		totalWeight: currentTotal,
	}, nil
}

func (w *weightedRandomImpl) get(generator *rand.Rand) interface{} {
	randomWeight := generator.Float32() * w.totalWeight
	for _, e := range w.entries {
		if e.currentTotal >= randomWeight && e.currentTotal > 0 {
			return e.entry.Item
		}
	}
	return w.entries[len(w.entries)-1].entry.Item
}

// Returns a random entry based on the weights
func (w *weightedRandomImpl) Get() interface{} {
	randGenerator := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return w.get(randGenerator)
}

// For a given seed, the same entry will be returned all the time.
func (w *weightedRandomImpl) GetWithSeed(seed string) (interface{}, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(seed))
	if err != nil {
		return nil, err
	}
	hashedSeed := int64(h.Sum64())
	randGenerator := rand.New(rand.NewSource(hashedSeed))
	return w.get(randGenerator), nil
}
