package ko

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

// Methods in this file return intermediate sequences.

// Filter keeps items for which predicate returns true.
func (c Seq[T]) Filter(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			if predicate(item, i) && !yield(item) {
				return
			}
			i++
		}
	})
}

// Reject drops items for which predicate returns true.
func (c Seq[T]) Reject(predicate func(item T, index int) bool) Seq[T] {
	return c.Filter(func(item T, index int) bool {
		return !predicate(item, index)
	})
}

// DistinctBy removes duplicate items by key, keeping the first occurrence.
func (c Seq[T]) DistinctBy[K comparable](mapper func(item T, index int) K) Seq[T] {
	return Seq[T](distinctBySeq(iter.Seq[T](c), mapper))
}

// Map transforms items and may change the element type.
func (c Seq[T]) Map[R any](mapper func(item T, index int) R) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			if !yield(mapper(item, i)) {
				return
			}
			i++
		}
	})
}

// FilterMap maps items and keeps only accepted results.
func (c Seq[T]) FilterMap[R any](mapper func(item T, index int) (R, bool)) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			if mapped, ok := mapper(item, i); ok && !yield(mapped) {
				return
			}
			i++
		}
	})
}

// FlatMap transforms each item into a sequence and flattens one level.
func (c Seq[T]) FlatMap[R any](mapper func(item T, index int) iter.Seq[R]) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			for mapped := range iter.Seq[R](mapper(item, i)) {
				if !yield(mapped) {
					return
				}
			}
			i++
		}
	})
}

// Enumerate indexes items by their zero-based position.
func (c Seq[T]) Enumerate() Seq2[int, T] {
	return Seq2[int, T](func(yield func(int, T) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			if !yield(i, item) {
				return
			}
			i++
		}
	})
}

// GroupBy groups items by key, preserving first key order.
func (c Seq[T]) GroupBy[K comparable](mapper func(item T, index int) K) groupedSeq[K, T] {
	return groupedSeq[K, T](func(yield func(K, []T) bool) {
		result := make(map[K][]T)
		order := make([]K, 0)
		i := 0
		for item := range iter.Seq[T](c) {
			key := mapper(item, i)
			if _, ok := result[key]; !ok {
				order = append(order, key)
			}
			result[key] = append(result[key], item)
			i++
		}
		for _, key := range order {
			if !yield(key, slices.Clone(result[key])) {
				return
			}
		}
	})
}

// Concat appends another sequence after the current sequence.
func (c Seq[T]) Concat(other iter.Seq[T]) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		for item := range iter.Seq[T](c) {
			if !yield(item) {
				return
			}
		}
		for item := range iter.Seq[T](other) {
			if !yield(item) {
				return
			}
		}
	})
}

// Chunk splits items into non-overlapping chunks of size n.
func (c Seq[T]) Chunk(n int) seqSeq[T] {
	return seqSeq[T](func(yield func([]T) bool) {
		for chunk := range chunkSeq(iter.Seq[T](c), n) {
			if !yield(chunk) {
				return
			}
		}
	})
}

// Window splits items into overlapping windows of size n, advancing by step.
func (c Seq[T]) Window(n, step int) seqSeq[T] {
	return seqSeq[T](func(yield func([]T) bool) {
		for window := range windowSeq(iter.Seq[T](c), n, step) {
			if !yield(window) {
				return
			}
		}
	})
}

// Take keeps the first n items.
func (c Seq[T]) Take(n int) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		i := 0
		for item := range iter.Seq[T](c) {
			if !yield(item) {
				return
			}
			i++
			if i >= n {
				return
			}
		}
	})
}

// TakeRight keeps the last n items.
func (c Seq[T]) TakeRight(n int) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		items := c.Collect()
		if n < len(items) {
			items = items[len(items)-n:]
		}
		for item := range Slice(items) {
			if !yield(item) {
				return
			}
		}
	})
}

// TakeWhile keeps items from the beginning while predicate returns true.
func (c Seq[T]) TakeWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			if !predicate(item, i) || !yield(item) {
				return
			}
			i++
		}
	})
}

// Drop skips the first n items.
func (c Seq[T]) Drop(n int) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			if i < n {
				i++
				continue
			}
			if !yield(item) {
				return
			}
		}
	})
}

// DropRight skips the last n items.
func (c Seq[T]) DropRight(n int) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		if n <= 0 {
			for item := range iter.Seq[T](c) {
				if !yield(item) {
					return
				}
			}
			return
		}
		buf := make([]T, 0, n+1)
		for item := range iter.Seq[T](c) {
			buf = append(buf, item)
			if len(buf) <= n {
				continue
			}
			if !yield(buf[0]) {
				return
			}
			buf = slices.Delete(buf, 0, 1)
		}
	})
}

