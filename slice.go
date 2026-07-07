package ko

import (
	"iter"
	"slices"
)

// Slice starts a chain over collection.
func Slice[T any](collection []T) Seq[T] {
	return Seq[T](slices.Values(collection))
}

func collectSeq[T any](seq iter.Seq[T]) []T {
	items := slices.Collect(seq)
	if items == nil {
		return []T{}
	}
	return items
}

// Collect returns the current slice.
func (c Seq[T]) Collect() []T {
	return collectSeq(iter.Seq[T](c))
}

// Filter keeps items for which predicate returns true.
func (c Seq[T]) Filter(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](filterSeq(iter.Seq[T](c), predicate))
}

// Reject drops items for which predicate returns true.
func (c Seq[T]) Reject(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](rejectSeq(iter.Seq[T](c), predicate))
}

// FilterReject splits items by predicate.
func (c Seq[T]) FilterReject(predicate func(item T, index int) bool) (Seq[T], Seq[T]) {
	kept, rejected := filterReject(collectSeq(iter.Seq[T](c)), predicate)
	return Seq[T](slices.Values(kept)), Seq[T](slices.Values(rejected))
}

// Uniq removes duplicate comparable items, keeping the first occurrence.
func (c Seq[T]) Uniq() Seq[T] {
	return Seq[T](uniqSeq(iter.Seq[T](c)))
}

// UniqBy removes duplicate items by key, keeping the first occurrence.
func (c Seq[T]) UniqBy[K comparable](mapper func(item T, index int) K) Seq[T] {
	return Seq[T](uniqBySeq(iter.Seq[T](c), mapper))
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
	return Seq[R](mapSeq(iter.Seq[T](c), mapper))
}

// FilterMap maps items and keeps only accepted results.
func (c Seq[T]) FilterMap[R any](mapper func(item T, index int) (R, bool)) Seq[R] {
	return Seq[R](filterMapSeq(iter.Seq[T](c), mapper))
}

// FlatMap transforms each item into zero or more output items.
func (c Seq[T]) FlatMap[R any](mapper func(item T, index int) []R) Seq[R] {
	return Seq[R](flatMapSeq(iter.Seq[T](c), mapper))
}

// GroupBy groups items by a key.
func (c Seq[T]) GroupBy[K comparable](mapper func(item T, index int) K) iter.Seq2[K, Seq[T]] {
	return func(yield func(K, Seq[T]) bool) {
		result := make(map[K][]T)
		i := 0
		for item := range iter.Seq[T](c) {
			key := mapper(item, i)
			result[key] = append(result[key], item)
			i++
		}
		for key, value := range result {
			if !yield(key, Seq[T](slices.Values(value))) {
				return
			}
		}
	}
}

// GroupByMap groups mapped values by a key.
func (c Seq[T]) GroupByMap[K comparable, V any](mapper func(item T, index int) (K, V)) iter.Seq2[K, Seq[V]] {
	return func(yield func(K, Seq[V]) bool) {
		result := make(map[K][]V)
		i := 0
		for item := range iter.Seq[T](c) {
			key, value := mapper(item, i)
			result[key] = append(result[key], value)
			i++
		}
		for key, value := range result {
			if !yield(key, Seq[V](slices.Values(value))) {
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

// KeyBy indexes items by a key. Later items replace earlier items with the same key.
func (c Seq[T]) KeyBy[K comparable](mapper func(item T, index int) K) Seq2[K, T] {
	return Seq2[K, T](func(yield func(K, T) bool) {
		result := make(map[K]T)
		i := 0
		for item := range iter.Seq[T](c) {
			result[mapper(item, i)] = item
			i++
		}
		for key, value := range result {
			if !yield(key, value) {
				return
			}
		}
	})
}

// CountBy counts items by a key.
func (c Seq[T]) CountBy[K comparable](mapper func(item T, index int) K) Seq2[K, int] {
	return Seq2[K, int](func(yield func(K, int) bool) {
		result := make(map[K]int)
		i := 0
		for item := range iter.Seq[T](c) {
			result[mapper(item, i)]++
			i++
		}
		for key, value := range result {
			if !yield(key, value) {
				return
			}
		}
	})
}

// ForEach calls iteratee for each item and returns the unchanged chain.
func (c Seq[T]) ForEach(iteratee func(item T, index int)) Seq[T] {
	forEachSeq(iter.Seq[T](c), iteratee)
	return c
}

// ForEachWhile calls predicate until it returns false and returns the unchanged chain.
func (c Seq[T]) ForEachWhile(predicate func(item T, index int) bool) Seq[T] {
	forEachWhileSeq(iter.Seq[T](c), predicate)
	return c
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

// Window splits items into overlapping windows of size n.
func (c Seq[T]) Window(n int) iter.Seq[Seq[T]] {
	return func(yield func(Seq[T]) bool) {
		for window := range windowSeq(iter.Seq[T](c), n) {
			if !yield(Seq[T](slices.Values(window))) {
				return
			}
		}
	}
}

// Take keeps the first n items.
func (c Seq[T]) Take(n int) Seq[T] {
	return Seq[T](takeSeq(iter.Seq[T](c), n))
}

// TakeRight keeps the last n items.
func (c Seq[T]) TakeRight(n int) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return takeRight(items, n)
	}))
}

// Subset returns up to length items starting at offset. Negative offsets count from the end.
func (c Seq[T]) Subset(offset, length int) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return subset(items, offset, length)
	}))
}

