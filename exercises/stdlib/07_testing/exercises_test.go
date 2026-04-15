package testing_pkg

import (
"errors"
"fmt"
"strings"
"testing"
)

// Exercise 1: Table-driven tests with t.Run subtests

func TestAdd(t *testing.T) {
tests := []struct {
name string
a, b int
want int
}{
{"positive", 2, 3, 5},
{"negative", -1, -2, -3},
{"zero", 0, 0, 0},
{"mixed", 10, -4, 6},
}
for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
t.Parallel()
if got := AddEx(tt.a, tt.b); got != tt.want {
t.Errorf("❌ Add(%d, %d) = %d, want %d\n\t\t"+
"Hint: return a + b. This test demonstrates table-driven "+
"tests with t.Run subtests and t.Parallel()",
tt.a, tt.b, got, tt.want)
} else {
t.Logf("✅ Add(%d, %d) = %d", tt.a, tt.b, got)
}
})
}
}

// Exercise 2: Testing error cases

func TestDivide(t *testing.T) {
t.Run("happy_path", func(t *testing.T) {
got, err := DivideEx(10, 2)
if err != nil {
t.Fatalf("❌ unexpected error: %v", err)
}
if got != 5 {
t.Errorf("❌ Divide(10, 2) = %v, want 5\n\t\t"+
"Hint: return a/b when b != 0", got)
} else {
t.Logf("✅ Divide(10, 2) = %v", got)
}
})

t.Run("division_by_zero", func(t *testing.T) {
_, err := DivideEx(10, 0)
if err == nil {
t.Error("❌ Divide(10, 0) should return error\n\t\t" +
"Hint: if b == 0 { return 0, errors.New(\"division by zero\") }")
} else {
t.Logf("✅ Divide(10, 0) returned error: %v", err)
}
})
}

// Exercise 3: Testing panic behavior

func TestMaxPanicsOnEmpty(t *testing.T) {
defer func() {
if r := recover(); r == nil {
t.Error("❌ Max(empty) should panic\n\t\t" +
"Hint: if len(nums) == 0 { panic(\"empty slice\") }. " +
"Test panics with defer + recover")
} else {
t.Logf("✅ Max(empty) panicked: %v", r)
}
}()
MaxEx([]int{})
}

func TestMax(t *testing.T) {
tests := []struct {
nums []int
want int
}{
{[]int{3, 1, 4, 1, 5, 9}, 9},
{[]int{-5, -1, -3}, -1},
{[]int{42}, 42},
}
for _, tt := range tests {
t.Run(fmt.Sprintf("Max(%v)", tt.nums), func(t *testing.T) {
if got := MaxEx(tt.nums); got != tt.want {
t.Errorf("❌ Max(%v) = %d, want %d", tt.nums, got, tt.want)
} else {
t.Logf("✅ Max(%v) = %d", tt.nums, got)
}
})
}
}

// Exercise 4: Benchmark

func BenchmarkContains(b *testing.B) {
s := make([]int, 1000)
for i := range s {
s[i] = i
}
b.ResetTimer()
for i := 0; i < b.N; i++ {
ContainsEx(s, 999)
}
}

func TestContains(t *testing.T) {
s := []int{1, 2, 3, 4, 5}
if !ContainsEx(s, 3) {
t.Error("❌ Contains([1..5], 3) = false, want true\n\t\t" +
"Hint: for _, v := range s { if v == target { return true } }; return false")
} else {
t.Logf("✅ Contains([1..5], 3) = true")
}
if ContainsEx(s, 99) {
t.Error("❌ Contains([1..5], 99) = true, want false")
} else {
t.Logf("✅ Contains([1..5], 99) = false")
}
}

// Exercise 5: Parallel subtests

func TestFizzBuzz(t *testing.T) {
tests := []struct {
n    int
want string
}{
{1, "1"}, {3, "Fizz"}, {5, "Buzz"}, {15, "FizzBuzz"}, {7, "7"},
}
for _, tt := range tests {
t.Run(tt.want, func(t *testing.T) {
t.Parallel()
if got := FizzBuzzEx(tt.n); got != tt.want {
t.Errorf("❌ FizzBuzz(%d) = %q, want %q\n\t\t"+
"Hint: check div 15 first, then 3, then 5, then default. "+
"t.Parallel() makes subtests run concurrently — good for "+
"detecting shared-state bugs",
tt.n, got, tt.want)
} else {
t.Logf("✅ FizzBuzz(%d) = %q", tt.n, got)
}
})
}
}

