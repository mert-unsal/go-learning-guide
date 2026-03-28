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
	c := NewMemCache()
	_ = c.Store("lang", "Go")
	lookupOrStore(c, "lang")    // hit: Go
	lookupOrStore(c, "missing") // miss: key not found
}
