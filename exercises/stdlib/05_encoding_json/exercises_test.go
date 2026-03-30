package encoding_json
import (
"strings"
"testing"
)
func TestRoundTrip(t *testing.T) {
e := Employee{Name: "Alice", Department: "Engineering", Salary: 100000}
got, err := RoundTripSolution(e)
if err != nil { t.Fatalf("RoundTrip error: %v", err) }
if got.Name != "Alice" || got.Department != "Engineering" {
t.Errorf("RoundTrip = %+v, want {Alice Engineering}", got)
}
}
func TestParseJSON(t *testing.T) {
name, age, err := ParseJSONSolution(`{"name":"Bob","age":25}`)
if err != nil || name != "Bob" || age != 25 {
t.Errorf("ParseJSON = (%q,%v,%v), want (Bob,25,nil)", name, age, err)
}
}
func TestMarshalSlice(t *testing.T) {
got, err := MarshalSliceSolution([]int{1, 2, 3})
if err != nil || got != "[1,2,3]" {
t.Errorf("MarshalSlice = (%q,%v), want ([1,2,3],nil)", got, err)
}
}
func TestMarshalAppConfig(t *testing.T) {
c := AppConfig{Host: "localhost", Port: 8080, Secret: "hunter2"}
data, err := MarshalAppConfigSolution(c)
if err != nil { t.Fatalf("MarshalAppConfig error: %v", err) }
s := string(data)
if strings.Contains(s, "hunter2") {
t.Error(`Secret field with json:"-" must never appear in JSON`)
}
if strings.Contains(s, "debug") {
t.Error("Debug=false with omitempty must not appear in JSON")
}
if !strings.Contains(s, "localhost") {
t.Error("Host must appear in JSON output")
}
}