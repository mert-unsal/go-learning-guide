package math_pkg
import (
"math"
"testing"
)
func TestHypotenuse(t *testing.T) {
got := HypotenuseSolution(3, 4)
if math.Abs(got-5.0) > 1e-9 {
t.Errorf("Hypotenuse(3,4) = %f, want 5", got)
}
}
func TestIsPowerOfTwo(t *testing.T) {
tests := []struct{ n int; want bool }{
{1, true}, {2, true}, {4, true}, {16, true},
{0, false}, {3, false}, {6, false}, {-4, false},
}
for _, tt := range tests {
if got := IsPowerOfTwoSolution(tt.n); got != tt.want {
t.Errorf("IsPowerOfTwo(%d) = %v, want %v", tt.n, got, tt.want)
}
}
}
func TestClampEx(t *testing.T) {
tests := []struct{ v, lo, hi, want float64 }{
{5, 0, 10, 5},
{-1, 0, 10, 0},
{15, 0, 10, 10},
}
for _, tt := range tests {
if got := ClampExSolution(tt.v, tt.lo, tt.hi); got != tt.want {
t.Errorf("Clamp(%v,%v,%v) = %v, want %v", tt.v, tt.lo, tt.hi, got, tt.want)
}
}
}
func TestRoundToN(t *testing.T) {
if got := RoundToNSolution(3.14159, 2); math.Abs(got-3.14) > 1e-9 {
t.Errorf("RoundToN(3.14159, 2) = %f, want 3.14", got)
}
if got := RoundToNSolution(2.5, 0); math.Abs(got-3.0) > 1e-9 {
t.Errorf("RoundToN(2.5, 0) = %f, want 3", got)
}
}
func TestGCD(t *testing.T) {
tests := []struct{ a, b, want int }{
{12, 8, 4}, {100, 75, 25}, {7, 3, 1}, {0, 5, 5},
}
for _, tt := range tests {
if got := GCDSolution(tt.a, tt.b); got != tt.want {
t.Errorf("GCD(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.want)
}
}
}
func TestLCM(t *testing.T) {
tests := []struct{ a, b, want int }{
{4, 6, 12}, {3, 5, 15}, {7, 7, 7},
}
for _, tt := range tests {
if got := LCMSolution(tt.a, tt.b); got != tt.want {
t.Errorf("LCM(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.want)
}
}
}