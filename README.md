# ko

`ko` is a small lodash-style helper package for Go collection chains.

It focuses on direct, type-safe helpers for slices and maps. There are no runtime dependencies outside the Go standard library.

## Install

```sh
go get github.com/go-kod/ko
```

`ko` currently targets Go `1.27`, matching the module's `go.mod`.

## Seq Chains

Start with `ko.Slice`, `ko.Of`, `ko.Range`, `ko.RangeStep`, or `ko.FromChannel`, chain operations on `ko.Seq[T]`, then call `Collect` for the final slice. `Seq[T]` is iterator-backed and can be ranged directly or converted to `iter.Seq[T]` for standard helpers such as `slices.Collect`.

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

Seq methods:

| Method | Returns | Notes |
| --- | --- | --- |
| `Collect()` | `[]T` | Materializes the sequence into a slice. |
| `Filter(predicate)` | `Seq[T]` | Keeps matching items. |
| `Reject(predicate)` | `Seq[T]` | Drops matching items. |
| `DistinctBy(mapper)` | `Seq[T]` | Removes duplicate items by a comparable key, keeping the first occurrence. |
| `Map(mapper)` | `Seq[R]` | Changes item type. |
| `FilterMap(mapper)` | `Seq[R]` | Maps items and keeps accepted results. |
| `FlatMap(mapper)` | `Seq[R]` | Maps each item to an `iter.Seq[R]` and flattens one level. |
| `Concat(other)` | `Seq[T]` | Appends an `iter.Seq[T]` after the current sequence. |
| `Take(n)` | `Seq[T]` | Keeps the first `n` items. |
| `TakeRight(n)` | `Seq[T]` | Keeps the last `n` items. |
| `TakeWhile(predicate)` | `Seq[T]` | Keeps the leading items while the predicate is true. |
| `Drop(n)` | `Seq[T]` | Drops the first `n` items. |
| `DropRight(n)` | `Seq[T]` | Drops the last `n` items. |
| `DropWhile(predicate)` | `Seq[T]` | Drops the leading items while the predicate is true. |
| `Sort(less)` | `Seq[T]` | Returns a stably sorted copy using a less predicate. |
| `SortBy(mapper)` | `Seq[T]` | Returns a stably sorted copy using mapped ordered keys. |
| `Reverse()` | `Seq[T]` | Returns a reversed copy. |
| `Reduce(accumulator, initial)` | `R` | Folds the slice into one value. |
| `MaxBy(mapper)` | `(T, bool)` | Returns the item with the greatest mapped ordered key. |
| `MinBy(mapper)` | `(T, bool)` | Returns the item with the least mapped ordered key. |
| `SumBy(mapper)` | `N` | Sums mapped numeric values. |
| `MeanBy(mapper)` | `float64` | Returns the arithmetic mean of mapped numeric values, or `0` for an empty sequence. |
| `Join(sep, mapper)` | `string` | Maps items to strings and joins them with `sep`. |
| `Scan(accumulator, initial)` | `Seq[R]` | Returns running accumulations, including the initial value. |
| `Find(predicate)` | `(T, bool)` | Returns the first match and whether it was found. |
| `FindIndex(predicate)` | `(int, bool)` | Returns the first matching index and whether it was found. |
| `First()` | `(T, bool)` | Returns the first item and whether it exists. |
| `Last()` | `(T, bool)` | Returns the last item and whether it exists. |
| `Nth(index)` | `(T, bool)` | Returns the item at index and whether it exists. Negative indexes count from the end. |
| `IsEmpty()` | `bool` | Reports whether the sequence yields no items. |
| `Some(predicate)` | `bool` | Reports whether any item matches. |
| `Count(predicate)` | `int` | Counts matching items. |
| `Every(predicate)` | `bool` | Reports whether all items match. |
| `Chunk(n)` | unexported adapter | Splits into chunks. `n <= 0` yields nothing. |
| `Window(n, step)` | unexported adapter | Splits into overlapping windows. `n <= 0` or `step <= 0` yields nothing. |

`Chunk` and `Window` return unexported adapter types. They can be ranged directly, collected to `[][]T`, or mapped as `[]T` values without naming the adapter type.

