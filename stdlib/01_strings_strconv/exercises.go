package strings_strconv
// ============================================================
// EXERCISES â€” 01 strings & strconv
// ============================================================
// Exercise 1:
// IsPalindrome returns true if s reads the same forwards and backwards.
// Ignore case. Only consider alphanumeric characters.
// Example: "A man a plan a canal Panama" â†’ true, "race a car" â†’ false
func IsPalindromeEx(s string) bool {
// TODO: clean s (keep only alphanumeric, lowercase), then check
return false
}
// Exercise 2:
// ReverseWords reverses the ORDER of words in a sentence.
// Example: "the sky is blue" â†’ "blue is sky the"
// Words are separated by single spaces, no leading/trailing spaces.
func ReverseWords(s string) string {
// TODO: strings.Fields + reverse the slice + strings.Join
return ""
}
// Exercise 3:
// CountOccurrences returns how many times substr appears in s.
// Example: "hello world hello", "hello" â†’ 2
func CountOccurrences(s, substr string) int {
// TODO: strings.Count
return 0
}
// Exercise 4:
// TitleCase converts a sentence to Title Case.
// Example: "the quick brown fox" â†’ "The Quick Brown Fox"
func TitleCase(s string) string {
// TODO: strings.Fields, capitalize each word, strings.Join
return ""
}
// Exercise 5:
// ParseCSVLine parses a comma-separated line into a slice of strings.
// Trim spaces from each field.
// Example: "alice, 30 , engineer" â†’ ["alice", "30", "engineer"]
func ParseCSVLine(line string) []string {
// TODO: strings.Split then strings.TrimSpace each field
return nil
}
// Exercise 6:
// IntToBase converts n to a string in the given base (2, 8, 10, 16).
// Example: IntToBase(255, 16) â†’ "ff"
func IntToBase(n, base int) string {
// TODO: strconv.FormatInt(int64(n), base)
return ""
}
// Exercise 7:
// ParseInts parses a slice of numeric strings into []int.
// Return an error if any string is not a valid integer.
// Example: ["1","2","3"] â†’ [1,2,3], nil
func ParseInts(strs []string) ([]int, error) {
// TODO: strconv.Atoi each, collect errors
return nil, nil
}