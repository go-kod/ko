package ko

import (
	"reflect"
	"testing"
)

func TestSeqOrderedMethods(t *testing.T) {
	seq := Ordered(Slice([]int{0, 3, 1, 2, 2}))

	if got := seq.Sort().Collect(); !reflect.DeepEqual(got, []int{0, 1, 2, 2, 3}) {
		t.Fatalf("sort: %#v", got)
	}
	max, ok := seq.Max()
	if !ok || max != 3 {
		t.Fatalf("max: %d %v", max, ok)
	}
	min, ok := seq.Min()
	if !ok || min != 0 {
		t.Fatalf("min: %d %v", min, ok)
	}
	if got := seq.Distinct().Drop(1).Collect(); !reflect.DeepEqual(got, []int{3, 1, 2}) {
		t.Fatalf("distinct/drop: %#v", got)
	}
	if got := seq.Compact().Collect(); !reflect.DeepEqual(got, []int{3, 1, 2, 2}) {
		t.Fatalf("compact: %#v", got)
	}
	if got := seq.Without(2).Collect(); !reflect.DeepEqual(got, []int{0, 3, 1}) {
		t.Fatalf("without: %#v", got)
	}
	if got := seq.Filter(func(item int, _ int) bool {
		return item > 0
	}).Reject(func(item int, _ int) bool {
		return item == 3
	}).Take(2).Collect(); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("filter chain: %#v", got)
	}
	if got := seq.TakeWhile(func(item int, _ int) bool {
		return item < 3
	}).Collect(); !reflect.DeepEqual(got, []int{0}) {
		t.Fatalf("takeWhile: %#v", got)
	}
	if got := seq.DropWhile(func(item int, _ int) bool {
		return item < 3
	}).Collect(); !reflect.DeepEqual(got, []int{3, 1, 2, 2}) {
		t.Fatalf("dropWhile: %#v", got)
	}
	if got := seq.Seq().Take(2).Collect(); !reflect.DeepEqual(got, []int{0, 3}) {
		t.Fatalf("seq: %#v", got)
	}
	if got := seq.Comparable().Distinct().Collect(); !reflect.DeepEqual(got, []int{0, 3, 1, 2}) {
		t.Fatalf("comparable: %#v", got)
	}
}

func TestSeqOrderedComparableMethods(t *testing.T) {
	seq := Ordered(Slice([]string{"go", "ko", "go"}))

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
	if !seq.Equal(Ordered(Slice([]string{"go", "ko", "go"}))) {
		t.Fatal("equal true")
	}
	if seq.Equal(Ordered(Slice([]string{"go"}))) {
		t.Fatal("equal false")
	}
}

func TestSeqOrderedSetOperations(t *testing.T) {
	a := Ordered(Slice([]int{1, 2, 2, 3}))
	b := Ordered(Slice([]int{2, 3, 4}))

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

func TestSeqOrderedReturningMethodsAreLazyUntilConsumed(t *testing.T) {
	assertLazySeq(t, "seqOrdered seq", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqOrdered comparable", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Comparable().Seq()
	}, 1, 0)
	assertLazySeq(t, "seqOrdered sort", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Sort().Seq()
	}, 4, 0)
	assertLazySeq(t, "seqOrdered distinct", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Distinct().Seq()
	}, 1, 0)
	assertLazySeq(t, "seqOrdered compact", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Compact().Seq()
	}, 2, 1)
	assertLazySeq(t, "seqOrdered without", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Without(0).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqOrdered union", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Union(Ordered(Slice([]int{3}))).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqOrdered intersect", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Intersect(Ordered(Slice([]int{1, 2}))).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqOrdered difference", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Difference(Ordered(Slice([]int{0, 2}))).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqOrdered symmetricDifference", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).SymmetricDifference(Ordered(Slice([]int{2, 3}))).Seq()
	}, 4, 0)
	assertLazySeq(t, "seqOrdered filter", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Filter(func(item int, _ int) bool { return item > 0 }).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqOrdered reject", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Reject(func(item int, _ int) bool { return item == 0 }).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqOrdered take", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Take(2).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqOrdered drop", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).Drop(1).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqOrdered takeWhile", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).TakeWhile(func(item int, _ int) bool { return item < 2 }).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqOrdered dropWhile", func(seq Seq[int]) Seq[int] {
		return Ordered(seq).DropWhile(func(item int, _ int) bool { return item < 1 }).Seq()
	}, 2, 1)
}

