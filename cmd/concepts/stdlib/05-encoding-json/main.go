// Encoding/JSON in Go — demonstrates JSON encoding and decoding using
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
//
// Run: go run cmd/concepts/stdlib/05-encoding-json/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

// ============================================================
// STRUCT DEFINITIONS — Used across demonstrations
// ============================================================

// User represents a user in our system.
type User struct {
	ID        int    `json:"id"`                // renamed: "id"
	FirstName string `json:"firstName"`         // renamed: "firstName"
	LastName  string `json:"lastName"`          // renamed: "lastName"
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

// Config represents application configuration.
type Config struct {
	Host     string          `json:"host"`
	Port     int             `json:"port"`
	Debug    bool            `json:"debug"`
	Database DBConfig        `json:"database"`
	Features map[string]bool `json:"features,omitempty"`
}

// DBConfig represents database configuration.
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Password string `json:"-"` // never save password to file
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Encoding/JSON                           %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	demonstrateStructTags()
	demonstrateMarshal()
	demonstrateUnmarshal()
	demonstrateEncoderDecoder()
	demonstrateConfigFile()
	demonstrateNumberPitfalls()
}

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

func demonstrateStructTags() {
	fmt.Printf("%s▸ 1. Struct Tags — Control JSON Field Names%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ `json:\"fieldName\"` — rename fields for JSON%s\n", green, reset)
	fmt.Printf("  %s✔ `json:\"name,omitempty\"` — omit if zero value%s\n", green, reset)
	fmt.Printf("  %s✔ `json:\"-\"` — always exclude from JSON%s\n", green, reset)
	fmt.Printf("  %s⚠ Only exported (capitalized) fields are marshalled%s\n\n", yellow, reset)
}

// ============================================================
// 2. json.Marshal — Go Struct → JSON bytes
// ============================================================
// json.Marshal(v) returns ([]byte, error)
// json.MarshalIndent(v, prefix, indent) returns pretty-printed JSON

func demonstrateMarshal() {
	fmt.Printf("%s▸ 2. json.Marshal — Go Struct → JSON%s\n", cyan+bold, reset)

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
		fmt.Println("  marshal error:", err)
		return
	}
	fmt.Println("  Compact JSON:")
	fmt.Println(" ", string(data))

	// --- Pretty-printed JSON ---
	pretty, _ := json.MarshalIndent(u, "  ", "  ") // prefix="  ", indent="  " (2 spaces)
	fmt.Println("\n  Pretty JSON:")
	fmt.Println(" ", string(pretty))

	// --- Slice of structs ---
	users := []User{
		{ID: 1, FirstName: "Alice", Email: "alice@example.com"},
		{ID: 2, FirstName: "Bob", Email: "bob@example.com", Age: 25},
	}
	usersJSON, _ := json.MarshalIndent(users, "  ", "  ")
	fmt.Println("\n  Users array:")
	fmt.Println(" ", string(usersJSON))

	// --- Map to JSON ---
	config := map[string]any{
		"host":  "localhost",
		"port":  8080,
		"debug": true,
		"tags":  []string{"web", "api"},
	}
	configJSON, _ := json.MarshalIndent(config, "  ", "  ")
	fmt.Println("\n  Config map:")
	fmt.Println(" ", string(configJSON))
	fmt.Println()
}

// ============================================================
// 3. json.Unmarshal — JSON bytes → Go Struct
// ============================================================
// json.Unmarshal(data []byte, v any) error
// v must be a POINTER — unmarshal writes into the pointed-to value.

