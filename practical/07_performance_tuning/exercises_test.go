package performance_tuning

import (
	"runtime"
	"testing"
	"unsafe"
)

// ============================================================
// BENCHMARKS — Exercise 1: Escape Analysis
// ============================================================
// Run: go test -bench=BenchmarkSum -benchmem ./practical/07_performance_tuning/

func BenchmarkSumToString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SumToString(42, 58)
	}
}

func BenchmarkSumToStringFixed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SumToStringFixed(42, 58)
	}
}

func BenchmarkNewPoint(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewPoint(1.0, 2.0)
	}
}

func BenchmarkNewPointFixed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewPointFixed(1.0, 2.0)
	}
}

func BenchmarkDistanceLabel(b *testing.B) {
	p1 := Point{1.0, 2.0}
	p2 := Point{4.0, 6.0}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DistanceLabel(p1, p2)
	}
}

func BenchmarkDistanceLabelFixed(b *testing.B) {
	p1 := Point{1.0, 2.0}
	p2 := Point{4.0, 6.0}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DistanceLabelFixed(p1, p2)
	}
}

// ============================================================
// BENCHMARKS — Exercise 2: String Building
// ============================================================
// Run: go test -bench=BenchmarkBuild -benchmem ./practical/07_performance_tuning/

func BenchmarkBuildCSVConcat10(b *testing.B) {
	fields := GenerateFields(10)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildCSVConcat(fields)
	}
}

func BenchmarkBuildCSVConcat100(b *testing.B) {
	fields := GenerateFields(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildCSVConcat(fields)
	}
}

func BenchmarkBuildCSVBuilder10(b *testing.B) {
	fields := GenerateFields(10)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildCSVBuilder(fields)
	}
}

func BenchmarkBuildCSVBuilder100(b *testing.B) {
	fields := GenerateFields(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildCSVBuilder(fields)
	}
}

func BenchmarkBuildCSVBuilderPrealloc10(b *testing.B) {
	fields := GenerateFields(10)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildCSVBuilderPrealloc(fields)
	}
}

func BenchmarkBuildCSVBuilderPrealloc100(b *testing.B) {
	fields := GenerateFields(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildCSVBuilderPrealloc(fields)
	}
}

// ============================================================
// BENCHMARKS — Exercise 3: Slice Pre-allocation
// ============================================================
// Run: go test -bench=BenchmarkCollect -benchmem ./practical/07_performance_tuning/

func BenchmarkCollectEvens1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = CollectEvens(1000)
	}
}

func BenchmarkCollectEvensPrealloc1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = CollectEvensPrealloc(1000)
	}
}

func BenchmarkCollectEvens100000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = CollectEvens(100000)
	}
}

func BenchmarkCollectEvensPrealloc100000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = CollectEvensPrealloc(100000)
	}
}

// ============================================================
// BENCHMARKS — Exercise 4: sync.Pool
// ============================================================
// Run: go test -bench=BenchmarkFormat -benchmem ./practical/07_performance_tuning/

func BenchmarkFormatRecords(b *testing.B) {
	records := GenerateRecords(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatRecords(records)
	}
}

func BenchmarkFormatRecordsPooled(b *testing.B) {
	records := GenerateRecords(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatRecordsPooled(records)
	}
}

// ============================================================
// BENCHMARKS — Exercise 5: Interface Boxing
// ============================================================
// Run: go test -bench=BenchmarkProcess -benchmem ./practical/07_performance_tuning/

func BenchmarkProcessValues(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ProcessValues(nums)
	}
}

func BenchmarkProcessValuesDirect(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ProcessValuesDirect(nums)
	}
}

// ============================================================
// BENCHMARKS — Exercise 6: Map Pre-allocation
// ============================================================
// Run: go test -bench=BenchmarkGroupBy -benchmem ./practical/07_performance_tuning/

func BenchmarkGroupByPrefix(b *testing.B) {
	words := GenerateWords(1000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GroupByPrefix(words, 3)
	}
}

func BenchmarkGroupByPrefixOptimized(b *testing.B) {
	words := GenerateWords(1000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GroupByPrefixOptimized(words, 3)
	}
}

// ============================================================
// BENCHMARKS — Exercise 7: GC Pressure (Pointers vs Values)
// ============================================================
// Run: go test -bench=BenchmarkCreateUsers -benchmem ./practical/07_performance_tuning/
//
// Watch for: allocs/op difference — fewer pointers = less GC scanning work.
// For a real-world comparison, also run with GODEBUG=gctrace=1.

func BenchmarkCreateUsersWithPointers(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		users := CreateUsersWithPointers(10000)
		runtime.KeepAlive(users)
	}
}

func BenchmarkCreateUsersWithValues(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		users := CreateUsersWithValues(10000)
		runtime.KeepAlive(users)
	}
}

// ============================================================
// TEST — Exercise 8: Struct Padding
// ============================================================
// Run: go test -run TestStructSizes -v ./practical/07_performance_tuning/

func TestStructSizes(t *testing.T) {
	bad := unsafe.Sizeof(BadLayout{})
	good := unsafe.Sizeof(GoodLayout{})

	t.Logf("BadLayout  size: %d bytes (fields: bool, int64, bool, int32, bool)", bad)
	t.Logf("GoodLayout size: %d bytes (same fields, reordered)", good)

	if bad <= good {
		t.Errorf("❌ GoodLayout (%d bytes) should be smaller than BadLayout (%d bytes)\n"+
			"   Hint: order fields from largest to smallest alignment", good, bad)
	} else {
		saved := bad - good
		t.Logf("✅ Saved %d bytes per struct by reordering fields", saved)
		t.Logf("   At 1M structs, that's %d MB saved", saved*1_000_000/1024/1024)
	}
}
