package ko

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestSeq2(t *testing.T) {
	got := Map(map[string]int{"a": 1, "bb": 2, "ccc": 3}).
		Filter(func(_ string, value int) bool {
			return value > 1
		}).
		Map(func(key string, value int) (int, string) {
			return len(key), strconv.Itoa(value * 10)
		}).
		Collect()

	want := map[int]string{2: "20", 3: "30"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeq2CollectEmptyReturnsEmptyMap(t *testing.T) {
	got := Map(map[string]int{}).Collect()
	if got == nil || len(got) != 0 {
		t.Fatalf("got %#v", got)
	}

	got = Map(map[string]int{"a": 1}).PickKeys().Collect()
	if got == nil || len(got) != 0 {
		t.Fatalf("picked none: %#v", got)
	}
}

func TestSeq2CollectDuplicateKeysKeepLast(t *testing.T) {
	got := Seq2[string, int](func(yield func(string, int) bool) {
		if !yield("x", 1) {
			return
		}
		if !yield("x", 2) {
			return
		}
		yield("y", 3)
	}).Collect()

	want := map[string]int{"x": 2, "y": 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeq2ToSlice(t *testing.T) {
	got := Map(map[string]int{"a": 1, "b": 2}).
		MapValues(func(_ string, value int) string {
			return strconv.Itoa(value)
		}).
		ToSlice(func(key string, value string) string {
			return key + value
		}).
		Collect()
	sort.Strings(got)

	want := []string{"a1", "b2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeq2Iteration(t *testing.T) {
	got := map[string]int{}
	for key, value := range Map(map[string]int{"a": 1, "b": 2}) {
		got[key] = value
	}

	want := map[string]int{"a": 1, "b": 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeq2LazyOperationsStopEarly(t *testing.T) {
	filterCalls := 0
	for range Map(map[string]int{"a": 1, "b": 2}).
		Filter(func(_ string, _ int) bool {
			filterCalls++
			return true
		}) {
		break
	}
	if filterCalls != 1 {
		t.Fatalf("filter calls: %d", filterCalls)
	}

	for range Map(map[string]int{"a": 1, "b": 2}).PickKeys("a", "b") {
		break
	}

	for range Map(map[string]int{"a": 1, "b": 2}).OmitKeys() {
		break
	}

	mapCalls := 0
	for range Map(map[string]int{"a": 1, "b": 2}).
		Map(func(key string, value int) (string, int) {
			mapCalls++
			return key, value * 10
		}) {
		break
	}
	if mapCalls != 1 {
		t.Fatalf("map calls: %d", mapCalls)
	}

	mapKeysCalls := 0
	for range Map(map[string]int{"a": 1, "b": 2}).
		MapKeys(func(key string, _ int) string {
			mapKeysCalls++
			return key + "!"
		}) {
		break
	}
	if mapKeysCalls != 1 {
		t.Fatalf("mapKeys calls: %d", mapKeysCalls)
	}

	mapValuesCalls := 0
	for range Map(map[string]int{"a": 1, "b": 2}).
		MapValues(func(_ string, value int) int {
			mapValuesCalls++
			return value * 10
		}) {
		break
	}
	if mapValuesCalls != 1 {
		t.Fatalf("mapValues calls: %d", mapValuesCalls)
	}

	for range Map(map[string]int{"a": 1, "b": 2}).Keys() {
		break
	}

	for range Map(map[string]int{"a": 1, "b": 2}).Values() {
		break
	}

	filterKeysCalls := 0
	for range Map(map[string]int{"a": 1, "b": 2}).
		Filter(func(_ string, _ int) bool {
			filterKeysCalls++
			return true
		}).
		Keys() {
		break
	}
	if filterKeysCalls != 1 {
		t.Fatalf("filterKeys calls: %d", filterKeysCalls)
	}

	filterValuesCalls := 0
	for range Map(map[string]int{"a": 1, "b": 2}).
		Filter(func(_ string, _ int) bool {
			filterValuesCalls++
			return true
		}).
		Values() {
		break
	}
	if filterValuesCalls != 1 {
		t.Fatalf("filterValues calls: %d", filterValuesCalls)
	}

	toSliceCalls := 0
	for range Map(map[string]int{"a": 1, "b": 2}).
		ToSlice(func(key string, value int) string {
			toSliceCalls++
			return key + strconv.Itoa(value)
		}) {
		break
	}
	if toSliceCalls != 1 {
		t.Fatalf("toSlice calls: %d", toSliceCalls)
	}

	toMapCalls := 0
	for range Slice([]int{1, 2, 3}).
		ToMap(func(item int, _ int) (int, int) {
			toMapCalls++
			return item, item * 10
		}) {
		break
	}
	if toMapCalls != 1 {
		t.Fatalf("toMap calls: %d", toMapCalls)
	}
}

func TestSeq2ReturningMethodsAreLazyUntilConsumed(t *testing.T) {
	assertLazySeq2(t, "pickKeys", func(seq Seq2[string, int]) Seq2[string, int] {
		return seq.PickKeys("a", "b")
	}, 1, "a", 1)
	assertLazySeq2(t, "omitKeys", func(seq Seq2[string, int]) Seq2[string, int] {
		return seq.OmitKeys("b")
	}, 1, "a", 1)
	assertLazySeq2(t, "filter", func(seq Seq2[string, int]) Seq2[string, int] {
		return seq.Filter(func(_ string, value int) bool { return value > 0 })
	}, 1, "a", 1)
	assertLazySeq2(t, "reject", func(seq Seq2[string, int]) Seq2[string, int] {
		return seq.Reject(func(_ string, value int) bool { return value == 0 })
	}, 1, "a", 1)
	assertLazySeq2(t, "map", func(seq Seq2[string, int]) Seq2[string, int] {
		return seq.Map(func(key string, value int) (string, int) { return key, value * 10 })
	}, 1, "a", 10)
	assertLazySeq2(t, "mapKeys", func(seq Seq2[string, int]) Seq2[string, int] {
		return seq.MapKeys(func(key string, _ int) string { return key + "!" })
	}, 1, "a!", 1)
	assertLazySeq2(t, "mapValues", func(seq Seq2[string, int]) Seq2[string, int] {
		return seq.MapValues(func(_ string, value int) int { return value * 10 })
	}, 1, "a", 10)
	calls := 0
	empty := lazySeq2(&calls).PickKeys()
	if calls != 0 {
		t.Fatalf("pickKeys empty called before consumption: %d", calls)
	}
	for key, value := range empty {
		t.Fatalf("pickKeys empty yielded %q %d", key, value)
	}
	if calls != 0 {
		t.Fatalf("pickKeys empty calls: %d", calls)
	}
}

func TestSeq2SliceReturningMethodsAreLazyUntilConsumed(t *testing.T) {
	assertLazySeq2ToSeq(t, "keys", func(seq Seq2[string, int]) Seq[string] {
		return seq.Keys()
	}, 1, "a")
	assertLazySeq2ToSeq(t, "values", func(seq Seq2[string, int]) Seq[string] {
		return seq.Values().Map(func(item int, _ int) string { return strconv.Itoa(item) })
	}, 1, "1")
	assertLazySeq2ToSeq(t, "toSlice", func(seq Seq2[string, int]) Seq[string] {
		return seq.ToSlice(func(key string, value int) string { return key + strconv.Itoa(value) })
	}, 1, "a1")
}

func lazySeq2(calls *int) Seq2[string, int] {
	return Seq2[string, int](func(yield func(string, int) bool) {
		(*calls)++
		if !yield("a", 1) {
			return
		}
		(*calls)++
		yield("b", 2)
	})
}

func assertLazySeq2(t *testing.T, name string, build func(Seq2[string, int]) Seq2[string, int], wantCalls int, wantKey string, wantValue int) {
	t.Helper()

	calls := 0
	seq := build(lazySeq2(&calls))
	if calls != 0 {
		t.Fatalf("%s called before consumption: %d", name, calls)
	}
	for key, value := range seq {
		if key != wantKey || value != wantValue {
			t.Fatalf("%s entry: %q %d", name, key, value)
		}
		break
	}
	if calls != wantCalls {
		t.Fatalf("%s calls: %d", name, calls)
	}
}

func assertLazySeq2ToSeq(t *testing.T, name string, build func(Seq2[string, int]) Seq[string], wantCalls int, wantFirst string) {
	t.Helper()

	calls := 0
	seq := build(lazySeq2(&calls))
	if calls != 0 {
		t.Fatalf("%s called before consumption: %d", name, calls)
	}
	for item := range seq {
		if item != wantFirst {
			t.Fatalf("%s item: %q", name, item)
		}
		break
	}
	if calls != wantCalls {
		t.Fatalf("%s calls: %d", name, calls)
	}
}

func TestSeq2Keys(t *testing.T) {
	got := Map(map[string]int{"b": 2, "a": 1}).Keys().Collect()
	sort.Strings(got)

	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeq2Values(t *testing.T) {
	got := Map(map[string]int{"a": 1, "b": 2}).Values().Collect()
	sort.Ints(got)

	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestSeq2FilterThenKeys(t *testing.T) {
	got := Map(map[string]int{"a": 1, "bb": 2, "ccc": 3}).
		Filter(func(key string, value int) bool {
			return len(key) == value
		}).
		Keys().
		Collect()
	sort.Strings(got)

	want := []string{"a", "bb", "ccc"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Map(map[string]int{"a": 2}).
		Filter(func(key string, value int) bool {
			return len(key) == value
		}).
		Keys().
		Collect()
	if len(got) != 0 {
		t.Fatalf("filter then keys miss: %#v", got)
	}
}

func TestSeq2FilterThenValues(t *testing.T) {
	got := Map(map[string]int{"a": 1, "bb": 2, "ccc": 4}).
		Filter(func(key string, value int) bool {
			return len(key) == value
		}).
		Values().
		Collect()
	sort.Ints(got)

	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Map(map[string]int{"a": 2}).
		Filter(func(key string, value int) bool {
			return len(key) == value
		}).
		Values().
		Collect()
	if len(got) != 0 {
		t.Fatalf("filter then values miss: %#v", got)
	}
}

func TestSeq2PickKeys(t *testing.T) {
	got := Map(map[string]int{"a": 1, "b": 2, "c": 3}).
		PickKeys("a", "c").
		Collect()

	want := map[string]int{"a": 1, "c": 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Map(map[string]int{"a": 1}).PickKeys().Collect()
	if len(got) != 0 {
		t.Fatalf("pick none: %#v", got)
	}

	calls := 0
	got = Seq2[string, int](func(yield func(string, int) bool) {
		calls++
		yield("a", 1)
	}).PickKeys().Collect()
	if len(got) != 0 {
		t.Fatalf("pick none seq: %#v", got)
	}
	if calls != 0 {
		t.Fatalf("pick none calls: %d", calls)
	}
}

func TestSeq2OmitKeys(t *testing.T) {
	got := Map(map[string]int{"a": 1, "b": 2, "c": 3}).
		OmitKeys("a", "c").
		Collect()

	want := map[string]int{"b": 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	got = Map(map[string]int{"a": 1}).OmitKeys().Collect()
	want = map[string]int{"a": 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("omit none: %#v", got)
	}

	for key, value := range Map(map[string]int{"a": 1, "b": 2}).OmitKeys("b") {
		if key != "a" || value != 1 {
			t.Fatalf("omit entry: %q %d", key, value)
		}
		break
	}
}

func TestSeq2Predicates(t *testing.T) {
	if Map(map[string]int{"a": 1}).IsEmpty() {
		t.Fatal("isEmpty non-empty")
	}

	if !Map(map[string]int{}).IsEmpty() {
		t.Fatal("isEmpty empty")
	}

	calls := 0
	empty := Seq2[string, int](func(yield func(string, int) bool) {
		calls++
		if !yield("a", 1) {
			return
		}
		yield("b", 2)
	}).IsEmpty()
	if empty {
		t.Fatal("isEmpty seq non-empty")
	}
	if calls != 1 {
		t.Fatalf("isEmpty should stop after first entry, calls: %d", calls)
	}

	if got := Map(map[string]int{"a": 1, "b": 2}).Some(func(_ string, value int) bool {
		return value == 2
	}); !got {
		t.Fatal("some: false")
	}

	if got := Map(map[string]int{"a": 1, "b": 2}).Every(func(key string, value int) bool {
		return len(key) == 1 && value > 0
	}); !got {
		t.Fatal("every: false")
	}

	if got := !Map(map[string]int{"a": 1, "b": 2}).Some(func(_ string, value int) bool {
		return value > 9
	}); !got {
		t.Fatal("none: false")
	}

	if got := !Map(map[string]int{"a": 1, "b": 2}).Some(func(_ string, value int) bool {
		return value == 2
	}); got {
		t.Fatal("none hit: true")
	}

	if got := Map(map[string]int{"a": 1, "b": 2, "c": 4}).Count(func(_ string, value int) bool {
		return value%2 == 0
	}); got != 2 {
		t.Fatalf("count: %d", got)
	}

	if got := Map(map[string]int{"a": 1}).Count(func(_ string, value int) bool {
		return value%2 == 0
	}); got != 0 {
		t.Fatalf("count miss: %d", got)
	}
}

func TestSeq2Find(t *testing.T) {
	key, value, ok := Map(map[string]int{"a": 1, "b": 2}).Find(func(_ string, value int) bool {
		return value == 2
	})

	if !ok || key != "b" || value != 2 {
		t.Fatalf("got %q, %d, %v", key, value, ok)
	}
}

func TestSeq2MoreHelpers(t *testing.T) {
	got := Map(map[string]int{"a": 1, "b": 2}).
		Reject(func(_ string, value int) bool {
			return value == 1
		}).
		MapKeys(func(key string, _ int) string {
			return key + key
		}).
		Collect()

	want := map[string]int{"bb": 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	if got := Map(map[string]int{"a": 1}).Some(func(_ string, value int) bool {
		return value == 2
	}); got {
		t.Fatal("some: true")
	}

	if got := Map(map[string]int{"a": 1, "b": 0}).Every(func(_ string, value int) bool {
		return value > 0
	}); got {
		t.Fatal("every: true")
	}

	key, value, ok := Map(map[string]int{"a": 1}).Find(func(_ string, value int) bool {
		return value == 2
	})
	if ok || key != "" || value != 0 {
		t.Fatalf("find miss: %q, %d, %v", key, value, ok)
	}
}