func demonstrateUnmarshal() {
	fmt.Printf("%s▸ 3. json.Unmarshal — JSON → Go Struct%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ v must be a POINTER — unmarshal writes into the pointed-to value%s\n", green, reset)

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
		fmt.Println("  unmarshal error:", err)
		return
	}
	fmt.Printf("  Unmarshalled: %+v\n", u)

	// --- Unmarshal array into slice ---
	arrayJSON := `[
		{"id":1,"firstName":"Alice","email":"alice@example.com"},
		{"id":2,"firstName":"Bob","email":"bob@example.com","age":25}
	]`
	var users []User
	json.Unmarshal([]byte(arrayJSON), &users)
	fmt.Printf("  Users: %d loaded\n", len(users))
	for _, u := range users {
		fmt.Printf("    %d: %s (%s)\n", u.ID, u.FirstName, u.Email)
	}

	// --- Unmarshal into map (for unknown structure) ---
	unknownJSON := `{"service":"payments","version":3,"active":true}`
	var m map[string]any // any = interface{} in older Go
	json.Unmarshal([]byte(unknownJSON), &m)
	fmt.Println("\n  Dynamic map:")
	for k, v := range m {
		fmt.Printf("    %s: %v (%T)\n", k, v, v)
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
	fmt.Printf("\n  Partial: id=%d email=%s\n", p.ID, p.Email)
	fmt.Println()
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

func demonstrateEncoderDecoder() {
	fmt.Printf("%s▸ 4. json.Encoder / json.Decoder — Streaming JSON%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Use for files, HTTP bodies, and large data streams%s\n", green, reset)

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
			fmt.Println("  encode error:", err)
		}
	}
	f.Close()
	fmt.Println("  Written users.json")

	// --- Decoder: read JSON from a file line by line ---
	rf, _ := os.Open("users.json")
	defer rf.Close()

	decoder := json.NewDecoder(rf)
	fmt.Println("  Reading users.json:")
	for decoder.More() { // More() returns true if there's another value to decode
		var u User
		if err := decoder.Decode(&u); err != nil {
			fmt.Println("  decode error:", err)
			break
		}
		fmt.Printf("    Loaded: %s %s\n", u.FirstName, u.LastName)
	}
	rf.Close()

	// --- Decoder from string (simulate HTTP body) ---
	body := strings.NewReader(`{"id":99,"firstName":"Carol","email":"carol@example.com"}`)
	dec := json.NewDecoder(body)
	var carol User
	dec.Decode(&carol)
	fmt.Printf("  From body: %+v\n", carol)

	os.Remove("users.json")
	fmt.Println()
}

// ============================================================
// 5. WRITE AND READ A JSON FILE — Full Workflow
// ============================================================

// writeConfig saves config to a JSON file.
func writeConfig(path string, cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// readConfig loads config from a JSON file.
func readConfig(path string) (Config, error) {
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

func demonstrateConfigFile() {
	fmt.Printf("%s▸ 5. JSON Config File — Full Workflow%s\n", cyan+bold, reset)

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
	if err := writeConfig("config.json", cfg); err != nil {
		fmt.Println("  write config error:", err)
		return
	}
	fmt.Println("  config.json written")

	// Read back
	loaded, err := readConfig("config.json")
	if err != nil {
		fmt.Println("  read config error:", err)
		return
	}
	fmt.Printf("  Loaded config: host=%s port=%d debug=%v\n",
		loaded.Host, loaded.Port, loaded.Debug)
	fmt.Printf("  DB password (should be empty): %q\n", loaded.Database.Password)

	// Show the file content
	data, _ := os.ReadFile("config.json")
	fmt.Println("  File content:")
	fmt.Println(" ", string(data))

	os.Remove("config.json")
}

// ============================================================
// 6. COMMON PITFALLS AND PATTERNS
// ============================================================

func demonstrateNumberPitfalls() {
	fmt.Printf("%s▸ 6. Common Pitfalls — JSON Number Types%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ JSON numbers unmarshal to float64 by default in map[string]any%s\n", yellow, reset)

	raw := `{"count":42,"ratio":3.14}`

	// Pitfall: count is float64, not int!
	var m map[string]any
	json.Unmarshal([]byte(raw), &m)
	count := m["count"]
	fmt.Printf("  count type: %T, value: %v\n", count, count) // float64 42

	// Solution 1: use UseNumber() on decoder
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.UseNumber() // numbers stay as json.Number (a string-based type)
	var m2 map[string]any
	dec.Decode(&m2)
	num := m2["count"].(json.Number)
	intVal, _ := num.Int64()
	fmt.Println("  With UseNumber:", intVal) // 42

	// Solution 2 (best): use a typed struct
	fmt.Printf("  %s✔ Best solution: use a typed struct%s\n", green, reset)
	type Stats struct {
		Count int     `json:"count"`
		Ratio float64 `json:"ratio"`
	}
	var s Stats
	json.Unmarshal([]byte(raw), &s)
	fmt.Printf("  Typed struct: count=%d ratio=%.2f\n", s.Count, s.Ratio)
}