// TakeWhile keeps items from the beginning while predicate returns true.
func (c Seq[T]) TakeWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](takeWhileSeq(iter.Seq[T](c), predicate))
}

// TakeFilter keeps the first n matching items.
func (c Seq[T]) TakeFilter(n int, predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](takeFilterSeq(iter.Seq[T](c), n, predicate))
}

// TakeRightWhile keeps items from the end while predicate returns true.
func (c Seq[T]) TakeRightWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return takeRightWhile(items, predicate)
	}))
}

// Drop skips the first n items.
func (c Seq[T]) Drop(n int) Seq[T] {
	return Seq[T](dropSeq(iter.Seq[T](c), n))
}

// DropRight skips the last n items.
func (c Seq[T]) DropRight(n int) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return dropRight(items, n)
	}))
}

// DropByIndex drops items by index. Negative indexes count from the end.
func (c Seq[T]) DropByIndex(indexes ...int) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return dropByIndex(items, indexes...)
	}))
}

// DropWhile drops items from the beginning while predicate returns true.
func (c Seq[T]) DropWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](dropWhileSeq(iter.Seq[T](c), predicate))
}

// DropRightWhile drops items from the end while predicate returns true.
func (c Seq[T]) DropRightWhile(predicate func(item T, index int) bool) Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), func(items []T) []T {
		return dropRightWhile(items, predicate)
	}))
}

// Reverse returns a reversed copy.
func (c Seq[T]) Reverse() Seq[T] {
	return Seq[T](materializeSeq(iter.Seq[T](c), reverse))
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
	return reduceRight(collectSeq(iter.Seq[T](c)), accumulator, initial)
}

// Some reports whether any item matches predicate.
func (c Seq[T]) Some(predicate func(item T, index int) bool) bool {
	_, _, ok := findSeq(iter.Seq[T](c), predicate)
	return ok
}

// None reports whether no item matches predicate.
func (c Seq[T]) None(predicate func(item T, index int) bool) bool {
	_, _, ok := findSeq(iter.Seq[T](c), predicate)
	return !ok
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
	i := 0
	for item := range iter.Seq[T](c) {
		if !predicate(item, i) {
			return false
		}
		i++
	}
	return true
}

// Find returns the first matching item.
func (c Seq[T]) Find(predicate func(item T, index int) bool) (T, bool) {
	item, _, ok := findSeq(iter.Seq[T](c), predicate)
	return item, ok
}

// FindOrElse returns the first matching item or fallback.
func (c Seq[T]) FindOrElse(fallback T, predicate func(item T, index int) bool) T {
	item, _, ok := findSeq(iter.Seq[T](c), predicate)
	if !ok {
		return fallback
	}
	return item
}

