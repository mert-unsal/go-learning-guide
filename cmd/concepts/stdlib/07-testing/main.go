// Testing Patterns in Go — demonstrates Go's testing idioms and best practices.
//
// This is a concepts demonstration file. Since testing functions require
// *testing.T and *testing.B, this program shows the patterns as educational
// output with runnable examples where possible.
//
// Topics:
//   - Table-driven tests (idiomatic Go)
//   - Subtests with t.Run
//   - t.Helper for clean error reporting
//   - Benchmark tests
//   - Test coverage
//   - Testing patterns quick reference
//
// Run: go run cmd/concepts/stdlib/07-testing/main.go
package main

import "fmt"

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

// Add is a trivial function used to demonstrate testing patterns.
func Add(a, b int) int { return a + b }

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Testing Patterns                        %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	demonstrateTableDriven()
	demonstrateSubtests()
	demonstrateHelper()
	demonstrateBenchmarks()
	demonstrateCoverage()
	demonstratePatterns()
}

// ============================================================
// 1. TABLE-DRIVEN TESTS — The Idiomatic Go Style
// ============================================================
// Instead of one test per case, define a slice of structs,
// each representing one input/expected pair.
// Benefits: easy to add cases, readable output, parallel-friendly.

func demonstrateTableDriven() {
	fmt.Printf("%s▸ 1. Table-Driven Tests — The Idiomatic Go Style%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Define a slice of structs, each with input/expected pairs%s\n", green, reset)
	fmt.Printf("  %s✔ Easy to add cases, readable output, parallel-friendly%s\n", green, reset)

	fmt.Println()
	fmt.Printf("  %sPattern (write in a _test.go file):%s\n", dim, reset)
	fmt.Println("  ┌──────────────────────────────────────────────────────┐")
	fmt.Println("  │ func TestAdd(t *testing.T) {                        │")
	fmt.Println("  │     tests := []struct {                             │")
	fmt.Println("  │         name     string                             │")
	fmt.Println("  │         a, b     int                                │")
	fmt.Println("  │         want     int                                │")
	fmt.Println("  │     }{                                              │")
	fmt.Println("  │         {\"positive\", 1, 2, 3},                      │")
	fmt.Println("  │         {\"negative\", -1, -1, -2},                   │")
	fmt.Println("  │         {\"zero\", 0, 0, 0},                          │")
	fmt.Println("  │     }                                               │")
	fmt.Println("  │     for _, tt := range tests {                      │")
	fmt.Println("  │         t.Run(tt.name, func(t *testing.T) {         │")
	fmt.Println("  │             got := Add(tt.a, tt.b)                  │")
	fmt.Println("  │             if got != tt.want {                     │")
	fmt.Println("  │                 t.Errorf(\"Add(%d,%d) = %d, want %d\",│")
	fmt.Println("  │                     tt.a, tt.b, got, tt.want)       │")
	fmt.Println("  │             }                                       │")
	fmt.Println("  │         })                                          │")
	fmt.Println("  │     }                                               │")
	fmt.Println("  │ }                                                   │")
	fmt.Println("  └──────────────────────────────────────────────────────┘")

	// Live demo with Add function
	fmt.Println()
	fmt.Printf("  %sLive demo with Add():%s\n", bold, reset)
	tests := []struct {
		name     string
		a, b     int
		want     int
	}{
		{"positive", 1, 2, 3},
		{"negative", -1, -1, -2},
		{"zero", 0, 0, 0},
		{"mixed", 5, -3, 2},
	}
	for _, tt := range tests {
		got := Add(tt.a, tt.b)
		status := green + "PASS" + reset
		if got != tt.want {
			status = red + "FAIL" + reset
		}
		fmt.Printf("    [%s] Add(%d, %d) = %s%d%s (want %d)\n",
			status, tt.a, tt.b, magenta, got, reset, tt.want)
	}
	fmt.Println()
}

// ============================================================
// 2. SUBTESTS — t.Run()
// ============================================================
// t.Run(name, func) creates a subtest. Benefits:
//   - Individual cases can fail without stopping others
//   - Run a specific case: go test -run TestAdd/zero
//   - Can be parallelized: call t.Parallel() inside the subtest