// Exercise 6: Golden value tests

func TestReverse(t *testing.T) {
tests := []struct {
input string
want  string
}{
{"hello", "olleh"},
{"", ""},
{"a", "a"},
{"racecar", "racecar"},
{"Hello, 世界", "界世 ,olleH"},
}
for _, tt := range tests {
t.Run(tt.input, func(t *testing.T) {
got := ReverseEx(tt.input)
if got != tt.want {
t.Errorf("❌ Reverse(%q) = %q, want %q\n\t\t"+
"Hint: Convert to []rune, reverse, convert back. "+
"This tests multi-byte rune handling (世界 = 2 runes, 6 bytes)",
tt.input, got, tt.want)
} else {
t.Logf("✅ Reverse(%q) = %q", tt.input, got)
}
})
}
}

// Exercise 7: Test helper function

func assertSorted(t *testing.T, nums []int, want bool) {
t.Helper()
got := IsSortedEx(nums)
if got != want {
t.Errorf("❌ IsSorted(%v) = %v, want %v\n\t\t"+
"Hint: for i := 1; i < len(nums); i++ { if nums[i] < nums[i-1] { return false } }. "+
"t.Helper() marks this as a helper — errors show the CALLER's line, not this one",
nums, got, want)
} else {
t.Logf("✅ IsSorted(%v) = %v", nums, got)
}
}

func TestIsSorted(t *testing.T) {
assertSorted(t, []int{1, 2, 3, 4, 5}, true)
assertSorted(t, []int{5, 3, 1}, false)
assertSorted(t, []int{}, true)
assertSorted(t, []int{42}, true)
assertSorted(t, []int{1, 1, 1}, true)
}

// Exercise 8: Testing error types with errors.As

func TestParseKeyValue(t *testing.T) {
t.Run("valid", func(t *testing.T) {
k, v, err := ParseKeyValue("host=localhost")
if err != nil {
t.Fatalf("❌ unexpected error: %v", err)
}
if k != "host" || v != "localhost" {
t.Errorf("❌ ParseKeyValue(\"host=localhost\") = (%q, %q), want (\"host\", \"localhost\")\n\t\t"+
"Hint: strings.SplitN(s, \"=\", 2)", k, v)
} else {
t.Logf("✅ ParseKeyValue(\"host=localhost\") = (%q, %q)", k, v)
}
})

t.Run("missing_equals", func(t *testing.T) {
_, _, err := ParseKeyValue("invalid")
if err == nil {
t.Fatal("❌ expected error for missing '='\n\t\t" +
"Hint: return &ParseError{Input: s, Message: \"missing '='\"}")
}
var pe *ParseError
if !errors.As(err, &pe) {
t.Errorf("❌ error type = %T, want *ParseError\n\t\t"+
"Hint: errors.As checks the error chain for a matching type. "+
"Return &ParseError{} (pointer) to implement the error interface",
err)
} else {
t.Logf("✅ ParseError: %v", pe)
}
})

t.Run("empty_key", func(t *testing.T) {
_, _, err := ParseKeyValue("=value")
if err == nil {
t.Error("❌ expected error for empty key")
} else {
t.Logf("✅ empty key returns error: %v", err)
}
})
}

// Exercise 9: Testing that original slice is not modified

func TestSortStrings(t *testing.T) {
input := []string{"banana", "apple", "cherry"}
original := make([]string, len(input))
copy(original, input)

got := SortStringsEx(input)
want := []string{"apple", "banana", "cherry"}

if got == nil {
t.Fatal("❌ SortStringsEx returned nil\n\t\t" +
"Hint: Make a copy with copy() or append([]string{}, input...), " +
"then sort.Strings(copied)")
}

for i, w := range want {
if i >= len(got) || got[i] != w {
t.Errorf("❌ SortStrings result[%d] = %q, want %q", i, got[i], w)
}
}

// Verify original was NOT modified
for i, o := range original {
if input[i] != o {
t.Errorf("❌ original input[%d] was modified: %q -> %q\n\t\t"+
"Hint: sort.Strings modifies in-place. You MUST copy first. "+
"This tests a critical Go gotcha: slices share backing arrays",
i, o, input[i])
}
}
t.Logf("✅ SortStrings: sorted correctly, original unchanged")
}

