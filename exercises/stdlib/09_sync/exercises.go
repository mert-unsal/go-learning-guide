package sync_exercises

import (
	"io"
	"sync"
)

// ============================================================
// EXERCISES — 09 sync: Mutexes, Pools, Once & Atomic Patterns
// ============================================================
// 12 exercises covering the sync package at production depth.
// Focus: when to use each primitive, internal tradeoffs, gotchas.

// ────────────────────────────────────────────────────────────
// Exercise 1: SafeCounter — protect a map with sync.Mutex
// ────────────────────────────────────────────────────────────
// Implement a goroutine-safe counter backed by map[string]int.
// Increment(key) adds 1, Value(key) returns current count.

type SafeCounter struct {
	mu sync.Mutex
	m  map[string]int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{m: make(map[string]int)}
}

func (c *SafeCounter) Increment(key string) {
	// TODO: Lock, increment m[key], Unlock
}

func (c *SafeCounter) Value(key string) int {
	// TODO: Lock, read m[key], Unlock, return
	return 0
}

// ────────────────────────────────────────────────────────────
// Exercise 2: RWCache — sync.RWMutex for read-heavy workloads
// ────────────────────────────────────────────────────────────
// Multiple goroutines read concurrently; writes are exclusive.
// Get returns (value, ok). Set stores a key-value pair.

type RWCache struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewRWCache() *RWCache {
	return &RWCache{m: make(map[string]string)}
}

func (c *RWCache) Get(key string) (string, bool) {
	// TODO: RLock (shared read lock), read, RUnlock
	return "", false
}

func (c *RWCache) Set(key, value string) {
	// TODO: Lock (exclusive write lock), write, Unlock
}

// ────────────────────────────────────────────────────────────
// Exercise 3: WaitForAll — sync.WaitGroup to fan-out work
// ────────────────────────────────────────────────────────────
// Launch n goroutines, each calling fn(i). Wait for all to finish.
// Return the total count of completed goroutines.
// Use sync.WaitGroup + atomic counter.

func WaitForAll(n int, fn func(i int)) int {
	// TODO: wg.Add(n), launch n goroutines, wg.Wait(), return n
	return 0
}

// ────────────────────────────────────────────────────────────
// Exercise 4: InitOnce — sync.Once for one-time initialization
// ────────────────────────────────────────────────────────────
// Create a Service that initializes its connection string only once,
// even if Init() is called from multiple goroutines concurrently.
// Return the connection string from ConnString().

type Service struct {
	once       sync.Once
	connString string
}

func (s *Service) Init(dsn string) {
	// TODO: s.once.Do(func() { s.connString = dsn })
}

func (s *Service) ConnString() string {
	return s.connString
}

// ────────────────────────────────────────────────────────────
// Exercise 5: BufferPool — sync.Pool for buffer reuse
// ────────────────────────────────────────────────────────────
// Create a pool of *bytes.Buffer. GetBuffer retrieves one (resetting it),
// PutBuffer returns it. This reduces GC pressure in hot paths.

var bufferPool = sync.Pool{
	New: func() any {
		// TODO: return new(bytes.Buffer)
		return nil
	},
}

func GetBuffer() io.Writer {
	// TODO: pool.Get().(*bytes.Buffer), reset it, return
	return nil
}

func PutBuffer(w io.Writer) {
	// TODO: type-assert to *bytes.Buffer, Reset(), pool.Put()
}

// ────────────────────────────────────────────────────────────
// Exercise 6: ConcurrentSum — mutex vs channel for aggregation
// ────────────────────────────────────────────────────────────
// Sum elements of a slice using n goroutines, each processing a chunk.
// Use sync.Mutex to aggregate partial sums.

func ConcurrentSum(nums []int, workers int) int {
	// TODO: split nums into chunks, each goroutine sums its chunk,
	// use mutex to add to total, WaitGroup to synchronize
	return 0
}

// ────────────────────────────────────────────────────────────
// Exercise 7: SyncMapStore — sync.Map for concurrent key-value
// ────────────────────────────────────────────────────────────
// Store n key-value pairs concurrently, then count total entries.
// sync.Map is optimized for write-once-read-many or disjoint keys.