func TestSeqNumericMethods(t *testing.T) {
	seq := Numbers(Slice([]int{0, 3, 1, 2, 2}))

	if got := seq.Sum(); got != 8 {
		t.Fatalf("sum: %d", got)
	}
	if got := seq.Product(); got != 0 {
		t.Fatalf("product: %d", got)
	}
	if got := Numbers(Slice([]int{3, 2})).Product(); got != 6 {
		t.Fatalf("product non-zero: %d", got)
	}
	if got := Numbers(Empty[int]()).Product(); got != 1 {
		t.Fatalf("product empty: %d", got)
	}
	if got := seq.Mean(); got != 1.6 {
		t.Fatalf("mean: %v", got)
	}
	if got := Numbers(Empty[int]()).Mean(); got != 0 {
		t.Fatalf("mean empty: %v", got)
	}
	if got := Numbers(Slice([]float64{1, 2, 3, 4})).Mean(); got != 2.5 {
		t.Fatalf("mean float: %v", got)
	}

	if got := seq.Sort().Collect(); !reflect.DeepEqual(got, []int{0, 1, 2, 2, 3}) {
		t.Fatalf("sort: %#v", got)
	}
	max, ok := seq.Max()
	if !ok || max != 3 {
		t.Fatalf("max: %d %v", max, ok)
	}
	min, ok := seq.Min()
	if !ok || min != 0 {
		t.Fatalf("min: %d %v", min, ok)
	}
	if got := seq.Distinct().Drop(1).Collect(); !reflect.DeepEqual(got, []int{3, 1, 2}) {
		t.Fatalf("distinct/drop: %#v", got)
	}
	if got := seq.Compact().Collect(); !reflect.DeepEqual(got, []int{3, 1, 2, 2}) {
		t.Fatalf("compact: %#v", got)
	}
	if got := seq.Without(2).Collect(); !reflect.DeepEqual(got, []int{0, 3, 1}) {
		t.Fatalf("without: %#v", got)
	}
	if got := seq.Filter(func(item int, _ int) bool {
		return item > 0
	}).Reject(func(item int, _ int) bool {
		return item == 3
	}).Take(2).Collect(); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("filter chain: %#v", got)
	}
	if got := seq.TakeWhile(func(item int, _ int) bool {
		return item < 3
	}).Collect(); !reflect.DeepEqual(got, []int{0}) {
		t.Fatalf("takeWhile: %#v", got)
	}
	if got := seq.DropWhile(func(item int, _ int) bool {
		return item < 3
	}).Collect(); !reflect.DeepEqual(got, []int{3, 1, 2, 2}) {
		t.Fatalf("dropWhile: %#v", got)
	}
	if got := seq.Seq().Take(2).Collect(); !reflect.DeepEqual(got, []int{0, 3}) {
		t.Fatalf("seq: %#v", got)
	}
	if got := seq.Ordered().Sort().Collect(); !reflect.DeepEqual(got, []int{0, 1, 2, 2, 3}) {
		t.Fatalf("ordered: %#v", got)
	}
	if got := seq.Comparable().Distinct().Collect(); !reflect.DeepEqual(got, []int{0, 3, 1, 2}) {
		t.Fatalf("comparable: %#v", got)
	}
}

func TestSeqNumericComparableMethods(t *testing.T) {
	seq := Numbers(Slice([]int{1, 2, 1}))

	if !seq.Contains(2) || seq.Contains(9) {
		t.Fatal("contains")
	}
	index, ok := seq.IndexOf(1)
	if !ok || index != 0 {
		t.Fatalf("indexOf: %d %v", index, ok)
	}
	index, ok = seq.LastIndexOf(1)
	if !ok || index != 2 {
		t.Fatalf("lastIndexOf: %d %v", index, ok)
	}
	if got := seq.CountValues(); !reflect.DeepEqual(got, map[int]int{1: 2, 2: 1}) {
		t.Fatalf("countValues: %#v", got)
	}
	if got := seq.ToSet(); !reflect.DeepEqual(got, map[int]struct{}{1: {}, 2: {}}) {
		t.Fatalf("toSet: %#v", got)
	}
	if !seq.Equal(Numbers(Slice([]int{1, 2, 1}))) {
		t.Fatal("equal true")
	}
	if seq.Equal(Numbers(Slice([]int{1}))) {
		t.Fatal("equal false")
	}
}

func TestSeqNumericSetOperations(t *testing.T) {
	a := Numbers(Slice([]int{1, 2, 2, 3}))
	b := Numbers(Slice([]int{2, 3, 4}))

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

func TestSeqNumericReturningMethodsAreLazyUntilConsumed(t *testing.T) {
	assertLazySeq(t, "seqNumeric seq", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqNumeric ordered", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Ordered().Seq()
	}, 1, 0)
	assertLazySeq(t, "seqNumeric comparable", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Comparable().Seq()
	}, 1, 0)
	assertLazySeq(t, "seqNumeric sort", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Sort().Seq()
	}, 4, 0)
	assertLazySeq(t, "seqNumeric distinct", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Distinct().Seq()
	}, 1, 0)
	assertLazySeq(t, "seqNumeric compact", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Compact().Seq()
	}, 2, 1)
	assertLazySeq(t, "seqNumeric without", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Without(0).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqNumeric union", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Union(Numbers(Slice([]int{3}))).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqNumeric intersect", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Intersect(Numbers(Slice([]int{1, 2}))).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqNumeric difference", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Difference(Numbers(Slice([]int{0, 2}))).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqNumeric symmetricDifference", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).SymmetricDifference(Numbers(Slice([]int{2, 3}))).Seq()
	}, 4, 0)
	assertLazySeq(t, "seqNumeric filter", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Filter(func(item int, _ int) bool { return item > 0 }).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqNumeric reject", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Reject(func(item int, _ int) bool { return item == 0 }).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqNumeric take", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Take(2).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqNumeric drop", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).Drop(1).Seq()
	}, 2, 1)
	assertLazySeq(t, "seqNumeric takeWhile", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).TakeWhile(func(item int, _ int) bool { return item < 2 }).Seq()
	}, 1, 0)
	assertLazySeq(t, "seqNumeric dropWhile", func(seq Seq[int]) Seq[int] {
		return Numbers(seq).DropWhile(func(item int, _ int) bool { return item < 1 }).Seq()
	}, 2, 1)
}
