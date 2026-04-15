package encoding_json

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: RoundTrip
// ────────────────────────────────────────────────────────────

func TestRoundTrip(t *testing.T) {
	e := Employee{Name: "Alice", Department: "Engineering", Salary: 100000}
	got, err := RoundTrip(e)
	if err != nil {
		t.Fatalf("❌ RoundTrip error: %v\n\t\tHint: json.Marshal(e) returns []byte. "+
			"json.Unmarshal(data, &result) fills it back. Two lines", err)
	}
	if got != e {
		t.Errorf("❌ RoundTrip = %+v, want %+v", got, e)
	} else {
		t.Logf("✅ RoundTrip round-trip OK: %+v", got)
	}

	// omitempty: zero Salary should be absent
	e2 := Employee{Name: "Bob", Department: "Sales"}
	data, _ := json.Marshal(e2)
	if strings.Contains(string(data), "salary") {
		t.Errorf("❌ Salary=0 with omitempty should be absent from JSON, got %s\n\t\t"+
			"Hint: omitempty omits zero-value fields. 0 is the zero value for int", string(data))
	} else {
		t.Logf("✅ omitempty correctly omits zero Salary")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: ParseJSON
// ────────────────────────────────────────────────────────────

func TestParseJSON(t *testing.T) {
	name, age, err := ParseJSON(`{"name":"Bob","age":25}`)
	if err != nil {
		t.Fatalf("❌ error: %v\n\t\tHint: var m map[string]any; json.Unmarshal([]byte(s), &m); "+
			"then m[\"name\"].(string), m[\"age\"].(float64)", err)
	}
	if name != "Bob" {
		t.Errorf("❌ name = %q, want \"Bob\"\n\t\tHint: JSON strings → Go string via type assertion", name)
	} else {
		t.Logf("✅ name = %q", name)
	}
	if age != 25 {
		t.Errorf("❌ age = %v, want 25\n\t\tHint: JSON numbers → Go float64 (not int!). "+
			"This is why json.Number exists for precision", age)
	} else {
		t.Logf("✅ age = %v", age)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: MarshalSlice
// ────────────────────────────────────────────────────────────

func TestMarshalSlice(t *testing.T) {
	got, err := MarshalSlice([]int{1, 2, 3})
	if err != nil {
		t.Fatalf("❌ error: %v", err)
	}
	if got != "[1,2,3]" {
		t.Errorf("❌ MarshalSlice = %q, want \"[1,2,3]\"\n\t\t"+
			"Hint: json.Marshal(nums) returns compact JSON. string(data) converts", got)
	} else {
		t.Logf("✅ MarshalSlice = %q", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: MarshalAppConfig
// ────────────────────────────────────────────────────────────

func TestMarshalAppConfig(t *testing.T) {
	c := AppConfig{Host: "localhost", Port: 8080, Secret: "hunter2"}
	data, err := MarshalAppConfig(c)
	if err != nil {
		t.Fatalf("❌ error: %v\n\t\tHint: json.Marshal(c) — struct tags control the output", err)
	}
	s := string(data)
	if s == "" {
		t.Fatal("❌ returned empty bytes\n\t\tHint: return json.Marshal(c)")
	}
	if strings.Contains(s, "hunter2") {
		t.Error("❌ Secret with json:\"-\" must NEVER appear in JSON\n\t\t" +
			"Hint: json:\"-\" means the field is completely excluded from marshal/unmarshal")
	} else {
		t.Logf("✅ Secret correctly excluded")
	}
	if strings.Contains(s, "debug") {
		t.Error("❌ Debug=false with omitempty should not appear\n\t\t" +
			"Hint: omitempty skips false, 0, \"\", nil, empty slice/map")
	} else {
		t.Logf("✅ Debug correctly omitted (omitempty + false)")
	}
	if !strings.Contains(s, "localhost") {
		t.Error("❌ Host must appear in JSON output")
	} else {
		t.Logf("✅ Host present in output")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: DecodeStream
// ────────────────────────────────────────────────────────────

func TestDecodeStream(t *testing.T) {
	input := `{"name":"Charlie","department":"Ops","salary":90000}`
	r := strings.NewReader(input)
	got, err := DecodeStream(r)
	if err != nil {
		t.Fatalf("❌ error: %v\n\t\tHint: json.NewDecoder(r).Decode(&emp) — "+
			"streams from io.Reader, doesn't buffer entire input", err)
	}
	if got.Name != "Charlie" || got.Department != "Ops" || got.Salary != 90000 {
		t.Errorf("❌ DecodeStream = %+v, want {Charlie Ops 90000}\n\t\t"+
			"Hint: Decoder reads tokens incrementally — better for HTTP bodies, "+
			"files, and pipes than json.Unmarshal(readAll)", got)
	} else {
		t.Logf("✅ DecodeStream = %+v", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: EncodeStream
// ────────────────────────────────────────────────────────────

func TestEncodeStream(t *testing.T) {
	var buf bytes.Buffer
	e := Employee{Name: "Dana", Department: "HR", Salary: 85000}
	err := EncodeStream(&buf, e)
	if err != nil {
		t.Fatalf("❌ error: %v\n\t\tHint: json.NewEncoder(w).Encode(v) — "+
			"writes JSON + trailing newline to the writer", err)
	}
	out := strings.TrimSpace(buf.String())
	if !strings.Contains(out, `"name":"Dana"`) {
		t.Errorf("❌ output = %q, missing name field\n\t\t"+
			"Hint: Encoder writes directly to io.Writer — no intermediate []byte", out)
	} else {
		t.Logf("✅ EncodeStream output = %s", out)
	}
	// Verify it's valid JSON
	var check Employee
	if err := json.Unmarshal([]byte(out), &check); err != nil {
		t.Errorf("❌ output is not valid JSON: %v", err)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: StatusCode custom marshaler
// ────────────────────────────────────────────────────────────

func TestStatusCodeMarshal(t *testing.T) {
	tests := []struct {
		code StatusCode
		want string
	}{
		{200, `"OK"`},
		{404, `"Not Found"`},
		{500, `"Internal Server Error"`},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			data, err := json.Marshal(tt.code)
			if err != nil {
				t.Fatalf("❌ Marshal error: %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("❌ Marshal(%d) = %s, want %s\n\t\t"+
					"Hint: MarshalJSON must return json.Marshal(stringRepr). "+
					"The json package calls your MarshalJSON instead of default int encoding",
					tt.code, data, tt.want)
			} else {
				t.Logf("✅ Marshal(%d) = %s", tt.code, data)
			}
		})
	}
}

func TestStatusCodeUnmarshal(t *testing.T) {
	tests := []struct {
		input string
		want  StatusCode
	}{
		{`"OK"`, 200},
		{`"Not Found"`, 404},
		{`"Internal Server Error"`, 500},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var got StatusCode
			err := json.Unmarshal([]byte(tt.input), &got)
			if err != nil {
				t.Fatalf("❌ Unmarshal error: %v\n\t\t"+
					"Hint: UnmarshalJSON receives raw bytes including quotes. "+
					"json.Unmarshal(data, &s) to get the string first, then map to int", err)
			}
			if got != tt.want {
				t.Errorf("❌ Unmarshal(%s) = %d, want %d", tt.input, got, tt.want)
			} else {
				t.Logf("✅ Unmarshal(%s) = %d", tt.input, got)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: DelayParse with json.RawMessage
// ────────────────────────────────────────────────────────────

func TestDelayParse(t *testing.T) {
	input := `{"type":"employee","data":{"name":"Eve","department":"Security","salary":120000}}`
	typeName, emp, err := DelayParse(input)
	if err != nil {
		t.Fatalf("❌ error: %v\n\t\t"+
			"Hint: Define an envelope struct with Type string + Data json.RawMessage. "+
			"Unmarshal envelope first. Then unmarshal Data based on Type value", err)
	}
	if typeName != "employee" {
		t.Errorf("❌ type = %q, want \"employee\"", typeName)
	} else {
		t.Logf("✅ type = %q", typeName)
	}
	if emp.Name != "Eve" || emp.Salary != 120000 {
		t.Errorf("❌ employee = %+v, want {Eve Security 120000}\n\t\t"+
			"Hint: json.RawMessage implements json.Marshaler/Unmarshaler. "+
			"It stays as raw []byte until you unmarshal it yourself", emp)
	} else {
		t.Logf("✅ employee = %+v", emp)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: NullableFields
// ────────────────────────────────────────────────────────────

func TestNullableFields(t *testing.T) {
	t.Run("explicit_null", func(t *testing.T) {
		input := `{"name":"Frank","nickname":null}`
		p, isNull, err := NullableFields(input)
		if err != nil {
			t.Fatalf("❌ error: %v", err)
		}
		if p.Name != "Frank" {
			t.Errorf("❌ Name = %q, want \"Frank\"", p.Name)
		}
		if !isNull {
			t.Errorf("❌ isNicknameNull = false, want true\n\t\t"+
				"Hint: Unmarshal into map[string]json.RawMessage first. "+
				"If key exists and value is 'null' (4 bytes), it's explicitly null. "+
				"*string pointer alone can't distinguish null from absent", )
		} else {
			t.Logf("✅ explicit null detected correctly")
		}
	})

	t.Run("absent_field", func(t *testing.T) {
		input := `{"name":"Grace","age":30}`
		p, isNull, err := NullableFields(input)
		if err != nil {
			t.Fatalf("❌ error: %v", err)
		}
		if p.Name != "Grace" {
			t.Errorf("❌ Name = %q, want \"Grace\"", p.Name)
		}
		if isNull {
			t.Errorf("❌ isNicknameNull = true, want false (field is absent, not null)")
		} else {
			t.Logf("✅ absent field correctly returns false")
		}
		if p.Age == nil || *p.Age != 30 {
			t.Errorf("❌ Age = %v, want *30", p.Age)
		} else {
			t.Logf("✅ Age = %d", *p.Age)
		}
	})
}

// ────────────────────────────────────────────────────────────
// Exercise 10: Event custom time format
// ────────────────────────────────────────────────────────────

func TestEventMarshal(t *testing.T) {
	ts := time.Date(2025, 6, 15, 14, 30, 0, 0, time.UTC)
	e := Event{Name: "deploy", Timestamp: ts}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("❌ Marshal error: %v\n\t\t"+
			"Hint: In MarshalJSON, create an helper struct with Timestamp as string, "+
			"format with time.Format(\"2006-01-02 15:04:05\")", err)
	}
	s := string(data)
	if !strings.Contains(s, "2025-06-15 14:30:00") {
		t.Errorf("❌ timestamp format wrong: %s\n\t\t"+
			"Hint: Go's reference time is Mon Jan 2 15:04:05 MST 2006. "+
			"Use layout \"2006-01-02 15:04:05\" for custom format", s)
	} else {
		t.Logf("✅ Event marshal = %s", s)
	}
}

func TestEventUnmarshal(t *testing.T) {
	input := `{"name":"rollback","timestamp":"2025-06-15 14:30:00"}`
	var e Event
	err := json.Unmarshal([]byte(input), &e)
	if err != nil {
		t.Fatalf("❌ Unmarshal error: %v\n\t\t"+
			"Hint: In UnmarshalJSON, parse into helper struct with Timestamp string, "+
			"then time.Parse(\"2006-01-02 15:04:05\", raw)", err)
	}
	if e.Name != "rollback" {
		t.Errorf("❌ Name = %q, want \"rollback\"", e.Name)
	}
	want := time.Date(2025, 6, 15, 14, 30, 0, 0, time.UTC)
	if !e.Timestamp.Equal(want) {
		t.Errorf("❌ Timestamp = %v, want %v", e.Timestamp, want)
	} else {
		t.Logf("✅ Event unmarshal: name=%s time=%v", e.Name, e.Timestamp)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: ValidateJSON
// ────────────────────────────────────────────────────────────

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{`{"name":"test"}`, true},
		{`[1,2,3]`, true},
		{`"hello"`, true},
		{`42`, true},
		{`null`, true},
		{`{invalid}`, false},
		{`{"name":}`, false},
		{``, false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ValidateJSON(tt.input)
			if got != tt.want {
				t.Errorf("❌ ValidateJSON(%q) = %v, want %v\n\t\t"+
					"Hint: json.Valid([]byte(s)) — one line. "+
					"Returns false for empty input, malformed JSON, trailing garbage",
					tt.input, got, tt.want)
			} else {
				t.Logf("✅ ValidateJSON(%q) = %v", tt.input, got)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: PreciseNumber
// ────────────────────────────────────────────────────────────

func TestPreciseNumber(t *testing.T) {
	// 2^53 + 1 = 9007199254740993 — float64 can't represent this exactly
	input := `{"id":9007199254740993}`
	raw, val, err := PreciseNumber(input)
	if err != nil {
		t.Fatalf("❌ error: %v\n\t\t"+
			"Hint: dec := json.NewDecoder(strings.NewReader(s)); "+
			"dec.UseNumber(); decode into map[string]any; "+
			"then m[\"id\"].(json.Number).String() and .Int64()", err)
	}
	if raw != "9007199254740993" {
		t.Errorf("❌ raw = %q, want \"9007199254740993\"\n\t\t"+
			"Hint: Without UseNumber(), this becomes 9007199254740992 (float64 truncation). "+
			"json.Number preserves the exact string from the JSON source", raw)
	} else {
		t.Logf("✅ raw = %q (exact)", raw)
	}
	if val != 9007199254740993 {
		t.Errorf("❌ val = %d, want 9007199254740993", val)
	} else {
		t.Logf("✅ val = %d", val)
	}
}
