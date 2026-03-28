package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof/* endpoints
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ══════════════════════════════════════════════════════════════
// ANSI color helpers
// ══════════════════════════════════════════════════════════════

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
	bgBlue  = "\033[44m"
)

func banner(text string) {
	fmt.Printf("\n%s%s %-70s %s\n", bold, bgBlue, text, reset)
}

func info(label, msg string) {
	fmt.Printf("  %s%s%-22s%s %s\n", cyan, bold, label, reset, msg)
}

func step(msg string) {
	fmt.Printf("  %s▸ %s%s\n", green, msg, reset)
}

func warn(msg string) {
	fmt.Printf("  %s⚠ %s%s\n", yellow, msg, reset)
}

func cmd(command string) {
	fmt.Printf("    %s$ %s%s\n", magenta, command, reset)
}

// ══════════════════════════════════════════════════════════════
// Workload generators — each targets a different profile type
// ══════════════════════════════════════════════════════════════

// cpuHog burns CPU doing useless string work.
// Shows up in: CPU profile
func cpuHog(done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			s := ""
			for i := 0; i < 1000; i++ {
				s += "x" // O(n²) — intentionally bad
			}
			_ = s
		}
	}
}

// heapSpammer allocates slices that escape to heap.
// Shows up in: heap profile (alloc_space, alloc_objects)
func heapSpammer(done <-chan struct{}) {
	var keep [][]byte
	for {
		select {
		case <-done:
			return
		default:
			buf := make([]byte, 4096) // 4KB per alloc
			buf[0] = byte(rand.IntN(256))
			keep = append(keep, buf)
			if len(keep) > 5000 {
				keep = keep[2500:] // let GC reclaim half
			}
			time.Sleep(100 * time.Microsecond)
		}
	}
}

// channelBlocker creates goroutines that block on channels.
// Shows up in: block profile, goroutine profile
func channelBlocker(done <-chan struct{}) {
	ch := make(chan int) // unbuffered — every send blocks until receive
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case ch <- id:
					// sent — other side will receive
				}
			}
		}(i)
	}

	// Slow consumer — deliberate bottleneck
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ch:
				time.Sleep(10 * time.Millisecond) // slow drain = senders block
			}
		}
	}()

	<-done
	wg.Wait()
}

// mutexContender creates goroutines fighting over a lock.
// Shows up in: mutex profile
func mutexContender(done <-chan struct{}, counter *int64) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					mu.Lock()
					atomic.AddInt64(counter, 1)
					time.Sleep(500 * time.Microsecond) // hold lock — creates contention
					mu.Unlock()
				}
			}
		}()
	}

	<-done
	wg.Wait()
}

// ══════════════════════════════════════════════════════════════
// Main — orchestrate all workloads and serve pprof
// ══════════════════════════════════════════════════════════════

func main() {
	banner("Go Profiling Demo — All Profile Types in One Program")
	fmt.Println()

	info("What this does:", "Runs 4 intentionally problematic workloads simultaneously")
	info("", "then serves pprof endpoints so you can diagnose each one.")
	fmt.Println()

	// ── Step 1: Enable block and mutex profiling ──
	step("Enabling block profiling (rate=1 — record all blocks)")
	runtime.SetBlockProfileRate(1)

	step("Enabling mutex profiling (fraction=1 — record all contention)")
	runtime.SetMutexProfileFraction(1)
	fmt.Println()

	// ── Step 2: Start pprof HTTP server ──
	pprofPort := "6060"
	if p := os.Getenv("PPROF_PORT"); p != "" {
		pprofPort = p
	}

	step(fmt.Sprintf("Starting pprof server on :%s", pprofPort))
	go func() {
		if err := http.ListenAndServe(":"+pprofPort, nil); err != nil {
			fmt.Printf("  %spprof server error: %v%s\n", red, err, reset)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	baseURL := "http://localhost:" + pprofPort
	info("pprof index:", baseURL+"/debug/pprof/")
	fmt.Println()

	// ── Step 3: Launch all workloads ──
	done := make(chan struct{})
	duration := 15 * time.Second
	var mutexOps int64

	banner("Launching Workloads")
	fmt.Println()

	step("CPU hog       — O(n²) string concatenation in tight loop")
	go cpuHog(done)

	step("Heap spammer  — 4KB allocations, keeps 5000 alive, lets GC churn")
	go heapSpammer(done)

	step("Channel block — 10 senders → 1 slow consumer on unbuffered channel")
	go channelBlocker(done)

	step("Mutex fight   — 20 goroutines contending on one sync.Mutex")
	go mutexContender(done, &mutexOps)

	fmt.Println()
	warn(fmt.Sprintf("Workloads running for %v — use the commands below NOW", duration))

	// ── Step 4: Print profiling commands ──
	banner("Profile Commands — Run These in Another Terminal")
	fmt.Println()

	profiles := []struct {
		name    string
		desc    string
		command string
	}{
		{
			"CPU Profile",
			"Where is CPU time going? → finds cpuHog's O(n²) string concat",
			fmt.Sprintf("go tool pprof %s/debug/pprof/profile?seconds=5", baseURL),
		},
		{
			"Heap (alloc_space)",
			"Total bytes allocated → finds heapSpammer's 4KB allocs",
			fmt.Sprintf("go tool pprof -alloc_space %s/debug/pprof/heap", baseURL),
		},
		{
			"Heap (alloc_objects)",
			"Object count → which functions create GC pressure",
			fmt.Sprintf("go tool pprof -alloc_objects %s/debug/pprof/heap", baseURL),
		},
		{
			"Heap (inuse_space)",
			"What's on heap RIGHT NOW → find memory leaks",
			fmt.Sprintf("go tool pprof -inuse_space %s/debug/pprof/heap", baseURL),
		},
		{
			"Goroutine Dump",
			"Where are goroutines? → finds 10 blocked senders + workers",
			fmt.Sprintf("go tool pprof %s/debug/pprof/goroutine", baseURL),
		},
		{
			"Block Profile",
			"Where do goroutines WAIT? → finds channel send blocks",
			fmt.Sprintf("go tool pprof %s/debug/pprof/block", baseURL),
		},
		{
			"Mutex Profile",
			"Lock contention → finds 20 goroutines fighting for mu",
			fmt.Sprintf("go tool pprof %s/debug/pprof/mutex", baseURL),
		},
		{
			"Execution Trace",
			"Visual timeline: scheduler, GC, goroutines → open in browser",
			fmt.Sprintf("curl -o trace.out %s/debug/pprof/trace?seconds=5 && go tool trace trace.out", baseURL),
		},
	}

	for _, p := range profiles {
		fmt.Printf("  %s%s%s — %s\n", bold+cyan, p.name, reset, p.desc)
		cmd(p.command)
		fmt.Println()
	}

	// ── Step 5: Also print the human-readable goroutine dump URL ──
	info("Quick goroutine dump:", fmt.Sprintf("curl %s/debug/pprof/goroutine?debug=2", baseURL))
	fmt.Println()

	// ── Step 6: Live status ticker ──
	banner("Live Status")
	fmt.Println()

	ticker := time.NewTicker(3 * time.Second)
	deadline := time.After(duration)
	startTime := time.Now()

	var memStats runtime.MemStats
	for {
		select {
		case <-ticker.C:
			runtime.ReadMemStats(&memStats)
			elapsed := time.Since(startTime).Round(time.Second)
			remaining := (duration - elapsed).Round(time.Second)

			numG := runtime.NumGoroutine()
			heapMB := float64(memStats.HeapAlloc) / 1024 / 1024
			gcRuns := memStats.NumGC
			ops := atomic.LoadInt64(&mutexOps)

			bar := strings.Repeat("█", int(elapsed.Seconds()))
			bar += strings.Repeat("░", int(remaining.Seconds()))

			fmt.Printf("  %s[%s]%s  Goroutines: %s%d%s | Heap: %s%.1f MB%s | GC runs: %s%d%s | Mutex ops: %s%d%s | %v remaining\n",
				blue, bar, reset,
				yellow, numG, reset,
				magenta, heapMB, reset,
				cyan, gcRuns, reset,
				green, ops, reset,
				remaining,
			)

		case <-deadline:
			ticker.Stop()
			close(done)
			time.Sleep(500 * time.Millisecond) // let goroutines drain

			fmt.Println()
			banner("Done — Workloads Stopped")
			fmt.Println()

			runtime.ReadMemStats(&memStats)
			info("Final goroutines:", fmt.Sprintf("%d", runtime.NumGoroutine()))
			info("Total GC runs:", fmt.Sprintf("%d", memStats.NumGC))
			info("Total mutex ops:", fmt.Sprintf("%d", atomic.LoadInt64(&mutexOps)))
			info("Total allocs:", fmt.Sprintf("%d objects, %.1f MB", memStats.Mallocs, float64(memStats.TotalAlloc)/1024/1024))
			fmt.Println()

			step("Key pprof commands inside the interactive shell:")
			cmd("top 10          — top 10 functions by cost")
			cmd("list FuncName   — source-annotated view of a function")
			cmd("web             — open flame graph in browser (needs graphviz)")
			cmd("peek FuncName   — callers and callees of a function")
			fmt.Println()
			return
		}
	}
}
