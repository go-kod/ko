package ko

import (
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

	if got := Of[int]().Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("of empty: %#v", got)
	}
}

func TestGenerate(t *testing.T) {
	calls := 0
	seq := Generate(1, func(item int) int {
		calls++
		return item * 2
	})
	if calls != 0 {
		t.Fatalf("generate should be lazy: %d", calls)
	}

	if got := seq.Take(4).Collect(); !reflect.DeepEqual(got, []int{1, 2, 4, 8}) {
		t.Fatalf("generate: %#v", got)
	}
	if calls != 3 {
		t.Fatalf("generate should only compute requested next values: %d", calls)
	}
}

func TestRange(t *testing.T) {
	if got := Range(1, 4).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("range ascending: %#v", got)
	}

	if got := Range(4, 1).Collect(); !reflect.DeepEqual(got, []int{4, 3, 2}) {
		t.Fatalf("range descending: %#v", got)
	}

	if got := Range(2, 2).Collect(); !reflect.DeepEqual(got, []int{}) {
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

	if got := RangeStep(1, 8, 0).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("rangeStep zero step: %#v", got)
	}

	if got := RangeStep(1, 8, -1).Collect(); !reflect.DeepEqual(got, []int{}) {
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

func TestTimes(t *testing.T) {
	if got := Times(4, func(index int) string {
		return strconv.Itoa(index * index)
	}).Collect(); !reflect.DeepEqual(got, []string{"0", "1", "4", "9"}) {
		t.Fatalf("times: %#v", got)
	}

	calls := 0
	if got := Times(0, func(index int) int {
		calls++
		return index
	}).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("times zero: %#v", got)
	}
	if calls != 0 {
		t.Fatalf("times zero should not call mapper: %d", calls)
	}

	seq := Times(5, func(index int) int {
		calls++
		return index + 10
	})
	if calls != 0 {
		t.Fatalf("times should be lazy: %d", calls)
	}

	var got []int
	for item := range seq {
		got = append(got, item)
		if len(got) == 2 {
			break
		}
	}
	if !reflect.DeepEqual(got, []int{10, 11}) {
		t.Fatalf("times should stop with range: %#v", got)
	}
	if calls != 2 {
		t.Fatalf("times calls: %d", calls)
	}
}

