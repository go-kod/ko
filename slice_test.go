package ko

import (
	"iter"
	"reflect"
	"strconv"
	"testing"
)

func TestSeqFilterMapChangesType(t *testing.T) {
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

func TestSeqFilterMap(t *testing.T) {
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

func TestOf(t *testing.T) {
	if got := Of(1, 2, 3).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("of: %#v", got)
	}

	if got := Of[int]().Collect(); len(got) != 0 {
		t.Fatalf("of empty: %#v", got)
	}
}

func TestRange(t *testing.T) {
	if got := Range(1, 4).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("range ascending: %#v", got)
	}

	if got := Range(4, 1).Collect(); !reflect.DeepEqual(got, []int{4, 3, 2}) {
		t.Fatalf("range descending: %#v", got)
	}

	if got := Range(2, 2).Collect(); len(got) != 0 {
		t.Fatalf("range empty: %#v", got)
	}

	var got []int
	for item := range Range(1, 100) {
		got = append(got, item)
		if len(got) == 2 {
			break
		}
	}
	if !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("range should stop with range: %#v", got)
	}

	got = nil
	for item := range Range(4, 1) {
		got = append(got, item)
		break
	}
	if !reflect.DeepEqual(got, []int{4}) {
		t.Fatalf("descending range should stop with range: %#v", got)
	}
}

func TestRangeStep(t *testing.T) {
	if got := RangeStep(1, 8, 2).Collect(); !reflect.DeepEqual(got, []int{1, 3, 5, 7}) {
		t.Fatalf("rangeStep ascending: %#v", got)
	}

	if got := RangeStep(8, 1, 3).Collect(); !reflect.DeepEqual(got, []int{8, 5, 2}) {
		t.Fatalf("rangeStep descending: %#v", got)
	}

	if got := RangeStep(1, 8, 0).Collect(); len(got) != 0 {
		t.Fatalf("rangeStep zero step: %#v", got)
	}

	if got := RangeStep(1, 8, -1).Collect(); len(got) != 0 {
		t.Fatalf("rangeStep negative step: %#v", got)
	}

	var got []int
	for item := range RangeStep(1, 100, 3) {
		got = append(got, item)
		if len(got) == 2 {
			break
		}
	}
	if !reflect.DeepEqual(got, []int{1, 4}) {
		t.Fatalf("rangeStep should stop with range: %#v", got)
	}

	got = nil
	for item := range RangeStep(10, 1, 4) {
		got = append(got, item)
		break
	}
	if !reflect.DeepEqual(got, []int{10}) {
		t.Fatalf("descending rangeStep should stop with range: %#v", got)
	}
}

func TestFromChannel(t *testing.T) {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)

	seq := FromChannel(ch)
	if got := seq.Collect(); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("fromChannel: %#v", got)
	}
	if got := seq.Collect(); len(got) != 0 {
		t.Fatalf("fromChannel should be one-shot: %#v", got)
	}

	ch = make(chan int, 2)
	ch <- 3
	ch <- 4
	close(ch)

	var got []int
	for item := range FromChannel(ch) {
		got = append(got, item)
		break
	}
	if !reflect.DeepEqual(got, []int{3}) {
		t.Fatalf("fromChannel should stop with range: %#v", got)
	}
	if len(ch) != 1 {
		t.Fatalf("fromChannel should not drain after stop, remaining: %d", len(ch))
	}
}

