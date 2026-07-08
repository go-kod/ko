package ko

import (
	"cmp"
	"iter"
	"slices"
)

// Methods in this file need whole-sequence knowledge before yielding, at least for some inputs.

// Collect materializes the sequence into a slice.
func (c Seq[T]) Collect() []T {
	return slices.Collect(iter.Seq[T](c))
}

// FindUniquesBy keeps items whose mapped key appears exactly once.
func (c Seq[T]) FindUniquesBy[K comparable](mapper func(item T, index int) K) Seq[T] {
	return Seq[T](findUniquesBySeq(iter.Seq[T](c), mapper))
}

// FindDuplicatesBy keeps the first item for each duplicated mapped key.
func (c Seq[T]) FindDuplicatesBy[K comparable](mapper func(item T, index int) K) Seq[T] {
	return Seq[T](findDuplicatesBySeq(iter.Seq[T](c), mapper))
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
			buf = append(slices.Delete(buf, 0, 1), item)
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

// TakeRightWhile keeps items from the end while predicate returns true.
func (c Seq[T]) TakeRightWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return takeRightWhile(items, predicate)
	}))
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

// DropRightWhile drops items from the end while predicate returns true.
func (c Seq[T]) DropRightWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return dropRightWhile(items, predicate)
	}))
}

// Sort returns a stably sorted copy using compare.
func (c Seq[T]) Sort(compare func(left, right T) int) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		slices.SortStableFunc(items, compare)
		return items
	}))
}

// SortBy returns a stably sorted copy using mapped ordered keys.
func (c Seq[T]) SortBy[K cmp.Ordered](mapper func(item T, index int) K) Seq[T] {
	return Seq[T](func(yield func(T) bool) {
		type entry struct {
			item T
			key  K
		}
		entries := []entry{}
		i := 0
		for item := range iter.Seq[T](c) {
			entries = append(entries, entry{item: item, key: mapper(item, i)})
			i++
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
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		result := slices.Clone(items)
		slices.Reverse(result)
		return result
	}))
}

// ReduceRight folds items from right to left.
func (c Seq[T]) ReduceRight[R any](accumulator func(agg R, item T, index int) R, initial R) R {
	items := slices.Collect(iter.Seq[T](c))
	result := initial
	for i := len(items) - 1; i >= 0; i-- {
		result = accumulator(result, items[i], i)
	}
	return result
}

// FindLast returns the last matching item.
func (c Seq[T]) FindLast(predicate func(item T, index int) bool) (T, bool) {
	candidate, _, ok := c.FindLastIndex(predicate)
	return candidate, ok
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
	items := slices.Collect(iter.Seq[T](c))
	index += len(items)
	if index < 0 || index >= len(items) {
		var zero T
		return zero, false
	}
	return items[index], true
}

func materializeSeq[T any](seq iter.Seq[T], transform func([]T) []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range transform(slices.Collect(seq)) {
			if !yield(item) {
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
	return slices.Clone(collection[offset:end])
}

func takeRightWhile[T any](collection []T, predicate func(item T, index int) bool) []T {
	for i := len(collection) - 1; i >= 0; i-- {
		if !predicate(collection[i], i) {
			return slices.Clone(collection[i+1:])
		}
	}
	return slices.Clone(collection)
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
	for i := len(collection) - 1; i >= 0; i-- {
		if !predicate(collection[i], i) {
			return slices.Clone(collection[:i+1])
		}
	}
	return []T{}
}
