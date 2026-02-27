package io_files
import (
"bufio"
"io"
"os"
"strings"
)
// SOLUTIONS — 04 io & files
func WriteAndReadFileSolution(filename, content string) (string, error) {
if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
return "", err
}
defer os.Remove(filename)
data, err := os.ReadFile(filename)
if err != nil {
return "", err
}
return string(data), nil
}
func CountLinesSolution(content string) int {
scanner := bufio.NewScanner(strings.NewReader(content))
count := 0
for scanner.Scan() {
count++
}
return count
}
func WordFrequencySolution(r io.Reader) map[string]int {
freq := make(map[string]int)
scanner := bufio.NewScanner(r)
for scanner.Scan() {
for _, word := range strings.Fields(scanner.Text()) {
freq[strings.ToLower(word)]++
}
}
return freq
}
func CopyReaderSolution(src io.Reader) (string, error) {
data, err := io.ReadAll(src)
if err != nil {
return "", err
}
return string(data), nil
}