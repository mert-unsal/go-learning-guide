// Package main demonstrates the power of small, focused interfaces.
//
// ============================================================
// 3. SMALL INTERFACES — THE io.Reader / io.Writer LESSON
// ============================================================
// Go's standard library is built on tiny, single-method interfaces.
// The smaller the interface, the more types can satisfy it,
// and the more places it can be used.
//
//   io.Reader  → Read(p []byte) (n int, err error)
//   io.Writer  → Write(p []byte) (n int, err error)
//   io.Closer  → Close() error
//   fmt.Stringer → String() string
//   error      → Error() string
//
// Because io.Reader is ONE method, strings, files, network connections,
// HTTP bodies, gzip readers, test buffers — ALL satisfy it.
// A 10-method interface would exclude most of them.
//
// Rule of thumb: if your interface has more than 3 methods, ask yourself
// if it's really one concern or several concerns bundled together.
package main

import (
	"fmt"
	"math"
)

// Measurer is a small, focused interface.
type Measurer interface {
	Measure() float64
}

// Notice: Circle and Rectangle are defined here NOT because they ARE shapes
// (that's OOP thinking), but because we need concrete types for the example.
// The insight is that ANYTHING with a Measure() method works — a distance,
// a weight, a duration, a file size.
type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }
type Segment struct{ Length float64 } // not a "shape" at all
type FileSize struct{ Bytes int64 }   // completely different domain

func (c Circle) Measure() float64    { return math.Pi * c.Radius * c.Radius }
func (r Rectangle) Measure() float64 { return r.Width * r.Height }
func (s Segment) Measure() float64   { return s.Length }
func (f FileSize) Measure() float64  { return float64(f.Bytes) }

func totalMeasure(items []Measurer) float64 {
	total := 0.0
	for _, m := range items {
		total += m.Measure()
	}
	return total
}

func main() {
	// Completely unrelated types used together — because behavior, not type.
	items := []Measurer{
		Circle{Radius: 2},
		Rectangle{Width: 3, Height: 4},
		Segment{Length: 10},
		FileSize{Bytes: 1024},
	}
	fmt.Printf("Total: %.2f\n", totalMeasure(items))
}
