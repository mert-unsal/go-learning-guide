package patterns
import "testing"
// These tests verify every template in templates.go compiles and produces
// correct output on representative inputs.
// 1. Binary Search
func TestBinarySearchExact(t *testing.T) {
nums := []int{1, 3, 5, 7, 9, 11}
if got := BinarySearchExact(nums, 7); got != 3 {
t.Errorf("BinarySearchExact(7) = %d, want 3", got)
}
if got := BinarySearchExact(nums, 4); got != -1 {
t.Errorf("BinarySearchExact(missing) = %d, want -1", got)
}
}
func TestBinarySearchLeftBound(t *testing.T) {
nums := []int{1, 2, 2, 2, 3}
if got := BinarySearchLeftBound(nums, 2); got != 1 {
t.Errorf("BinarySearchLeftBound(2) = %d, want 1 (leftmost)", got)
}
}
// 2. Sliding Window
func TestSlidingWindowMaxLen(t *testing.T) {
// Window sum <= 10: [2,1,5] sum=8 ✓, [2,1,5,1,3] sum=12 violates → shrink
got := SlidingWindowMaxLen([]int{2, 1, 5, 1, 3, 2})
if got < 1 {
t.Errorf("SlidingWindowMaxLen should return positive, got %d", got)
}
}
// 3. Two Pointers
func TestTwoPointers(t *testing.T) {
i, j := TwoPointers([]int{1, 2, 3, 4, 6}, 6)
if i != 1 || j != 3 {
t.Errorf("TwoPointers(target=6) = (%d,%d), want (1,3)", i, j)
}
i, j = TwoPointers([]int{1, 2, 3}, 99)
if i != -1 || j != -1 {
t.Errorf("TwoPointers(missing) = (%d,%d), want (-1,-1)", i, j)
}
}
// 4. BFS
func TestBFSGrid(t *testing.T) {
// Grid where 9 is the goal — place it at (2,2)
grid := [][]int{
{0, 0, 0},
{0, 0, 0},
{0, 0, 9},
}
got := BFSGrid(grid, 0, 0)
if got != 4 {
t.Errorf("BFSGrid to (2,2) = %d, want 4 steps", got)
}
}
// 5. DFS
func TestDFSRecursive(t *testing.T) {
graph := map[int][]int{
1: {2, 3},
2: {4},
3: {},
4: {},
}
result := DFSRecursive(graph, 1)
if len(result) != 4 {
t.Errorf("DFSRecursive visited %d nodes, want 4", len(result))
}
}
func TestDFSIterative(t *testing.T) {
graph := map[int][]int{
1: {2, 3},
2: {4},
3: {},
4: {},
}
result := DFSIterative(graph, 1)
if len(result) != 4 {
t.Errorf("DFSIterative visited %d nodes, want 4", len(result))
}
}
// 6. Dynamic Programming
func TestDPTopDown(t *testing.T) {
// fibonacci: 0,1,1,2,3,5,8,13,21
tests := []struct{ n, want int }{{0, 0}, {1, 1}, {7, 13}, {10, 55}}
for _, tt := range tests {
if got := DPTopDown(tt.n); got != tt.want {
t.Errorf("DPTopDown(%d) = %d, want %d", tt.n, got, tt.want)
}
}
}
func TestDPBottomUp(t *testing.T) {
tests := []struct{ n, want int }{{0, 0}, {1, 1}, {7, 13}, {10, 55}}
for _, tt := range tests {
if got := DPBottomUp(tt.n); got != tt.want {
t.Errorf("DPBottomUp(%d) = %d, want %d", tt.n, got, tt.want)
}
}
}
// 7. Monotonic Stack
func TestNextGreaterElement(t *testing.T) {
got := NextGreaterElement([]int{2, 1, 2, 4, 3})
want := []int{4, 2, 4, -1, -1}
for i := range want {
if got[i] != want[i] {
t.Errorf("NextGreaterElement[%d] = %d, want %d", i, got[i], want[i])
}
}
}
// 8. Union-Find
func TestUnionFind(t *testing.T) {
uf := NewUnionFind(5)
uf.Union(0, 1)
uf.Union(2, 3)
if !uf.Connected(0, 1) {
t.Error("0 and 1 should be connected")
}
if uf.Connected(0, 2) {
t.Error("0 and 2 should NOT be connected")
}
uf.Union(1, 2)
if !uf.Connected(0, 3) {
t.Error("after Union(1,2): 0 and 3 should be connected")
}
if uf.count != 2 { // {0,1,2,3} and {4}
t.Errorf("component count = %d, want 2", uf.count)
}
}
// 9. Heap
func TestHeapExample(t *testing.T) {
min := HeapExample()
if min != 0 {
t.Errorf("HeapExample min = %d, want 0 (pushed 0 into [5,3,8,1,2])", min)
}
}