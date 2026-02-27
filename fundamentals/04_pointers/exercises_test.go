package pointers
import "testing"
func TestIncrement(t *testing.T) {
n := 5
IncrementSolution(&n)
if n != 6 { t.Errorf("after Increment n=%d, want 6", n) }
}
func TestSwapPointers(t *testing.T) {
a, b := 10, 20
SwapPointersSolution(&a, &b)
if a != 20 || b != 10 { t.Errorf("swap: a=%d b=%d, want 20 10", a, b) }
}
func TestScoreBoard(t *testing.T) {
s := &ScoreBoard{}
s.AddPointsSolution(10)
s.AddPointsSolution(5)
if s.CurrentScoreSolution() != 15 {
t.Errorf("score = %d, want 15", s.CurrentScoreSolution())
}
}
func TestNewPlayer(t *testing.T) {
p := NewPlayerSolution("Alice", 5)
if p == nil { t.Fatal("NewPlayer returned nil") }
if p.Name != "Alice" || p.Level != 5 {
t.Errorf("got %+v, want {Alice 5}", p)
}
}
func TestDoubleValue(t *testing.T) {
n := 7
DoubleValueSolution(&n)
if n != 14 { t.Errorf("after DoubleValue n=%d, want 14", n) }
}