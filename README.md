# ko

`ko` is a small lodash-style helper package for Go collection chains.

It focuses on direct, type-safe helpers for slices and maps. There are no runtime dependencies outside the Go standard library.

## Install

```sh
go get github.com/go-kod/ko
```

`ko` currently targets Go `1.27`, matching the module's `go.mod`.

## Seq Chains

Start with `ko.Slice`, chain operations on `ko.Seq[T]`, then call `Collect` for the final slice. `Seq[T]` is iterator-backed and can be ranged directly; lazy operations stay lazy until range, `Collect`, or a terminal/materializing method consumes them.

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/go-kod/ko"
)

func main() {
	got := ko.Slice([]int{1, 2, 3, 4}).
		Filter(func(item int, _ int) bool {
			return item%2 == 0
		}).
		Map(func(item int, _ int) string {
			return strconv.Itoa(item * 10)
		}).
		Collect()

	fmt.Println(got) // [20 40]
}
```

Seq methods:

| Method | Returns | Notes |
| --- | --- | --- |
| `Collect()` | `[]T` | Returns the current slice. |
| `Filter(predicate)` | `Seq[T]` | Keeps matching items. |
| `Reject(predicate)` | `Seq[T]` | Drops matching items. |
| `FilterReject(predicate)` | `(Seq[T], Seq[T])` | Returns matching and non-matching items. |
| `Uniq()` | `Seq[T]` | Removes duplicate comparable items, keeping the first occurrence. Use `UniqBy` for non-comparable items. |
| `UniqBy(mapper)` | `Seq[T]` | Removes duplicate items by a comparable key, keeping the first occurrence. |
| `UniqMap(mapper)` | `Seq[R]` | Maps items and removes duplicate mapped values, keeping the first occurrence. |
| `FindUniquesBy(mapper)` | `Seq[T]` | Keeps items whose mapped key appears exactly once. |
| `FindDuplicatesBy(mapper)` | `Seq[T]` | Keeps the first item for each duplicated mapped key. |
| `Map(mapper)` | `Seq[R]` | Changes item type. |
| `FilterMap(mapper)` | `Seq[R]` | Maps items and keeps accepted results. |
| `FlatMap(mapper)` | `Seq[R]` | Maps each item to zero or more items. |
| `PartitionBy(mapper)` | `iter.Seq[ko.Seq[T]]` | Groups items by key, preserving first key order. |
| `ForEach(iteratee)` | `Seq[T]` | Calls the iteratee and keeps the chain unchanged. |
| `ForEachWhile(predicate)` | `Seq[T]` | Calls the predicate until it returns false and keeps the chain unchanged. |
| `Take(n)` | `Seq[T]` | Keeps the first `n` items. |
| `TakeRight(n)` | `Seq[T]` | Keeps the last `n` items. |
| `Subset(offset, length)` | `Seq[T]` | Keeps up to `length` items from `offset`. Negative offsets count from the end. |
| `TakeWhile(predicate)` | `Seq[T]` | Keeps the leading items while the predicate is true. |
| `TakeFilter(n, predicate)` | `Seq[T]` | Keeps the first `n` matching items. |
| `TakeRightWhile(predicate)` | `Seq[T]` | Keeps the trailing items while the predicate is true. |
| `Drop(n)` | `Seq[T]` | Drops the first `n` items. |
| `DropRight(n)` | `Seq[T]` | Drops the last `n` items. |
| `DropByIndex(indexes...)` | `Seq[T]` | Drops items by index. Negative indexes count from the end. |
| `DropWhile(predicate)` | `Seq[T]` | Drops the leading items while the predicate is true. |
| `DropRightWhile(predicate)` | `Seq[T]` | Drops the trailing items while the predicate is true. |
| `Reverse()` | `Seq[T]` | Returns a reversed copy. |
| `Reduce(accumulator, initial)` | `R` | Folds the slice into one value. |
| `ReduceRight(accumulator, initial)` | `R` | Folds the slice from right to left. |
| `Find(predicate)` | `(T, bool)` | Returns the first match and whether it was found. |
| `FindOrElse(fallback, predicate)` | `T` | Returns the first match or the fallback. |
| `FindIndex(predicate)` | `(T, int, bool)` | Returns the first match, its index, and whether it was found. |
| `FindLast(predicate)` | `(T, bool)` | Returns the last match and whether it was found. |
| `FindLastIndex(predicate)` | `(T, int, bool)` | Returns the last match, its index, and whether it was found. |
| `First()` | `(T, bool)` | Returns the first item and whether it exists. |
| `FirstOr(fallback)` | `T` | Returns the first item or the fallback. |
| `Last()` | `(T, bool)` | Returns the last item and whether it exists. |
| `LastOr(fallback)` | `T` | Returns the last item or the fallback. |
| `Nth(index)` | `(T, bool)` | Returns the item at index and whether it exists. Negative indexes count from the end. |
| `NthOr(index, fallback)` | `T` | Returns the item at index or the fallback. Negative indexes count from the end. |
| `Some(predicate)` | `bool` | Reports whether any item matches. |
| `None(predicate)` | `bool` | Reports whether no item matches. |
| `Count(predicate)` | `int` | Counts matching items. |
| `Every(predicate)` | `bool` | Reports whether all items match. |
| `ContainsBy(predicate)` | `bool` | Reports whether any item matches. |
| `WithoutBy(mapper, exclude...)` | `Seq[T]` | Drops items whose mapped key is excluded. |
| `Chunk(n)` | `iter.Seq[ko.Seq[T]]` | Splits into chunks. `n <= 0` yields nothing. |
| `Window(n)` | `iter.Seq[ko.Seq[T]]` | Splits into overlapping windows. `n <= 0` yields nothing. |

```go
var chunks [][]int
for chunk := range ko.Slice([]int{1, 2, 3, 4, 5}).Chunk(2) {
	chunks = append(chunks, chunk.Collect())
}

