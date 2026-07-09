package ko

import (
	"iter"
	"maps"
)

// Methods in this file consume Seq2 values and return end-state values.

// Collect materializes the entries into a map.
func (c Seq2[K, V]) Collect() map[K]V {
	return maps.Collect(iter.Seq2[K, V](c))
}

// IsEmpty reports whether the sequence yields no entries.
func (c Seq2[K, V]) IsEmpty() bool {
	return !c.Some(func(_ K, _ V) bool {
		return true
	})
}

// Some reports whether any entry matches predicate.
func (c Seq2[K, V]) Some(predicate func(key K, value V) bool) bool {
	_, _, ok := c.Find(predicate)
	return ok
}

// Count returns the number of entries matching predicate.
func (c Seq2[K, V]) Count(predicate func(key K, value V) bool) int {
	result := 0
	for key, value := range iter.Seq2[K, V](c) {
		if predicate(key, value) {
			result++
		}
	}
	return result
}

// Every reports whether all entries match predicate.
func (c Seq2[K, V]) Every(predicate func(key K, value V) bool) bool {
	return !c.Some(func(key K, value V) bool {
		return !predicate(key, value)
	})
}

// Find returns the first matching entry. Map iteration order is Go map order.
func (c Seq2[K, V]) Find(predicate func(key K, value V) bool) (K, V, bool) {
	for key, value := range iter.Seq2[K, V](c) {
		if predicate(key, value) {
			return key, value, true
		}
	}
	var zeroKey K
	var zeroValue V
	return zeroKey, zeroValue, false
}
