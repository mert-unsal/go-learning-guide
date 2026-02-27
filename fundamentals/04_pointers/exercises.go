package pointers
// ============================================================
// EXERCISES — 04 Pointers
// ============================================================
// Exercise 1:
// Write a function that increments an integer through a pointer.
func Increment(n *int) {
// TODO: *n++
}
// Exercise 2:
// Write a function that swaps two integers using pointers.
func SwapPointers(a, b *int) {
// TODO: *a, *b = *b, *a
}
// Exercise 3:
// Write a pointer receiver method on ScoreBoard that adds points.
// Write a value receiver method that returns the current score.
type ScoreBoard struct {
Score int
}
func (s *ScoreBoard) AddPoints(points int) {
// TODO: s.Score += points
}
func (s ScoreBoard) CurrentScore() int {
// TODO: return s.Score
return 0
}
// Exercise 4:
// NewPlayer returns a POINTER to a Player struct (constructor pattern).
type Player struct {
Name  string
Level int
}
func NewPlayer(name string, level int) *Player {
// TODO: return &Player{...}
return nil
}
// Exercise 5:
// DoubleValue doubles the value at the pointer address.
func DoubleValue(x *int) {
// TODO: *x = *x * 2
}