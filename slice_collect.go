package ko

import (
	"cmp"
	"iter"
	"slices"
	"strings"
)

// Methods in this file consume Seq values and return end-state values.

// Collect materializes the sequence into a slice.
func (c Seq[T]) Collect() []T {
	return slices.Collect(iter.Seq[T](c))
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

// FindIndex returns the first matching item's index.
func (c Seq[T]) FindIndex(predicate func(item T, index int) bool) (int, bool) {
	_, index, ok := findSeq(iter.Seq[T](c), predicate)
	return index, ok
}

// First returns the first item.
func (c Seq[T]) First() (T, bool) {
	return c.Find(func(_ T, _ int) bool {
		return true
	})
}

// MaxBy returns the item with the greatest mapped ordered key.
func (c Seq[T]) MaxBy[K cmp.Ordered](mapper func(item T, index int) K) (T, bool) {
	var best T
	var bestKey K
	ok := false
	i := 0
	for item := range iter.Seq[T](c) {
		key := mapper(item, i)
		if !ok || cmp.Less(bestKey, key) {
			best = item
			bestKey = key
			ok = true
		}
		i++
	}
	return best, ok
}

// MinBy returns the item with the least mapped ordered key.
func (c Seq[T]) MinBy[K cmp.Ordered](mapper func(item T, index int) K) (T, bool) {
	var best T
	var bestKey K
	ok := false
	i := 0
	for item := range iter.Seq[T](c) {
		key := mapper(item, i)
		if !ok || cmp.Less(key, bestKey) {
			best = item
			bestKey = key
			ok = true
		}
		i++
	}
	return best, ok
}

// SumBy sums the mapped numeric values.
func (c Seq[T]) SumBy[N numeric](mapper func(item T, index int) N) N {
	var sum N
	i := 0
	for item := range iter.Seq[T](c) {
		sum += mapper(item, i)
		i++
	}
	return sum
}

// MeanBy returns the arithmetic mean of mapped numeric values, or 0 for an empty sequence.
func (c Seq[T]) MeanBy[N numeric](mapper func(item T, index int) N) float64 {
	var sum N
	count := 0
	for item := range iter.Seq[T](c) {
		sum += mapper(item, count)
		count++
	}
	if count == 0 {
		return 0
	}
	return float64(sum) / float64(count)
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
	items := c.Collect()
	index += len(items)
	if index < 0 || index >= len(items) {
		var zero T
		return zero, false
	}
	return items[index], true
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