// [][]int{{1, 2}, {3, 4}, {5}}
```

`Chunk` yields inner `ko.Seq[T]` values. To chain the outer chunk stream, convert it explicitly to `ko.Seq[ko.Seq[T]]`.

```go
sums := ko.Seq[ko.Seq[int]](ko.Slice([]int{1, 2, 3, 4, 5}).Chunk(2)).
	Filter(func(chunk ko.Seq[int], _ int) bool {
		return len(chunk.Collect()) == 2
	}).
	Map(func(chunk ko.Seq[int], _ int) int {
		items := chunk.Collect()
		return items[0] + items[1]
	}).
	Collect()

// []int{3, 7}
```

`Window` has the same outer return shape, but each result overlaps the previous one.

```go
var windows [][]int
for window := range ko.Slice([]int{1, 2, 3, 4}).Window(3) {
	windows = append(windows, window.Collect())
}

// [][]int{{1, 2, 3}, {2, 3, 4}}
```

## Slice to Map

Use `ToMap`, `KeyBy`, or `CountBy` when a slice should become `Seq2`. `GroupBy` and `GroupByMap` return grouped `iter.Seq2` values whose grouped values are `ko.Seq`.

```go
byLength := map[int][]string{}
for key, group := range ko.Slice([]string{"go", "ko", "kod"}).
	GroupBy(func(item string, _ int) int {
		return len(item)
	}) {
	byLength[key] = group.Collect()
}

