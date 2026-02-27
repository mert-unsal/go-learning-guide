package testing_pkg
import (
"errors"
"strconv"
)
// SOLUTIONS â€” 07 testing
func AddExSolution(a, b int) int { return a + b }
func DivideExSolution(a, b float64) (float64, error) {
if b == 0 {
return 0, errors.New("division by zero")
}
return a / b, nil
}
func MaxExSolution(nums []int) int {
if len(nums) == 0 {
panic("Max called on empty slice")
}
m := nums[0]
for _, v := range nums[1:] {
if v > m {
m = v
}
}
return m
}
func ContainsExSolution(s []int, target int) bool {
for _, v := range s {
if v == target {
return true
}
}
return false
}
func FizzBuzzExSolution(n int) string {
switch {
case n%15 == 0:
return "FizzBuzz"
case n%3 == 0:
return "Fizz"
case n%5 == 0:
return "Buzz"
default:
return strconv.Itoa(n)
}
}