package builtins
import (
"errors"
"fmt"
)
// SOLUTIONS â€” 03 Builtins
func DeepCopySliceSolution(src []int) []int {
dst := make([]int, len(src))
copy(dst, src)
return dst
}
func DeepCopyMapSolution(src map[string]int) map[string]int {
dst := make(map[string]int, len(src))
for k, v := range src {
dst[k] = v
}
return dst
}
func SafeDivideExSolution(a, b int) (result int, err error) {
defer func() {
if r := recover(); r != nil {
err = errors.New(fmt.Sprintf("panic: %v", r))
result = 0
}
}()
return a / b, nil
}
func FlattenSolution(matrix [][]int) []int {
result := make([]int, 0)
for _, row := range matrix {
result = append(result, row...)
}
return result
}
func UniqueIntsSolution(nums []int) []int {
seen := make(map[int]bool)
result := make([]int, 0)
for _, n := range nums {
if !seen[n] {
seen[n] = true
result = append(result, n)
}
}
return result
}
func ChunkSliceSolution(s []int, n int) [][]int {
if n <= 0 {
return nil
}
var chunks [][]int
for i := 0; i < len(s); i += n {
end := i + n
if end > len(s) {
end = len(s)
}
chunks = append(chunks, s[i:end])
}
return chunks
}