package ko

import (
	"iter"
	"reflect"
	"slices"
	"strings"
)

// Slice returns a Seq over collection.
func Slice[T any](collection []T) Seq[T] {
	return Seq[T](slices.Values(collection))
}

// Of returns a Seq over items.
func Of[T any](items ...T) Seq[T] {
	return Slice(items)
}

// Generate creates an infinite sequence by repeatedly applying next.
func Generate[T any](seed T, next func(T) T) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		item := seed
		for {
			if !yield(item) {
				return
			}
			item = next(item)
		}
	})
}

// Range creates a sequence of integers in [start, end).
func Range(start, end int) Seq[int] {
	return Seq[int](func(yield func(int) bool) {
		if start < end {
			for i := start; i < end; i++ {
				if !yield(i) {
					return
				}
			}
			return
		}
		for i := start; i > end; i-- {
			if !yield(i) {
				return
			}
		}
	})
}

// RangeStep creates a sequence of integers in [start, end), advancing by step.
func RangeStep(start, end, step int) Seq[int] {
	return Seq[int](func(yield func(int) bool) {
		if step <= 0 {
			return
		}
		if start < end {
			for i := start; i < end; i += step {
				if !yield(i) {
					return
				}
			}
			return
		}
		for i := start; i > end; i -= step {
			if !yield(i) {
				return
			}
		}
	})
}

// Times creates a sequence by calling mapper for indexes [0, n).
func Times[T any](n int, mapper func(index int) T) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			if !yield(mapper(i)) {
				return
			}
		}
	})
}

// Repeat creates a sequence that yields value n times.
func Repeat[T any](n int, value T) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			if !yield(value) {
				return
			}
		}
	})
}

// FromChannel creates a one-shot sequence that drains channel.
func FromChannel[T any](channel <-chan T) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		for item := range channel {
			if !yield(item) {
				return
			}
		}
	})
}

func collectSeq[T any](seq iter.Seq[T]) []T {
	items := slices.Collect(seq)
	if items == nil {
		return []T{}
	}
	return items
}

// Collect materializes the sequence into a slice.
func (c Seq[T]) Collect() []T {
	return collectSeq(iter.Seq[T](c))
}

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

// Compact drops zero-value items.
func (c Seq[T]) Compact() Seq[T] {
	return c.Filter(func(item T, _ int) bool {
		return !reflect.ValueOf(&item).Elem().IsZero()
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

// Uniq removes duplicate comparable items, keeping the first occurrence.
func (c Seq[T]) Uniq() Seq[T] {
	return Seq[T](uniqSeq(iter.Seq[T](c)))
}

// UniqBy removes duplicate items by key, keeping the first occurrence.
func (c Seq[T]) UniqBy[K comparable](mapper func(item T, index int) K) Seq[T] {
	return Seq[T](uniqBySeq(iter.Seq[T](c), mapper))
}

// IsUniqBy reports whether all mapped keys are unique.
func (c Seq[T]) IsUniqBy[K comparable](mapper func(item T, index int) K) bool {
	seen := make(map[K]struct{})
	i := 0
	for item := range iter.Seq[T](c) {
		key := mapper(item, i)
		if _, ok := seen[key]; ok {
			return false
		}
		seen[key] = struct{}{}
		i++
	}
	return true
}

// UniqMap maps items and keeps the first occurrence of each mapped value.
func (c Seq[T]) UniqMap[R comparable](mapper func(item T, index int) R) Seq[R] {
	return Seq[R](uniqMapSeq(iter.Seq[T](c), mapper))
}

// FindUniquesBy keeps items whose mapped key appears exactly once.
func (c Seq[T]) FindUniquesBy[K comparable](mapper func(item T, index int) K) Seq[T] {
	return Seq[T](findUniquesBySeq(iter.Seq[T](c), mapper))
}

// FindDuplicatesBy keeps the first item for each duplicated mapped key.
func (c Seq[T]) FindDuplicatesBy[K comparable](mapper func(item T, index int) K) Seq[T] {
	return Seq[T](findDuplicatesBySeq(iter.Seq[T](c), mapper))
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
		buf := make([]T, 0, n)
		for item := range iter.Seq[T](c) {
			if len(buf) < n {
				buf = append(buf, item)
				continue
			}
			copy(buf, buf[1:])
			buf[n-1] = item
		}
		for _, item := range buf {
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
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return subset(items, offset, length)
	}))
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
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return takeRightWhile(items, predicate)
	}))
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
			copy(buf, buf[1:])
			buf = buf[:n]
		}
	})
}

