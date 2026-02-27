package sliding_window
import "testing"
func TestFindMaxAverage(t *testing.T) {
tests := []struct {
name string
nums []int
k    int
want float64
}{
{"basic", []int{1, 12, -5, -6, 50, 3}, 4, 12.75},
{"single window", []int{5, 5, 5}, 3, 5.0},
{"k=1", []int{3, 1, 4, 1, 5}, 1, 5.0},
}
for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := FindMaxAverage(tt.nums, tt.k)
if got != tt.want {
t.Errorf("FindMaxAverage(%v, %d) = %v, want %v", tt.nums, tt.k, got, tt.want)
}
})
}
}
func TestMinWindow(t *testing.T) {
tests := []struct {
name string
s, t string
want string
}{
{"classic", "ADOBECODEBANC", "ABC", "BANC"},
{"same string", "a", "a", "a"},
{"no match", "a", "aa", ""},
{"empty t", "a", "", ""},
}
for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := MinWindow(tt.s, tt.t)
if got != tt.want {
t.Errorf("MinWindow(%q, %q) = %q, want %q", tt.s, tt.t, got, tt.want)
}
})
}
}
func TestCheckInclusion(t *testing.T) {
tests := []struct {
name   string
s1, s2 string
want   bool
}{
{"contains perm", "ab", "eidbaooo", true},
{"no perm", "ab", "eidboaoo", false},
{"s1 longer", "abc", "ab", false},
{"exact match", "abc", "cba", true},
}
for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := CheckInclusion(tt.s1, tt.s2)
if got != tt.want {
t.Errorf("CheckInclusion(%q, %q) = %v, want %v", tt.s1, tt.s2, got, tt.want)
}
})
}
}
func TestTotalFruit(t *testing.T) {
tests := []struct {
name   string
fruits []int
want   int
}{
{"basic", []int{1, 2, 1}, 3},
{"three types", []int{0, 1, 2, 2}, 3},
{"longer", []int{1, 2, 3, 2, 2}, 4},
{"single type", []int{1, 1, 1}, 3},
}
for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := TotalFruit(tt.fruits)
if got != tt.want {
t.Errorf("TotalFruit(%v) = %d, want %d", tt.fruits, got, tt.want)
}
})
}
}
