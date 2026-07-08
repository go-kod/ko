# Repository Guidelines

## Project Structure & Module Organization

This repository is a small Go module, `github.com/go-kod/ko`, that provides lodash-style generic helpers for Go collections.

- Default to using the `grill-with-docs` and `tdd` skills when changing API shape, behavior, or tests.
- `seq.go` contains the public `Seq` and `Seq2` iterator-backed sequence types.
- `slice.go` contains ordered-value sequence constructors.
- `slice_stream.go` contains `Seq` methods that can yield while consuming the source.
- `slice_collect.go` contains `Seq` methods that need whole-sequence knowledge before yielding a derived sequence or final result, at least for some inputs.
- `map.go` contains key/value sequence constructors.
- `map_stream.go` contains `Seq2` methods that can yield while consuming entries.
- `map_collect.go` contains `Seq2` methods that materialize entries.
- `slice_test.go` contains the current unit tests for Seq chains.
- `map_test.go` contains the current unit tests for Seq2 chains.
- `README.md` documents public usage examples.

Keep new source files at the repository root unless a real package split becomes necessary. Public API changes should be reflected in `README.md`.

## Build, Test, and Development Commands

- `go test ./...` runs all tests in the module.
- `go test -run TestSeq2 ./...` runs one targeted test group.
- `go test -cover ./...` reports package coverage.
- `/home/pilot/sdk/go1.27rc1/bin/gofmt -w *.go` formats Go files with the module's Go 1.27 toolchain.

The module currently targets Go `1.27` with toolchain `go1.27rc1`, as declared in `go.mod`.

## Coding Style & Naming Conventions

Use standard Go formatting via `gofmt`; tabs are expected for indentation in Go files. Keep package name `ko`. Prefer small, direct generic helpers over new abstractions. Ordered collection chains should use `Seq[T]`, a defined type over `iter.Seq[T]`; key/value collection chains should use `Seq2[K, V]`, a defined type over `iter.Seq2[K, V]`. Public collection operations should be reached through exported constructors and methods such as `Slice`, `Of`, `Generate`, `Map`, `Range`, `Collect`, and `Filter`.

Name chain methods as short verbs that match collection operations, such as `Map`, `Filter`, `Reject`, `Reduce`, `Take`, and `Drop`. Internal helpers should stay unexported and lower camel case, for example `mapSeq` and `flatMapSeq`.

Collection operations should be exposed as chain methods, not package-level helper functions. Keep package-level exported functions limited to chain constructors such as `Slice`, `Of`, `Generate`, `Map`, `Range`, `RangeStep`, `Times`, `Repeat`, and `FromChannel`; do not add top-level helpers like `Uniq(collection)` when the operation belongs on `Slice(...).Uniq()` or `Map(...).Keys()`.

Avoid intermediate chain types and conversion-only helpers such as `chunkChain`, grouped Seq2 types, `sliceFromSeq`, or `mapFromSeq`. Return `Seq[T]`, `Seq2[K, V]`, or raw `iter.Seq`/`iter.Seq2` directly when Go 1.27 generic-method instantiation cycles require it, as with `Chunk`, `Window`, and grouped sequence results.

## Testing Guidelines

Use Go's standard `testing` package. Place tests in `*_test.go` files in package `ko`, and name new tests with current public vocabulary, such as `TestSeqWindow` or `TestSeq2ChunkEntries`.

For map-to-slice behavior, sort results before comparison because Go map iteration order is not stable. Use focused table tests only when they reduce repetition; otherwise keep tests simple and readable.

Keep statement coverage at `100%` for `go test -cover ./...`. When adding chain methods, cover both hit and miss/empty boundary behavior where it affects branch coverage.

## Commit & Pull Request Guidelines

This repository has no commit history yet, so use concise imperative commit messages such as `Add slice reject helper` or `Fix map value mapping`.

Pull requests should include a short description, note any public API changes, and mention the test command run. Include README updates when examples or supported methods change.
