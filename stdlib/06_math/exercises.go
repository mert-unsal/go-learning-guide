package math_pkg
// ============================================================
// EXERCISES â€” 06 math
// ============================================================
// Exercise 1:
// Hypotenuse computes the length of the hypotenuse of a right triangle.
// Formula: sqrt(aÂ² + bÂ²)
func Hypotenuse(a, b float64) float64 {
// TODO: math.Sqrt(a*a + b*b)
return 0
}
// Exercise 2:
// IsPowerOfTwo returns true if n is a power of 2 (n > 0).
// Example: 1,2,4,8,16 â†’ true   3,5,6 â†’ false
// Hint: use bits â€” n & (n-1) == 0
func IsPowerOfTwo(n int) bool {
// TODO: n > 0 && n&(n-1) == 0
return false
}
// Exercise 3:
// Clamp returns value clamped to the range [min, max].
// If value < min return min; if value > max return max; else return value.
func ClampEx(value, min, max float64) float64 {
// TODO: math.Max(min, math.Min(max, value))
return 0
}
// Exercise 4:
// RoundToN rounds f to n decimal places.
// Example: RoundToN(3.14159, 2) â†’ 3.14
func RoundToN(f float64, n int) float64 {
// TODO: use math.Pow(10, n) as scale factor, math.Round
return 0
}
// Exercise 5:
// GCD returns the greatest common divisor of a and b (Euclidean algorithm).
func GCD(a, b int) int {
// TODO: for b != 0 { a, b = b, a%b }; return a
return 0
}
// Exercise 6:
// LCM returns the least common multiple of a and b.
// Formula: LCM(a,b) = |a*b| / GCD(a,b)
func LCM(a, b int) int {
// TODO: use GCD
return 0
}