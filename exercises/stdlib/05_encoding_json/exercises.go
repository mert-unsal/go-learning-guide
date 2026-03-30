package encoding_json
// ============================================================
// EXERCISES — 05 encoding/json
// ============================================================
// Exercise 1: Round-trip an Employee struct (different name from concepts.go Person)
type Employee struct {
Name       string `json:"name"`
Department string `json:"department"`
Salary     int    `json:"salary,omitempty"`
}
func RoundTrip(e Employee) (Employee, error) {
// TODO: json.Marshal → json.Unmarshal → return result
return Employee{}, nil
}
// Exercise 2: Parse arbitrary JSON into map[string]any
func ParseJSON(jsonStr string) (name string, age float64, err error) {
// TODO: json.Unmarshal into map[string]any, type-assert name and age
return "", 0, nil
}
// Exercise 3: Marshal a []int to a JSON array string
func MarshalSlice(nums []int) (string, error) {
// TODO: json.Marshal, return string(data), nil
return "", nil
}
// Exercise 4: AppConfig demonstrates json struct tags
// omitempty — omit field when it is the zero value
// json:"-"  — never include this field in JSON output
type AppConfig struct {
Host     string `json:"host"`
Port     int    `json:"port"`
Debug    bool   `json:"debug,omitempty"`
Secret   string `json:"-"`
}
func MarshalAppConfig(c AppConfig) ([]byte, error) {
// TODO: json.Marshal(c)
return nil, nil
}