package strings_strconv
import (
"fmt"
"strconv"
"strings"
"unicode"
)
// SOLUTIONS â€” 01 strings & strconv
func IsPalindromeExSolution(s string) bool {
// Keep only alphanumeric, lowercase
var clean []rune
for _, ch := range strings.ToLower(s) {
if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
clean = append(clean, ch)
}
}
for i, j := 0, len(clean)-1; i < j; i, j = i+1, j-1 {
if clean[i] != clean[j] {
return false
}
}
return true
}
func ReverseWordsSolution(s string) string {
words := strings.Fields(s)
for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
words[i], words[j] = words[j], words[i]
}
return strings.Join(words, " ")
}
func CountOccurrencesSolution(s, substr string) int {
return strings.Count(s, substr)
}
func TitleCaseSolution(s string) string {
words := strings.Fields(s)
for i, w := range words {
if len(w) > 0 {
words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
}
}
return strings.Join(words, " ")
}
func ParseCSVLineSolution(line string) []string {
parts := strings.Split(line, ",")
for i, p := range parts {
parts[i] = strings.TrimSpace(p)
}
return parts
}
func IntToBaseSolution(n, base int) string {
return strconv.FormatInt(int64(n), base)
}
func ParseIntsSolution(strs []string) ([]int, error) {
result := make([]int, 0, len(strs))
for _, s := range strs {
n, err := strconv.Atoi(s)
if err != nil {
return nil, fmt.Errorf("ParseInts: %q is not a valid integer: %w", s, err)
}
result = append(result, n)
}
return result, nil
}