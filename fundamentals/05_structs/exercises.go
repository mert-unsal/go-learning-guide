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
	return r.Width * r.Height
}

func (r ExRectangle) Perimeter() float64 {
	// TODO: 2*(Width+Height)
	return 2 * (r.Width + r.Height)
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
	sum := 0.0
	for _, shape := range shapes {
		sum += shape.Area()
	}
	return sum
}

// Exercise 3:
// ExStack — a stack data structure backed by a slice.
type ExStack struct {
	// TODO: add items []int
	Items []int
}

func (s *ExStack) Push(val int) {
	// TODO: append to items
	s.Items = append(s.Items, val)
}

func (s *ExStack) Pop() (int, bool) {
	// TODO: remove and return top
	if len(s.Items) == 0 {
		return 0, false
	}
	n := len(s.Items)
	val := s.Items[n-1]
	s.Items = s.Items[:n-1]
	return val, true
}

func (s *ExStack) IsEmpty() bool {
	// TODO: return len == 0
	return len(s.Items) == 0
}
