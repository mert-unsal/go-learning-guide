package encoding_json

import (
	"encoding/json"
	"io"
	"time"
)

// ============================================================
// EXERCISES — 05 encoding/json
// ============================================================
// 12 exercises covering encoding/json from struct tags to
// streaming decoders, custom marshalers, and production patterns.

// ────────────────────────────────────────────────────────────
// Exercise 1: RoundTrip — marshal then unmarshal an Employee
// ────────────────────────────────────────────────────────────
// json.Marshal(e) → json.Unmarshal(data, &result) → return result

type Employee struct {
	Name       string `json:"name"`
	Department string `json:"department"`
	Salary     int    `json:"salary,omitempty"`
}

func RoundTrip(e Employee) (Employee, error) {
	return Employee{}, nil
}

// ────────────────────────────────────────────────────────────
// Exercise 2: ParseJSON — decode into map[string]any
// ────────────────────────────────────────────────────────────
// json.Unmarshal into map[string]any. Numbers become float64.
// Type-assert "name" as string, "age" as float64.

func ParseJSON(jsonStr string) (name string, age float64, err error) {
	return "", 0, nil
}

// ────────────────────────────────────────────────────────────
// Exercise 3: MarshalSlice — marshal []int to JSON array string
// ────────────────────────────────────────────────────────────

func MarshalSlice(nums []int) (string, error) {
	return "", nil
}

// ────────────────────────────────────────────────────────────
// Exercise 4: MarshalAppConfig — struct tags: omitempty, "-"
// ────────────────────────────────────────────────────────────
// omitempty skips zero-value fields. json:"-" never serializes.

type AppConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Debug  bool   `json:"debug,omitempty"`
	Secret string `json:"-"`
}

func MarshalAppConfig(c AppConfig) ([]byte, error) {
	return nil, nil
}

// ────────────────────────────────────────────────────────────
// Exercise 5: DecodeStream — use json.NewDecoder on io.Reader
// ────────────────────────────────────────────────────────────
// json.NewDecoder(r).Decode(&v) streams — doesn't read all into memory.
// Decode one Employee from the reader.

func DecodeStream(r io.Reader) (Employee, error) {
	return Employee{}, nil
}

// ────────────────────────────────────────────────────────────
// Exercise 6: EncodeStream — use json.NewEncoder on io.Writer
// ────────────────────────────────────────────────────────────
// json.NewEncoder(w).Encode(v) writes JSON + newline to w.

func EncodeStream(w io.Writer, e Employee) error {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 7: Custom Marshaler — StatusCode with string repr
// ────────────────────────────────────────────────────────────
// StatusCode is an int that serializes as a descriptive string:
//   200 → "OK", 404 → "Not Found", 500 → "Internal Server Error"
// Implement MarshalJSON and UnmarshalJSON on StatusCode.

type StatusCode int

func (s StatusCode) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (s *StatusCode) UnmarshalJSON(data []byte) error {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 8: DelayParse — use json.RawMessage to defer parsing
// ────────────────────────────────────────────────────────────
// Given JSON: {"type":"employee","data":{...}}
// Parse "type" first, then based on type, parse "data" into
// the correct struct. Return the type string and the decoded Employee.
//
// json.RawMessage keeps "data" as raw bytes until you know the type.

func DelayParse(jsonStr string) (typeName string, emp Employee, err error) {
	return "", Employee{}, nil
}

// ────────────────────────────────────────────────────────────
// Exercise 9: NullableFields — pointer fields distinguish null/absent
// ────────────────────────────────────────────────────────────
// Given a Profile struct, unmarshal JSON and report which fields
// were explicitly set to null vs simply absent.
//
// After unmarshal: nil pointer = absent OR null. You need the raw
// JSON to distinguish. Return: (profile, isNicknameNull bool)

type Profile struct {
	Name     string  `json:"name"`
	Nickname *string `json:"nickname,omitempty"`
	Age      *int    `json:"age,omitempty"`
}

func NullableFields(jsonStr string) (Profile, bool, error) {
	// Return: (profile, true if "nickname" is explicitly null, error)
	return Profile{}, false, nil
}

// ────────────────────────────────────────────────────────────
// Exercise 10: TimeEvent — custom time format in JSON
// ────────────────────────────────────────────────────────────
// time.Time marshals to RFC3339 by default. Implement MarshalJSON
// and UnmarshalJSON on Event so that "timestamp" uses the format
// "2006-01-02 15:04:05" instead.

type Event struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e Event) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (e *Event) UnmarshalJSON(data []byte) error {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 11: ValidateJSON — check if a string is valid JSON
// ────────────────────────────────────────────────────────────
// Return true if the input is syntactically valid JSON.
// Hint: json.Valid([]byte(s))

func ValidateJSON(s string) bool {
	return false
}

// ────────────────────────────────────────────────────────────
// Exercise 12: PreciseNumber — use json.Number for exact numbers
// ────────────────────────────────────────────────────────────
// Default JSON decoding turns numbers into float64, which loses
// precision for large integers (>2^53). Use json.Decoder with
// UseNumber() to preserve the exact number string.
//
// Return the raw number string and its int64 value.

func PreciseNumber(jsonStr string) (raw string, val int64, err error) {
	return "", 0, nil
}

// Ensure json import is used
var _ = json.Marshal
