package interfaces
import (
"fmt"
"strconv"
)
// SOLUTIONS — 06 Interfaces
func (b ExBook) StringSolution() string {
return fmt.Sprintf("%q by %s", b.Title, b.Author)
}
func (m ExMovie) StringSolution() string {
return fmt.Sprintf("%s (%d)", m.Title, m.Year)
}
func PrintAllSolution(items []ExStringer) {
for _, item := range items {
fmt.Println(item.String())
}
}
func (bw *ExBufferWriter) WriteSolution(data string) error {
bw.Buffer = append(bw.Buffer, data)
return nil
}
func WriteAllSolution(w ExWriter, items []string) error {
for _, item := range items {
if err := w.Write(item); err != nil {
return err
}
}
return nil
}
func DescribeSolution(i interface{}) string {
switch v := i.(type) {
case int:
return "int: " + strconv.Itoa(v)
case string:
return "string: " + v
case bool:
if v { return "bool: true" }
return "bool: false"
default:
return "unknown"
}
}