func demonstrateSubtests() {
	fmt.Printf("%s▸ 2. Subtests — t.Run()%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ t.Run(name, func) creates a named subtest%s\n", green, reset)
	fmt.Printf("  %s✔ Individual cases can fail without stopping others%s\n", green, reset)
	fmt.Printf("  %s✔ Run specific case: go test -run TestAdd/zero%s\n", green, reset)
	fmt.Printf("  %s✔ Parallelize: call t.Parallel() inside the subtest%s\n", green, reset)

	fmt.Println()
	fmt.Printf("  %sPattern:%s\n", dim, reset)
	fmt.Println("  ┌──────────────────────────────────────────────────────┐")
	fmt.Println("  │ for _, tc := range cases {                          │")
	fmt.Println("  │     tc := tc // capture for t.Parallel()            │")
	fmt.Println("  │     t.Run(tc.name, func(t *testing.T) {             │")
	fmt.Println("  │         t.Parallel() // run subtests in parallel    │")
	fmt.Println("  │         got := Add(tc.a, tc.b)                      │")
	fmt.Println("  │         assertEqual(t, got, tc.want)                │")
	fmt.Println("  │     })                                              │")
	fmt.Println("  │ }                                                   │")
	fmt.Println("  └──────────────────────────────────────────────────────┘")
	fmt.Println()
}

// ============================================================
// 3. TEST HELPER FUNCTIONS — t.Helper()
// ============================================================
// When you factor out assertion logic into a helper function,
// call t.Helper() so that failure messages show the CALLER's line,
// not the helper's line.

func demonstrateHelper() {
	fmt.Printf("%s▸ 3. t.Helper() — Clean Error Reporting%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Call t.Helper() in assertion helpers%s\n", green, reset)
	fmt.Printf("  %s✔ Errors point to the CALLER's line, not the helper's line%s\n", green, reset)

	fmt.Println()
	fmt.Printf("  %sPattern:%s\n", dim, reset)
	fmt.Println("  ┌──────────────────────────────────────────────────────┐")
	fmt.Println("  │ func assertEqual(t *testing.T, got, want int) {     │")
	fmt.Println("  │     t.Helper() // marks this as a helper            │")
	fmt.Println("  │     if got != want {                                │")
	fmt.Println("  │         t.Errorf(\"got %d, want %d\", got, want)      │")
	fmt.Println("  │     }                                               │")
	fmt.Println("  │ }                                                   │")
	fmt.Println("  └──────────────────────────────────────────────────────┘")
	fmt.Println()
}

// ============================================================
// 4. BENCHMARK TESTS
// ============================================================
// Benchmarks measure performance. They live in _test.go files and use *testing.B.
// Run: go test -bench=. -benchmem ./...
//
// Rules:
//   - Prefix function name with Benchmark
//   - Loop from 0 to b.N — the framework adjusts N to get stable timing
//   - Use b.ResetTimer() after setup that shouldn't be counted
//   - Use b.ReportAllocs() or pass -benchmem flag to see allocations

func demonstrateBenchmarks() {
	fmt.Printf("%s▸ 4. Benchmark Tests%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Run: go test -bench=. -benchmem ./...%s\n", green, reset)
	fmt.Printf("  %s✔ Loop 0 to b.N — framework adjusts N for stable timing%s\n", green, reset)
	fmt.Printf("  %s✔ Use b.ResetTimer() after setup%s\n", green, reset)

	fmt.Println()
	fmt.Printf("  %sPattern:%s\n", dim, reset)
	fmt.Println("  ┌──────────────────────────────────────────────────────┐")
	fmt.Println("  │ func BenchmarkAdd(b *testing.B) {                   │")
	fmt.Println("  │     for i := 0; i < b.N; i++ {                     │")
	fmt.Println("  │         Add(123456, 654321)                         │")
	fmt.Println("  │     }                                               │")
	fmt.Println("  │ }                                                   │")
	fmt.Println("  │                                                     │")
	fmt.Println("  │ func BenchmarkSliceAppend(b *testing.B) {           │")
	fmt.Println("  │     b.ReportAllocs()                                │")
	fmt.Println("  │     for i := 0; i < b.N; i++ {                     │")
	fmt.Println("  │         s := make([]int, 0)                         │")
	fmt.Println("  │         for j := 0; j < 1000; j++ {                │")
	fmt.Println("  │             s = append(s, j)                        │")
	fmt.Println("  │         }                                           │")
	fmt.Println("  │         _ = s                                       │")
	fmt.Println("  │     }                                               │")
	fmt.Println("  │ }                                                   │")
	fmt.Println("  └──────────────────────────────────────────────────────┘")
	fmt.Println()
}

