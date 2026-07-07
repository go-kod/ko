package ko

import (
	"reflect"
	"strconv"
	"testing"
)

func TestSliceChainFilterMapChangesType(t *testing.T) {
	got := Slice([]int{1, 2, 3, 4}).
		Filter(func(item int, _ int) bool {
			return item%2 == 0
		}).
		Map(func(item int, _ int) string {
			return strconv.Itoa(item * 10)
		}).
		Collect()

	want := []string{"20", "40"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainFilterMap(t *testing.T) {
	got := Slice([]int{1, 2, 3, 4}).
		FilterMap(func(item int, _ int) (string, bool) {
			return strconv.Itoa(item * 10), item%2 == 0
		}).
		Collect()

	want := []string{"20", "40"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainIteration(t *testing.T) {
	var got []int
	for item := range Slice([]int{1, 2, 3}) {
		got = append(got, item)
	}

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainFilterMapIsLazy(t *testing.T) {
	calls := 0
	chain := Slice([]int{1, 2, 3}).
		Filter(func(item int, _ int) bool {
			calls++
			return item%2 == 1
		}).
		Map(func(item int, _ int) int {
			calls++
			return item * 10
		})

	if calls != 0 {
		t.Fatalf("called before consumption: %d", calls)
	}

	got := chain.Collect()
	want := []int{10, 30}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	if calls != 5 {
		t.Fatalf("calls: %d", calls)
	}
}

func TestSliceChainLazyIterationStopsEarly(t *testing.T) {
	filterCalls := 0
	for item := range Slice([]int{1, 2, 3}).
		Filter(func(item int, _ int) bool {
			filterCalls++
			return item%2 == 1
		}) {
		if item != 1 {
			t.Fatalf("filter item: %d", item)
		}
		break
	}
	if filterCalls != 1 {
		t.Fatalf("filter calls: %d", filterCalls)
	}

	mapCalls := 0
	for item := range Slice([]int{1, 2, 3}).
		Map(func(item int, _ int) int {
			mapCalls++
			return item * 10
		}) {
		if item != 10 {
			t.Fatalf("map item: %d", item)
		}
		break
	}
	if mapCalls != 1 {
		t.Fatalf("map calls: %d", mapCalls)
	}

	filterMapCalls := 0
	for item := range Slice([]int{1, 2, 3}).
		FilterMap(func(item int, _ int) (int, bool) {
			filterMapCalls++
			return item * 10, item%2 == 1
		}) {
		if item != 10 {
			t.Fatalf("filterMap item: %d", item)
		}
		break
	}
	if filterMapCalls != 1 {
		t.Fatalf("filterMap calls: %d", filterMapCalls)
	}

	flatMapCalls := 0
	for item := range Slice([]int{1, 2, 3}).
		FlatMap(func(item int, _ int) []int {
			flatMapCalls++
			return []int{item, item * 10}
		}) {
		if item != 1 {
			t.Fatalf("flatMap item: %d", item)
		}
		break
	}
	if flatMapCalls != 1 {
		t.Fatalf("flatMap calls: %d", flatMapCalls)
	}
}

func TestSliceChainLazyTerminalsStopEarly(t *testing.T) {
	someCalls := 0
	if ok := Slice([]int{1, 2, 3}).
		Map(func(item int, _ int) int {
			someCalls++
			return item
		}).
		Some(func(item int, _ int) bool {
			return item == 1
		}); !ok {
		t.Fatal("some: false")
	}
	if someCalls != 1 {
		t.Fatalf("some calls: %d", someCalls)
	}

	findCalls := 0
	item, ok := Slice([]int{1, 2, 3}).
		Map(func(item int, _ int) int {
			findCalls++
			return item
		}).
		Find(func(item int, _ int) bool {
			return item == 1
		})
	if !ok || item != 1 {
		t.Fatalf("find: %d, %v", item, ok)
	}
	if findCalls != 1 {
		t.Fatalf("find calls: %d", findCalls)
	}

	firstCalls := 0
	item, ok = Slice([]int{1, 2, 3}).
		Map(func(item int, _ int) int {
			firstCalls++
			return item
		}).
		First()
	if !ok || item != 1 {
		t.Fatalf("first: %d, %v", item, ok)
	}
	if firstCalls != 1 {
		t.Fatalf("first calls: %d", firstCalls)
	}
}

func TestSliceChainLazyMiddleOperationsStopEarly(t *testing.T) {
	for item := range Slice([]int{1, 1, 2}).Uniq() {
		if item != 1 {
			t.Fatalf("uniq item: %d", item)
		}
		break
	}

	for item := range Slice([]string{"go", "ko", "kod"}).
		UniqBy(func(item string, _ int) int {
			return len(item)
		}) {
		if item != "go" {
			t.Fatalf("uniqBy item: %q", item)
		}
		break
	}

	for item := range Slice([]string{"go", "ko", "kod"}).
		UniqMap(func(item string, _ int) int {
			return len(item)
		}) {
		if item != 2 {
			t.Fatalf("uniqMap item: %d", item)
		}
		break
	}

	for item := range Slice([]int{1, 2, 3}).Take(2) {
		if item != 1 {
			t.Fatalf("take item: %d", item)
		}
		break
	}

	for item := range Slice([]int{1, 2, 3, 4}).
		TakeFilter(2, func(item int, _ int) bool {
			return item%2 == 0
		}) {
		if item != 2 {
			t.Fatalf("takeFilter item: %d", item)
		}
		break
	}

	for item := range Slice([]int{1, 2, 3}).Drop(1) {
		if item != 2 {
			t.Fatalf("drop item: %d", item)
		}
		break
	}

	for item := range Slice([]int{1, 2, 3}).
		DropWhile(func(item int, _ int) bool {
			return item < 2
		}) {
		if item != 2 {
			t.Fatalf("dropWhile item: %d", item)
		}
		break
	}

	for item := range Slice([]int{1, 2, 3}).Reverse() {
		if item != 3 {
			t.Fatalf("reverse item: %d", item)
		}
		break
	}
}

func TestSliceChainHelpers(t *testing.T) {
	if got := Slice([]int{1, 2, 3}).Reverse().Collect(); !reflect.DeepEqual(got, []int{3, 2, 1}) {
		t.Fatalf("reverse: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Reduce(func(sum int, item int, _ int) int {
		return sum + item
	}, 0); got != 6 {
		t.Fatalf("reduce: %d", got)
	}

	if got := Slice([]string{"a", "b", "c"}).ReduceRight(func(out string, item string, index int) string {
		return out + item + strconv.Itoa(index)
	}, ""); got != "c2b1a0" {
		t.Fatalf("reduceRight: %q", got)
	}
}

func TestSliceChainTakeWhile(t *testing.T) {
	got := Slice([]int{1, 2, 3, 2}).
		TakeWhile(func(item int, _ int) bool {
			return item < 3
		}).
		Collect()

	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainTakeFilter(t *testing.T) {
	got := Slice([]int{1, 2, 3, 4, 5, 6}).
		TakeFilter(2, func(item int, _ int) bool {
			return item%2 == 0
		}).
		Collect()

	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Slice([]int{1, 2, 3}).TakeFilter(0, func(item int, _ int) bool {
		return true
	}).Collect()

	if !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("zero: %#v", got)
	}
}

func TestSliceChainSubset(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod", "go-kod"}).Subset(1, 2).Collect()
	if !reflect.DeepEqual(got, []string{"ko", "kod"}) {
		t.Fatalf("subset: %#v", got)
	}

	got = Slice([]string{"go", "ko", "kod", "go-kod"}).Subset(-2, 9).Collect()
	if !reflect.DeepEqual(got, []string{"kod", "go-kod"}) {
		t.Fatalf("subset negative: %#v", got)
	}

	got = Slice([]string{"go", "ko"}).Subset(-9, 1).Collect()
	if !reflect.DeepEqual(got, []string{"go"}) {
		t.Fatalf("subset negative overflow: %#v", got)
	}

	got = Slice([]string{"go", "ko"}).Subset(9, 1).Collect()
	if !reflect.DeepEqual(got, []string{}) {
		t.Fatalf("subset overflow: %#v", got)
	}

	got = Slice([]string{"go", "ko"}).Subset(0, 0).Collect()
	if !reflect.DeepEqual(got, []string{}) {
		t.Fatalf("subset zero: %#v", got)
	}
}

func TestSliceChainDropWhile(t *testing.T) {
	got := Slice([]int{1, 2, 3, 2}).
		DropWhile(func(item int, _ int) bool {
			return item < 3
		}).
		Collect()

	want := []int{3, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainTakeRightWhile(t *testing.T) {
	got := Slice([]int{1, 2, 3, 2}).
		TakeRightWhile(func(item int, _ int) bool {
			return item < 3
		}).
		Collect()

	want := []int{2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainDropRightWhile(t *testing.T) {
	got := Slice([]int{1, 2, 3, 2}).
		DropRightWhile(func(item int, _ int) bool {
			return item < 3
		}).
		Collect()

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainDropByIndex(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod", "go-kod"}).
		DropByIndex(1, -1, 99).
		Collect()

	want := []string{"go", "kod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainUniq(t *testing.T) {
	got := Slice([]int{1, 2, 1, 3, 2}).Uniq().Collect()

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainUniqBy(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod"}).
		UniqBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want := []string{"go", "kod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainUniqBySupportsNonComparableItems(t *testing.T) {
	type item struct {
		id     int
		values []int
	}

	got := Slice([]item{
		{id: 1, values: []int{1}},
		{id: 1, values: []int{2}},
		{id: 2, values: []int{3}},
	}).UniqBy(func(item item, _ int) int {
		return item.id
	}).Collect()

	want := []item{
		{id: 1, values: []int{1}},
		{id: 2, values: []int{3}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainUniqMap(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod", "go-kod"}).
		UniqMap(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want := []int{2, 3, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainFindUniquesBy(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod", "x"}).
		FindUniquesBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want := []string{"kod", "x"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Slice([]string{"go", "ko"}).
		FindUniquesBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	if !reflect.DeepEqual(got, []string{}) {
		t.Fatalf("all duplicated: %#v", got)
	}
}

func TestSliceChainFindDuplicatesBy(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod", "go-kod", "x"}).
		FindDuplicatesBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want := []string{"go"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Slice([]string{"go", "kod"}).
		FindDuplicatesBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	if !reflect.DeepEqual(got, []string{}) {
		t.Fatalf("no duplicates: %#v", got)
	}
}

func TestSliceChainPredicates(t *testing.T) {
	if got := Slice([]int{1, 2, 3}).Some(func(item int, _ int) bool {
		return item%2 == 0
	}); !got {
		t.Fatal("some: false")
	}

	if got := Slice([]int{1, 2, 3}).Every(func(item int, _ int) bool {
		return item > 0
	}); !got {
		t.Fatal("every: false")
	}

	if got := Slice([]int{1, 2, 3}).Every(func(item int, _ int) bool {
		return item < 3
	}); got {
		t.Fatal("every miss: true")
	}

	if got := Slice([]int{}).Every(func(item int, _ int) bool {
		return item > 0
	}); !got {
		t.Fatal("every empty: false")
	}

	if got := Slice([]int{1, 2, 3}).None(func(item int, _ int) bool {
		return item > 9
	}); !got {
		t.Fatal("none: false")
	}

	if got := Slice([]int{1, 2, 3}).None(func(item int, _ int) bool {
		return item == 2
	}); got {
		t.Fatal("none hit: true")
	}

	if got := Slice([]int{1, 2, 3, 4}).Count(func(item int, _ int) bool {
		return item%2 == 0
	}); got != 2 {
		t.Fatalf("count: %d", got)
	}

	if got := Slice([]int{1, 3}).Count(func(item int, _ int) bool {
		return item%2 == 0
	}); got != 0 {
		t.Fatalf("count miss: %d", got)
	}
}

func TestSliceChainFilterReject(t *testing.T) {
	kept, rejected := Slice([]int{1, 2, 3, 4}).FilterReject(func(item int, _ int) bool {
		return item%2 == 0
	})

	if got := kept.Collect(); !reflect.DeepEqual(got, []int{2, 4}) {
		t.Fatalf("kept: %#v", got)
	}

	if got := rejected.Collect(); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("rejected: %#v", got)
	}
}

func TestSliceChainElementAccess(t *testing.T) {
	first, ok := Slice([]int{1, 2, 3}).First()
	if !ok || first != 1 {
		t.Fatalf("first: %d, %v", first, ok)
	}

	last, ok := Slice([]int{1, 2, 3}).Last()
	if !ok || last != 3 {
		t.Fatalf("last: %d, %v", last, ok)
	}

	first, ok = Slice([]int{}).First()
	if ok || first != 0 {
		t.Fatalf("first empty: %d, %v", first, ok)
	}

	last, ok = Slice([]int{}).Last()
	if ok || last != 0 {
		t.Fatalf("last empty: %d, %v", last, ok)
	}
}

func TestSliceChainFirstOr(t *testing.T) {
	if got := Slice([]int{1, 2, 3}).FirstOr(9); got != 1 {
		t.Fatalf("firstOr: %d", got)
	}

	if got := Slice([]int{}).FirstOr(9); got != 9 {
		t.Fatalf("firstOr empty: %d", got)
	}
}

func TestSliceChainLastOr(t *testing.T) {
	if got := Slice([]int{1, 2, 3}).LastOr(9); got != 3 {
		t.Fatalf("lastOr: %d", got)
	}

	if got := Slice([]int{}).LastOr(9); got != 9 {
		t.Fatalf("lastOr empty: %d", got)
	}
}

func TestSliceChainNth(t *testing.T) {
	item, ok := Slice([]string{"go", "ko", "kod"}).Nth(1)
	if !ok || item != "ko" {
		t.Fatalf("nth: %q, %v", item, ok)
	}

	item, ok = Slice([]string{"go", "ko", "kod"}).Nth(-1)
	if !ok || item != "kod" {
		t.Fatalf("nth negative: %q, %v", item, ok)
	}

	item, ok = Slice([]string{"go", "ko", "kod"}).Nth(9)
	if ok || item != "" {
		t.Fatalf("nth miss: %q, %v", item, ok)
	}

	item, ok = Slice([]string{"go", "ko", "kod"}).Nth(-4)
	if ok || item != "" {
		t.Fatalf("nth negative miss: %q, %v", item, ok)
	}
}

func TestSliceChainNthOr(t *testing.T) {
	if got := Slice([]string{"go", "ko", "kod"}).NthOr(1, "fallback"); got != "ko" {
		t.Fatalf("nthOr: %q", got)
	}

	if got := Slice([]string{"go", "ko", "kod"}).NthOr(-1, "fallback"); got != "kod" {
		t.Fatalf("nthOr negative: %q", got)
	}

	if got := Slice([]string{"go"}).NthOr(9, "fallback"); got != "fallback" {
		t.Fatalf("nthOr fallback: %q", got)
	}
}

func TestSliceChainContainsBy(t *testing.T) {
	if got := Slice([]string{"go", "ko", "kod"}).ContainsBy(func(item string, _ int) bool {
		return len(item) == 3
	}); !got {
		t.Fatal("containsBy: false")
	}

	if got := Slice([]string{"go", "ko"}).ContainsBy(func(item string, _ int) bool {
		return len(item) == 3
	}); got {
		t.Fatal("containsBy miss: true")
	}
}

func TestSliceChainFindIndex(t *testing.T) {
	item, index, ok := Slice([]string{"go", "ko", "kod"}).FindIndex(func(item string, _ int) bool {
		return len(item) == 3
	})
	if !ok || item != "kod" || index != 2 {
		t.Fatalf("findIndex: %q, %d, %v", item, index, ok)
	}

	item, index, ok = Slice([]string{"go"}).FindIndex(func(item string, _ int) bool {
		return len(item) == 3
	})
	if ok || item != "" || index != -1 {
		t.Fatalf("findIndex miss: %q, %d, %v", item, index, ok)
	}
}

func TestSliceChainFindOrElse(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod"}).FindOrElse("fallback", func(item string, _ int) bool {
		return len(item) == 3
	})
	if got != "kod" {
		t.Fatalf("got %q", got)
	}

	got = Slice([]string{"go", "ko"}).FindOrElse("fallback", func(item string, _ int) bool {
		return len(item) == 3
	})
	if got != "fallback" {
		t.Fatalf("fallback: %q", got)
	}
}

func TestSliceChainFindLast(t *testing.T) {
	item, ok := Slice([]string{"go", "kod", "ko"}).FindLast(func(item string, _ int) bool {
		return len(item) == 2
	})
	if !ok || item != "ko" {
		t.Fatalf("findLast: %q, %v", item, ok)
	}

	item, ok = Slice([]string{"go"}).FindLast(func(item string, _ int) bool {
		return len(item) == 3
	})
	if ok || item != "" {
		t.Fatalf("findLast miss: %q, %v", item, ok)
	}
}

func TestSliceChainFindLastIndex(t *testing.T) {
	item, index, ok := Slice([]string{"go", "kod", "ko"}).FindLastIndex(func(item string, _ int) bool {
		return len(item) == 2
	})
	if !ok || item != "ko" || index != 2 {
		t.Fatalf("findLastIndex: %q, %d, %v", item, index, ok)
	}

	item, index, ok = Slice([]string{"go"}).FindLastIndex(func(item string, _ int) bool {
		return len(item) == 3
	})
	if ok || item != "" || index != -1 {
		t.Fatalf("findLastIndex miss: %q, %d, %v", item, index, ok)
	}
}

func TestSliceChainWithoutBy(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod"}).
		WithoutBy(func(item string, _ int) int {
			return len(item)
		}, 2).
		Collect()

	want := []string{"kod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Slice([]string{"go", "ko"}).
		WithoutBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want = []string{"go", "ko"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainMoreHelpers(t *testing.T) {
	if got := Slice([]int{1, 2, 3}).Reject(func(item int, _ int) bool {
		return item == 2
	}).Collect(); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("reject: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).FlatMap(func(item int, _ int) []int {
		return []int{item, item * 10}
	}).Collect(); !reflect.DeepEqual(got, []int{1, 10, 2, 20, 3, 30}) {
		t.Fatalf("flatMap: %#v", got)
	}

	seen := 0
	Slice([]int{1, 2, 3}).ForEach(func(item int, _ int) {
		seen += item
	})
	if seen != 6 {
		t.Fatalf("forEach: %d", seen)
	}

	seen = 0
	got := Slice([]int{1, 2, 3}).ForEachWhile(func(item int, _ int) bool {
		seen += item
		return item < 2
	}).Collect()
	if seen != 3 {
		t.Fatalf("forEachWhile: %d", seen)
	}
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("forEachWhile value: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Take(0).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("take zero: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Take(9).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("take all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Take(1).Collect(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("take one: %#v", got)
	}

	if got := Slice([]int{}).Take(2).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("take empty: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Drop(0).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("drop zero: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Drop(9).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("drop all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).TakeRight(2).Collect(); !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("takeRight: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).TakeRight(0).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("takeRight zero: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).TakeRight(9).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("takeRight all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).DropRight(2).Collect(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("dropRight: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).DropRight(0).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("dropRight zero: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).DropRight(9).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("dropRight all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).TakeWhile(func(item int, _ int) bool {
		return item < 9
	}).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("takeWhile all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).DropWhile(func(item int, _ int) bool {
		return item < 9
	}).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("dropWhile all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).TakeRightWhile(func(item int, _ int) bool {
		return item < 9
	}).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("takeRightWhile all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).TakeRightWhile(func(item int, _ int) bool {
		return item > 9
	}).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("takeRightWhile none: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).DropRightWhile(func(item int, _ int) bool {
		return item < 9
	}).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("dropRightWhile all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).DropRightWhile(func(item int, _ int) bool {
		return item > 9
	}).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("dropRightWhile none: %#v", got)
	}

	item, ok := Slice([]int{1, 2, 3}).Find(func(item int, _ int) bool {
		return item == 9
	})
	if ok || item != 0 {
		t.Fatalf("find miss: %d, %v", item, ok)
	}
}

func TestSliceChainToMap(t *testing.T) {
	got := Slice([]string{"go", "kod"}).
		ToMap(func(item string, _ int) (int, string) {
			return len(item), item
		}).
		Collect()

	want := map[int]string{2: "go", 3: "kod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainGroupBy(t *testing.T) {
	got := map[int][]string{}
	for key, group := range Slice([]string{"go", "ko", "kod"}).
		GroupBy(func(item string, _ int) int {
			return len(item)
		}) {
		got[key] = group.Collect()
	}

	want := map[int][]string{2: {"go", "ko"}, 3: {"kod"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainGroupByMap(t *testing.T) {
	got := map[int][]string{}
	for key, group := range Slice([]string{"go", "ko", "kod"}).
		GroupByMap(func(item string, _ int) (int, string) {
			return len(item), item + "!"
		}) {
		got[key] = group.Collect()
	}

	want := map[int][]string{2: {"go!", "ko!"}, 3: {"kod!"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainPartitionBy(t *testing.T) {
	got := Slice([]string{"go", "kod", "ko", "go-kod"}).
		PartitionBy(func(item string, _ int) int {
			return len(item)
		})

	want := [][]string{{"go", "ko"}, {"kod"}, {"go-kod"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainKeyBy(t *testing.T) {
	got := Slice([]string{"g", "go", "ko"}).
		KeyBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want := map[int]string{1: "g", 2: "ko"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainCountBy(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod"}).
		CountBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want := map[int]int{2: 2, 3: 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainMapConversionsStopEarly(t *testing.T) {
	for range Slice([]string{"go", "ko", "kod"}).
		GroupBy(func(item string, _ int) int {
			return len(item)
		}) {
		break
	}

	for range Slice([]string{"go", "ko", "kod"}).
		GroupByMap(func(item string, _ int) (int, string) {
			return len(item), item + "!"
		}) {
		break
	}

	for range Slice([]string{"go", "ko", "kod"}).
		KeyBy(func(item string, _ int) int {
			return len(item)
		}) {
		break
	}

	for range Slice([]string{"go", "ko", "kod"}).
		CountBy(func(item string, _ int) int {
			return len(item)
		}) {
		break
	}
}

func TestSliceChainMapConversionsAreLazyUntilConsumed(t *testing.T) {
	calls := 0
	grouped := Slice([]string{"go", "ko", "kod"}).
		GroupBy(func(item string, _ int) int {
			calls++
			return len(item)
		})
	if calls != 0 {
		t.Fatalf("groupBy called before consumption: %d", calls)
	}
	for range grouped {
		break
	}
	if calls != 3 {
		t.Fatalf("groupBy calls: %d", calls)
	}

	calls = 0
	keyed := Slice([]string{"go", "ko", "kod"}).
		KeyBy(func(item string, _ int) int {
			calls++
			return len(item)
		})
	if calls != 0 {
		t.Fatalf("keyBy called before consumption: %d", calls)
	}
	for range keyed {
		break
	}
	if calls != 3 {
		t.Fatalf("keyBy calls: %d", calls)
	}
}

func TestSliceChainChunk(t *testing.T) {
	var got [][]int
	for chunk := range Slice([]int{1, 2, 3, 4, 5}).Chunk(2) {
		got = append(got, chunk.Collect())
	}

	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainChunkOuterCanBeChainedByConversion(t *testing.T) {
	got := Seq[Seq[int]](Slice([]int{1, 2, 3, 4, 5}).Chunk(2)).
		Filter(func(chunk Seq[int], _ int) bool {
			return len(chunk.Collect()) == 2
		}).
		Map(func(chunk Seq[int], _ int) int {
			items := chunk.Collect()
			return items[0] + items[1]
		}).
		Collect()

	want := []int{3, 7}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSliceChainChunkStopsEarly(t *testing.T) {
	calls := 0
	for chunk := range Slice([]int{1, 2, 3, 4}).Map(func(item int, _ int) int {
		calls++
		return item
	}).Chunk(2) {
		if got := chunk.Collect(); !reflect.DeepEqual(got, []int{1, 2}) {
			t.Fatalf("chunk: %#v", got)
		}
		break
	}
	if calls != 2 {
		t.Fatalf("calls: %d", calls)
	}
}

func TestSliceChainChunkZeroSizeYieldsNothing(t *testing.T) {
	for chunk := range Slice([]int{1, 2, 3}).Chunk(0) {
		t.Fatalf("unexpected chunk: %#v", chunk.Collect())
	}
}