// Exercise 10: Testing time-dependent behavior with injection

func TestRetry(t *testing.T) {
t.Run("succeeds_first_try", func(t *testing.T) {
err := Retry(3, func() {}, func() error { return nil })
if err != nil {
t.Errorf("❌ Retry returned error on success: %v\n\t\t"+
"Hint: call fn(), if nil return nil. No sleep before first attempt", err)
} else {
t.Logf("✅ Retry: succeeded on first try")
}
})

t.Run("succeeds_on_third_try", func(t *testing.T) {
attempt := 0
sleepCount := 0
err := Retry(3, func() { sleepCount++ }, func() error {
attempt++
if attempt < 3 {
return fmt.Errorf("attempt %d failed", attempt)
}
return nil
})
if err != nil {
t.Errorf("❌ Retry should succeed on attempt 3: %v\n\t\t"+
"Hint: loop maxAttempts times. Call sleep() between retries (not before first). "+
"The sleep func is injected — in real code it's time.Sleep, in tests it's a no-op",
err)
} else {
t.Logf("✅ Retry: succeeded on attempt 3")
}
if sleepCount != 2 {
t.Errorf("❌ sleep called %d times, want 2 (between retries only)", sleepCount)
}
})

t.Run("all_fail", func(t *testing.T) {
err := Retry(3, func() {}, func() error {
return fmt.Errorf("always fails")
})
if err == nil {
t.Error("❌ Retry should return last error when all attempts fail")
} else {
t.Logf("✅ Retry: returned error after 3 failures: %v", err)
}
})
}

// Exercise 11: Map-based test cases

func TestHTTPStatusText(t *testing.T) {
cases := map[int]string{
200: "OK",
201: "Created",
400: "Bad Request",
404: "Not Found",
500: "Internal Server Error",
}
for code, want := range cases {
t.Run(fmt.Sprintf("%d", code), func(t *testing.T) {
got := HTTPStatusText(code)
if got != want {
t.Errorf("❌ HTTPStatusText(%d) = %q, want %q\n\t\t"+
"Hint: map[int]string{200: \"OK\", ...}[code]. "+
"Map-based tests are concise for lookup-table functions",
code, got, want)
} else {
t.Logf("✅ HTTPStatusText(%d) = %q", code, got)
}
})
}

// Unknown code
if got := HTTPStatusText(999); got != "" {
t.Errorf("❌ HTTPStatusText(999) = %q, want \"\"", got)
} else {
t.Logf("✅ HTTPStatusText(999) = \"\" (unknown)")
}
}

// Exercise 12: Testing with function injection

func TestTransform(t *testing.T) {
t.Run("upper", func(t *testing.T) {
input := []string{"hello", "world"}
got := Transform(input, strings.ToUpper)
if got == nil || len(got) != 2 {
t.Fatal("❌ Transform returned nil or wrong length\n\t\t" +
"Hint: result := make([]string, len(input)); for i, s := range input { result[i] = fn(s) }")
}
if got[0] != "HELLO" || got[1] != "WORLD" {
t.Errorf("❌ Transform(upper) = %v, want [HELLO WORLD]\n\t\t"+
"Hint: fn is injected — any func(string) string works. "+
"This pattern makes functions testable without mocking",
got)
} else {
t.Logf("✅ Transform(upper) = %v", got)
}
})

t.Run("trim", func(t *testing.T) {
input := []string{"  hello  ", "\tworld\t"}
got := Transform(input, strings.TrimSpace)
if got == nil || got[0] != "hello" || got[1] != "world" {
t.Errorf("❌ Transform(trim) = %v, want [hello world]", got)
} else {
t.Logf("✅ Transform(trim) = %v", got)
}
})

t.Run("nil_input", func(t *testing.T) {
got := Transform(nil, strings.ToUpper)
if got != nil && len(got) != 0 {
t.Errorf("❌ Transform(nil) = %v, want nil or empty", got)
} else {
t.Logf("✅ Transform(nil) = nil/empty")
}
})
}
