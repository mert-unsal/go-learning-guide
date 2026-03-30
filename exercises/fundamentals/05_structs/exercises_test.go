package structs

import (
	"math"
	"testing"
)

func TestExRectangle(t *testing.T) {
	r := ExRectangle{Width: 3, Height: 4}

	area := r.Area()
	if area != 12 {
		t.Errorf("❌ Area() = %.1f, want 12", area)
	} else {
		t.Logf("✅ Rectangle{3,4}.Area() = %.1f", area)
	}

	perim := r.Perimeter()
	if perim != 14 {
		t.Errorf("❌ Perimeter() = %.1f, want 14", perim)
	} else {
		t.Logf("✅ Rectangle{3,4}.Perimeter() = %.1f", perim)
	}
}

func TestTotalArea(t *testing.T) {
	shapes := []ExShape{
		ExCircle{Radius: 1}, // pi
		ExCircle{Radius: 2}, // 4*pi
	}
	want := 5 * math.Pi
	got := TotalArea(shapes)
	if diff := math.Abs(got - want); diff > 0.0001 {
		t.Errorf("❌ TotalArea = %f, want %f", got, want)
	} else {
		t.Logf("✅ TotalArea([circle(r=1), circle(r=2)]) = %.4f", got)
	}
}

func TestExStack(t *testing.T) {
	s := &ExStack{}

	if !s.IsEmpty() {
		t.Error("❌ new stack should be empty  ← Hint: check len(items)==0")
	} else {
		t.Logf("✅ new ExStack is empty")
	}

	s.Push(10)
	s.Push(20)
	t.Logf("✅ Pushed 10 and 20")

	v, ok := s.Pop()
	if !ok || v != 20 {
		t.Errorf("❌ Pop() = (%d,%v), want (20,true)  ← Hint: pop from end of slice", v, ok)
	} else {
		t.Logf("✅ Pop() = (%d,%v)", v, ok)
	}

	if s.IsEmpty() {
		t.Error("❌ stack should still have 1 item after one pop")
	} else {
		t.Logf("✅ stack still has 1 item remaining")
	}
}
