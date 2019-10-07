package utils

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"reflect"
	"sort"
	"time"
)

var reflectInt = reflect.TypeOf(int(1))
var reflectStr = reflect.TypeOf("a")
var reflectFloat = reflect.TypeOf(float64(1))

type WeightedRandom interface {
	Get() interface{}
	GetWithSeed(seed string) interface{}
}

type Entry struct {
	Item    interface{}
	Weight  float32
}

type internalEntry struct {
	entry        Entry
	currentTotal float32
}

// WeightedRandom selects elements randomly from the list taking into account individual weights.
// Weight has to be assigned between 0 and 1.
// Support deterministic result given a particular seed and sortKey
type weightedRandomImpl struct {
	entries     []internalEntry
	totalWeight float32
}

func validateEntries(entries []Entry, sortKey string) error {
	if len(entries) == 0 {
		return fmt.Errorf("entries is empty")
	}
	for _, entry := range entries {
		if entry.Weight < 0 || entry.Weight > float32(1) {
			return fmt.Errorf("invalid weight %f", entry.Weight)
		}

		item := reflect.ValueOf(entry.Item)
		f := reflect.Indirect(item).FieldByName(sortKey)
		if !f.IsValid() {
			return fmt.Errorf("invalid sort key")
		}
		switch f.Type() {
		case reflectInt:
		case reflectStr:
		case reflectFloat:
			continue
		default:
			return fmt.Errorf("unsupported type")
		}
	}
	return nil
}

// Given a list of entries and sortKey, return WeightedRandom
// The sortKey indicates the field in the object to be used for sorting.
// This enables deterministic results for same seed and sortKey
func NewWeightedRandom(entries []Entry, sortKey string) (WeightedRandom, error) {
	err := validateEntries(entries, sortKey)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		item1 := reflect.ValueOf(entries[i].Item)
		field1 := reflect.Indirect(item1).FieldByName(sortKey)

		item2 := reflect.ValueOf(entries[j].Item)
		field2 := reflect.Indirect(item2).FieldByName(sortKey)
		switch field1.Type() {
		case reflectInt:
			return field1.Int() < field2.Int()
		case reflectStr:
			return field1.String() < field2.String()
		case reflectFloat:
			return field1.Float() < field2.Float()
		}
		// Should not reach here.
		return true
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
			currentTotal += 1.0/ float32(numberOfEntries)
		} else if e.Weight == 0 {
			// Entries which have zero weight are ignored
			continue
		}

		currentTotal += e.Weight
		internalEntries = append(internalEntries,internalEntry{
			entry:        e,
			currentTotal: currentTotal,
		})
	}

	return &weightedRandomImpl{
		entries:     internalEntries,
		totalWeight: currentTotal,
	}, nil
}

func (w weightedRandomImpl) get() interface{} {
	randomWeight := rand.Float32() * w.totalWeight
	for _, e := range w.entries {
		if e.currentTotal >= randomWeight && e.currentTotal > 0 {
			return e.entry.Item
		}
	}
	return w.entries[len(w.entries) -1].entry.Item
}

func (w weightedRandomImpl) Get() interface{} {
	rand.Seed(time.Now().UTC().UnixNano())
	return w.get()
}

func (w weightedRandomImpl) GetWithSeed(seed string) interface{} {
	h := fnv.New64a()
	h.Write([]byte(seed))
	hashedSeed := int64(h.Sum64())
	rand.Seed(hashedSeed)
	return w.get()
}
