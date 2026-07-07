# ko Collections

This context names the public collection-chain concepts exposed by `ko`.

## Language

**Seq**:
A public chaining concept over an ordered iterator that materializes to a slice. `Seq[T]` is the canonical ordered-value chain type; raw `iter.Seq` is used when Go generic-method cycles prevent returning `Seq`, such as the outer stream from chunking.
_Avoid_: SliceChain, ListChain, ChunkChain, slice helper, top-level collection helper

**Seq2**:
A public chaining concept over key/value entries that materializes to a map. `Seq2[K, V]` is the canonical key/value chain type; raw `iter.Seq2` is used when Go generic-method cycles prevent returning `Seq2`, such as grouped sequence results.
_Avoid_: grouped Seq2 chain, map helper, dictionary helper