func TestRepeat(t *testing.T) {
	if got := Repeat(3, "go").Collect(); !reflect.DeepEqual(got, []string{"go", "go", "go"}) {
		t.Fatalf("repeat: %#v", got)
	}

	if got := Repeat(0, "go").Collect(); !reflect.DeepEqual(got, []string{}) {
		t.Fatalf("repeat zero: %#v", got)
	}

	if got := Repeat(-1, "go").Collect(); !reflect.DeepEqual(got, []string{}) {
		t.Fatalf("repeat negative: %#v", got)
	}

	var got []int
	for item := range Repeat(5, 7) {
		got = append(got, item)
		if len(got) == 2 {
			break
		}
	}
	if !reflect.DeepEqual(got, []int{7, 7}) {
		t.Fatalf("repeat should stop with range: %#v", got)
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
	if got := seq.Collect(); !reflect.DeepEqual(got, []int{}) {
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

func TestSeqCompact(t *testing.T) {
	if got := Slice([]int{0, 1, 0, 2}).Compact().Collect(); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("compact ints: %#v", got)
	}

	type item struct {
		Name string
		Tags []string
	}

	got := Slice([]item{{}, {Name: "go"}, {Tags: []string{"ko"}}}).Compact().Collect()
	want := []item{{Name: "go"}, {Tags: []string{"ko"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("compact non-comparable: %#v", got)
	}

	calls := 0
	seq := Slice([]int{0, 0, 3, 4}).
		Map(func(item int, _ int) int {
			calls++
			return item
		}).
		Compact()
	if calls != 0 {
		t.Fatalf("compact should be lazy, calls: %d", calls)
	}

	for item := range seq {
		if item != 3 {
			t.Fatalf("compact first item: %d", item)
		}
		break
	}
	if calls != 3 {
		t.Fatalf("compact should stop after first non-zero item, calls: %d", calls)
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
		FlatMap(func(item int, _ int) Seq[int] {
			flatMapCalls++
			return Slice([]int{item, item * 10})
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
}

func TestSeqHelpers(t *testing.T) {
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

	if !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("zero: %#v", got)
	}
}

func TestSeqSubset(t *testing.T) {
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

	got = Slice([]string{"go", "ko"}).Subset(-1, 0).Collect()
	if !reflect.DeepEqual(got, []string{}) {
		t.Fatalf("subset negative zero: %#v", got)
	}

	got = Slice([]string{}).Subset(-1, 1).Collect()
	if !reflect.DeepEqual(got, []string{}) {
		t.Fatalf("subset negative empty: %#v", got)
	}
}

func TestSeqSubsetPositiveOffsetStopsEarly(t *testing.T) {
	calls := 0
	got := Slice([]int{1, 2, 3, 4}).Map(func(item int, _ int) int {
		calls++
		return item
	}).Subset(1, 2).Collect()

	if !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("subset: %#v", got)
	}
	if calls != 3 {
		t.Fatalf("calls: %d", calls)
	}

	for item := range Slice([]int{1, 2, 3}).Subset(1, 2) {
		if item != 2 {
			t.Fatalf("item: %d", item)
		}
		break
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

func TestSeqTakeRightWhile(t *testing.T) {
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

func TestSeqDropRightWhile(t *testing.T) {
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

func TestSeqDropByIndex(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod", "go-kod"}).
		DropByIndex(1, -1, 99).
		Collect()

	want := []string{"go", "kod"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	if got := Slice([]int{1, 2, 3}).DropByIndex(1).Collect(); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("drop positive: %#v", got)
	}
}

func TestSeqDropByPositiveIndexStopsEarly(t *testing.T) {
	calls := 0
	for item := range Slice([]int{1, 2, 3}).Map(func(item int, _ int) int {
		calls++
		return item
	}).DropByIndex(1) {
		if item != 1 {
			t.Fatalf("item: %d", item)
		}
		break
	}
	if calls != 1 {
		t.Fatalf("calls: %d", calls)
	}
}

func TestSeqUniq(t *testing.T) {
	got := Slice([]int{1, 2, 1, 3, 2}).Uniq().Collect()

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeqUniqBy(t *testing.T) {
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

func TestSeqUniqBySupportsNonComparableItems(t *testing.T) {
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

func TestSeqIsUniqBy(t *testing.T) {
	type item struct {
		id     int
		values []int
	}

	if !Slice([]item{{id: 1}, {id: 2}}).IsUniqBy(func(item item, _ int) int {
		return item.id
	}) {
		t.Fatalf("unique items should be unique")
	}

	if Slice([]item{{id: 1, values: []int{1}}, {id: 1, values: []int{2}}}).IsUniqBy(func(item item, _ int) int {
		return item.id
	}) {
		t.Fatalf("duplicate keys should not be unique")
	}

	calls := 0
	ok := Slice([]int{1, 2, 1, 3}).IsUniqBy(func(item int, _ int) int {
		calls++
		return item
	})
	if ok {
		t.Fatalf("duplicate values should not be unique")
	}
	if calls != 3 {
		t.Fatalf("isUniqBy should stop at first duplicate, calls: %d", calls)
	}

	if !Slice([]int{}).IsUniqBy(func(item int, _ int) int {
		return item
	}) {
		t.Fatalf("empty sequence should be unique")
	}
}

func TestSeqUniqMap(t *testing.T) {
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

func TestSeqFindUniquesBy(t *testing.T) {
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

	for item := range Slice([]string{"go", "kod", "x"}).
		FindUniquesBy(func(item string, _ int) int {
			return len(item)
		}) {
		if item != "go" {
			t.Fatalf("unique item: %q", item)
		}
		break
	}
}

func TestSeqFindDuplicatesBy(t *testing.T) {
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

	for item := range Slice([]string{"go", "ko", "x"}).
		FindDuplicatesBy(func(item string, _ int) int {
			return len(item)
		}) {
		if item != "go" {
			t.Fatalf("duplicate item: %q", item)
		}
		break
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

func TestSeqFilterReject(t *testing.T) {
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

func TestSeqFilterRejectIsLazyUntilConsumed(t *testing.T) {
	calls := 0
	kept, rejected := Slice([]int{1, 2, 3, 4}).FilterReject(func(item int, _ int) bool {
		calls++
		return item%2 == 0
	})
	if calls != 0 {
		t.Fatalf("filterReject called before consumption: %d", calls)
	}

	if got := kept.Collect(); !reflect.DeepEqual(got, []int{2, 4}) {
		t.Fatalf("kept: %#v", got)
	}
	if calls != 4 {
		t.Fatalf("filterReject calls after kept: %d", calls)
	}

	if got := rejected.Collect(); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("rejected: %#v", got)
	}
	if calls != 4 {
		t.Fatalf("filterReject called again: %d", calls)
	}
}

func TestSeqFilterRejectStopsEarly(t *testing.T) {
	kept, rejected := Slice([]int{1, 2, 3, 4}).FilterReject(func(item int, _ int) bool {
		return item%2 == 0
	})

	for item := range kept {
		if item != 2 {
			t.Fatalf("kept: %d", item)
		}
		break
	}

	for item := range rejected {
		if item != 1 {
			t.Fatalf("rejected: %d", item)
		}
		break
	}
}

func TestSeqFilterRejectAdvancesOnDemand(t *testing.T) {
	calls := 0
	kept, rejected := Slice([]int{1, 2, 3, 4}).Map(func(item int, _ int) int {
		calls++
		return item
	}).FilterReject(func(item int, _ int) bool {
		return item%2 == 0
	})

	for item := range kept {
		if item != 2 {
			t.Fatalf("kept: %d", item)
		}
		break
	}
	if calls != 2 {
		t.Fatalf("kept should only advance to first match, calls: %d", calls)
	}

	for item := range rejected {
		if item != 1 {
			t.Fatalf("rejected: %d", item)
		}
		break
	}
	if calls != 2 {
		t.Fatalf("rejected should use cached item, calls: %d", calls)
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

func TestSeqFindLast(t *testing.T) {
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

func TestSeqFindLastIndex(t *testing.T) {
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

func TestSeqWithoutBy(t *testing.T) {
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

	calls := 0
	got = Slice([]string{"go", "ko"}).
		WithoutBy(func(item string, _ int) int {
			calls++
			return len(item)
		}).
		Collect()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("without empty: %#v", got)
	}
	if calls != 0 {
		t.Fatalf("without empty calls: %d", calls)
	}
}

func TestSeqMoreHelpers(t *testing.T) {
	if got := Slice([]int{1, 2, 3}).Reject(func(item int, _ int) bool {
		return item == 2
	}).Collect(); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("reject: %#v", got)
	}

	if got := Slice([]int{1, 2, 3}).FlatMap(func(item int, _ int) Seq[int] {
		return Slice([]int{item, item * 10})
	}).Collect(); !reflect.DeepEqual(got, []int{1, 10, 2, 20, 3, 30}) {
		t.Fatalf("flatMap: %#v", got)
	}

	if got := Slice([]int{1, 2}).Concat(Slice([]int{3, 4})).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("concat: %#v", got)
	}

	if got := Slice([]int{1, 2}).Concat(Slice([]int{})).Collect(); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("concat empty: %#v", got)
	}

	concatCalls := 0
	concatSeq := Slice([]int{1, 2}).Map(func(item int, _ int) int {
		concatCalls++
		return item
	}).Concat(Slice([]int{3, 4}))
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
	for item := range Slice([]int{}).Concat(Slice([]int{3, 4}).Map(func(item int, _ int) int {
		concatCalls++
		return item
	})) {
		if item != 3 {
			t.Fatalf("concat right item: %d", item)
		}
		break
	}
	if concatCalls != 1 {
		t.Fatalf("concat should stop in right sequence: %d", concatCalls)
	}

	if got := Slice([]int{1, 2, 3}).Intersperse(0).Collect(); !reflect.DeepEqual(got, []int{1, 0, 2, 0, 3}) {
		t.Fatalf("intersperse: %#v", got)
	}

	if got := Slice([]int{}).Intersperse(0).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("intersperse empty: %#v", got)
	}

	if got := Slice([]int{1}).Intersperse(0).Collect(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("intersperse single: %#v", got)
	}

	intersperseCalls := 0
	intersperseSeq := Slice([]int{1, 2}).Map(func(item int, _ int) int {
		intersperseCalls++
		return item
	}).Intersperse(0)
	if intersperseCalls != 0 {
		t.Fatalf("intersperse should be lazy: %d", intersperseCalls)
	}
	for item := range intersperseSeq {
		if item != 1 {
			t.Fatalf("intersperse first item: %d", item)
		}
		break
	}
	if intersperseCalls != 1 {
		t.Fatalf("intersperse should stop with range: %d", intersperseCalls)
	}

	var interspersed []int
	for item := range Slice([]int{1, 2, 3}).Intersperse(0) {
		interspersed = append(interspersed, item)
		if len(interspersed) == 2 {
			break
		}
	}
	if !reflect.DeepEqual(interspersed, []int{1, 0}) {
		t.Fatalf("intersperse break after separator: %#v", interspersed)
	}

	seen := 0
	seq := Slice([]int{1, 2, 3}).ForEach(func(item int, _ int) {
		seen += item
	})
	if seen != 0 {
		t.Fatalf("forEach should be lazy: %d", seen)
	}
	if got := seq.Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("forEach value: %#v", got)
	}
	if seen != 6 {
		t.Fatalf("forEach: %d", seen)
	}

	seen = 0
	got := Slice([]int{1, 2, 3}).ForEachWhile(func(item int, _ int) bool {
		seen += item
		return item < 2
	})
	if seen != 0 {
		t.Fatalf("forEachWhile should be lazy: %d", seen)
	}
	values := got.Collect()
	if seen != 3 {
		t.Fatalf("forEachWhile: %d", seen)
	}
	if !reflect.DeepEqual(values, []int{1, 2, 3}) {
		t.Fatalf("forEachWhile value: %#v", values)
	}

	seen = 0
	for item := range Slice([]int{1, 2, 3}).ForEach(func(item int, _ int) {
		seen += item
	}) {
		if item != 1 {
			t.Fatalf("forEach item: %d", item)
		}
		break
	}
	if seen != 1 {
		t.Fatalf("forEach should stop with range: %d", seen)
	}

	seen = 0
	for item := range Slice([]int{1, 2, 3}).ForEachWhile(func(item int, _ int) bool {
		seen += item
		return true
	}) {
		if item != 1 {
			t.Fatalf("forEachWhile item: %d", item)
		}
		break
	}
	if seen != 1 {
		t.Fatalf("forEachWhile should stop with range: %d", seen)
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

	for item := range Slice([]int{1, 2, 3}).TakeRight(2) {
		if item != 2 {
			t.Fatalf("takeRight item: %d", item)
		}
		break
	}

	if got := Slice([]int{1, 2, 3}).TakeRight(0).Collect(); !reflect.DeepEqual(got, []int{}) {
		t.Fatalf("takeRight zero: %#v", got)
	}

	takeRightCalls := 0
	if got := Slice([]int{1, 2, 3}).Map(func(item int, _ int) int {
		takeRightCalls++
		return item
	}).TakeRight(0).Collect(); !reflect.DeepEqual(got, []int{}) {
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

	if got := Slice([]int{1, 2, 3}).DropRight(9).Collect(); !reflect.DeepEqual(got, []int{}) {
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

func TestSeqGroupByMap(t *testing.T) {
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

	for range 64 {
		var keys []int
		for key := range Slice([]string{"kod", "go", "ko", "go-kod"}).
			GroupByMap(func(item string, _ int) (int, string) {
				return len(item), item + "!"
			}) {
			keys = append(keys, key)
		}
		if !reflect.DeepEqual(keys, []int{3, 2, 6}) {
			t.Fatalf("groupByMap keys: %#v", keys)
		}
	}
}

func TestSeqPartitionBy(t *testing.T) {
	got := [][]string{}
	for group := range Slice([]string{"go", "kod", "ko", "go-kod"}).
		PartitionBy(func(item string, _ int) int {
			return len(item)
		}) {
		got = append(got, group.Collect())
	}

	want := [][]string{{"go", "ko"}, {"kod"}, {"go-kod"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeqPartitionByIsLazyUntilConsumed(t *testing.T) {
	calls := 0
	groups := Slice([]string{"go", "kod", "ko"}).
		PartitionBy(func(item string, _ int) int {
			calls++
			return len(item)
		})
	if calls != 0 {
		t.Fatalf("partitionBy called before consumption: %d", calls)
	}
	for group := range groups {
		if got := group.Collect(); !reflect.DeepEqual(got, []string{"go", "ko"}) {
			t.Fatalf("group: %#v", got)
		}
		break
	}
	if calls != 3 {
		t.Fatalf("partitionBy calls: %d", calls)
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

func TestSeqCountBy(t *testing.T) {
	got := Slice([]string{"go", "ko", "kod"}).
		CountBy(func(item string, _ int) int {
			return len(item)
		}).
		Collect()

	want := map[int]int{2: 2, 3: 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	for range 64 {
		var keys []int
		for key := range Slice([]string{"kod", "go", "ko", "go-kod"}).
			CountBy(func(item string, _ int) int {
				return len(item)
			}) {
			keys = append(keys, key)
		}
		if !reflect.DeepEqual(keys, []int{3, 2, 6}) {
			t.Fatalf("countBy keys: %#v", keys)
		}
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

func TestSeqChunkOuterCanBeChainedByConversion(t *testing.T) {
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

func TestSeqChunkStopsEarly(t *testing.T) {
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

func TestSeqWindowOuterCanBeChainedByConversion(t *testing.T) {
	got := Seq[Seq[int]](Slice([]int{1, 2, 3, 4}).Window(3, 1)).
		Map(func(window Seq[int], _ int) int {
			items := window.Collect()
			return items[0] + items[2]
		}).
		Collect()

	want := []int{4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeqWindowStopsEarly(t *testing.T) {
	calls := 0
	for window := range Slice([]int{1, 2, 3, 4}).Map(func(item int, _ int) int {
		calls++
		return item
	}).Window(3, 1) {
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