// DropByIndex drops items by index. Negative indexes count from the end.
func (c Seq[T]) DropByIndex(indexes ...int) Seq[T] {
	for _, index := range indexes {
		if index < 0 {
			return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
				return dropByIndex(items, indexes...)
			}))
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

// DropRightWhile drops items from the end while predicate returns true.
func (c Seq[T]) DropRightWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return dropRightWhile(items, predicate)
	}))
}

// Reverse returns a reversed copy.
func (c Seq[T]) Reverse() Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		result := make([]T, len(items))
		for i, item := range items {
			result[len(items)-1-i] = item
		}
		return result
	}))
}

// Reduce folds collection into one value.
func (c Seq[T]) Reduce[R any](accumulator func(agg R, item T, index int) R, initial R) R {
	result := initial
	i := 0
	for item := range iter.Seq[T](c) {
		result = accumulator(result, item, i)
		i++
	}
	return result
}

// ReduceRight folds items from right to left.
func (c Seq[T]) ReduceRight[R any](accumulator func(agg R, item T, index int) R, initial R) R {
	items := collectSeq(iter.Seq[T](c))
	result := initial
	for i := len(items) - 1; i >= 0; i-- {
		result = accumulator(result, items[i], i)
	}
	return result
}

// Join maps items to strings and joins them with sep.
func (c Seq[T]) Join(sep string, mapper func(item T, index int) string) string {
	var builder strings.Builder
	i := 0
	first := true
	for item := range iter.Seq[T](c) {
		if !first {
			builder.WriteString(sep)
		}
		first = false
		builder.WriteString(mapper(item, i))
		i++
	}
	return builder.String()
}

