package testing_pkg
// ============================================================
// EXERCISES â€” 07 testing
// ============================================================
// This package teaches Go testing patterns by BEING tested.
// All functions below have known, testable behaviors.
// Exercise 1:
// Add returns the sum of a and b.
// The test demonstrates the TABLE-DRIVEN pattern with t.Run subtests.
func AddEx(a, b int) int {
// TODO: return a + b
return 0
}
// Exercise 2:
// Divide returns a/b and an error if b==0.
// The test demonstrates testing BOTH the happy path AND error cases.
func DivideEx(a, b float64) (float64, error) {
// TODO: return error if b==0, else a/b
return 0, nil
}
// Exercise 3:
// Max returns the largest value in nums. Panics if nums is empty.
// The test demonstrates testing PANIC behavior with recover.
func MaxEx(nums []int) int {
// TODO: panic if empty, else find max
return 0
}
// Exercise 4:
// Contains reports whether target exists in s.
// The test demonstrates a BENCHMARK using b.N.
func ContainsEx(s []int, target int) bool {
// TODO: linear search
return false
}
// Exercise 5:
// FizzBuzz returns the FizzBuzz string for n.
// The test demonstrates a PARALLEL subtest with t.Parallel().
func FizzBuzzEx(n int) string {
// TODO: "FizzBuzz", "Fizz", "Buzz", or strconv.Itoa(n)
return ""
}