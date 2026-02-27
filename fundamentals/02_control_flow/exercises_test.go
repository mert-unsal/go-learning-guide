package control_flow
import "testing"
func TestFizzBuzzSwitch(t *testing.T) {
tests := []struct {
n    int
want string
}{
{1, "1"}, {3, "Fizz"}, {5, "Buzz"}, {15, "FizzBuzz"},
{9, "Fizz"}, {10, "Buzz"}, {30, "FizzBuzz"}, {7, "7"},
}
for _, tt := range tests {
got := FizzBuzzSwitchSolution(tt.n)
if got != tt.want {
t.Errorf("FizzBuzzSwitch(%d) = %q, want %q", tt.n, got, tt.want)
}
}
}
func TestSumTo(t *testing.T) {
tests := []struct{ n, want int }{
{1, 1}, {5, 15}, {10, 55}, {100, 5050},
}
for _, tt := range tests {
got := SumToSolution(tt.n)
if got != tt.want {
t.Errorf("SumTo(%d) = %d, want %d", tt.n, got, tt.want)
}
}
}
func TestCountVowels(t *testing.T) {
tests := []struct {
s    string
want int
}{
{"hello", 2}, {"AEIOU", 5}, {"rhythm", 0}, {"Go is fun", 3},
}
for _, tt := range tests {
got := CountVowelsSolution(tt.s)
if got != tt.want {
t.Errorf("CountVowels(%q) = %d, want %d", tt.s, got, tt.want)
}
}
}
func TestIsPrime(t *testing.T) {
tests := []struct {
n    int
want bool
}{
{2, true}, {3, true}, {4, false}, {17, true},
{1, false}, {0, false}, {97, true}, {100, false},
}
for _, tt := range tests {
got := IsPrimeSolution(tt.n)
if got != tt.want {
t.Errorf("IsPrime(%d) = %v, want %v", tt.n, got, tt.want)
}
}
}
func TestDeferOrder(t *testing.T) {
got := DeferOrderSolution()
want := []string{"third", "second", "first"}
if len(got) != 3 {
t.Fatalf("DeferOrder() returned %d elements, want 3", len(got))
}
for i := range want {
if got[i] != want[i] {
t.Errorf("DeferOrder()[%d] = %q, want %q", i, got[i], want[i])
}
}
}