// IsEmpty reports whether the sequence yields no items.
func (c Seq[T]) IsEmpty() bool {
	return !c.Some(func(_ T, _ int) bool {
		return true
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

// Some reports whether any item matches predicate.
func (c Seq[T]) Some(predicate func(item T, index int) bool) bool {
	_, _, ok := findSeq(iter.Seq[T](c), predicate)
	return ok
}

// Count returns the number of items matching predicate.
func (c Seq[T]) Count(predicate func(item T, index int) bool) int {
	result := 0
	i := 0
	for item := range iter.Seq[T](c) {
		if predicate(item, i) {
			result++
		}
		i++
	}
	return result
}

// Every reports whether all items match predicate.
func (c Seq[T]) Every(predicate func(item T, index int) bool) bool {
	return !c.Some(func(item T, index int) bool {
		return !predicate(item, index)
	})
}

// Find returns the first matching item.
func (c Seq[T]) Find(predicate func(item T, index int) bool) (T, bool) {
	item, _, ok := findSeq(iter.Seq[T](c), predicate)
	return item, ok
}

// FindIndex returns the first matching item, its index, and whether it was found.
func (c Seq[T]) FindIndex(predicate func(item T, index int) bool) (T, int, bool) {
	return findSeq(iter.Seq[T](c), predicate)
}

// FindLast returns the last matching item.
func (c Seq[T]) FindLast(predicate func(item T, index int) bool) (T, bool) {
	return c.Filter(predicate).Last()
}

// FindLastIndex returns the last matching item, its index, and whether it was found.
func (c Seq[T]) FindLastIndex(predicate func(item T, index int) bool) (T, int, bool) {
	var result T
	resultIndex := -1
	i := 0
	for item := range iter.Seq[T](c) {
		if predicate(item, i) {
			result = item
			resultIndex = i
		}
		i++
	}
	return result, resultIndex, resultIndex >= 0
}

// First returns the first item.
func (c Seq[T]) First() (T, bool) {
	return c.Find(func(_ T, _ int) bool {
		return true
	})
}

// Last returns the last item.
func (c Seq[T]) Last() (T, bool) {
	var result T
	ok := false
	for item := range iter.Seq[T](c) {
		result = item
		ok = true
	}
	return result, ok
}

// Nth returns the item at index. Negative indexes count from the end.
func (c Seq[T]) Nth(index int) (T, bool) {
	if index >= 0 {
		return c.Find(func(_ T, itemIndex int) bool {
			return itemIndex == index
		})
	}
	items := collectSeq(iter.Seq[T](c))
	index += len(items)
	if index < 0 || index >= len(items) {
		var zero T
		return zero, false
	}
	return items[index], true
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

func materializeSeq[T any](seq iter.Seq[T], transform func([]T) []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range transform(collectSeq(seq)) {
			if !yield(item) {
				return
			}
		}
	}
}

func uniqSeq[T any](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := map[any]struct{}{}
		for item := range seq {
			key := any(item)
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

func uniqBySeq[T any, K comparable](seq iter.Seq[T], mapper func(item T, index int) K) iter.Seq[T] {
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

func uniqMapSeq[T any, R comparable](seq iter.Seq[T], mapper func(item T, index int) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		seen := map[R]struct{}{}
		i := 0
		for item := range seq {
			mapped := mapper(item, i)
			i++
			if _, ok := seen[mapped]; ok {
				continue
			}
			seen[mapped] = struct{}{}
			if !yield(mapped) {
				return
			}
		}
	}
}

func findUniquesBySeq[T any, K comparable](seq iter.Seq[T], mapper func(item T, index int) K) iter.Seq[T] {
	return func(yield func(T) bool) {
		counts := make(map[K]int)
		first := make(map[K]T)
		order := make([]K, 0)
		i := 0
		for item := range seq {
			key := mapper(item, i)
			if _, ok := counts[key]; !ok {
				first[key] = item
				order = append(order, key)
			}
			counts[key]++
			i++
		}
		for _, key := range order {
			if counts[key] == 1 && !yield(first[key]) {
				return
			}
		}
	}
}

func findDuplicatesBySeq[T any, K comparable](seq iter.Seq[T], mapper func(item T, index int) K) iter.Seq[T] {
	return func(yield func(T) bool) {
		counts := make(map[K]int)
		first := make(map[K]T)
		order := make([]K, 0)
		i := 0
		for item := range seq {
			key := mapper(item, i)
			if _, ok := counts[key]; !ok {
				first[key] = item
				order = append(order, key)
			}
			counts[key]++
			i++
		}
		for _, key := range order {
			if counts[key] > 1 && !yield(first[key]) {
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
			if !yield(append([]T(nil), chunk...)) {
				return
			}
			chunk = chunk[:0]
		}
		if len(chunk) > 0 {
			yield(append([]T(nil), chunk...))
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
			if !yield(append([]T(nil), window...)) {
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

func findSeq[T any](seq iter.Seq[T], predicate func(item T, index int) bool) (T, int, bool) {
	i := 0
	for item := range seq {
		if predicate(item, i) {
			return item, i, true
		}
		i++
	}
	var zero T
	return zero, -1, false
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
	return append([]T(nil), collection[offset:end]...)
}

func takeRightWhile[T any](collection []T, predicate func(item T, index int) bool) []T {
	for i := len(collection) - 1; i >= 0; i-- {
		if !predicate(collection[i], i) {
			return append([]T{}, collection[i+1:]...)
		}
	}
	return append([]T(nil), collection...)
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

	result := make([]T, 0, len(collection)-len(drop))
	for i, item := range collection {
		if _, ok := drop[i]; !ok {
			result = append(result, item)
		}
	}
	return result
}

func dropRightWhile[T any](collection []T, predicate func(item T, index int) bool) []T {
	for i := len(collection) - 1; i >= 0; i-- {
		if !predicate(collection[i], i) {
			return append([]T(nil), collection[:i+1]...)
		}
	}
	return []T{}
}