func TestSeqIteration(t *testing.T) {
	var got []int
	for item := range Slice([]int{1, 2, 3}) {
		got = append(got, item)
	}

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeqFilterMapIsLazy(t *testing.T) {
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

func TestSeqLazyIterationStopsEarly(t *testing.T) {
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
		FlatMap(func(item int, _ int) iter.Seq[int] {
			flatMapCalls++
			return iter.Seq[int](Slice([]int{item, item * 10}))
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

func TestSeqLazyTerminalsStopEarly(t *testing.T) {
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

func TestSeqLazyMiddleOperationsStopEarly(t *testing.T) {
	for item := range Slice([]int{1, 1, 2}).DistinctBy(func(item int, _ int) int {
		return item
	}) {
		if item != 1 {
			t.Fatalf("distinct item: %d", item)
		}
		break
	}

	for item := range Slice([]string{"go", "ko", "kod"}).
		DistinctBy(func(item string, _ int) int {
			return len(item)
		}) {
		if item != "go" {
			t.Fatalf("distinctBy item: %q", item)
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
		Filter(func(item int, _ int) bool {
			return item%2 == 0
		}).
		Take(2) {
		if item != 2 {
			t.Fatalf("filter take item: %d", item)
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

	for item := range Slice([]int{3, 1, 2}).Sort(func(left, right int) bool {
		return left < right
	}) {
		if item != 1 {
			t.Fatalf("sort item: %d", item)
		}
		break
	}
}

func TestSeqHelpers(t *testing.T) {
	if got := Slice([]int{1, 2, 3}).Reverse().Collect(); !reflect.DeepEqual(got, []int{3, 2, 1}) {
		t.Fatalf("reverse: %#v", got)
	}

	if got := Slice([]int{3, 1, 2}).Sort(func(left, right int) bool {
		return left < right
	}).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("sort: %#v", got)
	}

	if got := Slice([]int{2, 1, 1}).Sort(func(left, right int) bool {
		return left < right
	}).Collect(); !reflect.DeepEqual(got, []int{1, 1, 2}) {
		t.Fatalf("sort equal: %#v", got)
	}

	got := Slice([]string{"aaa", "b", "cc"}).
		SortBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()
	if !reflect.DeepEqual(got, []string{"b", "cc", "aaa"}) {
		t.Fatalf("sortBy: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Reduce(func(sum int, item int, _ int) int {
		return sum + item
	}, 0); got != 6 {
		t.Fatalf("reduce: %d", got)
	}

	if got := Slice([]int{1, 2, 3}).Join(":", func(item int, index int) string {
		return strconv.Itoa(index) + "=" + strconv.Itoa(item)
	}); got != "0=1:1=2:2=3" {
		t.Fatalf("join: %q", got)
	}

	if got := Slice([]int{}).Join(",", func(item int, _ int) string {
		return strconv.Itoa(item)
	}); got != "" {
		t.Fatalf("join empty: %q", got)
	}
}

func TestSeqTerminalAggregates(t *testing.T) {
	if got, ok := Slice([]string{"go", "gopher", "ko"}).MaxBy(func(item string, _ int) int {
		return len(item)
	}); !ok || got != "gopher" {
		t.Fatalf("maxBy: %q, %v", got, ok)
	}

	if got, ok := Slice([]string{"go", "gopher", "ko"}).MaxBy(func(_ string, index int) int {
		return index
	}); !ok || got != "ko" {
		t.Fatalf("maxBy index: %q, %v", got, ok)
	}

	if got, ok := Slice([]string{"go", "gopher", "ko"}).MinBy(func(item string, _ int) int {
		return len(item)
	}); !ok || got != "go" {
		t.Fatalf("minBy: %q, %v", got, ok)
	}

	if got := Slice([]string{"go", "gopher", "ko"}).SumBy(func(item string, _ int) int {
		return len(item)
	}); got != 10 {
		t.Fatalf("sumBy: %d", got)
	}

	if got := Slice([]string{"go", "gopher", "ko"}).SumBy(func(_ string, index int) int {
		return index
	}); got != 3 {
		t.Fatalf("sumBy index: %d", got)
	}

	if got := Slice([]string{"go", "gopher", "ko"}).MeanBy(func(item string, _ int) int {
		return len(item)
	}); got != float64(10)/3 {
		t.Fatalf("meanBy: %f", got)
	}

	if got := Of[string]().MeanBy(func(item string, _ int) int {
		return len(item)
	}); got != 0 {
		t.Fatalf("meanBy empty: %f", got)
	}
}

func TestSeqScan(t *testing.T) {
	got := Slice([]int{1, 2, 3}).Scan(func(sum int, item int, _ int) int {
		return sum + item
	}, 0).Collect()
	if !reflect.DeepEqual(got, []int{0, 1, 3, 6}) {
		t.Fatalf("scan: %#v", got)
	}

	gotStrings := Slice([]int{1, 2}).Scan(func(out string, item int, _ int) string {
		return out + strconv.Itoa(item)
	}, "").Collect()
	if !reflect.DeepEqual(gotStrings, []string{"", "1", "12"}) {
		t.Fatalf("scan strings: %#v", gotStrings)
	}

	if got := Slice([]int{}).Scan(func(sum int, item int, _ int) int {
		return sum + item
	}, 10).Collect(); !reflect.DeepEqual(got, []int{10}) {
		t.Fatalf("scan empty: %#v", got)
	}
}

func TestSeqScanIsLazyAndStopsEarly(t *testing.T) {
	calls := 0
	seq := Slice([]int{1, 2, 3}).Map(func(item int, _ int) int {
		calls++
		return item
	}).Scan(func(sum int, item int, _ int) int {
		return sum + item
	}, 0)

	if calls != 0 {
		t.Fatalf("scan should be lazy, calls: %d", calls)
	}

	for item := range seq {
		if item != 0 {
			t.Fatalf("scan first item: %d", item)
		}
		break
	}
	if calls != 0 {
		t.Fatalf("scan should yield initial without consuming source, calls: %d", calls)
	}

	for item := range seq {
		if item != 0 {
			t.Fatalf("scan second range first item: %d", item)
		}
		break
	}
	if calls != 0 {
		t.Fatalf("scan should remain lazy across ranges, calls: %d", calls)
	}

	for item := range seq.Drop(1) {
		if item != 1 {
			t.Fatalf("scan second item: %d", item)
		}
		break
	}
	if calls != 1 {
		t.Fatalf("scan should stop after one source item, calls: %d", calls)
	}
}

func TestSeqSortByIsLazy(t *testing.T) {
	calls := 0
	seq := Slice([]string{"aaa", "b", "cc"}).SortBy(func(item string, _ int) int {
		calls++
		return len(item)
	})
	if calls != 0 {
		t.Fatalf("sortBy should be lazy, calls: %d", calls)
	}

	for item := range seq {
		if item != "b" {
			t.Fatalf("sortBy first item: %q", item)
		}
		break
	}
	if calls != 3 {
		t.Fatalf("sortBy needs the whole source once consumed, calls: %d", calls)
	}
}

func TestSeqCollectBackedOperationsAreLazyUntilConsumed(t *testing.T) {
	assertLazy := func(name string, build func(Seq[int]) Seq[int], wantFirst int) {
		calls := 0
		seq := build(Slice([]int{1, 2, 3}).Map(func(item int, _ int) int {
			calls++
			return item
		}))
		if calls != 0 {
			t.Fatalf("%s called before consumption: %d", name, calls)
		}
		for item := range seq {
			if item != wantFirst {
				t.Fatalf("%s item: %d", name, item)
			}
			break
		}
		if calls != 3 {
			t.Fatalf("%s calls: %d", name, calls)
		}
	}

	assertLazy("sort", func(seq Seq[int]) Seq[int] {
		return seq.Sort(func(left, right int) bool { return left < right })
	}, 1)
	assertLazy("reverse", func(seq Seq[int]) Seq[int] { return seq.Reverse() }, 3)
	assertLazy("takeRight", func(seq Seq[int]) Seq[int] { return seq.TakeRight(2) }, 2)
}

func TestSeqCountBackedOperationsAreLazyUntilConsumed(t *testing.T) {
	calls := 0
	grouped := Slice([]string{"go", "ko", "kod"}).Map(func(item string, _ int) string {
		calls++
		return item
	}).GroupBy(func(item string, _ int) int {
		return len(item)
	})
	if calls != 0 {
		t.Fatalf("groupBy called before consumption: %d", calls)
	}
	for key, group := range grouped {
		if key != 2 {
			t.Fatalf("groupBy key: %d", key)
		}
		if got := group.Collect(); !reflect.DeepEqual(got, []string{"go", "ko"}) {
			t.Fatalf("groupBy group: %#v", got)
		}
		break
	}
	if calls != 3 {
		t.Fatalf("groupBy calls: %d", calls)
	}
}

func TestSeqTakeWhile(t *testing.T) {
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

func TestSeqFilterThenTake(t *testing.T) {
	got := Slice([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(item int, _ int) bool {
			return item%2 == 0
		}).
		Take(2).
		Collect()

	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Slice([]int{1, 2, 3}).
		Filter(func(item int, _ int) bool {
			return true
		}).
		Take(0).
		Collect()

	if len(got) != 0 {
		t.Fatalf("zero: %#v", got)
	}
}

func TestSeqDropWhile(t *testing.T) {
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

func TestSeqDistinctBy(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod"}).
		DistinctBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want := []string{"go", "kod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeqDistinctBySupportsNonComparableItems(t *testing.T) {
	type item struct {
		id     int
		values []int
	}

	got := Slice([]item{
		{id: 1, values: []int{1}},
		{id: 1, values: []int{2}},
		{id: 2, values: []int{3}},
	}).DistinctBy(func(item item, _ int) int {
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

func TestSeqPredicates(t *testing.T) {
	if Slice([]int{1}).IsEmpty() {
		t.Fatal("isEmpty non-empty")
	}

	if !Slice([]int{}).IsEmpty() {
		t.Fatal("isEmpty empty")
	}

	calls := 0
	empty := Slice([]int{1, 2, 3}).Map(func(item int, _ int) int {
		calls++
		return item
	}).IsEmpty()
	if empty {
		t.Fatal("isEmpty mapped non-empty")
	}
	if calls != 1 {
		t.Fatalf("isEmpty should stop after first item, calls: %d", calls)
	}

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

	if got := !Slice([]int{1, 2, 3}).Some(func(item int, _ int) bool {
		return item > 9
	}); !got {
		t.Fatal("none: false")
	}

	if got := !Slice([]int{1, 2, 3}).Some(func(item int, _ int) bool {
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

func TestSeqElementAccess(t *testing.T) {
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

func TestSeqNth(t *testing.T) {
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

func TestSeqFindIndex(t *testing.T) {
	index, ok := Slice([]string{"go", "ko", "kod"}).FindIndex(func(item string, _ int) bool {
		return len(item) == 3
	})
	if !ok || index != 2 {
		t.Fatalf("findIndex: %d, %v", index, ok)
	}

	index, ok = Slice([]string{"go"}).FindIndex(func(item string, _ int) bool {
		return len(item) == 3
	})
	if ok || index != -1 {
		t.Fatalf("findIndex miss: %d, %v", index, ok)
	}
}

func TestSeqMoreHelpers(t *testing.T) {
	if got := Slice([]int{1, 2, 3}).Reject(func(item int, _ int) bool {
		return item == 2
	}).Collect(); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("reject: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).FlatMap(func(item int, _ int) iter.Seq[int] {
		return iter.Seq[int](Slice([]int{item, item * 10}))
	}).Collect(); !reflect.DeepEqual(got, []int{1, 10, 2, 20, 3, 30}) {
		t.Fatalf("flatMap: %#v", got)
	}

	if got := Slice([]int{1, 2}).Concat(iter.Seq[int](Slice([]int{3, 4}))).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("concat: %#v", got)
	}

	if got := Slice([]int{1, 2}).Concat(iter.Seq[int](Slice([]int{}))).Collect(); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("concat empty: %#v", got)
	}

	concatCalls := 0
	concatSeq := Slice([]int{1, 2}).Map(func(item int, _ int) int {
		concatCalls++
		return item
	}).Concat(iter.Seq[int](Slice([]int{3, 4})))
	if concatCalls != 0 {
		t.Fatalf("concat should be lazy: %d", concatCalls)
	}
	for item := range concatSeq {
		if item != 1 {
			t.Fatalf("concat first item: %d", item)
		}
		break
	}
	if concatCalls != 1 {
		t.Fatalf("concat should stop with range: %d", concatCalls)
	}

	concatCalls = 0
	for item := range Slice([]int{}).Concat(iter.Seq[int](Slice([]int{3, 4}).Map(func(item int, _ int) int {
		concatCalls++
		return item
	}))) {
		if item != 3 {
			t.Fatalf("concat right item: %d", item)
		}
		break
	}
	if concatCalls != 1 {
		t.Fatalf("concat should stop in right sequence: %d", concatCalls)
	}

	if got := Slice([]int{1, 2, 3}).Take(0).Collect(); len(got) != 0 {
		t.Fatalf("take zero: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Take(9).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("take all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Take(1).Collect(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("take one: %#v", got)
	}

	if got := Slice([]int{}).Take(2).Collect(); len(got) != 0 {
		t.Fatalf("take empty: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Drop(0).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("drop zero: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).Drop(9).Collect(); len(got) != 0 {
		t.Fatalf("drop all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).TakeRight(2).Collect(); !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("takeRight: %#v", got)
	}

	for item := range Slice([]int{1, 2, 3}).TakeRight(2) {
		if item != 2 {
			t.Fatalf("takeRight item: %d", item)
		}
		break
	}

	if got := Slice([]int{1, 2, 3}).TakeRight(0).Collect(); len(got) != 0 {
		t.Fatalf("takeRight zero: %#v", got)
	}

	takeRightCalls := 0
	if got := Slice([]int{1, 2, 3}).Map(func(item int, _ int) int {
		takeRightCalls++
		return item
	}).TakeRight(0).Collect(); len(got) != 0 {
		t.Fatalf("takeRight zero mapped: %#v", got)
	}
	if takeRightCalls != 0 {
		t.Fatalf("takeRight zero calls: %d", takeRightCalls)
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

	for item := range Slice([]int{1, 2, 3}).DropRight(0) {
		if item != 1 {
			t.Fatalf("dropRight zero item: %d", item)
		}
		break
	}

	if got := Slice([]int{1, 2, 3}).DropRight(9).Collect(); len(got) != 0 {
		t.Fatalf("dropRight all: %#v", got)
	}

	calls := 0
	for item := range Slice([]int{1, 2, 3, 4}).Map(func(item int, _ int) int {
		calls++
		return item
	}).DropRight(2) {
		if item != 1 {
			t.Fatalf("dropRight item: %d", item)
		}
		break
	}
	if calls != 3 {
		t.Fatalf("dropRight calls: %d", calls)
	}

	if got := Slice([]int{1, 2, 3}).TakeWhile(func(item int, _ int) bool {
		return item < 9
	}).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("takeWhile all: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).DropWhile(func(item int, _ int) bool {
		return item < 9
	}).Collect(); len(got) != 0 {
		t.Fatalf("dropWhile all: %#v", got)
	}

	item, ok := Slice([]int{1, 2, 3}).Find(func(item int, _ int) bool {
		return item == 9
	})
	if ok || item != 0 {
		t.Fatalf("find miss: %d, %v", item, ok)
	}
}

func TestSeqToMap(t *testing.T) {
	got := Slice([]string{"go", "kod"}).
		ToMap(func(item string, _ int) (int, string) {
			return len(item), item
		}).
		Collect()

	want := map[int]string{2: "go", 3: "kod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	calls := 0
	seq := Slice([]int{1, 2}).Map(func(item int, _ int) int {
		calls++
		return item
	}).ToMap(func(item int, _ int) (int, int) {
		return item, item * 10
	})
	if calls != 0 {
		t.Fatalf("toMap should be lazy: %d", calls)
	}
	for key, value := range seq {
		if key != 1 || value != 10 {
			t.Fatalf("toMap first entry: %d %d", key, value)
		}
		break
	}
	if calls != 1 {
		t.Fatalf("toMap calls: %d", calls)
	}
}

func TestSeqEnumerate(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod"}).Enumerate().Collect()
	want := map[int]string{0: "go", 1: "ko", 2: "kod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	calls := 0
	seq := Slice([]string{"go", "ko"}).
		Map(func(item string, _ int) string {
			calls++
			return item
		}).
		Enumerate()
	if calls != 0 {
		t.Fatalf("enumerate should be lazy: %d", calls)
	}

	for index, item := range seq {
		if index != 0 || item != "go" {
			t.Fatalf("first pair: %d, %q", index, item)
		}
		break
	}
	if calls != 1 {
		t.Fatalf("enumerate should stop with range: %d", calls)
	}
}

func TestSeqGroupBy(t *testing.T) {
	grouped := Slice([]string{"go", "ko", "kod"}).
		GroupBy(func(item string, _ int) int {
			return len(item)
		})
	got := grouped.Collect()

	want := map[int][]string{2: {"go", "ko"}, 3: {"kod"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	sizes := grouped.Map(func(key int, group Seq[string]) int {
		return key + group.Count(func(_ string, _ int) bool { return true })
	}).Collect()
	if !reflect.DeepEqual(sizes, []int{4, 4}) {
		t.Fatalf("group map: %#v", sizes)
	}

	groupMapCalls := 0
	for item := range grouped.Map(func(key int, _ Seq[string]) int {
		groupMapCalls++
		return key
	}) {
		if item != 2 {
			t.Fatalf("group map first: %d", item)
		}
		break
	}
	if groupMapCalls != 1 {
		t.Fatalf("group map calls: %d", groupMapCalls)
	}

	for range 64 {
		var keys []int
		for key := range Slice([]string{"kod", "go", "ko", "go-kod"}).
			GroupBy(func(item string, _ int) int {
				return len(item)
			}) {
			keys = append(keys, key)
		}
		if !reflect.DeepEqual(keys, []int{3, 2, 6}) {
			t.Fatalf("group keys: %#v", keys)
		}
	}
}

func TestSeqKeyBy(t *testing.T) {
	got := Slice([]string{"g", "go", "ko"}).
		KeyBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want := map[int]string{1: "g", 2: "ko"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Slice([]string{"go", "ko", "kod"}).
		KeyBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want = map[int]string{2: "ko", 3: "kod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("duplicate keys got %#v, want %#v", got, want)
	}
}

func TestSeqMapConversionsStopEarly(t *testing.T) {
	for range Slice([]string{"go", "ko", "kod"}).
		GroupBy(func(item string, _ int) int {
			return len(item)
		}) {
		break
	}

	for range Slice([]string{"go", "ko", "kod"}).
		KeyBy(func(item string, _ int) int {
			return len(item)
		}) {
		break
	}
}

func TestSeqMapConversionsAreLazyUntilConsumed(t *testing.T) {
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
	if calls != 1 {
		t.Fatalf("keyBy calls: %d", calls)
	}
}

func TestSeqChunk(t *testing.T) {
	var got [][]int
	for chunk := range Slice([]int{1, 2, 3, 4, 5}).Chunk(2) {
		got = append(got, chunk.Collect())
	}

	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeqChunkOuterCanBeChained(t *testing.T) {
	got := Slice([]int{1, 2, 3, 4}).Chunk(2).
		Map(func(chunk Seq[int], _ int) int {
			items := chunk.Collect()
			return items[0] + items[1]
		}).
		Collect()

	want := []int{3, 7}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	chunks := Slice([]int{1, 2, 3}).Chunk(2).Collect()
	if !reflect.DeepEqual(chunks, [][]int{{1, 2}, {3}}) {
		t.Fatalf("chunk collect: %#v", chunks)
	}

	mapCalls := 0
	for item := range Slice([]int{1, 2, 3, 4}).Chunk(2).Map(func(chunk Seq[int], _ int) int {
		mapCalls++
		return len(chunk.Collect())
	}) {
		if item != 2 {
			t.Fatalf("chunk map first: %d", item)
		}
		break
	}
	if mapCalls != 1 {
		t.Fatalf("chunk map calls: %d", mapCalls)
	}
}

func TestSeqChunkStopsEarly(t *testing.T) {
	calls := 0
	chunks := Slice([]int{1, 2, 3, 4}).Map(func(item int, _ int) int {
		calls++
		return item
	}).Chunk(2)
	if calls != 0 {
		t.Fatalf("chunk called before consumption: %d", calls)
	}
	for chunk := range chunks {
		if got := chunk.Collect(); !reflect.DeepEqual(got, []int{1, 2}) {
			t.Fatalf("chunk: %#v", got)
		}
		break
	}
	if calls != 2 {
		t.Fatalf("calls: %d", calls)
	}
}

func TestSeqChunkZeroSizeYieldsNothing(t *testing.T) {
	for chunk := range Slice([]int{1, 2, 3}).Chunk(0) {
		t.Fatalf("unexpected chunk: %#v", chunk.Collect())
	}
}

func TestSeqWindow(t *testing.T) {
	var got [][]int
	for window := range Slice([]int{1, 2, 3, 4}).Window(3, 1) {
		got = append(got, window.Collect())
	}

	want := [][]int{{1, 2, 3}, {2, 3, 4}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = nil
	for window := range Slice([]int{1, 2, 3, 4, 5}).Window(2, 3) {
		got = append(got, window.Collect())
	}

	want = [][]int{{1, 2}, {4, 5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeqWindowOuterCanBeChained(t *testing.T) {
	got := Slice([]int{1, 2, 3, 4}).Window(3, 1).
		Map(func(window Seq[int], _ int) int {
			items := window.Collect()
			return items[0] + items[2]
		}).
		Collect()

	want := []int{4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	windows := Slice([]int{1, 2, 3, 4}).Window(3, 1).Collect()
	if !reflect.DeepEqual(windows, [][]int{{1, 2, 3}, {2, 3, 4}}) {
		t.Fatalf("window collect: %#v", windows)
	}
}

func TestSeqWindowStopsEarly(t *testing.T) {
	calls := 0
	windows := Slice([]int{1, 2, 3, 4}).Map(func(item int, _ int) int {
		calls++
		return item
	}).Window(3, 1)
	if calls != 0 {
		t.Fatalf("window called before consumption: %d", calls)
	}
	for window := range windows {
		if got := window.Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
			t.Fatalf("window: %#v", got)
		}
		break
	}
	if calls != 3 {
		t.Fatalf("calls: %d", calls)
	}
}

func TestSeqWindowZeroSizeYieldsNothing(t *testing.T) {
	for window := range Slice([]int{1, 2, 3}).Window(0, 1) {
		t.Fatalf("unexpected window: %#v", window.Collect())
	}
	for window := range Slice([]int{1, 2, 3}).Window(2, 0) {
		t.Fatalf("unexpected window: %#v", window.Collect())
	}
}
