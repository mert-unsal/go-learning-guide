// Package encoding_json covers JSON encoding and decoding in Go using
// the encoding/json standard library package.
//
// Topics:
//   - Struct tags (json:"fieldName")
//   - Marshal (Go → JSON)
//   - Unmarshal (JSON → Go)
//   - Streaming: json.Encoder / json.Decoder
//   - Working with unknown structure: map[string]any
//   - json.RawMessage for deferred parsing
//   - Omitempty, nested structs, custom types
package encoding_json

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ============================================================
// 1. STRUCT TAGS — Control JSON Field Names
// ============================================================
// Struct tags tell the JSON encoder/decoder how to map fields.
//
// Syntax:  `json:"fieldName"`
//   - Use camelCase names to match typical JSON APIs
//   - "omitempty": omit the field if it has a zero value
//   - "-":         always exclude this field from JSON
//
// Go field names are exported (capitalized), but JSON typically uses camelCase.

// User represents a user in our system.
type User struct {
	ID        int    `json:"id"`        // renamed: "id"
	FirstName string `json:"firstName"` // renamed: "firstName"
	LastName  string `json:"lastName"`  // renamed: "lastName"
	Email     string `json:"email"`
	Age       int    `json:"age,omitempty"`     // omitted if Age == 0
	Password  string `json:"-"`                 // NEVER included in JSON
	IsAdmin   bool   `json:"isAdmin,omitempty"` // omitted if false
}

// Address demonstrates nested structs.
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
	Zip     string `json:"zip,omitempty"`
}

// Person demonstrates nested struct composition.
type Person struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Address Address  `json:"address"`         // nested object
	Tags    []string `json:"tags,omitempty"`  // array
	Score   *float64 `json:"score,omitempty"` // pointer: null if nil, otherwise value
}

// ============================================================
// 2. json.Marshal — Go Struct → JSON bytes
// ============================================================
// json.Marshal(v) returns ([]byte, error)
// json.MarshalIndent(v, prefix, indent) returns pretty-printed JSON

func DemonstrateMarshal() {
	u := User{
		ID:        1,
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice@example.com",
		Age:       30,
		Password:  "secret123", // will NOT appear in JSON (tag: "-")
		IsAdmin:   false,       // will NOT appear (omitempty + zero value)
	}

	// --- Compact JSON ---
	data, err := json.Marshal(u)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}
	fmt.Println("Compact JSON:")
	fmt.Println(string(data))
	// {"id":1,"firstName":"Alice","lastName":"Smith","email":"alice@example.com","age":30}

	// --- Pretty-printed JSON ---
	pretty, _ := json.MarshalIndent(u, "", "  ") // prefix="", indent="  " (2 spaces)
	fmt.Println("\nPretty JSON:")
	fmt.Println(string(pretty))

	// --- Slice of structs ---
	users := []User{
		{ID: 1, FirstName: "Alice", Email: "alice@example.com"},
		{ID: 2, FirstName: "Bob", Email: "bob@example.com", Age: 25},
	}
	usersJSON, _ := json.MarshalIndent(users, "", "  ")
	fmt.Println("\nUsers array:")
	fmt.Println(string(usersJSON))

	// --- Map to JSON ---
	config := map[string]any{
		"host":  "localhost",
		"port":  8080,
		"debug": true,
		"tags":  []string{"web", "api"},
	}
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println("\nConfig map:")
	fmt.Println(string(configJSON))
}

// ============================================================
// 3. json.Unmarshal — JSON bytes → Go Struct
// ============================================================
// json.Unmarshal(data []byte, v any) error
// v must be a POINTER — unmarshal writes into the pointed-to value.

func DemonstrateUnmarshal() {
	// --- Unmarshal object into struct ---
	jsonStr := `{
		"id": 42,
		"firstName": "Bob",
		"lastName": "Jones",
		"email": "bob@example.com",
		"age": 28
	}`

	var u User
	err := json.Unmarshal([]byte(jsonStr), &u) // pass pointer &u
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}
	fmt.Printf("Unmarshalled: %+v\n", u)
	// {ID:42 FirstName:Bob LastName:Jones Email:bob@example.com Age:28 Password: IsAdmin:false}

	// --- Unmarshal array into slice ---
	arrayJSON := `[
		{"id":1,"firstName":"Alice","email":"alice@example.com"},
		{"id":2,"firstName":"Bob","email":"bob@example.com","age":25}
	]`
	var users []User
	json.Unmarshal([]byte(arrayJSON), &users)
	fmt.Printf("Users: %d loaded\n", len(users))
	for _, u := range users {
		fmt.Printf("  %d: %s (%s)\n", u.ID, u.FirstName, u.Email)
	}

	// --- Unmarshal into map (for unknown structure) ---
	unknownJSON := `{"service":"payments","version":3,"active":true}`
	var m map[string]any // any = interface{} in older Go
	json.Unmarshal([]byte(unknownJSON), &m)
	fmt.Println("\nDynamic map:")
	for k, v := range m {
		fmt.Printf("  %s: %v (%T)\n", k, v, v)
	}
	// Note: numbers become float64 by default in map[string]any

	// --- Partial unmarshal: only specific fields ---
	type PartialUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		// Extra fields in JSON are silently ignored
	}
	var p PartialUser
	json.Unmarshal([]byte(jsonStr), &p)
	fmt.Printf("\nPartial: id=%d email=%s\n", p.ID, p.Email)
}

