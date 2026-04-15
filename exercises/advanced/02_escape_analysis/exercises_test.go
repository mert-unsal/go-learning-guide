package escape

import (
	"math"
	"strings"
	"testing"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: StackOnly
// ────────────────────────────────────────────────────────────

func TestStackOnly(t *testing.T) {
	if got := StackOnly(3, 4); got != 7 {
		t.Errorf("❌ StackOnly(3, 4) = %d, want 7\n\t\t"+
			"Hint: return a + b. Verify with go build -gcflags='-m': nothing should escape",
			got)
	} else {
		t.Logf("✅ StackOnly(3, 4) = %d", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: EscapeToHeap
// ────────────────────────────────────────────────────────────

func TestEscapeToHeap(t *testing.T) {
	p := EscapeToHeap(3, 4)
	if p == nil {
		t.Fatal("❌ EscapeToHeap returned nil\n\t\t" +
			"Hint: result := a + b; return &result. " +
			"go build -gcflags='-m' will show 'moved to heap: result'")
	}
	if *p != 7 {
		t.Errorf("❌ *EscapeToHeap(3, 4) = %d, want 7", *p)
	} else {
		t.Logf("✅ *EscapeToHeap(3, 4) = %d", *p)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: NoEscapeSlice
// ────────────────────────────────────────────────────────────

func TestNoEscapeSlice(t *testing.T) {
	if got := NoEscapeSlice(); got != 6 {
		t.Errorf("❌ NoEscapeSlice() = %d, want 6\n\t\t"+
			"Hint: s := []int{1, 2, 3}; return s[0]+s[1]+s[2]. "+
			"Small slices that don't escape stay on stack",
			got)
	} else {
		t.Logf("✅ NoEscapeSlice = %d (slice stayed on stack)", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: EscapeSlice
// ────────────────────────────────────────────────────────────

func TestEscapeSlice(t *testing.T) {
	s := EscapeSlice(3)
	if s == nil || len(s) != 3 {
		t.Fatal("❌ EscapeSlice(3) returned nil or wrong length\n\t\t" +
			"Hint: s := make([]int, n); for i := range s { s[i] = i }; return s. " +
			"The backing array escapes because it's returned")
	}
	if s[0] != 0 || s[1] != 1 || s[2] != 2 {
		t.Errorf("❌ EscapeSlice(3) = %v, want [0 1 2]", s)
	} else {
		t.Logf("✅ EscapeSlice(3) = %v", s)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: InterfaceEscape
// ────────────────────────────────────────────────────────────

func TestInterfaceEscape(t *testing.T) {
	v := InterfaceEscape(42)
	if v == nil {
		t.Fatal("❌ InterfaceEscape(42) returned nil\n\t\t" +
			"Hint: return n. The value is boxed into an interface, causing escape")
	}
	if v.(int) != 42 {
		t.Errorf("❌ InterfaceEscape(42) = %v, want 42", v)
	} else {
		t.Logf("✅ InterfaceEscape(42) = %v", v)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: ClosureCapture
// ────────────────────────────────────────────────────────────

func TestClosureCapture(t *testing.T) {
	fn := ClosureCapture()
	if fn == nil {
		t.Fatal("❌ ClosureCapture returned nil\n\t\t" +
			"Hint: counter := 0; return func() int { counter++; return counter }. " +
			"counter escapes to heap because the closure outlives the stack frame")
	}
	if fn() != 1 || fn() != 2 || fn() != 3 {
		t.Error("❌ closure should return 1, 2, 3 on successive calls")
	} else {
		t.Logf("✅ ClosureCapture increments: 1, 2, 3")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: PreallocateVsAppend
// ────────────────────────────────────────────────────────────

func TestPreallocateVsAppend(t *testing.T) {
	s := PreallocateVsAppend(5)
	if s == nil || len(s) != 5 {
		t.Fatal("❌ PreallocateVsAppend(5) returned nil or wrong length\n\t\t" +
			"Hint: s := make([]int, n); for i := 0; i < n; i++ { s[i] = i * i }; return s")
	}
	expected := []int{0, 1, 4, 9, 16}
	for i, want := range expected {
		if s[i] != want {
			t.Errorf("❌ s[%d] = %d, want %d", i, s[i], want)
		}
	}
	t.Logf("✅ PreallocateVsAppend(5) = %v", s)
}

// ────────────────────────────────────────────────────────────
// Exercise 8: JoinStrings
// ────────────────────────────────────────────────────────────

func TestJoinStrings(t *testing.T) {
	got := JoinStrings([]string{"a", "b", "c"}, ",")
	if got != "a,b,c" {
		t.Errorf("❌ JoinStrings = %q, want \"a,b,c\"\n\t\t"+
			"Hint: var b strings.Builder; for i, s := range parts { "+
			"if i > 0 { b.WriteString(sep) }; b.WriteString(s) }; return b.String()",
			got)
	} else {
		t.Logf("✅ JoinStrings = %q", got)
	}

	if got := JoinStrings(nil, ","); got != "" {
		t.Errorf("❌ JoinStrings(nil) = %q, want \"\"", got)
	}

	if got := JoinStrings([]string{"only"}, ","); got != "only" {
		t.Errorf("❌ JoinStrings([only]) = %q, want \"only\"", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: Distance
// ────────────────────────────────────────────────────────────

func TestDistance(t *testing.T) {
	got := Distance(3, 4)
	if math.Abs(got-5.0) > 0.001 {
		t.Errorf("❌ Distance(3, 4) = %v, want 5.0\n\t\t"+
			"Hint: p := Point{X: x, Y: y}; return math.Sqrt(p.X*p.X + p.Y*p.Y). "+
			"Point is passed by value (stays on stack)",
			got)
	} else {
		t.Logf("✅ Distance(3, 4) = %v", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 10: FormatRecord
// ────────────────────────────────────────────────────────────

func TestFormatRecord(t *testing.T) {
	got := FormatRecord("users", 42)
	if got != "users=42" {
		t.Errorf("❌ FormatRecord(\"users\", 42) = %q, want \"users=42\"\n\t\t"+
			"Hint: buf := bufferPool.Get().(*bytes.Buffer); buf.Reset(); "+
			"fmt.Fprintf(buf, \"%%s=%%d\", name, value); result := buf.String(); "+
			"bufferPool.Put(buf); return result",
			got)
	} else {
		t.Logf("✅ FormatRecord(users, 42) = %q", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: RemoveAt
// ────────────────────────────────────────────────────────────

func TestRemoveAt(t *testing.T) {
	s := []int{10, 20, 30, 40, 50}
	got := RemoveAt(s, 2)
	if len(got) != 4 || got[0] != 10 || got[1] != 20 || got[2] != 40 || got[3] != 50 {
		t.Errorf("❌ RemoveAt([10,20,30,40,50], 2) = %v, want [10 20 40 50]\n\t\t"+
			"Hint: return append(s[:i], s[i+1:]...)",
			got)
	} else {
		t.Logf("✅ RemoveAt(index 2) = %v", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: BytesEqualString
// ────────────────────────────────────────────────────────────

func TestBytesEqualString(t *testing.T) {
	if !BytesEqualString([]byte("hello"), "hello") {
		t.Error("❌ BytesEqualString(hello, hello) = false\n\t\t" +
			"Hint: if len(b) != len(s) { return false }; " +
			"for i := range b { if b[i] != s[i] { return false } }; return true")
	} else {
		t.Log("✅ BytesEqualString matches")
	}

	if BytesEqualString([]byte("hello"), "world") {
		t.Error("❌ BytesEqualString(hello, world) = true")
	}

	if BytesEqualString([]byte("hi"), "hello") {
		t.Error("❌ BytesEqualString(hi, hello) = true (different lengths)")
	}
}

// Keep imports used
var _ = strings.Builder{}
