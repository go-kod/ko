package ko

import (
	"iter"
	"reflect"
	"slices"
	"strings"
)

// Methods in this file can yield while consuming the source.

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

// None reports whether no item matches predicate.
func (c Seq[T]) None(predicate func(item T, index int) bool) bool {
	return !c.Some(predicate)
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

// First returns the first item.
func (c Seq[T]) First() (T, bool) {
	return c.Find(func(_ T, _ int) bool {
		return true
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
