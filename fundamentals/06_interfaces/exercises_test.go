package interfaces
import "testing"
func TestExStringers(t *testing.T) {
b := ExBook{Title: "The Go Programming Language", Author: "Donovan"}
want := `"The Go Programming Language" by Donovan`
if got := b.StringSolution(); got != want {
t.Errorf("Book.String() = %q, want %q", got, want)
}
m := ExMovie{Title: "Inception", Year: 2010}
if got := m.StringSolution(); got != "Inception (2010)" {
t.Errorf("Movie.String() = %q, want Inception (2010)", got)
}
}
func TestExBufferWriter(t *testing.T) {
bw := &ExBufferWriter{}
bw.WriteSolution("hello")
bw.WriteSolution("world")
if len(bw.Buffer) != 2 {
t.Errorf("buffer has %d items, want 2", len(bw.Buffer))
}
if bw.Buffer[0] != "hello" { t.Error("buffer[0] should be hello") }
}
func TestDescribe(t *testing.T) {
tests := []struct{ input interface{}; want string }{
{42, "int: 42"},
{"hello", "string: hello"},
{true, "bool: true"},
{3.14, "unknown"},
}
for _, tt := range tests {
got := DescribeSolution(tt.input)
if got != tt.want {
t.Errorf("Describe(%v) = %q, want %q", tt.input, got, tt.want)
}
}
}