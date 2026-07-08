package ko

import (
	"cmp"
	"iter"
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

// FilterReject splits items by predicate.
func (c Seq[T]) FilterReject(predicate func(item T, index int) bool) (Seq[T], Seq[T]) {
	var kept []T
	var rejected []T
	next, stop := iter.Pull(iter.Seq[T](c))
	done := false
	stopped := false
	i := 0
	stopOnce := func() {
		if !stopped {
			stop()
			stopped = true
		}
	}

	advance := func(wantKept bool, cursor int) bool {
		target := &rejected
		if wantKept {
			target = &kept
		}
		for len(*target) <= cursor {
			if done {
				return false
			}
			item, ok := next()
			if !ok {
				done = true
				stopOnce()
				return false
			}
			if predicate(item, i) {
				kept = append(kept, item)
			} else {
				rejected = append(rejected, item)
			}
			i++
		}
		return true
	}

	yieldFrom := func(wantKept bool, yield func(T) bool) {
		cursor := 0
		for advance(wantKept, cursor) {
			items := rejected
			if wantKept {
				items = kept
			}
			if !yield(items[cursor]) {
				return
			}
			cursor++
		}
		if done {
			stopOnce()
		}
	}

	return Seq[T](func(yield func(T) bool) {
			yieldFrom(true, yield)
		}), Seq[T](func(yield func(T) bool) {
			yieldFrom(false, yield)
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
func (c Seq[T]) FlatMap[R any](mapper func(item T, index int) Seq[R]) Seq[R] {
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

// KeyBy indexes items by a key. Later items replace earlier items with the same key.
func (c Seq[T]) KeyBy[K comparable](mapper func(item T, index int) K) Seq2[K, T] {
	return Seq2[K, T](func(yield func(K, T) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			key := mapper(item, i)
			i++
			if !yield(key, item) {
				return
			}
		}
	})
}

// GroupBy groups items by key, preserving first key order.
func (c Seq[T]) GroupBy[K comparable](mapper func(item T, index int) K) iter.Seq2[K, Seq[T]] {
	return func(yield func(K, Seq[T]) bool) {
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
			if !yield(key, Seq[T](slices.Values(result[key]))) {
				return
			}
		}
	}
}

// GroupByMap groups mapped values by key, preserving first key order.
func (c Seq[T]) GroupByMap[K comparable, V any](mapper func(item T, index int) (K, V)) iter.Seq2[K, Seq[V]] {
	return func(yield func(K, Seq[V]) bool) {
		result := make(map[K][]V)
		order := make([]K, 0)
		i := 0
		for item := range iter.Seq[T](c) {
			key, value := mapper(item, i)
			if _, ok := result[key]; !ok {
				order = append(order, key)
			}
			result[key] = append(result[key], value)
			i++
		}
		for _, key := range order {
			if !yield(key, Seq[V](slices.Values(result[key]))) {
				return
			}
		}
	}
}

// PartitionBy groups items by key, preserving first key order.
func (c Seq[T]) PartitionBy[K comparable](mapper func(item T, index int) K) iter.Seq[Seq[T]] {
	return func(yield func(Seq[T]) bool) {
		groups := make([][]T, 0)
		seen := make(map[K]int)
		i := 0
		for item := range iter.Seq[T](c) {
			key := mapper(item, i)
			if index, ok := seen[key]; ok {
				groups[index] = append(groups[index], item)
			} else {
				seen[key] = len(groups)
				groups = append(groups, []T{item})
			}
			i++
		}
		for _, group := range groups {
			if !yield(Seq[T](slices.Values(group))) {
				return
			}
		}
	}
}

// CountBy counts items by key, preserving first key order.
func (c Seq[T]) CountBy[K comparable](mapper func(item T, index int) K) Seq2[K, int] {
	return Seq2[K, int](func(yield func(K, int) bool) {
		result := make(map[K]int)
		order := make([]K, 0)
		i := 0
		for item := range iter.Seq[T](c) {
			key := mapper(item, i)
			if _, ok := result[key]; !ok {
				order = append(order, key)
			}
			result[key]++
			i++
		}
		for _, key := range order {
			if !yield(key, result[key]) {
				return
			}
		}
	})
}

// Concat appends another sequence after the current sequence.
func (c Seq[T]) Concat(other Seq[T]) Seq[T] {
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

// Intersperse inserts sep between adjacent items.
func (c Seq[T]) Intersperse(sep T) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		first := true
		for item := range iter.Seq[T](c) {
			if !first && !yield(sep) {
				return
			}
			first = false
			if !yield(item) {
				return
			}
		}
	})
}

// ForEach calls iteratee for each item and returns the unchanged Seq.
func (c Seq[T]) ForEach(iteratee func(item T, index int)) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			iteratee(item, i)
			i++
			if !yield(item) {
				return
			}
		}
	})
}

