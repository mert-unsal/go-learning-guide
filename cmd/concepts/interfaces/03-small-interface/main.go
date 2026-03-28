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
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Small Interfaces — Power of Minimalism  %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Measurer has ONE method: Measure() float64%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Small interface = more types can satisfy it%s\n", green, reset)
	fmt.Printf("  %s✔ io.Reader has 1 method → files, strings, HTTP bodies, gzip all satisfy it%s\n", green, reset)
	fmt.Printf("  %s✔ A 10-method interface would exclude most types%s\n\n", green, reset)

	// Completely unrelated types used together — because behavior, not type.
	c := Circle{Radius: 2}
	r := Rectangle{Width: 3, Height: 4}
	seg := Segment{Length: 10}
	fs := FileSize{Bytes: 1024}

	items := []Measurer{c, r, seg, fs}

	fmt.Printf("%s▸ Four completely unrelated types — all satisfy Measurer%s\n", cyan+bold, reset)
	fmt.Printf("  Circle{Radius: 2}        → Measure() = %s%.2f%s  (π·r²)\n", magenta, c.Measure(), reset)
	fmt.Printf("  Rectangle{3, 4}          → Measure() = %s%.2f%s  (w·h)\n", magenta, r.Measure(), reset)
	fmt.Printf("  Segment{Length: 10}      → Measure() = %s%.2f%s  (not a shape!)\n", magenta, seg.Measure(), reset)
	fmt.Printf("  FileSize{Bytes: 1024}    → Measure() = %s%.2f%s  (different domain entirely)\n\n", magenta, fs.Measure(), reset)

	total := totalMeasure(items)
	fmt.Printf("%s▸ totalMeasure([]Measurer) works on ALL of them%s\n", cyan+bold, reset)
	fmt.Printf("  Total = %s%.2f%s\n\n", magenta, total, reset)

	fmt.Printf("  %s⚠ Go groups by WHAT THINGS DO, not what they ARE%s\n", yellow, reset)
	fmt.Printf("  %s⚠ OOP would force Circle and FileSize into separate hierarchies%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Rule: >3 methods? Ask if it's one concern or several bundled together%s\n", yellow, reset)
}
