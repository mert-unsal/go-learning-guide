package math_pkg
import "math"
// SOLUTIONS â€” 06 math
func HypotenuseSolution(a, b float64) float64 {
return math.Sqrt(a*a + b*b)
}
func IsPowerOfTwoSolution(n int) bool {
return n > 0 && n&(n-1) == 0
}
func ClampExSolution(value, min, max float64) float64 {
return math.Max(min, math.Min(max, value))
}
func RoundToNSolution(f float64, n int) float64 {
scale := math.Pow(10, float64(n))
return math.Round(f*scale) / scale
}
func GCDSolution(a, b int) int {
for b != 0 {
a, b = b, a%b
}
return a
}
func LCMSolution(a, b int) int {
if a == 0 || b == 0 {
return 0
}
// Use abs to handle negative inputs
product := a * b
if product < 0 {
product = -product
}
return product / GCDSolution(a, b)
}