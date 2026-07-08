package ko

import (
	"cmp"
	"iter"
)

// Seq is an iterator-backed sequence type for ordered values.
type Seq[T any] iter.Seq[T]

// SeqComparable is an iterator-backed sequence type for comparable values.
type SeqComparable[T comparable] iter.Seq[T]

// SeqOrdered is an iterator-backed sequence type for ordered values.
type SeqOrdered[T cmp.Ordered] iter.Seq[T]

// SeqNumeric is an iterator-backed sequence type for numeric values.
type SeqNumeric[T Numeric] iter.Seq[T]

// Seq2 is an iterator-backed sequence type for key/value entries.
type Seq2[K comparable, V any] iter.Seq2[K, V]

// Numeric is any integer or floating-point type.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}
