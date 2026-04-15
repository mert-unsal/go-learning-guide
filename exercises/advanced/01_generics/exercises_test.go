package generics

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: Min
// ────────────────────────────────────────────────────────────

func TestMin(t *testing.T) {
	if got := Min(3, 5); got != 3 {
		t.Errorf("❌ Min(3, 5) = %d, want 3\n\t\t"+
			"Hint: if a < b { return a }; return b. "+
			"cmp.Ordered covers int, float64, string — anything with < operator", got)
	} else {
		t.Logf("✅ Min(3, 5) = %d", got)
	}

	if got := Min("banana", "apple"); got != "apple" {
		t.Errorf("❌ Min(\"banana\", \"apple\") = %q, want \"apple\"", got)
	} else {
		t.Logf("✅ Min(\"banana\", \"apple\") = %q", got)
	}

	if got := Min(3.14, 2.71); got != 2.71 {
		t.Errorf("❌ Min(3.14, 2.71) = %v, want 2.71", got)
	} else {
		t.Logf("✅ Min(3.14, 2.71) = %v", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: Contains
// ────────────────────────────────────────────────────────────

func TestContains(t *testing.T) {
	if !Contains([]int{1, 2, 3}, 2) {
		t.Error("❌ Contains([1,2,3], 2) = false, want true\n\t\t" +
			"Hint: comparable constraint allows ==. for _, v := range s { if v == target { ... } }")
	} else {
		t.Logf("✅ Contains([1,2,3], 2) = true")
	}

	if Contains([]string{"a", "b"}, "c") {
		t.Error("❌ Contains([a,b], c) = true, want false")
	} else {
		t.Logf("✅ Contains([a,b], c) = false")
	}

	if Contains[int](nil, 1) {
		t.Error("❌ Contains(nil, 1) = true, want false")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: Map
// ────────────────────────────────────────────────────────────

func TestMap(t *testing.T) {
	// int → string
	got := Map([]int{1, 2, 3}, func(n int) string {
		return fmt.Sprintf("%d!", n)
	})
	if got == nil || len(got) != 3 {
		t.Fatal("❌ Map returned nil or wrong length\n\t\t" +
			"Hint: result := make([]U, len(s)); for i, v := range s { result[i] = fn(v) }; return result")
	}
	if got[0] != "1!" || got[2] != "3!" {
		t.Errorf("❌ Map = %v, want [1! 2! 3!]\n\t\t"+
			"Hint: Two type params: Map[T any, U any]. T is input, U is output", got)
	} else {
		t.Logf("✅ Map([1,2,3], n→n!) = %v", got)
	}

	// string → int
	lengths := Map([]string{"hi", "hello"}, func(s string) int { return len(s) })
	if lengths[0] != 2 || lengths[1] != 5 {
		t.Errorf("❌ Map(lengths) = %v, want [2 5]", lengths)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: Filter
// ────────────────────────────────────────────────────────────

func TestFilter(t *testing.T) {
	evens := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	if len(evens) != 3 || evens[0] != 2 || evens[1] != 4 || evens[2] != 6 {
		t.Errorf("❌ Filter(evens) = %v, want [2 4 6]\n\t\t"+
			"Hint: var result []T; for _, v := range s { if keep(v) { result = append(result, v) } }",
			evens)
	} else {
		t.Logf("✅ Filter(evens) = %v", evens)
	}

	long := Filter([]string{"go", "rust", "c", "python"}, func(s string) bool { return len(s) > 2 })
	if len(long) != 3 {
		t.Errorf("❌ Filter(long strings) = %v, want [rust python]... wait, [go? no...] %v", long, long)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: Reduce
// ────────────────────────────────────────────────────────────

func TestReduce(t *testing.T) {
	sum := Reduce([]int{1, 2, 3, 4}, 0, func(acc, n int) int { return acc + n })
	if sum != 10 {
		t.Errorf("❌ Reduce(sum) = %d, want 10\n\t\t"+
			"Hint: acc := initial; for _, v := range s { acc = fn(acc, v) }; return acc",
			sum)
	} else {
		t.Logf("✅ Reduce(sum [1..4]) = %d", sum)
	}

	// Reduce strings into concatenation
	joined := Reduce([]string{"a", "b", "c"}, "", func(acc, s string) string {
		if acc == "" {
			return s
		}
		return acc + "," + s
	})
	if joined != "a,b,c" {
		t.Errorf("❌ Reduce(join) = %q, want \"a,b,c\"", joined)
	} else {
		t.Logf("✅ Reduce(join) = %q", joined)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: Keys and Values
// ────────────────────────────────────────────────────────────

func TestKeysValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := Keys(m)
	if keys == nil || len(keys) != 3 {
		t.Fatal("❌ Keys returned nil or wrong length\n\t\t" +
			"Hint: result := make([]K, 0, len(m)); for k := range m { result = append(result, k) }")
	}
	sort.Strings(keys) // map iteration order is random
	if keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
		t.Errorf("❌ Keys = %v, want [a b c]", keys)
	} else {
		t.Logf("✅ Keys = %v", keys)
	}

	vals := Values(m)
	if vals == nil || len(vals) != 3 {
		t.Fatal("❌ Values returned nil or wrong length")
	}
	sort.Ints(vals)
	if vals[0] != 1 || vals[2] != 3 {
		t.Errorf("❌ Values = %v, want [1 2 3]", vals)
	} else {
		t.Logf("✅ Values = %v", vals)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: Stack
// ────────────────────────────────────────────────────────────

func TestStack(t *testing.T) {
	var s Stack[int]

	if s.Len() != 0 {
		t.Errorf("❌ empty stack Len = %d, want 0", s.Len())
	}

	s.Push(10)
	s.Push(20)
	s.Push(30)

	if s.Len() != 3 {
		t.Errorf("❌ Len = %d after 3 pushes, want 3\n\t\t"+
			"Hint: Push appends, Pop removes last, Peek reads last. "+
			"Return (zero, false) when empty",
			s.Len())
	}

	if v, ok := s.Peek(); !ok || v != 30 {
		t.Errorf("❌ Peek = (%d, %v), want (30, true)", v, ok)
	} else {
		t.Logf("✅ Peek = %d", v)
	}

	if v, ok := s.Pop(); !ok || v != 30 {
		t.Errorf("❌ Pop = (%d, %v), want (30, true)", v, ok)
	} else {
		t.Logf("✅ Pop = %d", v)
	}

	if s.Len() != 2 {
		t.Errorf("❌ Len after pop = %d, want 2", s.Len())
	}

	// Pop until empty
	s.Pop()
	s.Pop()
	if _, ok := s.Pop(); ok {
		t.Error("❌ Pop on empty stack should return (_, false)")
	} else {
		t.Logf("✅ Pop on empty returns false")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: Pair
// ────────────────────────────────────────────────────────────

func TestPair(t *testing.T) {
	p := Pair[string, int]{Key: "age", Value: 30}
	swapped := p.Swap()

	if swapped.Key != 30 || swapped.Value != "age" {
		t.Errorf("❌ Swap = {%v, %v}, want {30, \"age\"}\n\t\t"+
			"Hint: return Pair[V, K]{Key: p.Value, Value: p.Key}",
			swapped.Key, swapped.Value)
	} else {
		t.Logf("✅ Pair{\"age\", 30}.Swap() = {%v, %q}", swapped.Key, swapped.Value)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: MaxBy
// ────────────────────────────────────────────────────────────

func TestMaxBy(t *testing.T) {
	words := []string{"go", "rust", "c", "python"}
	longest := MaxBy(words, func(s string) int { return len(s) })
	if longest != "python" {
		t.Errorf("❌ MaxBy(longest) = %q, want \"python\"\n\t\t"+
			"Hint: track the element with highest fn(elem) value. "+
			"Panic if slice is empty",
			longest)
	} else {
		t.Logf("✅ MaxBy(longest word) = %q", longest)
	}
}

func TestMaxByPanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("❌ MaxBy(empty) should panic")
		} else {
			t.Logf("✅ MaxBy(empty) panicked: %v", r)
		}
	}()
	MaxBy([]int{}, func(n int) int { return n })
}

// ────────────────────────────────────────────────────────────
// Exercise 10: GroupBy
// ────────────────────────────────────────────────────────────

func TestGroupBy(t *testing.T) {
	words := []string{"go", "rust", "c", "java", "py"}
	groups := GroupBy(words, func(s string) int { return len(s) })
	if groups == nil {
		t.Fatal("❌ GroupBy returned nil\n\t\t" +
			"Hint: result := make(map[K][]T); for _, v := range s { k := keyFn(v); result[k] = append(result[k], v) }")
	}
	if len(groups[2]) != 2 { // "go", "py"
		t.Errorf("❌ groups[2] = %v, want [go py]", groups[2])
	} else {
		t.Logf("✅ groups[2] = %v", groups[2])
	}
	if len(groups[4]) != 2 { // "rust", "java"
		t.Errorf("❌ groups[4] = %v, want [rust java]", groups[4])
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: Result
// ────────────────────────────────────────────────────────────

func TestResult(t *testing.T) {
	okResult := Ok[string]("hello")
	if !okResult.IsOk() {
		t.Error("❌ Ok(\"hello\").IsOk() = false\n\t\t" +
			"Hint: Ok sets ok=true, value=v. Err sets ok=false, err=e")
	}
	val, err := okResult.Unwrap()
	if err != nil || val != "hello" {
		t.Errorf("❌ Unwrap = (%q, %v), want (\"hello\", nil)\n\t\t"+
			"Hint: if r.ok { return r.value, nil } else { return zero, r.err }",
			val, err)
	} else {
		t.Logf("✅ Ok(\"hello\").Unwrap() = %q", val)
	}

	errResult := Err[int](errors.New("boom"))
	if errResult.IsOk() {
		t.Error("❌ Err(\"boom\").IsOk() = true, want false")
	}
	_, err = errResult.Unwrap()
	if err == nil || err.Error() != "boom" {
		t.Errorf("❌ Err.Unwrap error = %v, want \"boom\"", err)
	} else {
		t.Logf("✅ Err(\"boom\").Unwrap() returns error: %v", err)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: Set
// ────────────────────────────────────────────────────────────

func TestSet(t *testing.T) {
	s := NewSet[string]()
	s.Add("go")
	s.Add("rust")
	s.Add("go") // duplicate

	if s.Len() != 2 {
		t.Errorf("❌ Set.Len = %d after adding go,rust,go — want 2\n\t\t"+
			"Hint: map[T]struct{}{} — struct{} is zero bytes. "+
			"This is the idiomatic Go set implementation",
			s.Len())
	} else {
		t.Logf("✅ Set.Len = 2")
	}

	if !s.Has("go") {
		t.Error("❌ Has(\"go\") = false")
	}
	if s.Has("python") {
		t.Error("❌ Has(\"python\") = true")
	}

	s.Remove("rust")
	if s.Has("rust") || s.Len() != 1 {
		t.Error("❌ Remove(\"rust\") didn't work")
	} else {
		t.Logf("✅ Remove(\"rust\") OK, Len = %d", s.Len())
	}
}

func TestSetUnion(t *testing.T) {
	a := NewSet[int]()
	a.Add(1)
	a.Add(2)
	b := NewSet[int]()
	b.Add(2)
	b.Add(3)

	u := a.Union(b)
	if u.Len() != 3 || !u.Has(1) || !u.Has(2) || !u.Has(3) {
		t.Errorf("❌ Union = len %d, want 3 with {1,2,3}\n\t\t"+
			"Hint: create new set, add all from s, add all from other",
			u.Len())
	} else {
		t.Logf("✅ Union({1,2}, {2,3}) = {1,2,3}")
	}
}

func TestSetIntersection(t *testing.T) {
	a := NewSet[int]()
	a.Add(1)
	a.Add(2)
	a.Add(3)
	b := NewSet[int]()
	b.Add(2)
	b.Add(3)
	b.Add(4)

	inter := a.Intersection(b)
	if inter.Len() != 2 || !inter.Has(2) || !inter.Has(3) {
		t.Errorf("❌ Intersection = len %d, want 2 with {2,3}\n\t\t"+
			"Hint: for each in s, if other.Has(it), add to result",
			inter.Len())
	} else {
		t.Logf("✅ Intersection({1,2,3}, {2,3,4}) = {2,3}")
	}
}

// Keep imports used
var _ = strings.ToUpper
