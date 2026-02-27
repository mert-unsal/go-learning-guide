package strings_strconv
import (
"reflect"
"testing"
)
func TestIsPalindromeEx(t *testing.T) {
tests := []struct{ s string; want bool }{
{"A man a plan a canal Panama", true},
{"race a car", false},
{"", true},
{"Was it a car or a cat I saw", true},
{"hello", false},
}
for _, tt := range tests {
if got := IsPalindromeExSolution(tt.s); got != tt.want {
t.Errorf("IsPalindrome(%q) = %v, want %v", tt.s, got, tt.want)
}
}
}
func TestReverseWords(t *testing.T) {
tests := []struct{ s, want string }{
{"the sky is blue", "blue is sky the"},
{"hello", "hello"},
{"Go is fun", "fun is Go"},
}
for _, tt := range tests {
if got := ReverseWordsSolution(tt.s); got != tt.want {
t.Errorf("ReverseWords(%q) = %q, want %q", tt.s, got, tt.want)
}
}
}
func TestCountOccurrences(t *testing.T) {
tests := []struct{ s, sub string; want int }{
{"hello world hello", "hello", 2},
{"aaa", "aa", 1},
{"abc", "xyz", 0},
}
for _, tt := range tests {
if got := CountOccurrencesSolution(tt.s, tt.sub); got != tt.want {
t.Errorf("CountOccurrences(%q,%q) = %d, want %d", tt.s, tt.sub, got, tt.want)
}
}
}
func TestTitleCase(t *testing.T) {
tests := []struct{ s, want string }{
{"the quick brown fox", "The Quick Brown Fox"},
{"hello world", "Hello World"},
{"GO", "Go"},
}
for _, tt := range tests {
if got := TitleCaseSolution(tt.s); got != tt.want {
t.Errorf("TitleCase(%q) = %q, want %q", tt.s, got, tt.want)
}
}
}
func TestParseCSVLine(t *testing.T) {
got := ParseCSVLineSolution("alice, 30 , engineer")
want := []string{"alice", "30", "engineer"}
if !reflect.DeepEqual(got, want) {
t.Errorf("ParseCSVLine = %v, want %v", got, want)
}
}
func TestIntToBase(t *testing.T) {
tests := []struct{ n, base int; want string }{
{255, 16, "ff"},
{8, 2, "1000"},
{10, 10, "10"},
{255, 8, "377"},
}
for _, tt := range tests {
if got := IntToBaseSolution(tt.n, tt.base); got != tt.want {
t.Errorf("IntToBase(%d,%d) = %q, want %q", tt.n, tt.base, got, tt.want)
}
}
}
func TestParseInts(t *testing.T) {
got, err := ParseIntsSolution([]string{"1", "2", "3"})
if err != nil || !reflect.DeepEqual(got, []int{1, 2, 3}) {
t.Errorf("ParseInts ok case: got %v err %v", got, err)
}
_, err = ParseIntsSolution([]string{"1", "abc", "3"})
if err == nil {
t.Error("ParseInts should error on non-integer string")
}
}