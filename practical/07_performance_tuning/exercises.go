// Package performance_tuning contains deliberately suboptimal code for
// profiling practice. Each file targets a different performance issue
// and profiling tool. Your job: benchmark → profile → diagnose → fix → verify.
package performance_tuning

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// ============================================================
// EXERCISE 1: Escape Analysis — Stop Unnecessary Heap Allocations
// ============================================================
// Run: go build -gcflags='-m' ./practical/07_performance_tuning/
// Look for: "escapes to heap" lines
//
// These functions have HIDDEN heap allocations that could be avoided.
// Your job: identify WHY each one escapes, then write a fixed version
// that keeps values on the stack.

// SumToString converts a sum result to a string.
// PROBLEM: This allocates more than necessary. Why?
func SumToString(a, b int) string {
	result := a + b
	return fmt.Sprintf("sum=%d", result)
}

// SumToStringFixed is your optimized version.
// HINT: fmt.Sprintf uses interface{} args → escape. Use strconv instead.
func SumToStringFixed(a, b int) string {
	// TODO: implement without fmt.Sprintf
	return ""
}

// ContainsAny checks if any of the needles exist in the haystack.
// PROBLEM: The slice literal in the loop creates allocations.
func ContainsAny(haystack string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// Point represents a 2D point.
type Point struct {
	X, Y float64
}

// NewPoint creates a new point.
// PROBLEM: Returns a pointer → forces heap allocation.
// Is the pointer necessary here? Point is only 16 bytes.
func NewPoint(x, y float64) *Point {
	return &Point{X: x, Y: y}
}

// NewPointFixed is your optimized version.
// HINT: Return by value for small structs (≤ ~64 bytes).
func NewPointFixed(x, y float64) Point {
	// TODO: implement — return by value
	return Point{}
}

// DistanceLabel returns a formatted label for two points.
// PROBLEM: Multiple allocations hidden in fmt and string concat.
func DistanceLabel(p1, p2 Point) string {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	dist := fmt.Sprintf("%.2f", dx*dx+dy*dy)
	return "distance=" + dist
}

// DistanceLabelFixed is your optimized version.
// HINT: Use strconv.AppendFloat with a pre-allocated buffer.
func DistanceLabelFixed(p1, p2 Point) string {
	// TODO: implement with strconv.AppendFloat
	return ""
}

// ============================================================
// EXERCISE 2: String Building — The O(n²) Trap
// ============================================================
// Benchmark: go test -bench=BenchmarkBuild -benchmem ./practical/07_performance_tuning/
//
// Three implementations of the same thing: build a CSV line.
// Compare their allocation counts and throughput.

// BuildCSVConcat builds a CSV line using string concatenation.
// PROBLEM: Each += allocates a new string. O(n²) total allocations.
func BuildCSVConcat(fields []string) string {
	result := ""
	for i, f := range fields {
		if i > 0 {
			result += ","
		}
		result += f
	}
	return result
}

// BuildCSVBuilder builds a CSV line using strings.Builder.
// BETTER: strings.Builder avoids repeated allocations.
func BuildCSVBuilder(fields []string) string {
	var b strings.Builder
	for i, f := range fields {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(f)
	}
	return b.String()
}

// BuildCSVBuilderPrealloc builds a CSV line with pre-allocated builder.
// BEST: Pre-calculate total size, allocate once.
func BuildCSVBuilderPrealloc(fields []string) string {
	// TODO: implement
	// HINT: Calculate total length first, then b.Grow(totalLen)
	return ""
}

// ============================================================
// EXERCISE 3: Slice Growth — Pre-allocation Matters
// ============================================================
// Benchmark: go test -bench=BenchmarkCollect -benchmem ./practical/07_performance_tuning/

// CollectEvens returns all even numbers from 0 to n.
// PROBLEM: append triggers multiple growslice calls.
func CollectEvens(n int) []int {
	var result []int // nil slice, cap=0
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			result = append(result, i)
		}
	}
	return result
}

// CollectEvensPrealloc is your optimized version.
// HINT: You know roughly how many evens there will be.
func CollectEvensPrealloc(n int) []int {
	// TODO: implement with make([]int, 0, expectedSize)
	return nil
}

// ============================================================
// EXERCISE 4: sync.Pool — Reuse Expensive Buffers
// ============================================================
// Benchmark: go test -bench=BenchmarkFormat -benchmem ./practical/07_performance_tuning/

// FormatRecords formats a batch of records into a single string.
// PROBLEM: Allocates a new bytes.Buffer for every call.
// In a hot path (10k+ calls/sec), this creates massive GC pressure.
func FormatRecords(records []Record) string {
	buf := new(bytes.Buffer) // fresh allocation every call
	for _, r := range records {
		fmt.Fprintf(buf, "%d:%s;", r.ID, r.Name)
	}
	return buf.String()
}

// Record is a simple data record.
type Record struct {
	ID   int
	Name string
}

// bufferPool is a sync.Pool for reusing bytes.Buffer.
// TODO: initialize this pool with a New function that returns *bytes.Buffer
var bufferPool = sync.Pool{
	// TODO: set New function
}

// FormatRecordsPooled is your optimized version using sync.Pool.
// HINT: Get from pool, Reset, use, convert to string, Put back.
// CAREFUL: buf.String() copies the bytes — safe to reuse after.
func FormatRecordsPooled(records []Record) string {
	// TODO: implement using bufferPool
	return ""
}