// ForEachWhile calls predicate until it returns false and returns the unchanged Seq.
func (c Seq[T]) ForEachWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		i := 0
		running := true
		for item := range iter.Seq[T](c) {
			if running {
				running = predicate(item, i)
			}
			i++
			if !yield(item) {
				return
			}
		}
	})
}

// Chunk splits items into non-overlapping chunks of size n.
func (c Seq[T]) Chunk(n int) iter.Seq[Seq[T]] {
	return func(yield func(Seq[T]) bool) {
		for chunk := range chunkSeq(iter.Seq[T](c), n) {
			if !yield(Seq[T](slices.Values(chunk))) {
				return
			}
		}
	}
}

// Window splits items into overlapping windows of size n, advancing by step.
func (c Seq[T]) Window(n, step int) iter.Seq[Seq[T]] {
	return func(yield func(Seq[T]) bool) {
		for window := range windowSeq(iter.Seq[T](c), n, step) {
			if !yield(Seq[T](slices.Values(window))) {
				return
			}
		}
	}
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

// Subset returns up to length items starting at offset. Negative offsets count from the end.
func (c Seq[T]) Subset(offset, length int) Seq[T] {
	if offset >= 0 {
		return Seq[T](func(yield func(T) bool) {
			if length <= 0 {
				return
			}
			i := 0
			taken := 0
			for item := range iter.Seq[T](c) {
				if i < offset {
					i++
					continue
				}
				if !yield(item) {
					return
				}
				taken++
				if taken >= length {
					return
				}
			}
		})
	}
	return Seq[T](func(yield func(T) bool) {
		for item := range Slice(subset(c.Collect(), offset, length)) {
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

// TakeRightWhile keeps items from the end while predicate returns true.
func (c Seq[T]) TakeRightWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		for item := range Slice(takeRightWhile(c.Collect(), predicate)) {
			if !yield(item) {
				return
			}
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

// DropByIndex drops items by index. Negative indexes count from the end.
func (c Seq[T]) DropByIndex(indexes ...int) Seq[T] {
	for _, index := range indexes {
		if index < 0 {
			return Seq[T](func(yield func(T) bool) {
				for item := range Slice(dropByIndex(c.Collect(), indexes...)) {
					if !yield(item) {
						return
					}
				}
			})
		}
	}
	drop := make(map[int]struct{}, len(indexes))
	for _, index := range indexes {
		drop[index] = struct{}{}
	}
	return Seq[T](func(yield func(T) bool) {
		i := 0
		for item := range iter.Seq[T](c) {
			if _, ok := drop[i]; ok {
				i++
				continue
			}
			if !yield(item) {
				return
			}
			i++
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

// DropRightWhile drops items from the end while predicate returns true.
func (c Seq[T]) DropRightWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		for item := range Slice(dropRightWhile(c.Collect(), predicate)) {
			if !yield(item) {
				return
			}
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

// WithoutBy drops items whose mapped key is excluded.
func (c Seq[T]) WithoutBy[K comparable](mapper func(item T, index int) K, exclude ...K) Seq[T] {
	if len(exclude) == 0 {
		return c
	}
	excluded := make(map[K]struct{}, len(exclude))
	for _, key := range exclude {
		excluded[key] = struct{}{}
	}
	return c.Filter(func(item T, index int) bool {
		_, ok := excluded[mapper(item, index)]
		return !ok
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

// Sort returns a stably sorted copy using compare.
func (c Seq[T]) Sort(compare func(left, right T) int) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		for item := range Slice(slices.SortedStableFunc(iter.Seq[T](c), compare)) {
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

func subset[T any](collection []T, offset, length int) []T {
	if length <= 0 {
		return []T{}
	}
	offset += len(collection)
	if offset < 0 {
		offset = 0
	}
	if offset >= len(collection) {
		return []T{}
	}

	end := offset + length
	if end > len(collection) {
		end = len(collection)
	}
	return collection[offset:end]
}

func takeRightWhile[T any](collection []T, predicate func(item T, index int) bool) []T {
	for i, item := range slices.Backward(collection) {
		if !predicate(item, i) {
			return collection[i+1:]
		}
	}
	return collection
}

func dropByIndex[T any](collection []T, indexes ...int) []T {
	drop := make(map[int]struct{}, len(indexes))
	for _, index := range indexes {
		if index < 0 {
			index += len(collection)
		}
		if index >= 0 && index < len(collection) {
			drop[index] = struct{}{}
		}
	}

	i := -1
	return slices.DeleteFunc(collection, func(T) bool {
		i++
		_, ok := drop[i]
		return ok
	})
}

func dropRightWhile[T any](collection []T, predicate func(item T, index int) bool) []T {
	for i, item := range slices.Backward(collection) {
		if !predicate(item, i) {
			return collection[:i+1]
		}
	}
	return []T{}
}
