# Reddit Post Draft

## Subreddit

- `r/golang`

## Flair

- `Show and Tell`
- `Package`
- `Feedback`

## Recommended Title

I built `ko`, a small `iter.Seq` collection helper package for Go

## Alternative Titles

- Looking for API feedback on a small `iter.Seq` helper package
- `ko`: typed `Seq` and `Seq2` chains backed by Go iterators
- Is this `Seq` / `Seq2` API idiomatic, or too much abstraction for Go?
- A tiny lodash-style helper package for Go collections, built on `iter.Seq`

## Post Body

Hi r/golang,

I have been working on a small Go package called `ko`:

https://github.com/go-kod/ko

It is a typed collection helper package for common slice and map pipelines. The current design is built around Go iterators:

- `Seq[T]` is based on `iter.Seq[T]`
- `Seq2[K, V]` is based on `iter.Seq2[K, V]`
- intermediate operations stay sequence-based where possible
- `Collect()` is the explicit materialization boundary
- there are no runtime dependencies outside the standard library

The package is not meant to replace normal loops. I still reach for loops when they are clearer. The goal is to make common transformation pipelines compact without repeatedly converting between slices and maps in the middle of the chain.

Example:

```go
got := ko.Of(1, 2, 3, 4).
	Filter(func(item int, _ int) bool {
		return item%2 == 0
	}).
	Map(func(item int, _ int) string {
		return strconv.Itoa(item * 10)
	}).
	Collect()

// []string{"20", "40"}
```

Map-style pipelines use `Seq2`:

```go
got := ko.Map(map[string]int{"a": 1, "bb": 2}).
	Filter(func(_ string, value int) bool {
		return value > 1
	}).
	Map(func(key string, value int) (int, string) {
		return len(key), strconv.Itoa(value)
	}).
	Collect()

// map[int]string{2: "2"}
```

Seq-to-Seq2 operations also stay lazy until `Collect()`:

```go
byLength := ko.Slice([]string{"go", "ko", "kod"}).
	ToMap(func(item string, _ int) (int, string) {
		return len(item), item
	}).
	Collect()

// map[int]string{2: "ko", 3: "kod"}
```

In that example, duplicate keys are resolved only when the `Seq2` is collected into a Go map, so the last value for key `2` wins.

There are also small sequence constructors:

```go
got := ko.Range(1, 6).
	Filter(func(item int, _ int) bool {
		return item%2 == 1
	}).
	Collect()

// []int{1, 3, 5}
```

For infinite generated sequences, the caller uses normal chain methods to bound consumption:

```go
got := ko.Generate(1, func(item int) int {
	return item * 2
}).Take(4).Collect()

// []int{1, 2, 4, 8}
```

Some of the API surface includes `Filter`, `Reject`, `Map`, `FilterMap`, `FlatMap`, `UniqBy`, `Chunk`, `Window`, `GroupBy`, `CountBy`, `Keys`, `Values`, and `ToSlice`.

I would appreciate critical feedback on the API shape:

- Does `Seq[T]` / `Seq2[K, V]` feel natural for iterator-backed Go code?
- Is `Collect()` the right name for the materialization boundary?
- Are methods like `ToMap`, `GroupBy`, `Chunk`, and `Window` worth having in a small package?
- Are any return types surprising?
- Does this feel useful, or would you keep helpers like this local to each project?
- Are there methods that should be removed to keep the package more idiomatic?

I am especially interested in feedback from people already using Go iterators in real code.

## Short Version

I built `ko`, a small typed collection helper package for Go:

https://github.com/go-kod/ko

It provides `Seq[T]` and `Seq2[K, V]` chains backed by `iter.Seq` / `iter.Seq2`, keeps intermediate operations sequence-based where possible, and uses `Collect()` as the explicit materialization point.

I am looking for API feedback from Go developers, especially around whether this style feels idiomatic or like unnecessary abstraction.

## First Comment Draft

A bit more context: this started as a small lodash-style helper package, but I refactored it to use `iter.Seq` more thoroughly. The main design constraint is that intermediate operations should remain sequence-based instead of collecting into slices or maps between steps.

The implementation is intentionally small and dependency-free. I am trying to find the line between useful collection helpers and an API that feels too clever for Go.
