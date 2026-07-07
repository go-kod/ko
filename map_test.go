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

func TestSeq2FilterMapToSlice(t *testing.T) {
	got := Map(map[string]int{"a": 1, "b": 2, "c": 3}).
		FilterMapToSlice(func(key string, value int) (string, bool) {
			return key + strconv.Itoa(value), value%2 == 1
		}).
		Collect()
	sort.Strings(got)

	want := []string{"a1", "c3"}
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

	for range Map(map[string]int{"a": 1}).Assign(map[string]int{"b": 2}) {
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

	filterMapCalls := 0
	for range Map(map[string]int{"a": 1, "b": 2}).
		FilterMapToSlice(func(key string, value int) (string, bool) {
			filterMapCalls++
			return key + strconv.Itoa(value), true
		}) {
		break
	}
	if filterMapCalls != 1 {
		t.Fatalf("filterMap calls: %d", filterMapCalls)
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

func TestSeq2ChunkEntries(t *testing.T) {
	chunks := []map[string]int{}
	for chunk := range Map(map[string]int{"a": 1, "b": 2, "c": 3}).ChunkEntries(2) {
		chunks = append(chunks, chunk.Collect())
	}
	if len(chunks) != 2 {
		t.Fatalf("chunks: %#v", chunks)
	}

	merged := map[string]int{}
	for _, chunk := range chunks {
		if len(chunk) > 2 {
			t.Fatalf("chunk too large: %#v", chunk)
		}
		for key, value := range chunk {
			merged[key] = value
		}
	}

	want := map[string]int{"a": 1, "b": 2, "c": 3}
	if !reflect.DeepEqual(merged, want) {
		t.Fatalf("merged %#v, want %#v", merged, want)
	}

	for chunk := range Map(map[string]int{}).ChunkEntries(2) {
		t.Fatalf("empty chunk: %#v", chunk.Collect())
	}
}

func TestSeq2ChunkEntriesStopsEarly(t *testing.T) {
	calls := 0
	for chunk := range Map(map[int]int{1: 1, 2: 2, 3: 3}).
		MapValues(func(_ int, value int) int {
			calls++
			return value
		}).
		ChunkEntries(2) {
		if got := chunk.Collect(); len(got) != 2 {
			t.Fatalf("chunk: %#v", got)
		}
		break
	}
	if calls != 2 {
		t.Fatalf("calls: %d", calls)
	}
}

func TestSeq2ChunkEntriesCountsDuplicateKeysAsEntries(t *testing.T) {
	chunks := []map[string]int{}
	var firstChunkValues []int
	for chunk := range Seq2[string, int](func(yield func(string, int) bool) {
		if !yield("x", 1) {
			return
		}
		if !yield("x", 2) {
			return
		}
		yield("y", 3)
	}).ChunkEntries(2) {
		if len(chunks) == 0 {
			for key, value := range chunk {
				if key != "x" {
					t.Fatalf("first chunk key: %q", key)
				}
				firstChunkValues = append(firstChunkValues, value)
			}
		}
		chunks = append(chunks, chunk.Collect())
	}

	if !reflect.DeepEqual(firstChunkValues, []int{1, 2}) {
		t.Fatalf("first chunk values: %#v", firstChunkValues)
	}

	for chunk := range Seq2[string, int](func(yield func(string, int) bool) {
		if !yield("x", 1) {
			return
		}
		yield("x", 2)
	}).ChunkEntries(2) {
		for _, value := range chunk {
			if value != 1 {
				t.Fatalf("first chunk value: %d", value)
			}
			break
		}
		break
	}

	want := []map[string]int{{"x": 2}, {"y": 3}}
	if !reflect.DeepEqual(chunks, want) {
		t.Fatalf("got %#v, want %#v", chunks, want)
	}
}

func TestSeq2ChunkEntriesZeroSizeYieldsNothing(t *testing.T) {
	for chunk := range Map(map[string]int{"a": 1}).ChunkEntries(0) {
		t.Fatalf("unexpected chunk: %#v", chunk.Collect())
	}
}

func TestSeq2Assign(t *testing.T) {
	base := map[string]int{"a": 1, "b": 2}
	got := Map(base).
		Assign(map[string]int{"b": 20}, map[string]int{"c": 3}).
		Collect()

	want := map[string]int{"a": 1, "b": 20, "c": 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	if !reflect.DeepEqual(base, map[string]int{"a": 1, "b": 2}) {
		t.Fatalf("base mutated: %#v", base)
	}
}

func TestSeq2AssignStopsEarly(t *testing.T) {
	calls := 0
	for key, value := range Seq2[string, int](func(yield func(string, int) bool) {
		calls++
		if !yield("a", 1) {
			return
		}
		calls++
		yield("b", 2)
	}).Assign(map[string]int{"b": 20}) {
		if key != "a" || value != 1 {
			t.Fatalf("entry: %q %d", key, value)
		}
		break
	}
	if calls != 1 {
		t.Fatalf("calls: %d", calls)
	}
}

func TestSeq2AssignLaterMapsReplaceEarlierMaps(t *testing.T) {
	seen := 0
	for key, value := range Map(map[string]int{}).
		Assign(map[string]int{"a": 1}, map[string]int{"a": 2}) {
		seen++
		if key != "a" || value != 2 {
			t.Fatalf("entry: %q %d", key, value)
		}
		break
	}
	if seen != 1 {
		t.Fatalf("seen: %d", seen)
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

func TestSeq2Lookup(t *testing.T) {
	chain := Map(map[string]int{"a": 1})

	if !chain.HasKey("a") {
		t.Fatal("hasKey: false")
	}

	if chain.HasKey("b") {
		t.Fatal("hasKey miss: true")
	}

	if got := chain.ValueOr("a", 9); got != 1 {
		t.Fatalf("valueOr: %d", got)
	}

	if got := chain.ValueOr("b", 9); got != 9 {
		t.Fatalf("valueOr miss: %d", got)
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

	sum := 0
	seq := Map(map[string]int{"a": 1, "b": 2}).ForEach(func(_ string, value int) {
		sum += value
	})
	if sum != 0 {
		t.Fatalf("forEach should be lazy: %d", sum)
	}
	if got := seq.Collect(); !reflect.DeepEqual(got, map[string]int{"a": 1, "b": 2}) {
		t.Fatalf("forEach value: %#v", got)
	}
	if sum != 3 {
		t.Fatalf("forEach: %d", sum)
	}

	sum = 0
	for range Map(map[string]int{"a": 1, "b": 2}).ForEach(func(_ string, value int) {
		sum += value
	}) {
		break
	}
	if sum == 0 || sum == 3 {
		t.Fatalf("forEach should stop after one entry, sum: %d", sum)
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
