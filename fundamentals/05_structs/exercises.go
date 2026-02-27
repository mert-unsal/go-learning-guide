package structs
// ============================================================
// EXERCISES — 05 Structs
// ============================================================
// Exercise 1:
// ExRectangle has Width and Height. Implement Area, Perimeter, String.
type ExRectangle struct {
Width, Height float64
}
func (r ExRectangle) Area() float64 {
// TODO: Width * Height
return 0
}
func (r ExRectangle) Perimeter() float64 {
// TODO: 2*(Width+Height)
return 0
}
// Exercise 2:
// ExShape interface — any type with Area and Perimeter satisfies it.
type ExShape interface {
Area() float64
Perimeter() float64
}
// TotalArea sums areas of all shapes.
func TotalArea(shapes []ExShape) float64 {
// TODO: iterate and sum
return 0
}
// Exercise 3:
// ExStack — a stack data structure backed by a slice.
type ExStack struct {
// TODO: add items []int
}
func (s *ExStack) Push(val int) {
// TODO: append to items
}
func (s *ExStack) Pop() (int, bool) {
// TODO: remove and return top
return 0, false
}
func (s *ExStack) IsEmpty() bool {
// TODO: return len == 0
return true
}