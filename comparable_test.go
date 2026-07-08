package ko

import (
	"reflect"
	"testing"
)

func TestComparableDistinctCompactWithout(t *testing.T) {
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

func TestSeqComparableMethods(t *testing.T) {
	seq := Comparable(Slice([]int{0, 1, 2, 2, 3}))

	if got := seq.Distinct().Drop(1).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("distinct drop: %#v", got)
	}
	if got := seq.Compact().Collect(); !reflect.DeepEqual(got, []int{1, 2, 2, 3}) {
		t.Fatalf("compact: %#v", got)
	}
	if got := seq.Without(2).Collect(); !reflect.DeepEqual(got, []int{0, 1, 3}) {
		t.Fatalf("without: %#v", got)
	}
	if got := seq.Filter(func(item int, _ int) bool {
		return item > 0
	}).Reject(func(item int, _ int) bool {
		return item == 2
	}).Take(1).Collect(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("filter chain: %#v", got)
	}
	if got := seq.TakeWhile(func(item int, _ int) bool {
		return item < 2
	}).Collect(); !reflect.DeepEqual(got, []int{0, 1}) {
		t.Fatalf("takeWhile: %#v", got)
	}
	if got := seq.DropWhile(func(item int, _ int) bool {
		return item < 2
	}).Collect(); !reflect.DeepEqual(got, []int{2, 2, 3}) {
		t.Fatalf("dropWhile: %#v", got)
	}
	if got := seq.Seq().Take(2).Collect(); !reflect.DeepEqual(got, []int{0, 1}) {
		t.Fatalf("seq: %#v", got)
	}
}

func TestSeqComparableTerminals(t *testing.T) {
	seq := Comparable(Slice([]string{"go", "ko", "go"}))

	if !seq.Contains("ko") || seq.Contains("kod") {
		t.Fatal("contains")
	}
	index, ok := seq.IndexOf("go")
	if !ok || index != 0 {
		t.Fatalf("indexOf: %d %v", index, ok)
	}
	index, ok = seq.LastIndexOf("go")
	if !ok || index != 2 {
		t.Fatalf("lastIndexOf: %d %v", index, ok)
	}
	if got := seq.CountValues(); !reflect.DeepEqual(got, map[string]int{"go": 2, "ko": 1}) {
		t.Fatalf("countValues: %#v", got)
	}
	if got := seq.ToSet(); !reflect.DeepEqual(got, map[string]struct{}{"go": {}, "ko": {}}) {
		t.Fatalf("toSet: %#v", got)
	}
	if !seq.Equal(Comparable(Slice([]string{"go", "ko", "go"}))) {
		t.Fatal("equal true")
	}
	if seq.Equal(Comparable(Slice([]string{"go"}))) {
		t.Fatal("equal false")
	}
}

func TestSeqComparableSetOperations(t *testing.T) {
	a := Comparable(Slice([]int{1, 2, 2, 3}))
	b := Comparable(Slice([]int{2, 3, 4}))

	if got := a.Union(b).Collect(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("union: %#v", got)
	}
	if got := a.Intersect(b).Collect(); !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("intersect: %#v", got)
	}
	if got := a.Difference(b).Collect(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("difference: %#v", got)
	}
	if got := a.SymmetricDifference(b).Collect(); !reflect.DeepEqual(got, []int{1, 4}) {
		t.Fatalf("symmetricDifference: %#v", got)
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
	assertLazySeq(t, "distinct", Distinct[int], 1, 0)
	assertLazySeq(t, "compact", Compact[int], 2, 1)
	assertLazySeq(t, "without", func(seq Seq[int]) Seq[int] { return Without(seq, 0) }, 2, 1)
	assertLazySeq(t, "union", func(seq Seq[int]) Seq[int] { return Union(seq, Slice([]int{3})) }, 1, 0)
	assertLazySeq(t, "intersect", func(seq Seq[int]) Seq[int] { return Intersect(seq, Slice([]int{1, 2})) }, 2, 1)
	assertLazySeq(t, "difference", func(seq Seq[int]) Seq[int] { return Difference(seq, Slice([]int{0, 2})) }, 2, 1)
	assertLazySeq(t, "symmetricDifference", func(seq Seq[int]) Seq[int] { return SymmetricDifference(seq, Slice([]int{2, 3})) }, 4, 0)
}

func TestSeqComparableReturningMethodsAreLazyUntilConsumed(t *testing.T) {
	assertLazySeq(t, "seqComparable seq", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqComparable distinct", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Distinct().Seq()
	}, 1, 0)
	assertLazySeq(t, "seqComparable compact", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Compact().Seq()
	}, 2, 1)
	assertLazySeq(t, "seqComparable without", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Without(0).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqComparable union", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Union(Comparable(Slice([]int{3}))).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqComparable intersect", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Intersect(Comparable(Slice([]int{1, 2}))).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqComparable difference", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Difference(Comparable(Slice([]int{0, 2}))).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqComparable symmetricDifference", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).SymmetricDifference(Comparable(Slice([]int{2, 3}))).Seq()
	}, 4, 0)
	assertLazySeq(t, "seqComparable filter", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Filter(func(item int, _ int) bool { return item > 0 }).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqComparable reject", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Reject(func(item int, _ int) bool { return item == 0 }).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqComparable take", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Take(2).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqComparable drop", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).Drop(1).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqComparable takeWhile", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).TakeWhile(func(item int, _ int) bool { return item < 2 }).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqComparable dropWhile", func(seq Seq[int]) Seq[int] {
		return Comparable(seq).DropWhile(func(item int, _ int) bool { return item < 1 }).Seq()
	}, 2, 1)
}

func assertLazySeq(t *testing.T, name string, build func(Seq[int]) Seq[int], wantCalls int, wantFirst int) {
	t.Helper()

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
