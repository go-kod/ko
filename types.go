package ko

import (
	"iter"
)

// Seq is an iterator-backed sequence type for ordered values.
type Seq[T any] iter.Seq[T]

// Seq2 is an iterator-backed sequence type for key/value entries.
type Seq2[K comparable, V any] iter.Seq2[K, V]

type seqSeq[T any] iter.Seq[Seq[T]]

type groupedSeq[K comparable, V any] iter.Seq2[K, Seq[V]]

// Numeric is any integer or floating-point type.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}
