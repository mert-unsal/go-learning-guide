package pointers
import "fmt"
// SOLUTIONS — 04 Pointers
func IncrementSolution(n *int) { *n++ }
func SwapPointersSolution(a, b *int) { *a, *b = *b, *a }
func (s *ScoreBoard) AddPointsSolution(points int) { s.Score += points }
func (s ScoreBoard) CurrentScoreSolution() int     { return s.Score }
func NewPlayerSolution(name string, level int) *Player {
return &Player{Name: name, Level: level}
}
func (p *Player) String() string {
return fmt.Sprintf("%s (level %d)", p.Name, p.Level)
}
func DoubleValueSolution(x *int) { *x = *x * 2 }