// FindIndex returns the first matching item, its index, and whether it was found.
func (c Seq[T]) FindIndex(predicate func(item T, index int) bool) (T, int, bool) {
	return findSeq(iter.Seq[T](c), predicate)
}

// FindLast returns the last matching item.
func (c Seq[T]) FindLast(predicate func(item T, index int) bool) (T, bool) {
	item, _, ok := c.FindLastIndex(predicate)
	return item, ok
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
	for item := range iter.Seq[T](c) {
		return item, true
	}
	var zero T
	return zero, false
}

// FirstOr returns the first item or fallback.
func (c Seq[T]) FirstOr(fallback T) T {
	item, ok := c.First()
	if !ok {
		return fallback
	}
	return item
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

// LastOr returns the last item or fallback.
func (c Seq[T]) LastOr(fallback T) T {
	item, ok := c.Last()
	if !ok {
		return fallback
	}
	return item
}

// Nth returns the item at index. Negative indexes count from the end.
func (c Seq[T]) Nth(index int) (T, bool) {
	if index >= 0 {
		i := 0
		for item := range iter.Seq[T](c) {
			if i == index {
				return item, true
			}
			i++
		}
		var zero T
		return zero, false
	}
	return nth(collectSeq(iter.Seq[T](c)), index)
}

// NthOr returns the item at index or fallback. Negative indexes count from the end.
func (c Seq[T]) NthOr(index int, fallback T) T {
	item, ok := c.Nth(index)
	if !ok {
		return fallback
	}
	return item
}

// ContainsBy reports whether any item matches predicate.
func (c Seq[T]) ContainsBy(predicate func(item T, index int) bool) bool {
	_, _, ok := findSeq(iter.Seq[T](c), predicate)
	return ok
}

// WithoutBy drops items whose mapped key is excluded.
func (c Seq[T]) WithoutBy[K comparable](mapper func(item T, index int) K, exclude ...K) Seq[T] {
	excluded := make(map[K]struct{}, len(exclude))
	for _, key := range exclude {
		excluded[key] = struct{}{}
	}
	return Seq[T](filterSeq(iter.Seq[T](c), func(item T, index int) bool {
		_, ok := excluded[mapper(item, index)]
		return !ok
	}))
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

func filterSeq[T any](seq iter.Seq[T], predicate func(item T, index int) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for item := range seq {
			if predicate(item, i) && !yield(item) {
				return
			}
			i++
		}
	}
}

func rejectSeq[T any](seq iter.Seq[T], predicate func(item T, index int) bool) iter.Seq[T] {
	return filterSeq(seq, func(item T, index int) bool {
		return !predicate(item, index)
	})
}

func mapSeq[T, R any](seq iter.Seq[T], mapper func(item T, index int) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		i := 0
		for item := range seq {
			if !yield(mapper(item, i)) {
				return
			}
			i++
		}
	}
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
	return materializeSeq(seq, func(items []T) []T {
		return findUniquesBy(items, mapper)
	})
}

func findDuplicatesBySeq[T any, K comparable](seq iter.Seq[T], mapper func(item T, index int) K) iter.Seq[T] {
	return materializeSeq(seq, func(items []T) []T {
		return findDuplicatesBy(items, mapper)
	})
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

func windowSeq[T any](seq iter.Seq[T], n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if n <= 0 {
			return
		}
		window := make([]T, 0, n)
		for item := range seq {
			window = append(window, item)
			if len(window) < n {
				continue
			}
			if !yield(append([]T(nil), window...)) {
				return
			}
			// ponytail: copy-shift is enough; switch to a ring buffer if large windows matter.
			copy(window, window[1:])
			window = window[:n-1]
		}
	}
}

func takeSeq[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		i := 0
		for item := range seq {
			if !yield(item) {
				return
			}
			i++
			if i >= n {
				return
			}
		}
	}
}

func takeWhileSeq[T any](seq iter.Seq[T], predicate func(item T, index int) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for item := range seq {
			if !predicate(item, i) || !yield(item) {
				return
			}
			i++
		}
	}
}

func takeFilterSeq[T any](seq iter.Seq[T], n int, predicate func(item T, index int) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		i := 0
		taken := 0
		for item := range seq {
			if predicate(item, i) {
				if !yield(item) {
					return
				}
				taken++
				if taken >= n {
					return
				}
			}
			i++
		}
	}
}

