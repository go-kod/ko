package ko

import "slices"

// Slice returns a Seq over collection.
func Slice[T any](collection []T) Seq[T] {
	return Seq[T](slices.Values(collection))
}

// Of returns a Seq over items.
func Of[T any](items ...T) Seq[T] {
	return Slice(items)
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
