package ko

import (
	"iter"
	"maps"
)

// Map starts a chain over collection.
func Map[K comparable, V any](collection map[K]V) Seq2[K, V] {
	return Seq2[K, V](maps.All(collection))
}

func collectMap[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	return maps.Collect(seq)
}

// Collect returns the current map.
func (c Seq2[K, V]) Collect() map[K]V {
	return collectMap(iter.Seq2[K, V](c))
}

// HasKey reports whether key exists in the map.
func (c Seq2[K, V]) HasKey(key K) bool {
	for itemKey := range iter.Seq2[K, V](c) {
		if itemKey == key {
			return true
		}
	}
	return false
}

// ValueOr returns the value for key, or fallback when key is absent.
func (c Seq2[K, V]) ValueOr(key K, fallback V) V {
	for itemKey, value := range iter.Seq2[K, V](c) {
		if itemKey == key {
			return value
		}
	}
	return fallback
}

// PickKeys keeps entries whose key is listed.
func (c Seq2[K, V]) PickKeys(keys ...K) Seq2[K, V] {
	keep := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		keep[key] = struct{}{}
	}

	return Seq2[K, V](func(yield func(K, V) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if _, ok := keep[key]; ok && !yield(key, value) {
				return
			}
		}
	})
}

// OmitKeys drops entries whose key is listed.
func (c Seq2[K, V]) OmitKeys(keys ...K) Seq2[K, V] {
	drop := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		drop[key] = struct{}{}
	}

	return Seq2[K, V](func(yield func(K, V) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if _, ok := drop[key]; !ok && !yield(key, value) {
				return
			}
		}
	})
}

// Assign merges maps into a new Seq2. Later maps replace earlier keys.
func (c Seq2[K, V]) Assign(maps ...map[K]V) Seq2[K, V] {
	return Seq2[K, V](func(yield func(K, V) bool) {
		result := collectMap(iter.Seq2[K, V](c))
		for _, items := range maps {
			for key, value := range items {
				result[key] = value
			}
		}
		for key, value := range result {
			if !yield(key, value) {
				return
			}
		}
	})
}

// Filter keeps entries for which predicate returns true.
func (c Seq2[K, V]) Filter(predicate func(key K, value V) bool) Seq2[K, V] {
	return Seq2[K, V](func(yield func(K, V) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if predicate(key, value) && !yield(key, value) {
				return
			}
		}
	})
}

// Reject drops entries for which predicate returns true.
func (c Seq2[K, V]) Reject(predicate func(key K, value V) bool) Seq2[K, V] {
	return c.Filter(func(key K, value V) bool {
		return !predicate(key, value)
	})
}

// Map transforms entries and may change both key and value types.
func (c Seq2[K, V]) Map[RK comparable, RV any](mapper func(key K, value V) (RK, RV)) Seq2[RK, RV] {
	return Seq2[RK, RV](func(yield func(RK, RV) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			nextKey, nextValue := mapper(key, value)
			if !yield(nextKey, nextValue) {
				return
			}
		}
	})
}

// MapKeys transforms keys and keeps values.
func (c Seq2[K, V]) MapKeys[RK comparable](mapper func(key K, value V) RK) Seq2[RK, V] {
	return Seq2[RK, V](func(yield func(RK, V) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if !yield(mapper(key, value), value) {
				return
			}
		}
	})
}

// MapValues transforms values and keeps keys.
func (c Seq2[K, V]) MapValues[RV any](mapper func(key K, value V) RV) Seq2[K, RV] {
	return Seq2[K, RV](func(yield func(K, RV) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if !yield(key, mapper(key, value)) {
				return
			}
		}
	})
}

// ForEach calls iteratee for each entry and returns the unchanged chain.
func (c Seq2[K, V]) ForEach(iteratee func(key K, value V)) Seq2[K, V] {
	for key, value := range iter.Seq2[K, V](c) {
		iteratee(key, value)
	}
	return c
}

// Keys returns the map keys as a slice chain. Map iteration order is Go map order.
func (c Seq2[K, V]) Keys() Seq[K] {
	return Seq[K](func(yield func(K) bool) {
		for key := range iter.Seq2[K, V](c) {
			if !yield(key) {
				return
			}
		}
	})
}

// ChunkEntries splits map entries into Seq2 chunks of size n. Map iteration order is Go map order.
func (c Seq2[K, V]) ChunkEntries(n int) iter.Seq[Seq2[K, V]] {
	return func(yield func(Seq2[K, V]) bool) {
		if n <= 0 {
			return
		}
		chunk := make(map[K]V, n)
		count := 0
		for key, value := range iter.Seq2[K, V](c) {
			chunk[key] = value
			count++
			if count < n {
				continue
			}
			if !yield(Seq2[K, V](maps.All(chunk))) {
				return
			}
			chunk = make(map[K]V, n)
			count = 0
		}
		if count > 0 {
			yield(Seq2[K, V](maps.All(chunk)))
		}
	}
}

// FilterKeys returns matching keys as a slice chain. Map iteration order is Go map order.
func (c Seq2[K, V]) FilterKeys(predicate func(key K, value V) bool) Seq[K] {
	return Seq[K](func(yield func(K) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if predicate(key, value) && !yield(key) {
				return
			}
		}
	})
}

// Values returns the map values as a slice chain. Map iteration order is Go map order.
func (c Seq2[K, V]) Values() Seq[V] {
	return Seq[V](func(yield func(V) bool) {
		for _, value := range iter.Seq2[K, V](c) {
			if !yield(value) {
				return
			}
		}
	})
}

// FilterValues returns matching values as a slice chain. Map iteration order is Go map order.
func (c Seq2[K, V]) FilterValues(predicate func(key K, value V) bool) Seq[V] {
	return Seq[V](func(yield func(V) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if predicate(key, value) && !yield(value) {
				return
			}
		}
	})
}

// Some reports whether any entry matches predicate.
func (c Seq2[K, V]) Some(predicate func(key K, value V) bool) bool {
	_, _, ok := c.Find(predicate)
	return ok
}

// None reports whether no entry matches predicate.
func (c Seq2[K, V]) None(predicate func(key K, value V) bool) bool {
	_, _, ok := c.Find(predicate)
	return !ok
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
	for key, value := range iter.Seq2[K, V](c) {
		if !predicate(key, value) {
			return false
		}
	}
	return true
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

// ToSlice transforms the map into a slice chain. Map iteration order is Go map order.
func (c Seq2[K, V]) ToSlice[R any](mapper func(key K, value V) R) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if !yield(mapper(key, value)) {
				return
			}
		}
	})
}

// FilterMapToSlice transforms matching entries into a slice chain. Map iteration order is Go map order.
func (c Seq2[K, V]) FilterMapToSlice[R any](mapper func(key K, value V) (R, bool)) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if item, ok := mapper(key, value); ok && !yield(item) {
				return
			}
		}
	})
}
