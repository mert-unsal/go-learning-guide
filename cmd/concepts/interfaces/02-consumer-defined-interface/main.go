// Package main demonstrates the consumer-defined interface pattern.
//
// ============================================================
// 2. THE CONSUMER DEFINES THE INTERFACE
// ============================================================
// This is the key Go pattern. You do NOT go to the producer's package
// and make their type implement your interface. You define the interface
// in YOUR package describing only what YOU need.
//
// Imagine you're writing a function that saves data somewhere.
// You don't care if it's a file, a database, or an in-memory buffer.
// You only care that it can Write bytes.
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

// Saver is defined HERE, by the consumer. It describes the minimum
// behavior this package needs. It does NOT live next to File or DB.
type Saver interface {
	Save(data string) error
}

// FileSaver and DBSaver are defined independently — they know nothing about Saver.
type FileSaver struct{ Path string }

func (f FileSaver) Save(data string) error {
	// (pretend we write to a file)
	fmt.Printf("[FileSaver] writing %d bytes to %s\n", len(data), f.Path)
	return nil
}

type DBSaver struct{ Table string }

func (d DBSaver) Save(data string) error {
	// (pretend we write to a DB)
	fmt.Printf("[DBSaver] inserting into table %s: %q\n", d.Table, data)
	return nil
}

// persist only knows about Saver — it is fully decoupled from FileSaver and DBSaver.
// You can add a new storage backend without touching this function at all.
func persist(s Saver, data string) {
	if err := s.Save(data); err != nil {
		fmt.Println("save failed:", err)
	}
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Consumer-Defined Interface Pattern      %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ The interface is defined by the CONSUMER, not the producer%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Saver interface lives here — next to persist(), which needs it%s\n", green, reset)
	fmt.Printf("  %s✔ FileSaver and DBSaver know nothing about Saver%s\n", green, reset)
	fmt.Printf("  %s✔ New backends can be added without touching persist()%s\n\n", green, reset)

	fmt.Printf("%s▸ FileSaver satisfies Saver via Save(data string) error%s\n", cyan+bold, reset)
	fmt.Printf("  ")
	persist(FileSaver{Path: "/data/events.log"}, "hello world")

	fmt.Printf("\n%s▸ DBSaver satisfies Saver via Save(data string) error%s\n", cyan+bold, reset)
	fmt.Printf("  ")
	persist(DBSaver{Table: "events"}, "hello world")

	fmt.Printf("\n  %s⚠ In Java, FileSaver would 'implements Saver' — coupling at definition%s\n", yellow, reset)
	fmt.Printf("  %s⚠ In Go, the consumer defines what it needs — decoupling by default%s\n", yellow, reset)
	fmt.Printf("  %s⚠ This is why Go code is so easy to test: swap in a mock Saver trivially%s\n", yellow, reset)
}
