package encoding_json
import "encoding/json"
// SOLUTIONS — 05 encoding/json
func RoundTripSolution(e Employee) (Employee, error) {
data, err := json.Marshal(e)
if err != nil { return Employee{}, err }
var result Employee
if err := json.Unmarshal(data, &result); err != nil { return Employee{}, err }
return result, nil
}
func ParseJSONSolution(jsonStr string) (name string, age float64, err error) {
var m map[string]any
if err = json.Unmarshal([]byte(jsonStr), &m); err != nil { return }
name, _ = m["name"].(string)
age, _ = m["age"].(float64)
return
}
func MarshalSliceSolution(nums []int) (string, error) {
data, err := json.Marshal(nums)
if err != nil { return "", err }
return string(data), nil
}
func MarshalAppConfigSolution(c AppConfig) ([]byte, error) {
return json.Marshal(c)
}