package ko

import "maps"

// Map returns a Seq2 over collection.
func Map[K comparable, V any](collection map[K]V) Seq2[K, V] {
	return Seq2[K, V](maps.All(collection))
}
