package ko_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"sort"
	"strconv"

	"github.com/go-kod/ko"
)

func ExampleOf() {
	got := ko.Uniq(ko.Of(1, 2, 2, 3, 4)).
		Filter(func(item int, _ int) bool {
			return item%2 == 0
		}).
		Map(func(item int, _ int) string {
			return strconv.Itoa(item * 10)
		}).
		Collect()

	fmt.Println(got)
	// Output: [20 40]
}

func ExampleSlice_uniqBy() {
	got := ko.Slice([]string{"go", "ko", "kod"}).
		UniqBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	fmt.Println(got)
	// Output: [go kod]
}

func ExampleSeq_standardLibraryInterop() {
	items := slices.Collect(iter.Seq[int](ko.Range(1, 4)))
	entries := maps.Collect(iter.Seq2[string, int](ko.Map(map[string]int{"a": 1})))

	fmt.Println(items)
	fmt.Println(entries["a"])
	// Output:
	// [1 2 3]
	// 1
}

func ExampleSlice_chunk() {
	var got [][]int
	for chunk := range ko.Slice([]int{1, 2, 3, 4, 5}).Chunk(2) {
		got = append(got, chunk.Collect())
	}

	fmt.Println(got)
	// Output: [[1 2] [3 4] [5]]
}

func ExampleSlice_window() {
	var got [][]int
	for window := range ko.Slice([]int{1, 2, 3, 4}).Window(3, 1) {
		got = append(got, window.Collect())
	}

	fmt.Println(got)
	// Output: [[1 2 3] [2 3 4]]
}

func ExampleSlice_takeRight() {
	got := ko.Slice([]int{1, 2, 3, 4}).
		TakeRight(2).
		DropRight(1).
		Collect()

	fmt.Println(got)
	// Output: [3]
}

func ExampleSlice_mapChunks() {
	got := ko.Seq[ko.Seq[int]](ko.Slice([]int{1, 2, 3, 4, 5}).Chunk(2)).
		Filter(func(chunk ko.Seq[int], _ int) bool {
			return len(chunk.Collect()) == 2
		}).
		Map(func(chunk ko.Seq[int], _ int) int {
			items := chunk.Collect()
			return items[0] + items[1]
		}).
		Collect()

	fmt.Println(got)
	// Output: [3 7]
}

func ExampleMap() {
	got := ko.Map(map[string]int{"a": 1, "b": 2}).
		MapValues(func(_ string, value int) string {
			return strconv.Itoa(value)
		}).
		ToSlice(func(key string, value string) string {
			return key + value
		}).
		Collect()
	sort.Strings(got)

	fmt.Println(got)
	// Output: [a1 b2]
}

func ExampleSeq2_Keys() {
	got := ko.Map(map[string]int{"b": 2, "a": 1}).
		Keys().
		Collect()
	sort.Strings(got)

	fmt.Println(got)
	// Output: [a b]
}

func ExampleSeq2_Find() {
	key, value, ok := ko.Map(map[string]int{"a": 1, "b": 2}).
		Find(func(_ string, value int) bool {
			return value == 2
		})

	fmt.Println(key, value, ok)
	// Output: b 2 true
}
