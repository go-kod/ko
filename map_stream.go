package ko

import "iter"

// Methods in this file return intermediate sequences.

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

// Keys returns the entry keys as a Seq. Map-backed sources use Go map order.
func (c Seq2[K, V]) Keys() Seq[K] {
	return c.ToSlice(func(key K, _ V) K {
		return key
	})
}

// Values returns the entry values as a Seq. Map-backed sources use Go map order.
func (c Seq2[K, V]) Values() Seq[V] {
	return c.ToSlice(func(_ K, value V) V {
		return value
	})
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
