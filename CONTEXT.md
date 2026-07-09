# ko Collections

This context names the public collection-chain concepts exposed by `ko`.

## Language

**Seq**:
A public chaining concept over an ordered iterator that materializes to a slice. `Seq[T]` is the canonical ordered-value sequence type; second-order ordered results may use unexported adapters when `Seq[Seq[T]]` would create Go generic-method instantiation cycles.
_Avoid_: SliceChain, ListChain, ChunkChain, constrained wrapper chain, slice helper, top-level collection helper

**Seq2**:
A public chaining concept over key/value entries that materializes to a map. `Seq2[K, V]` is the canonical key/value sequence type; grouped results may use unexported adapters rather than exported grouped chain types.
_Avoid_: grouped Seq2 chain, map helper, dictionary helper

**Terminal Aggregate**:
A Seq method that consumes the ordered sequence and returns a scalar result or optional item, such as `MaxBy`, `MinBy`, `SumBy`, or `MeanBy`.
_Avoid_: map-returning duplicate of an existing chain method, no-predicate overload