func dropSeq[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for item := range seq {
			if i < n {
				i++
				continue
			}
			if !yield(item) {
				return
			}
		}
	}
}

func dropWhileSeq[T any](seq iter.Seq[T], predicate func(item T, index int) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		dropping := true
		for item := range seq {
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
	}
}

func forEachSeq[T any](seq iter.Seq[T], iteratee func(item T, index int)) {
	i := 0
	for item := range seq {
		iteratee(item, i)
		i++
	}
}

func forEachWhileSeq[T any](seq iter.Seq[T], predicate func(item T, index int) bool) {
	i := 0
	for item := range seq {
		if !predicate(item, i) {
			return
		}
		i++
	}
}

func filterMapSeq[T, R any](seq iter.Seq[T], mapper func(item T, index int) (R, bool)) iter.Seq[R] {
	return func(yield func(R) bool) {
		i := 0
		for item := range seq {
			if mapped, ok := mapper(item, i); ok && !yield(mapped) {
				return
			}
			i++
		}
	}
}

func flatMapSeq[T, R any](seq iter.Seq[T], mapper func(item T, index int) []R) iter.Seq[R] {
	return func(yield func(R) bool) {
		i := 0
		for item := range seq {
			for _, mapped := range mapper(item, i) {
				if !yield(mapped) {
					return
				}
			}
			i++
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

func filterReject[T any](collection []T, predicate func(item T, index int) bool) ([]T, []T) {
	kept := make([]T, 0, len(collection))
	rejected := make([]T, 0, len(collection))
	for i, item := range collection {
		if predicate(item, i) {
			kept = append(kept, item)
			continue
		}
		rejected = append(rejected, item)
	}
	return kept, rejected
}

func findUniquesBy[T any, K comparable](collection []T, mapper func(item T, index int) K) []T {
	counts := make(map[K]int, len(collection))
	keys := make([]K, len(collection))
	for i, item := range collection {
		key := mapper(item, i)
		keys[i] = key
		counts[key]++
	}

	result := make([]T, 0, len(collection))
	for i, item := range collection {
		if counts[keys[i]] == 1 {
			result = append(result, item)
		}
	}
	return result
}

func findDuplicatesBy[T any, K comparable](collection []T, mapper func(item T, index int) K) []T {
	counts := make(map[K]int, len(collection))
	keys := make([]K, len(collection))
	for i, item := range collection {
		key := mapper(item, i)
		keys[i] = key
		counts[key]++
	}

	result := make([]T, 0)
	added := make(map[K]struct{})
	for i, item := range collection {
		key := keys[i]
		if counts[key] < 2 {
			continue
		}
		if _, ok := added[key]; ok {
			continue
		}
		added[key] = struct{}{}
		result = append(result, item)
	}
	return result
}

func reduceRight[T, R any](collection []T, accumulator func(agg R, item T, index int) R, initial R) R {
	result := initial
	for i := len(collection) - 1; i >= 0; i-- {
		result = accumulator(result, collection[i], i)
	}
	return result
}

func nth[T any](collection []T, index int) (T, bool) {
	if index < 0 {
		index += len(collection)
	}
	if index < 0 || index >= len(collection) {
		var zero T
		return zero, false
	}
	return collection[index], true
}

func takeRight[T any](collection []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(collection) {
		return append([]T(nil), collection...)
	}
	return append([]T(nil), collection[len(collection)-n:]...)
}

func subset[T any](collection []T, offset, length int) []T {
	if length <= 0 {
		return []T{}
	}
	if offset < 0 {
		offset += len(collection)
		if offset < 0 {
			offset = 0
		}
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

func dropRight[T any](collection []T, n int) []T {
	if n <= 0 {
		return append([]T(nil), collection...)
	}
	if n >= len(collection) {
		return []T{}
	}
	return append([]T(nil), collection[:len(collection)-n]...)
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

func reverse[T any](collection []T) []T {
	result := make([]T, len(collection))
	for i, item := range collection {
		result[len(collection)-1-i] = item
	}
	return result
}
