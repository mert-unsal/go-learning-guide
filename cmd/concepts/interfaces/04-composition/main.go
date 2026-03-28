// Package main demonstrates interface composition — building bigger
// interfaces from smaller ones.
//
// ============================================================
// 4. INTERFACE COMPOSITION — BUILD BIGGER FROM SMALL
// ============================================================
// Instead of one big interface, compose small ones.
// This mirrors how Go's stdlib works (io.ReadWriter = Reader + Writer).
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

type Loader interface {
	Load(key string) (string, error)
}

type Storer interface {
	Store(key, value string) error
}

// Cache composes Loader and Storer.
// A type must satisfy BOTH to be used as a Cache.
type Cache interface {
	Loader
	Storer
}

type MemCache struct {
	data map[string]string
}

func NewMemCache() *MemCache {
	return &MemCache{data: make(map[string]string)}
}

func (m *MemCache) Load(key string) (string, error) {
	v, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found", key)
	}
	return v, nil
}

func (m *MemCache) Store(key, value string) error {
	m.data[key] = value
	return nil
}

// lookupOrStore only needs to load — it accepts the narrower Loader.
// This makes it easier to test and reuse.
func lookupOrStore(l Loader, key string) {
	v, err := l.Load(key)
	if err != nil {
		fmt.Println("miss:", err)
		return
	}
	fmt.Println("hit:", v)
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Interface Composition                   %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Building bigger interfaces from small ones%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Loader: Load(key string) (string, error) — single responsibility%s\n", green, reset)
	fmt.Printf("  %s✔ Storer: Store(key, value string) error   — single responsibility%s\n", green, reset)
	fmt.Printf("  %s✔ Cache = Loader + Storer — composed, not inherited%s\n\n", green, reset)

	fmt.Printf("  %s⚠ This mirrors stdlib: io.ReadWriter = io.Reader + io.Writer%s\n", yellow, reset)
	fmt.Printf("  %s⚠ io.ReadWriteCloser = Reader + Writer + Closer (3 one-method interfaces)%s\n\n", yellow, reset)

	c := NewMemCache()

	fmt.Printf("%s▸ MemCache satisfies Cache (both Loader and Storer)%s\n", cyan+bold, reset)
	_ = c.Store("lang", "Go")
	fmt.Printf("  Store(%s\"lang\"%s, %s\"Go\"%s) → stored\n", magenta, reset, magenta, reset)

	fmt.Printf("\n%s▸ lookupOrStore() accepts Loader — the narrower interface%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ MemCache satisfies Loader (it has Load method) — works!%s\n", green, reset)
	fmt.Printf("  %s✔ Function asks only for what it needs — easier to test and reuse%s\n\n", green, reset)

	fmt.Printf("  lookupOrStore(c, \"lang\")    → ")
	lookupOrStore(c, "lang")
	fmt.Printf("  lookupOrStore(c, \"missing\") → ")
	lookupOrStore(c, "missing")

	fmt.Printf("\n  %s⚠ Accept the narrowest interface possible — don't take Cache if you only Load%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Compose interfaces at the point of need, not speculatively%s\n", yellow, reset)
}