// map[int][]string{2: {"go", "ko"}, 3: {"kod"}}
```

| Method | Returns | Notes |
| --- | --- | --- |
| `ToMap(mapper)` | `Seq2[K, V]` | Mapper returns a key and value. Later duplicate keys replace earlier ones. |
| `GroupBy(mapper)` | `iter.Seq2[K, ko.Seq[T]]` | Groups items by key. |
| `GroupByMap(mapper)` | `iter.Seq2[K, ko.Seq[V]]` | Groups mapped values by key. |
| `KeyBy(mapper)` | `Seq2[K, T]` | Indexes items by key. Later duplicate keys replace earlier ones. |
| `CountBy(mapper)` | `Seq2[K, int]` | Counts items by key. |

## Seq2 Chains

Start with `ko.Map`, chain operations, then call `Collect` for the final map.

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

Seq2 methods:

| Method | Returns | Notes |
| --- | --- | --- |
| `Collect()` | `map[K]V` | Returns the current map. |
| `HasKey(key)` | `bool` | Reports whether the key exists. |
| `ValueOr(key, fallback)` | `V` | Returns the key's value or the fallback when absent. |
| `PickKeys(keys...)` | `Seq2[K, V]` | Keeps entries whose key is listed. |
| `OmitKeys(keys...)` | `Seq2[K, V]` | Drops entries whose key is listed. |
| `Assign(maps...)` | `Seq2[K, V]` | Merges maps. Later maps replace earlier keys. |
| `Filter(predicate)` | `Seq2[K, V]` | Keeps matching entries. |
| `Reject(predicate)` | `Seq2[K, V]` | Drops matching entries. |
| `Map(mapper)` | `Seq2[RK, RV]` | Changes key and value types. |
| `MapKeys(mapper)` | `Seq2[RK, V]` | Changes keys and keeps values. |
| `MapValues(mapper)` | `Seq2[K, RV]` | Changes values and keeps keys. |
| `ForEach(iteratee)` | `Seq2[K, V]` | Calls the iteratee and keeps the chain unchanged. |
| `ChunkEntries(n)` | `iter.Seq[ko.Seq2[K, V]]` | Splits entries into chunks. Map iteration order is not stable. `n <= 0` yields nothing. |
| `Keys()` | `Seq[K]` | Returns map keys. Map iteration order is not stable. |
| `FilterKeys(predicate)` | `Seq[K]` | Returns matching map keys. Map iteration order is not stable. |
| `Values()` | `Seq[V]` | Returns map values. Map iteration order is not stable. |
| `FilterValues(predicate)` | `Seq[V]` | Returns matching map values. Map iteration order is not stable. |
| `Find(predicate)` | `(K, V, bool)` | Returns the first matching entry and whether it was found. Map iteration order is not stable. |
| `Some(predicate)` | `bool` | Reports whether any entry matches. |
| `None(predicate)` | `bool` | Reports whether no entry matches. |
| `Count(predicate)` | `int` | Counts matching entries. |
| `Every(predicate)` | `bool` | Reports whether all entries match. |
| `ToSlice(mapper)` | `Seq[R]` | Converts entries to `Seq[R]`. Map iteration order is not stable. |
| `FilterMapToSlice(mapper)` | `Seq[R]` | Converts matching entries to `Seq[R]`. Map iteration order is not stable. |

```go
keys := ko.Map(map[string]int{"b": 2, "a": 1}).
	Keys().
	Collect()
sort.Strings(keys)

// []string{"a", "b"}
```

## Behavior Notes

- Predicates and mappers receive the item plus its index for slices.
- `Slice` accepts any element type. `Uniq` requires values that can be used as Go map keys; for non-comparable items, use `UniqBy`.
- Map predicates and mappers receive key and value.
- Map iteration order follows Go map iteration order, so sort `ToSlice` output in tests when order matters.
- Chain methods that transform collections return iterable chains; `Collect` materializes the current collection.
- This package intentionally stays small. Prefer Go's standard library when it already covers the job.

## Development

```sh
go test ./...
go test -cover ./...
/home/pilot/sdk/go1.27rc1/bin/gofmt -w *.go
```

Format Go files with the Go toolchain declared in `go.mod`; older `gofmt` binaries do not parse Go 1.27 generic methods.