func SyncMapStore(pairs map[string]int) int {
	// TODO: range pairs, store each in sync.Map from goroutines,
	// Range() to count entries after all goroutines complete
	return 0
}

// ────────────────────────────────────────────────────────────
// Exercise 8: OnceValue — compute an expensive value exactly once
// ────────────────────────────────────────────────────────────
// Use sync.OnceValue (Go 1.21+) to create a lazy singleton.
// The factory function should be called exactly once.
// Return the factory function wrapper.

func MakeLazy(factory func() string) func() string {
	// TODO: return sync.OnceValue(factory)
	return func() string { return "" }
}

// ────────────────────────────────────────────────────────────
// Exercise 9: MutexTimeout — try-lock pattern with channel
// ────────────────────────────────────────────────────────────
// Go's sync.Mutex has no TryLock timeout. Implement a "lock with
// timeout" using a channel-based semaphore (buffered chan of size 1).
// TryAcquire returns true if lock acquired within the timeout.

type TimedMutex struct {
	ch chan struct{}
}

func NewTimedMutex() *TimedMutex {
	ch := make(chan struct{}, 1)
	ch <- struct{}{} // initially unlocked
	return &TimedMutex{ch: ch}
}

func (m *TimedMutex) Lock() {
	// TODO: <-m.ch (blocks until available)
}

func (m *TimedMutex) Unlock() {
	// TODO: m.ch <- struct{}{}
}

func (m *TimedMutex) TryLock(timeout <-chan struct{}) bool {
	// TODO: select { case <-m.ch: return true; case <-timeout: return false }
	return false
}

// ────────────────────────────────────────────────────────────
// Exercise 10: CondBroadcast — sync.Cond for signaling waiters
// ────────────────────────────────────────────────────────────
// Implement a gate that blocks N goroutines until Open() is called.
// Open() broadcasts to all waiting goroutines.

type Gate struct {
	cond   *sync.Cond
	opened bool
}

func NewGate() *Gate {
	return &Gate{cond: sync.NewCond(&sync.Mutex{})}
}

func (g *Gate) Wait() {
	// TODO: g.cond.L.Lock(); for !g.opened { g.cond.Wait() }; g.cond.L.Unlock()
}

func (g *Gate) Open() {
	// TODO: g.cond.L.Lock(); g.opened = true; g.cond.Broadcast(); g.cond.L.Unlock()
}

// ────────────────────────────────────────────────────────────
// Exercise 11: AtomicConfig — atomic.Value for lock-free reads
// ────────────────────────────────────────────────────────────
// Store and load a Config atomically. Writers call Update(),
// readers call Current(). No mutex needed for reads.

type Config struct {
	Host string
	Port int
}

type AtomicConfig struct {
	v sync.Map // using sync.Map as container; real impl uses atomic.Value
}

func (ac *AtomicConfig) Update(cfg Config) {
	// TODO: ac.v.Store("cfg", cfg)
}

func (ac *AtomicConfig) Current() Config {
	// TODO: v, ok := ac.v.Load("cfg"); if ok { return v.(Config) }; return Config{}
	return Config{}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: Singleflight — deduplicate concurrent calls
// ────────────────────────────────────────────────────────────
// Implement a simplified singleflight: if multiple goroutines call
// Do(key, fn) concurrently with the same key, only ONE executes fn.
// All others wait and receive the same result.
// Use sync.Mutex + sync.WaitGroup + a map of in-flight calls.

type call struct {
	wg  sync.WaitGroup
	val string
	err error
}

type SingleFlight struct {
	mu sync.Mutex
	m  map[string]*call
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight{m: make(map[string]*call)}
}

func (sf *SingleFlight) Do(key string, fn func() (string, error)) (string, error) {
	// TODO:
	// 1. Lock, check if key already in sf.m
	// 2. If yes: unlock, wg.Wait(), return existing result
	// 3. If no: create call, add to map, wg.Add(1), unlock
	// 4. Execute fn, store result, wg.Done()
	// 5. Lock, delete from map, unlock
	// 6. Return result
	return "", nil
}

// Ensure sync import is used
var _ sync.Mutex
