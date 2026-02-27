package interfaces
// ============================================================
// EXERCISES — 06 Interfaces
// ============================================================
// Exercise 1:
// Stringer — any type with String() string method.
// Book and Movie both implement it.
type ExStringer interface {
String() string
}
type ExBook struct{ Title, Author string }
type ExMovie struct{ Title string; Year int }
func (b ExBook) String() string {
// TODO: return `"Title" by Author`
return ""
}
func (m ExMovie) String() string {
// TODO: return `Title (Year)`
return ""
}
func PrintAll(items []ExStringer) {
// TODO: print each item.String()
}
// Exercise 2:
// ExWriter interface — anything that can Write a string.
type ExWriter interface {
Write(data string) error
}
type ExBufferWriter struct {
Buffer []string
}
func (bw *ExBufferWriter) Write(data string) error {
// TODO: append data to Buffer
return nil
}
func WriteAll(w ExWriter, items []string) error {
// TODO: call w.Write for each item
return nil
}
// Exercise 3:
// Type switch — describe what kind of value is passed.
func Describe(i interface{}) string {
// TODO: type switch on int, string, bool, default
return "unknown"
}