// ============================================================
// 5. TEST COVERAGE
// ============================================================
// Run with coverage:
//   go test -cover ./...
//
// Generate HTML coverage report:
//   go test -coverprofile=coverage.out ./...
//   go tool cover -html=coverage.out
//
// Aim for high coverage on critical logic, but 100% isn't always practical.
// Focus on covering edge cases (empty input, single element, negatives).

func demonstrateCoverage() {
	fmt.Printf("%s▸ 5. Test Coverage%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Run: go test -cover ./...%s\n", green, reset)
	fmt.Printf("  %s✔ HTML report: go test -coverprofile=coverage.out ./...%s\n", green, reset)
	fmt.Printf("  %s✔ View: go tool cover -html=coverage.out%s\n", green, reset)
	fmt.Printf("  %s⚠ Focus on edge cases: empty input, single element, negatives%s\n", yellow, reset)
	fmt.Println()
}

// ============================================================
// 6. TESTING PATTERNS QUICK REFERENCE
// ============================================================

func demonstratePatterns() {
	fmt.Printf("%s▸ 6. Testing Patterns Quick Reference%s\n", cyan+bold, reset)

	fmt.Println()
	fmt.Printf("  %st.Error vs t.Fatal:%s\n", bold, reset)
	fmt.Printf("  %s✔ t.Error:  marks test as failed, continues execution%s\n", green, reset)
	fmt.Printf("  %s✔ t.Fatal:  marks test as failed, stops immediately%s\n", green, reset)
	fmt.Printf("  %s⚠ Use Fatal when further tests would panic (e.g., nil pointer)%s\n", yellow, reset)

	fmt.Println()
	fmt.Printf("  %st.Skip — Conditional Skipping:%s\n", bold, reset)
	fmt.Println("  ┌──────────────────────────────────────────────────────┐")
	fmt.Println("  │ if runtime.GOOS == \"windows\" {                      │")
	fmt.Println("  │     t.Skip(\"skipping on Windows\")                   │")
	fmt.Println("  │ }                                                   │")
	fmt.Println("  └──────────────────────────────────────────────────────┘")

	fmt.Println()
	fmt.Printf("  %st.Cleanup — Deferred Cleanup:%s\n", bold, reset)
	fmt.Printf("  %s✔ Register cleanup functions that run when the test ends%s\n", green, reset)
	fmt.Println("  ┌──────────────────────────────────────────────────────┐")
	fmt.Println("  │ t.Cleanup(func() {                                  │")
	fmt.Println("  │     // runs after test, even on failure             │")
	fmt.Println("  │     os.Remove(\"testfile.txt\")                       │")
	fmt.Println("  │ })                                                  │")
	fmt.Println("  └──────────────────────────────────────────────────────┘")

	fmt.Println()
	fmt.Printf("  %st.Log / t.Logf:%s\n", bold, reset)
	fmt.Printf("  %s✔ Output only shows with -v flag or on failure%s\n", green, reset)

	fmt.Println()
	fmt.Printf("  %sKey Commands:%s\n", bold, reset)
	fmt.Printf("    go test ./...                      %s# run all tests%s\n", dim, reset)
	fmt.Printf("    go test -v ./...                   %s# verbose output%s\n", dim, reset)
	fmt.Printf("    go test -race ./...                %s# race detector (NON-NEGOTIABLE)%s\n", dim, reset)
	fmt.Printf("    go test -run TestAdd/zero ./...    %s# run specific subtest%s\n", dim, reset)
	fmt.Printf("    go test -bench=. -benchmem ./...   %s# benchmarks with allocs%s\n", dim, reset)
	fmt.Printf("    go test -cover ./...               %s# coverage summary%s\n", dim, reset)
	fmt.Printf("    go test -count=1 ./...             %s# disable test caching%s\n", dim, reset)
}
