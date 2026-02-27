package pointers

import "testing"

func TestIncrement(t *testing.T) {
	n := 5
	Increment(&n)
	if n != 6 {
		t.Errorf("❌ after Increment n=%d, want 6", n)
	} else {
		t.Logf("✅ Increment(&5) → n = %d", n)
	}
}

func TestSwapPointers(t *testing.T) {
	a, b := 10, 20
	SwapPointers(&a, &b)
	if a != 20 || b != 10 {
		t.Errorf("❌ SwapPointers: a=%d b=%d, want a=20 b=10", a, b)
	} else {
		t.Logf("✅ SwapPointers(10, 20) → a=%d, b=%d", a, b)
	}
}

func TestScoreBoard(t *testing.T) {
	s := &ScoreBoard{}
	s.AddPoints(10)
	s.AddPoints(5)
	got := s.CurrentScore()
	if got != 15 {
		t.Errorf("❌ ScoreBoard score = %d, want 15", got)
	} else {
		t.Logf("✅ AddPoints(10)+AddPoints(5) → CurrentScore = %d", got)
	}
}

func TestNewPlayer(t *testing.T) {
	p := NewPlayer("Alice", 5)
	if p == nil {
		t.Fatal("❌ NewPlayer returned nil  ← Hint: return &Player{...}")
	}
	if p.Name != "Alice" || p.Level != 5 {
		t.Errorf("❌ NewPlayer = %+v, want {Alice 5}", p)
	} else {
		t.Logf("✅ NewPlayer(\"Alice\", 5) = %+v", p)
	}
}

func TestDoubleValue(t *testing.T) {
	n := 7
	DoubleValue(&n)
	if n != 14 {
		t.Errorf("❌ after DoubleValue n=%d, want 14", n)
	} else {
		t.Logf("✅ DoubleValue(&7) → n = %d", n)
	}
}
