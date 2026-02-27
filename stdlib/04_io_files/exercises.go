package io_files
import "io"
// ============================================================
// EXERCISES — 04 io & files
// ============================================================
// Exercise 1:
// WriteAndReadFile writes content to a temp file, reads it back, returns the content.
// Use os.WriteFile and os.ReadFile.
// Delete the file when done (defer os.Remove).
func WriteAndReadFile(filename, content string) (string, error) {
// TODO: os.WriteFile → os.ReadFile → return string(data), nil
return "", nil
}
// Exercise 2:
// CountLines counts the number of lines in a string using bufio.Scanner.
// Simulate reading from a string by wrapping it in strings.NewReader.
func CountLines(content string) int {
// TODO: strings.NewReader → bufio.NewScanner → scan loop → count
return 0
}
// Exercise 3:
// WordFrequency reads text line-by-line from r (an io.Reader)
// and returns a map of word→frequency.
// Words are split by whitespace. Case-insensitive ("Go" and "go" are the same).
func WordFrequency(r io.Reader) map[string]int {
// TODO: bufio.NewScanner, scanner.Scan loop, strings.Fields, strings.ToLower
return nil
}
// Exercise 4:
// CopyReader reads all bytes from src and returns them as a string.
// Use io.ReadAll.
func CopyReader(src io.Reader) (string, error) {
// TODO: io.ReadAll(src)
return "", nil
}