// ============================================================
// EXERCISE 5: Interface Allocation — The Hidden Cost
// ============================================================
// Benchmark: go test -bench=BenchmarkProcess -benchmem ./practical/07_performance_tuning/

// ProcessValues sums values by converting them through interface{}.
// PROBLEM: Each int → interface{} conversion allocates (boxing).
func ProcessValues(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += toInt(n) // n gets boxed into interface{}
	}
	return sum
}

func toInt(v interface{}) int {
	return v.(int) // type assertion unboxes
}

// ProcessValuesDirect is your optimized version.
// HINT: Don't use interface{} when you know the type.
func ProcessValuesDirect(nums []int) int {
	// TODO: implement without interface conversion
	return 0
}

// ============================================================
// EXERCISE 6: Map Pre-allocation & Efficient Key Building
// ============================================================
// Benchmark: go test -bench=BenchmarkGroupBy -benchmem ./practical/07_performance_tuning/

// GroupByPrefix groups strings by their first N characters.
// PROBLEM: (a) map not pre-allocated, (b) fmt.Sprintf for key extraction.
func GroupByPrefix(words []string, prefixLen int) map[string][]string {
	groups := map[string][]string{} // no size hint
	for _, w := range words {
		if len(w) >= prefixLen {
			key := fmt.Sprintf("%s", w[:prefixLen]) // unnecessary Sprintf!
			groups[key] = append(groups[key], w)
		}
	}
	return groups
}

// GroupByPrefixOptimized is your optimized version.
// HINT: (a) make(map, len/expectedBuckets), (b) direct string slicing.
func GroupByPrefixOptimized(words []string, prefixLen int) map[string][]string {
	// TODO: implement with pre-allocation and no fmt.Sprintf
	return nil
}

// ============================================================
// EXERCISE 7: GC Pressure — Struct of Pointers vs Values
// ============================================================
// The GC must scan every pointer in the heap. Reducing pointers
// reduces GC work. This is critical at scale.

// UserWithPointers has pointer fields — GC must scan all of them.
type UserWithPointers struct {
	Name   *string
	Email  *string
	Age    *int
	Active *bool
	Score  *float64
	Tags   []*string
}

// UserWithValues embeds values directly — GC scans fewer pointers.
type UserWithValues struct {
	Name   string
	Email  string
	Age    int
	Active bool
	Score  float64
	Tags   []string // slice header has one pointer, not N pointers
}

// CreateUsersWithPointers creates n users with pointer fields.
func CreateUsersWithPointers(n int) []*UserWithPointers {
	users := make([]*UserWithPointers, n)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("user_%d", i)
		email := fmt.Sprintf("user_%d@test.com", i)
		age := 20 + (i % 50)
		active := i%2 == 0
		score := float64(i) * 1.5
		users[i] = &UserWithPointers{
			Name:   &name,
			Email:  &email,
			Age:    &age,
			Active: &active,
			Score:  &score,
		}
	}
	return users
}

// CreateUsersWithValues creates n users with value fields.
func CreateUsersWithValues(n int) []UserWithValues {
	users := make([]UserWithValues, n)
	for i := 0; i < n; i++ {
		users[i] = UserWithValues{
			Name:   fmt.Sprintf("user_%d", i),
			Email:  fmt.Sprintf("user_%d@test.com", i),
			Age:    20 + (i % 50),
			Active: i%2 == 0,
			Score:  float64(i) * 1.5,
		}
	}
	return users
}

// ============================================================
// EXERCISE 8: Struct Padding — Field Ordering Matters
// ============================================================
// Run: go test -run TestStructSizes ./practical/07_performance_tuning/

// BadLayout wastes memory due to padding between fields.
// The compiler aligns fields to their natural alignment boundary.
type BadLayout struct {
	A bool  // 1 byte + 7 bytes padding (to align B)
	B int64 // 8 bytes
	C bool  // 1 byte + 3 bytes padding (to align D)
	D int32 // 4 bytes
	E bool  // 1 byte + 7 bytes padding (to align next field / struct size)
	// Total: 1+7 + 8 + 1+3 + 4 + 1+7 = 32 bytes
}

// GoodLayout minimizes padding by ordering fields largest-first.
// TODO: Reorder the SAME fields to minimize size.
// Target: 24 bytes (or less).
type GoodLayout struct {
	// TODO: reorder A, B, C, D, E fields to minimize padding
	B int64
	D int32
	A bool
	C bool
	E bool
	// Total should be: 8 + 4 + 1 + 1 + 1 + 3(padding) = 18 → rounded to 24
}

// ============================================================
// HELPERS for exercises
// ============================================================

// GenerateWords creates n random-ish words for benchmarks.
func GenerateWords(n int) []string {
	words := make([]string, n)
	prefixes := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	for i := 0; i < n; i++ {
		words[i] = prefixes[i%len(prefixes)] + "_" + strconv.Itoa(i)
	}
	return words
}

// GenerateRecords creates n records for benchmarks.
func GenerateRecords(n int) []Record {
	records := make([]Record, n)
	for i := 0; i < n; i++ {
		records[i] = Record{ID: i, Name: "item_" + strconv.Itoa(i)}
	}
	return records
}

// GenerateFields creates n string fields for CSV benchmarks.
func GenerateFields(n int) []string {
	fields := make([]string, n)
	for i := 0; i < n; i++ {
		fields[i] = "field_" + strconv.Itoa(i)
	}
	return fields
}
