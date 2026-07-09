# Simplify Collection API

`ko` keeps `Seq[T]` and `Seq2[K, V]` as the public chain concepts and removes constrained wrapper chains, top-level collection helpers, and combination helpers that duplicate ordinary chaining. Second-order results such as chunks and groups return unexported adapter types with only the small methods needed for direct chaining, avoiding exported `ChunkChain`-style types while also avoiding `Seq[Seq[T]]` instantiation cycles.
