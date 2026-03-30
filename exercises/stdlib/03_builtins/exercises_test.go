package builtins
import (
"reflect"
"testing"
)
func TestDeepCopySlice(t *testing.T) {
src := []int{1, 2, 3}
dst := DeepCopySliceSolution(src)
dst[0] = 99
if src[0] == 99 {
t.Error("DeepCopySlice: modifying copy should not affect src")
}
if !reflect.DeepEqual(DeepCopySliceSolution([]int{4, 5}), []int{4, 5}) {
t.Error("DeepCopySlice values wrong")
}
}
func TestDeepCopyMap(t *testing.T) {
src := map[string]int{"a": 1, "b": 2}
dst := DeepCopyMapSolution(src)
dst["a"] = 99
if src["a"] == 99 {
t.Error("DeepCopyMap: modifying copy should not affect src")
}
}
func TestSafeDivideEx(t *testing.T) {
r, err := SafeDivideExSolution(10, 2)
if err != nil || r != 5 {
t.Errorf("SafeDivide(10,2) = (%d,%v), want (5,nil)", r, err)
}
r, err = SafeDivideExSolution(10, 0)
if err == nil {
t.Errorf("SafeDivide(10,0) should return error, got result=%d", r)
}
}
func TestFlatten(t *testing.T) {
got := FlattenSolution([][]int{{1, 2}, {3}, {4, 5}})
want := []int{1, 2, 3, 4, 5}
if !reflect.DeepEqual(got, want) {
t.Errorf("Flatten = %v, want %v", got, want)
}
}
func TestUniqueInts(t *testing.T) {
got := UniqueIntsSolution([]int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3})
want := []int{3, 1, 4, 5, 9, 2, 6}
if !reflect.DeepEqual(got, want) {
t.Errorf("UniqueInts = %v, want %v", got, want)
}
}
func TestChunkSlice(t *testing.T) {
got := ChunkSliceSolution([]int{1, 2, 3, 4, 5}, 2)
if len(got) != 3 || got[0][0] != 1 || got[2][0] != 5 {
t.Errorf("ChunkSlice = %v, unexpected", got)
}
// single chunk
got2 := ChunkSliceSolution([]int{1, 2}, 5)
if len(got2) != 1 || len(got2[0]) != 2 {
t.Errorf("ChunkSlice single chunk = %v", got2)
}
}