// Package testing_pkg demonstrates Go's testing patterns.
// This is a concepts file — to see tests in action, look at any *_test.go file.
//
// Topics:
//   - Table-driven tests (idiomatic Go)
//   - Subtests with t.Run
//   - t.Helper for clean error reporting
//   - Benchmark tests
//   - Test coverage
package testing_pkg

import (
	"fmt"
	"testing"
)

// ============================================================
// 1. TABLE-DRIVEN TESTS — The Idiomatic Go Style
// ============================================================
// Instead of one test per case, define a slice of structs,
// each representing one input/expected pair.
// Benefits: easy to add cases, readable output, parallel-friendly.
//
// Example structure (write in a _test.go file):
//
//  func TestAdd(t *testing.T) {
//      tests := []struct {
//          name     string  // identifies the test case in output
//          a, b     int
//          want     int
//      }{
//          {"positive", 1, 2, 3},
//          {"negative", -1, -1, -2},
//          {"zero", 0, 0, 0},
//      }
//      for _, tt := range tests {
//          t.Run(tt.name, func(t *testing.T) {
//              got := Add(tt.a, tt.b)
//              if got != tt.want {
//                  t.Errorf("Add(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.want)
//              }
//          })
//      }
//  }

// Add is a trivial function used in testing examples below.
func Add(a, b int) int { return a + b }

// ============================================================
// 2. TEST HELPER FUNCTIONS — t.Helper()
// ============================================================
// When you factor out assertion logic into a helper function,
// call t.Helper() so that failure messages show the CALLER's line,
// not the helper's line.

// assertEqual is a reusable assertion helper.
// Call t.Helper() immediately — this makes the error point to the test line.
func assertEqual(t *testing.T, got, want int) {
	t.Helper() // marks this as a helper — errors point to the caller
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

// ============================================================
// 3. SUBTESTS — t.Run()
// ============================================================
// t.Run(name, func) creates a subtest. Benefits:
//   - Individual cases can fail without stopping others
//   - Run a specific case: go test -run TestAdd/zero
//   - Can be parallelized: call t.Parallel() inside the subtest

// DemonstrateSubtests shows the t.Run pattern (would live in a _test.go).
func DemonstrateSubtests(t *testing.T) {
	cases := []struct {
		name       string
		a, b, want int
	}{
		{"positive", 1, 2, 3},
		{"zero", 0, 0, 0},
	}
	for _, tc := range cases {
		tc := tc // capture for t.Parallel()
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // run subtests in parallel
			got := Add(tc.a, tc.b)
			assertEqual(t, got, tc.want)
		})
	}
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
//
// Example (would live in a _test.go file):
//
//  func BenchmarkAdd(b *testing.B) {
//      for i := 0; i < b.N; i++ {
//          Add(123456, 654321)
//      }
//  }
//
//  func BenchmarkSliceAppend(b *testing.B) {
//      b.ReportAllocs()
//      for i := 0; i < b.N; i++ {
//          s := make([]int, 0)
//          for j := 0; j < 1000; j++ {
//              s = append(s, j)
//          }
//          _ = s
//      }
//  }

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

// ============================================================
// 6. TESTING PATTERNS QUICK REFERENCE
// ============================================================

// TestingPatterns demonstrates common testing idioms (call from a test).
func TestingPatterns(t *testing.T) {
	// --- t.Error vs t.Fatal ---
	// t.Error:  marks test as failed, continues execution
	// t.Fatal:  marks test as failed, stops execution immediately (good for nil checks)
	result := Add(1, 1)
	if result != 2 {
		t.Fatalf("Add(1,1) = %d; cannot continue", result) // use Fatal if further tests would panic
	}

	// --- t.Skip ---
	// Skip a test conditionally (e.g., based on OS or environment variable)
	// if runtime.GOOS == "windows" {
	//     t.Skip("skipping on Windows")
	// }

	// --- t.Cleanup ---
	// Register cleanup functions that run when the test ends (like defer for tests).
	// Useful for releasing resources.
	t.Cleanup(func() {
		fmt.Println("cleanup runs after test, even on failure")
	})

	// --- t.Log / t.Logf ---
	// Logged output only shows up with -v flag or on test failure.
	t.Logf("result = %d", result)
}