```go
sums := ko.Slice([]int{1, 2, 3, 4}).Chunk(2).
	Map(func(chunk []int, _ int) int {
		return chunk[0] + chunk[1]
	}).
	Collect()

// []int{3, 7}
```

`Window` advances by `step`; use `step == 1` for sliding windows that overlap by one item.

## Seq to Seq2

Use `ToMap` or `Enumerate` when a `Seq` should become `Seq2`. `GroupBy` returns an unexported grouped adapter whose grouped values are slices; grouped keys are yielded in first-seen order.

```go
byLength := ko.Slice([]string{"go", "ko", "kod"}).
	GroupBy(func(item string, _ int) int {
		return len(item)
	}).
	Collect()

// map[int][]string{2: {"go", "ko"}, 3: {"kod"}}
```

| Method | Returns | Notes |
| --- | --- | --- |
| `ToMap(mapper)` | `Seq2[K, V]` | Mapper returns a key and value. When collected to a map, later duplicate keys replace earlier ones. |
| `GroupBy(mapper)` | unexported adapter | Groups items by key, preserving first key order. The adapter supports range, `Map`, and `Collect`. |
| `Enumerate()` | `Seq2[int, T]` | Indexes items by zero-based position. |

## Seq2 Chains

Start with `ko.Map`, chain operations, then call `Collect` for the final map. `Seq2[K, V]` can also be converted to `iter.Seq2[K, V]` for standard helpers such as `maps.Collect`.

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
| `Collect()` | `map[K]V` | Materializes entries into a map. |
| `PickKeys(keys...)` | `Seq2[K, V]` | Keeps entries whose key is listed. |
| `OmitKeys(keys...)` | `Seq2[K, V]` | Drops entries whose key is listed. |
| `Filter(predicate)` | `Seq2[K, V]` | Keeps matching entries. |
| `Reject(predicate)` | `Seq2[K, V]` | Drops matching entries. |
| `Map(mapper)` | `Seq2[RK, RV]` | Changes key and value types. |
| `MapKeys(mapper)` | `Seq2[RK, V]` | Changes keys and keeps values. |
| `MapValues(mapper)` | `Seq2[K, RV]` | Changes values and keeps keys. |
| `Keys()` | `Seq[K]` | Returns map keys. Map iteration order is not stable. |
| `Values()` | `Seq[V]` | Returns map values. Map iteration order is not stable. |
| `Find(predicate)` | `(K, V, bool)` | Returns the first matching entry and whether it was found. Map iteration order is not stable. |
| `IsEmpty()` | `bool` | Reports whether the sequence yields no entries. |
| `Some(predicate)` | `bool` | Reports whether any entry matches. |
| `Count(predicate)` | `int` | Counts matching entries. |
| `Every(predicate)` | `bool` | Reports whether all entries match. |
| `ToSlice(mapper)` | `Seq[R]` | Converts entries to `Seq[R]`. Map iteration order is not stable. |

## Behavior Notes

- Predicates and mappers receive the item plus its index for slices.
- `Slice` accepts any element type. Use `DistinctBy` to remove duplicates by a comparable key.
- `Seq2` is an entry stream; duplicate key replacement happens when `Collect` materializes a Go map.
- Map predicates and mappers receive key and value.
- Map iteration order follows Go map iteration order, so sort `ToSlice` output in tests when order matters.
- Chain methods that transform collections return iterable chains; `Collect` materializes the current collection.
- `FromChannel` is one-shot because it drains the source channel as it is consumed.
- Methods that need suffix, sort, reverse, or whole-sequence knowledge buffer only when consumed, such as `Sort`, `SortBy`, `Reverse`, negative `Nth`, `TakeRight`, `DropRight`, and grouping.
- This package intentionally stays small. Prefer Go's standard library when it already covers the job.

## Development

```sh
go test ./...
go test -cover ./...
/home/pilot/sdk/go1.27rc1/bin/gofmt -w *.go
```

Format Go files with the Go toolchain declared in `go.mod`; older `gofmt` binaries do not parse Go 1.27 generic methods.