// ============================================================
// 4. json.Encoder / json.Decoder — STREAMING JSON (FILES, HTTP)
// ============================================================
// Use Encoder/Decoder when working with:
//   - Files (encode directly to *os.File)
//   - HTTP request/response bodies
//   - Large data streams (avoids loading everything into memory)
//
// Encoder: encodes Go values to an io.Writer
// Decoder: decodes JSON from an io.Reader, one value at a time

func DemonstrateEncoderDecoder() {
	// --- Encoder: write JSON to a file ---
	users := []User{
		{ID: 1, FirstName: "Alice", Email: "alice@example.com", Age: 30},
		{ID: 2, FirstName: "Bob", Email: "bob@example.com", Age: 25},
	}

	f, _ := os.Create("users.json")
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ") // pretty-print (optional)

	for _, u := range users {
		if err := encoder.Encode(u); err != nil { // Encode adds a newline after each value
			fmt.Println("encode error:", err)
		}
	}
	f.Close()
	fmt.Println("Written users.json")

	// --- Decoder: read JSON from a file line by line ---
	rf, _ := os.Open("users.json")
	defer rf.Close()

	decoder := json.NewDecoder(rf)
	fmt.Println("Reading users.json:")
	for decoder.More() { // More() returns true if there's another value to decode
		var u User
		if err := decoder.Decode(&u); err != nil {
			fmt.Println("decode error:", err)
			break
		}
		fmt.Printf("  Loaded: %s %s\n", u.FirstName, u.LastName)
	}
	rf.Close()

	// --- Decoder from string (simulate HTTP body) ---
	body := strings.NewReader(`{"id":99,"firstName":"Carol","email":"carol@example.com"}`)
	dec := json.NewDecoder(body)
	var carol User
	dec.Decode(&carol)
	fmt.Printf("From body: %+v\n", carol)

	os.Remove("users.json")
}

// ============================================================
// 5. WRITE AND READ A JSON FILE — Full Workflow
// ============================================================

// Config represents application configuration.
type Config struct {
	Host     string          `json:"host"`
	Port     int             `json:"port"`
	Debug    bool            `json:"debug"`
	Database DBConfig        `json:"database"`
	Features map[string]bool `json:"features,omitempty"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Password string `json:"-"` // never save password to file
}

// WriteConfig saves config to a JSON file.
func WriteConfig(path string, cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// ReadConfig loads config from a JSON file.
func ReadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read file: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal: %w", err)
	}
	return cfg, nil
}

func DemonstrateConfigFile() {
	cfg := Config{
		Host:  "0.0.0.0",
		Port:  8080,
		Debug: true,
		Database: DBConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "myapp",
			Password: "supersecret", // will be excluded (json:"-")
		},
		Features: map[string]bool{
			"darkMode":   true,
			"betaAccess": false,
		},
	}

	// Write
	if err := WriteConfig("config.json", cfg); err != nil {
		fmt.Println("write config error:", err)
		return
	}
	fmt.Println("config.json written")

	// Read back
	loaded, err := ReadConfig("config.json")
	if err != nil {
		fmt.Println("read config error:", err)
		return
	}
	fmt.Printf("Loaded config: host=%s port=%d debug=%v\n",
		loaded.Host, loaded.Port, loaded.Debug)
	fmt.Printf("DB password (should be empty): %q\n", loaded.Database.Password)

	// Show the file content
	data, _ := os.ReadFile("config.json")
	fmt.Println("File content:\n" + string(data))

	os.Remove("config.json")
}

// ============================================================
// 6. COMMON PITFALLS AND PATTERNS
// ============================================================

// NumberFromJSON shows that JSON numbers unmarshal to float64 by default
// when using map[string]any. Use json.Number or a typed struct to avoid this.
func NumberFromJSON() {
	raw := `{"count":42,"ratio":3.14}`

	// Pitfall: count is float64, not int!
	var m map[string]any
	json.Unmarshal([]byte(raw), &m)
	count := m["count"]
	fmt.Printf("count type: %T, value: %v\n", count, count) // float64 42

	// Solution 1: use UseNumber() on decoder
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.UseNumber() // numbers stay as json.Number (a string-based type)
	var m2 map[string]any
	dec.Decode(&m2)
	num := m2["count"].(json.Number)
	intVal, _ := num.Int64()
	fmt.Println("With UseNumber:", intVal) // 42

	// Solution 2 (best): use a typed struct
	type Stats struct {
		Count int     `json:"count"`
		Ratio float64 `json:"ratio"`
	}
	var s Stats
	json.Unmarshal([]byte(raw), &s)
	fmt.Printf("Typed struct: count=%d ratio=%.2f\n", s.Count, s.Ratio)
}
