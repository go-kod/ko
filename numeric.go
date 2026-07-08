package ko

import "iter"

// Numbers converts seq into a SeqNumeric, enabling numeric-only methods.
func Numbers[T Numeric](s Seq[T]) SeqNumeric[T] {
	return SeqNumeric[T](iter.Seq[T](s))
}

// Seq converts s back to the unconstrained Seq type.
func (s SeqNumeric[T]) Seq() Seq[T] {
	return Seq[T](iter.Seq[T](s))
}

// Ordered converts s to SeqOrdered.
func (s SeqNumeric[T]) Ordered() SeqOrdered[T] {
	return SeqOrdered[T](iter.Seq[T](s))
}

// Comparable converts s to SeqComparable.
func (s SeqNumeric[T]) Comparable() SeqComparable[T] {
	return s.Ordered().Comparable()
}

// Collect materializes the sequence into a slice.
func (s SeqNumeric[T]) Collect() []T {
	return s.Seq().Collect()
}

// Sum returns the sum of all items.
func (s SeqNumeric[T]) Sum() T {
	var sum T
	for item := range iter.Seq[T](s) {
		sum += item
	}
	return sum
}

// Product returns the product of all items, or 1 for an empty sequence.
func (s SeqNumeric[T]) Product() T {
	var product T = 1
	for item := range iter.Seq[T](s) {
		product *= item
	}
	return product
}

// Mean returns the arithmetic mean, or 0 for an empty sequence.
func (s SeqNumeric[T]) Mean() float64 {
	var sum T
	count := 0
	for item := range iter.Seq[T](s) {
		sum += item
		count++
	}
	if count == 0 {
		return 0
	}
	return float64(sum) / float64(count)
}

// Sort returns a stably sorted copy.
func (s SeqNumeric[T]) Sort() SeqNumeric[T] {
	return SeqNumeric[T](s.Ordered().Sort())
}

// Max returns the greatest item.
func (s SeqNumeric[T]) Max() (T, bool) {
	return s.Ordered().Max()
}

// Min returns the least item.
func (s SeqNumeric[T]) Min() (T, bool) {
	return s.Ordered().Min()
}

// Distinct yields each distinct element once.
func (s SeqNumeric[T]) Distinct() SeqNumeric[T] {
	return SeqNumeric[T](s.Comparable().Distinct())
}

// Compact drops zero-value elements.
func (s SeqNumeric[T]) Compact() SeqNumeric[T] {
	return SeqNumeric[T](s.Comparable().Compact())
}

// Without excludes any element equal to one of vals.
func (s SeqNumeric[T]) Without(vals ...T) SeqNumeric[T] {
	return SeqNumeric[T](s.Comparable().Without(vals...))
}

// Contains reports whether v occurs in s.
func (s SeqNumeric[T]) Contains(v T) bool {
	return s.Comparable().Contains(v)
}

// IndexOf returns the first index of v in s.
func (s SeqNumeric[T]) IndexOf(v T) (int, bool) {
	return s.Comparable().IndexOf(v)
}

// LastIndexOf returns the last index of v in s.
func (s SeqNumeric[T]) LastIndexOf(v T) (int, bool) {
	return s.Comparable().LastIndexOf(v)
}

// CountValues returns a map from each element to the number of times it occurs.
func (s SeqNumeric[T]) CountValues() map[T]int {
	return s.Comparable().CountValues()
}

// ToSet returns a set of distinct elements.
func (s SeqNumeric[T]) ToSet() map[T]struct{} {
	return s.Comparable().ToSet()
}

// Equal reports whether s and other yield the same elements in the same order.
func (s SeqNumeric[T]) Equal(other SeqNumeric[T]) bool {
	return s.Comparable().Equal(other.Comparable())
}

// Union returns distinct elements across s and others.
func (s SeqNumeric[T]) Union(others ...SeqNumeric[T]) SeqNumeric[T] {
	seqs := make([]Seq[T], 0, len(others)+1)
	seqs = append(seqs, s.Seq())
	for _, other := range others {
		seqs = append(seqs, other.Seq())
	}
	return SeqNumeric[T](Union(seqs...))
}

// Intersect returns elements present in both s and other.
func (s SeqNumeric[T]) Intersect(other SeqNumeric[T]) SeqNumeric[T] {
	return SeqNumeric[T](Intersect(s.Seq(), other.Seq()))
}

// Difference returns elements in s but not in other.
func (s SeqNumeric[T]) Difference(other SeqNumeric[T]) SeqNumeric[T] {
	return SeqNumeric[T](Difference(s.Seq(), other.Seq()))
}

// SymmetricDifference returns elements in exactly one sequence.
func (s SeqNumeric[T]) SymmetricDifference(other SeqNumeric[T]) SeqNumeric[T] {
	return SeqNumeric[T](SymmetricDifference(s.Seq(), other.Seq()))
}

// Filter keeps items for which predicate returns true.
func (s SeqNumeric[T]) Filter(predicate func(item T, index int) bool) SeqNumeric[T] {
	return SeqNumeric[T](s.Seq().Filter(predicate))
}

// Reject drops items for which predicate returns true.
func (s SeqNumeric[T]) Reject(predicate func(item T, index int) bool) SeqNumeric[T] {
	return SeqNumeric[T](s.Seq().Reject(predicate))
}

// Take keeps the first n items.
func (s SeqNumeric[T]) Take(n int) SeqNumeric[T] {
	return SeqNumeric[T](s.Seq().Take(n))
}

// Drop skips the first n items.
func (s SeqNumeric[T]) Drop(n int) SeqNumeric[T] {
	return SeqNumeric[T](s.Seq().Drop(n))
}

// TakeWhile keeps items from the beginning while predicate returns true.
func (s SeqNumeric[T]) TakeWhile(predicate func(item T, index int) bool) SeqNumeric[T] {
	return SeqNumeric[T](s.Seq().TakeWhile(predicate))
}

// DropWhile drops items from the beginning while predicate returns true.
func (s SeqNumeric[T]) DropWhile(predicate func(item T, index int) bool) SeqNumeric[T] {
	return SeqNumeric[T](s.Seq().DropWhile(predicate))
}
