package ko

import (
	"reflect"
	"testing"
)

func TestComparableDistinctCompactWithout(t *testing.T) {
	if got := Uniq(Slice([]int{1, 2, 1, 3, 2})).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("uniq: %#v", got)
	}

	if got := Distinct(Slice([]int{1, 2, 1, 3, 2})).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("distinct: %#v", got)
	}

	if got := Compact(Slice([]string{"", "go", "", "ko"})).Collect(); !reflect.DeepEqual(got, []string{"go", "ko"}) {
		t.Fatalf("compact: %#v", got)
	}

	if got := Without(Slice([]int{1, 2, 3, 2}), 2, 9).Collect(); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("without: %#v", got)
	}
}

func TestComparableLookups(t *testing.T) {
	if !Contains(Slice([]string{"go", "ko"}), "ko") {
		t.Fatal("contains miss")
	}
	if Contains(Slice([]string{"go", "ko"}), "kod") {
		t.Fatal("contains false hit")
	}

	index, ok := IndexOf(Slice([]string{"go", "ko", "go"}), "go")
	if !ok || index != 0 {
		t.Fatalf("indexOf: %d %v", index, ok)
	}
	index, ok = IndexOf(Slice([]string{"go"}), "ko")
	if ok || index != -1 {
		t.Fatalf("indexOf miss: %d %v", index, ok)
	}

	index, ok = LastIndexOf(Slice([]string{"go", "ko", "go"}), "go")
	if !ok || index != 2 {
		t.Fatalf("lastIndexOf: %d %v", index, ok)
	}
	index, ok = LastIndexOf(Slice([]string{"go"}), "ko")
	if ok || index != -1 {
		t.Fatalf("lastIndexOf miss: %d %v", index, ok)
	}
}

func TestComparableAggregates(t *testing.T) {
	if got := CountValues(Slice([]string{"go", "ko", "go"})); !reflect.DeepEqual(got, map[string]int{"go": 2, "ko": 1}) {
		t.Fatalf("countValues: %#v", got)
	}
	if got := CountValues(Empty[string]()); got == nil || len(got) != 0 {
		t.Fatalf("countValues empty: %#v", got)
	}

	if got := ToSet(Slice([]int{1, 2, 1})); !reflect.DeepEqual(got, map[int]struct{}{1: {}, 2: {}}) {
		t.Fatalf("toSet: %#v", got)
	}
	if got := ToSet(Empty[int]()); got == nil || len(got) != 0 {
		t.Fatalf("toSet empty: %#v", got)
	}

	if got := JoinStrings(Slice([]string{"go", "ko"}), "/"); got != "go/ko" {
		t.Fatalf("joinStrings: %q", got)
	}
	if got := JoinStrings(Empty[string](), "/"); got != "" {
		t.Fatalf("joinStrings empty: %q", got)
	}
}

func TestComparableEqual(t *testing.T) {
	if !Equal(Slice([]int{1, 2}), Slice([]int{1, 2})) {
		t.Fatal("equal false")
	}
	if !Equal(Empty[int](), Empty[int]()) {
		t.Fatal("empty equal false")
	}
	if Equal(Slice([]int{1}), Slice([]int{1, 2})) {
		t.Fatal("different length equal")
	}
	if Equal(Slice([]int{1, 2}), Slice([]int{1, 3})) {
		t.Fatal("different value equal")
	}
}

func TestComparableSetOperations(t *testing.T) {
	if got := Union(Slice([]int{1, 2, 1}), Slice([]int{2, 3})).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("union: %#v", got)
	}
	if got := Union[int]().Collect(); len(got) != 0 {
		t.Fatalf("union empty: %#v", got)
	}

	if got := Intersect(Slice([]int{1, 2, 2, 3}), Slice([]int{2, 3, 4})).Collect(); !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("intersect: %#v", got)
	}

	if got := Difference(Slice([]int{1, 2, 2, 3, 4}), Slice([]int{2, 4})).Collect(); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("difference: %#v", got)
	}
	if got := Difference(Slice([]int{1, 1, 2}), Slice([]int{2})).Collect(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("difference duplicate: %#v", got)
	}

	if got := SymmetricDifference(Slice([]int{1, 2, 2, 3}), Slice([]int{2, 3, 4, 4})).Collect(); !reflect.DeepEqual(got, []int{1, 4}) {
		t.Fatalf("symmetricDifference: %#v", got)
	}
	for item := range SymmetricDifference(Slice([]int{1}), Slice([]int{1, 2})) {
		if item != 2 {
			t.Fatalf("symmetricDifference right item: %d", item)
		}
		break
	}
}

func TestComparableSequencesAreLazyUntilConsumed(t *testing.T) {
	assertLazy := func(name string, build func(Seq[int]) Seq[int], wantCalls int, wantFirst int) {
		calls := 0
		seq := build(Slice([]int{0, 1, 2, 1}).Map(func(item int, _ int) int {
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
		if calls != wantCalls {
			t.Fatalf("%s calls: %d", name, calls)
		}
	}

	assertLazy("uniq", Uniq[int], 1, 0)
	assertLazy("distinct", Distinct[int], 1, 0)
	assertLazy("compact", Compact[int], 2, 1)
	assertLazy("without", func(seq Seq[int]) Seq[int] { return Without(seq, 0) }, 2, 1)
	assertLazy("union", func(seq Seq[int]) Seq[int] { return Union(seq, Slice([]int{3})) }, 1, 0)
	assertLazy("intersect", func(seq Seq[int]) Seq[int] { return Intersect(seq, Slice([]int{1, 2})) }, 2, 1)
	assertLazy("difference", func(seq Seq[int]) Seq[int] { return Difference(seq, Slice([]int{0, 2})) }, 2, 1)
	assertLazy("symmetricDifference", func(seq Seq[int]) Seq[int] { return SymmetricDifference(seq, Slice([]int{2, 3})) }, 4, 0)
}
