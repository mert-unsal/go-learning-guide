package testing_pkg
import "testing"
// TABLE-DRIVEN TEST with t.Run subtests
func TestAdd(t *testing.T) {
tests := []struct {
name    string
a, b    int
want    int
}{
{"positive", 2, 3, 5},
{"negative", -1, -2, -3},
{"zero", 0, 0, 0},
{"mixed", 10, -4, 6},
}
for _, tt := range tests {
tt := tt // capture for t.Parallel()
t.Run(tt.name, func(t *testing.T) {
t.Parallel() // subtests run concurrently
if got := AddExSolution(tt.a, tt.b); got != tt.want {
t.Errorf("Add(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.want)
}
})
}
}
// TESTING ERROR CASES
func TestDivide(t *testing.T) {
t.Run("happy path", func(t *testing.T) {
got, err := DivideExSolution(10, 2)
if err != nil || got != 5 {
t.Errorf("Divide(10,2) = (%v,%v), want (5,nil)", got, err)
}
})
t.Run("division by zero", func(t *testing.T) {
_, err := DivideExSolution(10, 0)
if err == nil {
t.Error("expected error for division by zero, got nil")
}
})
}
// TESTING PANIC with recover
func TestMaxPanicsOnEmpty(t *testing.T) {
defer func() {
if r := recover(); r == nil {
t.Error("Max(empty) should have panicked")
}
}()
MaxExSolution([]int{}) // should panic
}
func TestMax(t *testing.T) {
tests := []struct{ nums []int; want int }{
{[]int{3, 1, 4, 1, 5, 9}, 9},
{[]int{-5, -1, -3}, -1},
{[]int{42}, 42},
}
for _, tt := range tests {
if got := MaxExSolution(tt.nums); got != tt.want {
t.Errorf("Max(%v) = %d, want %d", tt.nums, got, tt.want)
}
}
}
// BENCHMARK â€” run with: go test -bench=. ./stdlib/07_testing/
func BenchmarkContains(b *testing.B) {
s := make([]int, 1000)
for i := range s { s[i] = i }
b.ResetTimer()
for i := 0; i < b.N; i++ {
ContainsExSolution(s, 999)
}
}
// PARALLEL SUBTESTS
func TestFizzBuzz(t *testing.T) {
tests := []struct{ n int; want string }{
{1, "1"}, {3, "Fizz"}, {5, "Buzz"}, {15, "FizzBuzz"}, {7, "7"},
}
for _, tt := range tests {
tt := tt
t.Run(tt.want, func(t *testing.T) {
t.Parallel()
if got := FizzBuzzExSolution(tt.n); got != tt.want {
t.Errorf("FizzBuzz(%d) = %q, want %q", tt.n, got, tt.want)
}
})
}
}