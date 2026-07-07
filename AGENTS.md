# Repository Guidelines

## Project Structure & Module Organization

This repository is a small Go module, `github.com/go-kod/ko`, that provides lodash-style generic helpers for Go collections.

- Default to using the `grill-with-docs` and `tdd` skills when changing API shape, behavior, or tests.
- `slice.go` contains `Seq`, the iterator-backed ordered-value chain type, and slice helper functions.
- `map.go` contains `Seq2`, the iterator-backed key/value chain type, and map helper functions.
- `slice_test.go` contains the current unit tests for both slice and Seq2 chains.
- `README.md` documents public usage examples.

Keep new source files at the repository root unless a real package split becomes necessary. Public API changes should be reflected in `README.md`.

## Build, Test, and Development Commands

- `go test ./...` runs all tests in the module.
- `go test -run TestSliceChainHelpers` runs one targeted test.
- `go test -cover ./...` reports package coverage.
- `gofmt -w *.go` formats Go files in the current flat layout.

The module currently targets Go `1.27` with toolchain `go1.27rc1`, as declared in `go.mod`.

## Coding Style & Naming Conventions

Use standard Go formatting via `gofmt`; tabs are expected for indentation in Go files. Keep package name `ko`. Prefer small, direct generic helpers over new abstractions. Ordered collection chains should use `Seq[T]`, a defined type over `iter.Seq[T]`; key/value collection chains should use `Seq2[K, V]`, a defined type over `iter.Seq2[K, V]`. Public collection operations should be reached through exported constructors and methods such as `Slice`, `Map`, `Collect`, and `Filter`.

Name chain methods as short verbs that match collection operations, such as `Map`, `Filter`, `Reject`, `Reduce`, `Take`, and `Drop`. Internal helpers should stay unexported and lower camel case, for example `mapSlice` and `flatMap`.

Collection operations should be exposed as chain methods, not package-level helper functions. Keep package-level exported functions limited to chain constructors such as `Slice` and `Map`; do not add top-level helpers like `Uniq(collection)` when the operation belongs on `Slice(...).Uniq()` or `Map(...).Keys()`.

Avoid intermediate chain types and conversion-only helpers such as `chunkChain`, grouped Seq2 types, `sliceFromSeq`, or `mapFromSeq`. Return `Seq[T]`, `Seq2[K, V]`, or raw `iter.Seq`/`iter.Seq2` directly when Go 1.27 generic-method instantiation cycles require it, as with `Chunk` and grouped sequence results.

## Testing Guidelines

Use Go's standard `testing` package. Place tests in `*_test.go` files in package `ko`, and name tests `Test<FeatureOrMethod>`, as in `TestSliceChainToMap`.

For map-to-slice behavior, sort results before comparison because Go map iteration order is not stable. Use focused table tests only when they reduce repetition; otherwise keep tests simple and readable.

Keep statement coverage at `100%` for `go test -cover ./...`. When adding chain methods, cover both hit and miss/empty boundary behavior where it affects branch coverage.

## Commit & Pull Request Guidelines

This repository has no commit history yet, so use concise imperative commit messages such as `Add slice reject helper` or `Fix map value mapping`.

Pull requests should include a short description, note any public API changes, and mention the test command run. Include README updates when examples or supported methods change.
