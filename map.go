package ko

import (
	"iter"
	"maps"
)

// Map returns a Seq2 over collection.
func Map[K comparable, V any](collection map[K]V) Seq2[K, V] {
	return Seq2[K, V](maps.All(collection))
}

// Collect materializes the entries into a map.
func (c Seq2[K, V]) Collect() map[K]V {
	return maps.Collect(iter.Seq2[K, V](c))
}

// HasKey reports whether key exists in the map.
func (c Seq2[K, V]) HasKey(key K) bool {
	_, _, ok := c.Find(func(itemKey K, _ V) bool {
		return itemKey == key
	})
	return ok
}

// ValueOr returns the value for key, or fallback when key is absent.
func (c Seq2[K, V]) ValueOr(key K, fallback V) V {
	_, value, ok := c.Find(func(itemKey K, _ V) bool {
		return itemKey == key
	})
	if !ok {
		return fallback
	}
	return value
}

// PickKeys keeps entries whose key is listed.
func (c Seq2[K, V]) PickKeys(keys ...K) Seq2[K, V] {
	if len(keys) == 0 {
		return Seq2[K, V](func(yield func(K, V) bool) {})
	}
	keep := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		keep[key] = struct{}{}
	}

	return c.Filter(func(key K, _ V) bool {
		_, ok := keep[key]
		return ok
	})
}

// OmitKeys drops entries whose key is listed.
func (c Seq2[K, V]) OmitKeys(keys ...K) Seq2[K, V] {
	if len(keys) == 0 {
		return c
	}
	drop := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		drop[key] = struct{}{}
	}

	return c.Filter(func(key K, _ V) bool {
		_, ok := drop[key]
		return !ok
	})
}

// Assign merges maps into a new Seq2. Later maps replace earlier keys.
func (c Seq2[K, V]) Assign(others ...map[K]V) Seq2[K, V] {
	return Seq2[K, V](func(yield func(K, V) bool) {
		overrides := make(map[K]int)
		for _, items := range others {
			for key := range items {
				overrides[key]++
			}
		}

		for key, value := range iter.Seq2[K, V](c) {
			if overrides[key] > 0 {
				continue
			}
			if !yield(key, value) {
				return
			}
		}

		for _, items := range others {
			for key, value := range items {
				overrides[key]--
				if overrides[key] > 0 {
					continue
				}
				if !yield(key, value) {
					return
				}
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
	return c.Map(func(key K, value V) (RK, V) {
		return mapper(key, value), value
	})
}

// MapValues transforms values and keeps keys.
func (c Seq2[K, V]) MapValues[RV any](mapper func(key K, value V) RV) Seq2[K, RV] {
	return c.Map(func(key K, value V) (K, RV) {
		return key, mapper(key, value)
	})
}

// ForEach calls iteratee for each entry and returns the unchanged Seq2.
func (c Seq2[K, V]) ForEach(iteratee func(key K, value V)) Seq2[K, V] {
	return Seq2[K, V](func(yield func(K, V) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			iteratee(key, value)
			if !yield(key, value) {
				return
			}
		}
	})
}

// Keys returns the entry keys as a Seq. Map-backed sources use Go map order.
func (c Seq2[K, V]) Keys() Seq[K] {
	return c.ToSlice(func(key K, _ V) K {
		return key
	})
}

// ChunkEntries splits entries into Seq2 chunks of size n. Map-backed sources use Go map order.
func (c Seq2[K, V]) ChunkEntries(n int) iter.Seq[Seq2[K, V]] {
	return func(yield func(Seq2[K, V]) bool) {
		if n <= 0 {
			return
		}
		chunk := make([]struct {
			key   K
			value V
		}, 0, n)
		emit := func() bool {
			entries := append([]struct {
				key   K
				value V
			}(nil), chunk...)
			return yield(Seq2[K, V](func(yield func(K, V) bool) {
				for _, entry := range entries {
					if !yield(entry.key, entry.value) {
						return
					}
				}
			}))
		}

		for key, value := range iter.Seq2[K, V](c) {
			chunk = append(chunk, struct {
				key   K
				value V
			}{key: key, value: value})
			if len(chunk) < n {
				continue
			}
			if !emit() {
				return
			}
			chunk = chunk[:0]
		}
		if len(chunk) > 0 {
			emit()
		}
	}
}

// Values returns the entry values as a Seq. Map-backed sources use Go map order.
func (c Seq2[K, V]) Values() Seq[V] {
	return c.ToSlice(func(_ K, value V) V {
		return value
	})
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

// ToSlice transforms entries into a Seq. Map-backed sources use Go map order.
func (c Seq2[K, V]) ToSlice[R any](mapper func(key K, value V) R) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if !yield(mapper(key, value)) {
				return
			}
		}
	})
}

// FilterMapToSlice transforms matching entries into a Seq. Map-backed sources use Go map order.
func (c Seq2[K, V]) FilterMapToSlice[R any](mapper func(key K, value V) (R, bool)) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		for key, value := range iter.Seq2[K, V](c) {
			if item, ok := mapper(key, value); ok && !yield(item) {
				return
			}
		}
	})
}
