package ko

import (
	"cmp"
	"iter"
)

// Ordered converts seq into a SeqOrdered, enabling ordered-only methods.
func Ordered[T cmp.Ordered](s Seq[T]) SeqOrdered[T] {
	return SeqOrdered[T](iter.Seq[T](s))
}

// Seq converts s back to the unconstrained Seq type.
func (s SeqOrdered[T]) Seq() Seq[T] {
	return Seq[T](iter.Seq[T](s))
}

// Comparable converts s to SeqComparable.
func (s SeqOrdered[T]) Comparable() SeqComparable[T] {
	return SeqComparable[T](iter.Seq[T](s))
}

// Collect materializes the sequence into a slice.
func (s SeqOrdered[T]) Collect() []T {
	return s.Seq().Collect()
}

// Sort returns a stably sorted copy.
func (s SeqOrdered[T]) Sort() SeqOrdered[T] {
	return SeqOrdered[T](s.Seq().Sort(cmp.Compare[T]))
}

// Max returns the greatest item.
func (s SeqOrdered[T]) Max() (T, bool) {
	return s.Seq().Max(cmp.Compare[T])
}

// Min returns the least item.
func (s SeqOrdered[T]) Min() (T, bool) {
	return s.Seq().Min(cmp.Compare[T])
}

// Distinct yields each distinct element once.
func (s SeqOrdered[T]) Distinct() SeqOrdered[T] {
	return SeqOrdered[T](s.Comparable().Distinct())
}

// Compact drops zero-value elements.
func (s SeqOrdered[T]) Compact() SeqOrdered[T] {
	return SeqOrdered[T](s.Comparable().Compact())
}

// Without excludes any element equal to one of vals.
func (s SeqOrdered[T]) Without(vals ...T) SeqOrdered[T] {
	return SeqOrdered[T](s.Comparable().Without(vals...))
}

// Contains reports whether v occurs in s.
func (s SeqOrdered[T]) Contains(v T) bool {
	return s.Comparable().Contains(v)
}

// IndexOf returns the first index of v in s.
func (s SeqOrdered[T]) IndexOf(v T) (int, bool) {
	return s.Comparable().IndexOf(v)
}

// LastIndexOf returns the last index of v in s.
func (s SeqOrdered[T]) LastIndexOf(v T) (int, bool) {
	return s.Comparable().LastIndexOf(v)
}

// CountValues returns a map from each element to the number of times it occurs.
func (s SeqOrdered[T]) CountValues() map[T]int {
	return s.Comparable().CountValues()
}

// ToSet returns a set of distinct elements.
func (s SeqOrdered[T]) ToSet() map[T]struct{} {
	return s.Comparable().ToSet()
}

// Equal reports whether s and other yield the same elements in the same order.
func (s SeqOrdered[T]) Equal(other SeqOrdered[T]) bool {
	return s.Comparable().Equal(other.Comparable())
}

// Union returns distinct elements across s and others.
func (s SeqOrdered[T]) Union(others ...SeqOrdered[T]) SeqOrdered[T] {
	seqs := make([]Seq[T], 0, len(others)+1)
	seqs = append(seqs, s.Seq())
	for _, other := range others {
		seqs = append(seqs, other.Seq())
	}
	return SeqOrdered[T](Union(seqs...))
}

// Intersect returns elements present in both s and other.
func (s SeqOrdered[T]) Intersect(other SeqOrdered[T]) SeqOrdered[T] {
	return SeqOrdered[T](Intersect(s.Seq(), other.Seq()))
}

// Difference returns elements in s but not in other.
func (s SeqOrdered[T]) Difference(other SeqOrdered[T]) SeqOrdered[T] {
	return SeqOrdered[T](Difference(s.Seq(), other.Seq()))
}

// SymmetricDifference returns elements in exactly one sequence.
func (s SeqOrdered[T]) SymmetricDifference(other SeqOrdered[T]) SeqOrdered[T] {
	return SeqOrdered[T](SymmetricDifference(s.Seq(), other.Seq()))
}

// Filter keeps items for which predicate returns true.
func (s SeqOrdered[T]) Filter(predicate func(item T, index int) bool) SeqOrdered[T] {
	return SeqOrdered[T](s.Seq().Filter(predicate))
}

// Reject drops items for which predicate returns true.
func (s SeqOrdered[T]) Reject(predicate func(item T, index int) bool) SeqOrdered[T] {
	return SeqOrdered[T](s.Seq().Reject(predicate))
}

// Take keeps the first n items.
func (s SeqOrdered[T]) Take(n int) SeqOrdered[T] {
	return SeqOrdered[T](s.Seq().Take(n))
}

// Drop skips the first n items.
func (s SeqOrdered[T]) Drop(n int) SeqOrdered[T] {
	return SeqOrdered[T](s.Seq().Drop(n))
}

// TakeWhile keeps items from the beginning while predicate returns true.
func (s SeqOrdered[T]) TakeWhile(predicate func(item T, index int) bool) SeqOrdered[T] {
	return SeqOrdered[T](s.Seq().TakeWhile(predicate))
}

// DropWhile drops items from the beginning while predicate returns true.
func (s SeqOrdered[T]) DropWhile(predicate func(item T, index int) bool) SeqOrdered[T] {
	return SeqOrdered[T](s.Seq().DropWhile(predicate))
}
