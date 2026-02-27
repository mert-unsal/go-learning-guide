package structs
import (
"fmt"
"math"
)
// SOLUTIONS — 05 Structs
func (r ExRectangle) AreaSolution() float64      { return r.Width * r.Height }
func (r ExRectangle) PerimeterSolution() float64 { return 2 * (r.Width + r.Height) }
func (r ExRectangle) String() string {
return fmt.Sprintf("Rectangle(%.0fx%.0f)", r.Width, r.Height)
}
type ExCircle struct{ Radius float64 }
func (c ExCircle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c ExCircle) Perimeter() float64 { return 2 * math.Pi * c.Radius }
func TotalAreaSolution(shapes []ExShape) float64 {
total := 0.0
for _, s := range shapes {
total += s.Area()
}
return total
}
type ExStackSolution struct{ items []int }
func (s *ExStackSolution) Push(val int) { s.items = append(s.items, val) }
func (s *ExStackSolution) Pop() (int, bool) {
if len(s.items) == 0 { return 0, false }
top := s.items[len(s.items)-1]
s.items = s.items[:len(s.items)-1]
return top, true
}
func (s *ExStackSolution) IsEmpty() bool { return len(s.items) == 0 }