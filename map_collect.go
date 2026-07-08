package ko

import (
	"iter"
	"maps"
)

// Methods in this file materialize Seq2 entries.

// Collect materializes the entries into a map.
func (c Seq2[K, V]) Collect() map[K]V {
	return maps.Collect(iter.Seq2[K, V](c))
}
