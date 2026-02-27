package structs
import (
"math"
"testing"
)
func TestExRectangle(t *testing.T) {
r := ExRectangle{Width: 3, Height: 4}
if r.AreaSolution() != 12 {
t.Errorf("Area() = %.1f, want 12", r.AreaSolution())
}
if r.PerimeterSolution() != 14 {
t.Errorf("Perimeter() = %.1f, want 14", r.PerimeterSolution())
}
}
func TestTotalArea(t *testing.T) {
shapes := []ExShape{
ExCircle{Radius: 1},         // pi
ExCircle{Radius: 2},         // 4*pi
}
want := 5 * math.Pi
got := TotalAreaSolution(shapes)
if diff := math.Abs(got - want); diff > 0.0001 {
t.Errorf("TotalArea = %f, want %f", got, want)
}
}
func TestExStack(t *testing.T) {
s := &ExStackSolution{}
if !s.IsEmpty() { t.Error("new stack should be empty") }
s.Push(10); s.Push(20)
if v, ok := s.Pop(); !ok || v != 20 {
t.Errorf("Pop() = (%d,%v), want (20,true)", v, ok)
}
if s.IsEmpty() { t.Error("stack should still have 1 item") }
}