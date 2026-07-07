package ko

import "iter"

// Seq is an iterator-backed chain type for ordered values.
type Seq[T any] iter.Seq[T]

// Seq2 is an iterator-backed chain type for key/value entries.
type Seq2[K comparable, V any] iter.Seq2[K, V]