// DropWhile drops items from the beginning while predicate returns true.
func (c Seq[T]) DropWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		i := 0
		dropping := true
		for item := range iter.Seq[T](c) {
			if dropping && predicate(item, i) {
				i++
				continue
			}
			dropping = false
			if !yield(item) {
				return
			}
			i++
		}
	})
}

// Scan returns the running accumulation, including the initial value.
func (c Seq[T]) Scan[R any](accumulator func(agg R, item T, index int) R, initial R) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		result := initial
		if !yield(result) {
			return
		}
		i := 0
		for item := range iter.Seq[T](c) {
			result = accumulator(result, item, i)
			i++
			if !yield(result) {
				return
			}
		}
	})
}

// ToMap transforms the sequence into Seq2.
func (c Seq[T]) ToMap[K comparable, V any](mapper func(item T, index int) (K, V)) Seq2[K, V] {
	return Seq2[K, V](func(yield func(K, V) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			key, value := mapper(item, i)
			i++
			if !yield(key, value) {
				return
			}
		}
	})
}

// Sort returns a stably sorted copy using less.
func (c Seq[T]) Sort(less func(left, right T) bool) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		items := slices.SortedStableFunc(iter.Seq[T](c), func(left, right T) int {
			if less(left, right) {
				return -1
			}
			if less(right, left) {
				return 1
			}
			return 0
		})
		for item := range Slice(items) {
			if !yield(item) {
				return
			}
		}
	})
}

// SortBy returns a stably sorted copy using mapped ordered keys.
func (c Seq[T]) SortBy[K cmp.Ordered](mapper func(item T, index int) K) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		type entry struct {
			item T
			key  K
		}
		items := c.Collect()
		entries := make([]entry, 0, len(items))
		for i, item := range items {
			entries = append(entries, entry{item: item, key: mapper(item, i)})
		}
		slices.SortStableFunc(entries, func(left, right entry) int {
			return cmp.Compare(left.key, right.key)
		})
		for _, entry := range entries {
			if !yield(entry.item) {
				return
			}
		}
	})
}

// Reverse returns a reversed copy.
func (c Seq[T]) Reverse() Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		for _, item := range slices.Backward(c.Collect()) {
			if !yield(item) {
				return
			}
		}
	})
}

// Map transforms each inner slice and returns a normal Seq.
func (s seqSeq[T]) Map[R any](mapper func(item []T, index int) R) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		i := 0
		for item := range iter.Seq[[]T](s) {
			if !yield(mapper(item, i)) {
				return
			}
			i++
		}
	})
}

// Collect materializes the outer and inner sequences.
func (s seqSeq[T]) Collect() [][]T {
	return slices.Collect(iter.Seq[[]T](s))
}

// Map transforms each grouped entry and returns a normal Seq.
func (g groupedSeq[K, V]) Map[R any](mapper func(key K, group []V) R) Seq[R] {
	return Seq[R](func(yield func(R) bool) {
		for key, group := range iter.Seq2[K, []V](g) {
			if !yield(mapper(key, group)) {
				return
			}
		}
	})
}

// Collect materializes groups into slices.
func (g groupedSeq[K, V]) Collect() map[K][]V {
	return maps.Collect(iter.Seq2[K, []V](g))
}

func distinctBySeq[T any, K comparable](seq iter.Seq[T], mapper func(item T, index int) K) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := map[K]struct{}{}
		i := 0
		for item := range seq {
			key := mapper(item, i)
			i++
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			if !yield(item) {
				return
			}
		}
	}
}

func chunkSeq[T any](seq iter.Seq[T], n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if n <= 0 {
			return
		}
		chunk := make([]T, 0, n)
		for item := range seq {
			chunk = append(chunk, item)
			if len(chunk) < n {
				continue
			}
			if !yield(slices.Clone(chunk)) {
				return
			}
			chunk = chunk[:0]
		}
		if len(chunk) > 0 {
			yield(slices.Clone(chunk))
		}
	}
}

func windowSeq[T any](seq iter.Seq[T], n, step int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if n <= 0 || step <= 0 {
			return
		}
		window := make([]T, 0, n)
		skip := 0
		for item := range seq {
			if skip > 0 {
				skip--
				continue
			}
			window = append(window, item)
			if len(window) < n {
				continue
			}
			if !yield(slices.Clone(window)) {
				return
			}
			if step >= n {
				skip = step - n
				window = window[:0]
				continue
			}
			// ponytail: copy-shift is enough; switch to a ring buffer if large windows matter.
			copy(window, window[step:])
			window = window[:n-step]
		}
	}
}
