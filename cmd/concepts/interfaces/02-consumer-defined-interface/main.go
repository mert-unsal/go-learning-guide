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
	persist(FileSaver{Path: "/tmp/data.txt"}, "hello world")
	persist(DBSaver{Table: "events"}, "hello world")
}
