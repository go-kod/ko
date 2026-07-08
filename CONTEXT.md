# ko Collections

This context names the public collection-chain concepts exposed by `ko`.

## Language

**Seq**:
A public chaining concept over an ordered iterator that materializes to a slice. `Seq[T]` is the canonical ordered-value sequence type; raw `iter.Seq` is used when Go generic-method cycles prevent returning `Seq`, such as the outer stream from chunking or windowing.
_Avoid_: SliceChain, ListChain, ChunkChain, slice helper, top-level collection helper

**Seq2**:
A public chaining concept over key/value entries that materializes to a map. `Seq2[K, V]` is the canonical key/value sequence type; raw `iter.Seq2` is used when Go generic-method cycles prevent returning `Seq2`, such as grouped sequence results.
_Avoid_: grouped Seq2 chain, map helper, dictionary helper

**Terminal Aggregate**:
A Seq method that consumes the ordered sequence and returns a scalar result or optional item, such as `Max`, `Min`, `SumBy`, or `MeanBy`.
_Avoid_: map-returning duplicate of an existing chain method, no-predicate overload
