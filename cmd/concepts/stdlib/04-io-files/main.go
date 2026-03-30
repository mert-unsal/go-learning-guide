// I/O & Files in Go — demonstrates reading and writing files using
// the os, bufio, io, and filepath packages.
//
// Topics:
//   - os.ReadFile / os.WriteFile  (simple, whole-file operations)
//   - os.Open / os.Create         (fine-grained file control)
//   - bufio.Scanner               (line-by-line reading)
//   - bufio.Writer                (buffered writing)
//   - io.Reader / io.Writer       (interfaces behind everything)
//   - Working with paths          (filepath package)
//
// Run: go run cmd/concepts/stdlib/04-io-files/main.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  I/O & Files                             %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	demonstrateSimpleReadWrite()
	demonstrateOpenCreate()
	demonstrateScanner()
	demonstrateBufferedWrite()
	demonstrateIOInterfaces()
	demonstrateFilepath()
	demonstrateFileStat()
	demonstrateFastIOConcepts()
}

// ============================================================
// 1. SIMPLE READ/WRITE — os.ReadFile & os.WriteFile
// ============================================================
// These are the easiest APIs — load or save the ENTIRE file into memory.
// Perfect for configuration files, small data files, JSON payloads.

func demonstrateSimpleReadWrite() {
	fmt.Printf("%s▸ 1. os.ReadFile & os.WriteFile — Simple Read/Write%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Easiest API — loads/saves ENTIRE file into memory%s\n", green, reset)

	// --- Write a file ---
	content := []byte("Hello, Go!\nLine two.\nLine three.\n")
	err := os.WriteFile("example.txt", content, 0644)
	// 0644 = owner can read+write, group/others can read
	if err != nil {
		fmt.Println("  write error:", err)
		return
	}
	fmt.Println("  Written example.txt")

	// --- Read entire file into memory ---
	data, err := os.ReadFile("example.txt")
	if err != nil {
		fmt.Println("  read error:", err)
		return
	}
	fmt.Println("  File contents:\n  " + strings.ReplaceAll(strings.TrimRight(string(data), "\n"), "\n", "\n  "))

	// Cleanup
	os.Remove("example.txt")
	fmt.Println()
}

// ============================================================
// 2. os.Open / os.Create — Fine-Grained File Control
// ============================================================
// os.Open(path)   — open for reading only  (returns *os.File)
// os.Create(path) — create/truncate for writing (returns *os.File)
// os.OpenFile(path, flags, perm) — full control
//
// ALWAYS defer file.Close() to avoid resource leaks.
//
// Common flags:
//   os.O_RDONLY  — read only
//   os.O_WRONLY  — write only
//   os.O_RDWR   — read and write
//   os.O_CREATE  — create if not exists
//   os.O_APPEND  — append to existing file
//   os.O_TRUNC   — truncate to zero when opening

func demonstrateOpenCreate() {
	fmt.Printf("%s▸ 2. os.Open / os.Create — Fine-Grained Control%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ ALWAYS defer file.Close() to avoid resource leaks%s\n", green, reset)

	// --- Create and write ---
	f, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("  create error:", err)
		return
	}
	defer f.Close() // ALWAYS defer Close!

	// Write methods on *os.File
	f.WriteString("First line\n")
	fmt.Fprintln(f, "Second line")    // fmt.Fprintln writes to any io.Writer
	fmt.Fprintf(f, "Value: %d\n", 42) // formatted write

	f.Close() // close before re-opening (or defer handles it)

	// --- Open and read with os.Open ---
	rf, err := os.Open("output.txt")
	if err != nil {
		fmt.Println("  open error:", err)
		return
	}
	defer rf.Close()

	buf := make([]byte, 1024)
	n, _ := rf.Read(buf) // reads up to 1024 bytes
	fmt.Println("  Read bytes:", n)
	fmt.Print("  Content: ", string(buf[:n]))

	rf.Close()

	// --- Append to existing file ---
	af, _ := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, 0644)
	defer af.Close()
	af.WriteString("Appended line\n")
	af.Close()

	// Cleanup
	os.Remove("output.txt")
	fmt.Println()
}

// ============================================================
// 3. bufio.Scanner — Line-by-Line Reading (MOST COMMON PATTERN)
// ============================================================
// bufio.Scanner reads a file line by line without loading the whole file.
// Essential for large files. Also works with stdin (os.Stdin).
//
// Default split function: ScanLines (splits on \n)
// Other options: bufio.ScanWords, bufio.ScanRunes, bufio.ScanBytes

func demonstrateScanner() {
	fmt.Printf("%s▸ 3. bufio.Scanner — Line-by-Line Reading%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Reads line by line without loading the whole file%s\n", green, reset)

	// Create a test file with multiple lines
	os.WriteFile("data.txt", []byte("alice 30\nbob 25\ncharlie 35\n"), 0644)

	f, err := os.Open("data.txt")
	if err != nil {
		fmt.Println("  open error:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// --- Read line by line ---
	lineNum := 1
	for scanner.Scan() { // Scan() returns false at EOF or error
		line := scanner.Text() // current line (without newline)
		fmt.Printf("  Line %d: %q\n", lineNum, line)
		lineNum++
	}

	// Always check for scanner errors after the loop
	if err := scanner.Err(); err != nil {
		fmt.Println("  scanner error:", err)
	}

	// --- Scan words instead of lines ---
	f.Seek(0, io.SeekStart) // rewind to beginning
	wordScanner := bufio.NewScanner(f)
	wordScanner.Split(bufio.ScanWords) // split on whitespace
	var words []string
	for wordScanner.Scan() {
		words = append(words, wordScanner.Text())
	}
	fmt.Println("  Words:", words)

	os.Remove("data.txt")
	fmt.Println()
}

// ============================================================
// 4. bufio.Writer — Buffered Writing (FASTER for many small writes)
// ============================================================
// Buffered writing accumulates writes in memory and flushes them
// in larger batches — much faster than writing one line at a time.
//
// CRITICAL: always call Flush() before closing, or data may be lost!

func demonstrateBufferedWrite() {
	fmt.Printf("%s▸ 4. bufio.Writer — Buffered Writing%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ CRITICAL: always call Flush() before closing, or data may be lost!%s\n", yellow, reset)

	f, _ := os.Create("buffered.txt")
	defer f.Close()

	w := bufio.NewWriter(f) // wrap file in buffered writer
	// or: bufio.NewWriterSize(f, 65536) for custom buffer size

	for i := 1; i <= 5; i++ {
		fmt.Fprintf(w, "Line %d: hello\n", i)
		// Data is in the buffer, NOT written to disk yet
	}

	fmt.Printf("  Bytes in buffer before flush: %d\n", w.Buffered())
	err := w.Flush() // FLUSH: write buffered data to the underlying file
	if err != nil {
		fmt.Println("  flush error:", err)
	}
	fmt.Println("  Flushed successfully")

	os.Remove("buffered.txt")
	fmt.Println()
}

// ============================================================
// 5. io.Reader and io.Writer INTERFACES
// ============================================================
// Almost everything in Go that reads or writes implements these:
//
//   type Reader interface { Read(p []byte) (n int, err error) }
//   type Writer interface { Write(p []byte) (n int, err error) }
//
// This means you can pass a *os.File, bytes.Buffer, strings.Reader,
// http.Response.Body, etc. to any function accepting io.Reader.
//
// Key io helpers:
//   io.Copy(dst, src)          — copy all bytes from src to dst
//   io.ReadAll(r)              — read everything from r into []byte
//   strings.NewReader(s)       — create a Reader from a string
//   bytes.Buffer               — in-memory read/write buffer

func demonstrateIOInterfaces() {
	fmt.Printf("%s▸ 5. io.Reader & io.Writer Interfaces%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Everything in Go that reads/writes implements these interfaces%s\n", green, reset)

	// io.Copy — copy from any Reader to any Writer
	src := strings.NewReader("Hello from strings.Reader!")
	var dst strings.Builder // strings.Builder implements io.Writer
	n, err := io.Copy(&dst, src)
	fmt.Printf("  Copied %d bytes: %q, err=%v\n", n, dst.String(), err)

	// io.ReadAll — read everything from a Reader
	r := strings.NewReader("Read all of this")
	data, _ := io.ReadAll(r)
	fmt.Println("  ReadAll:", string(data))

	// Write to multiple destinations simultaneously
	f, _ := os.Create("tee.txt")
	defer f.Close()
	mw := io.MultiWriter(f, os.Stdout) // write to file AND stdout
	fmt.Fprint(mw, "  Written to both!\n")

	f.Close()
	os.Remove("tee.txt")
	fmt.Println()
}

// ============================================================
// 6. filepath — Working with File Paths (Cross-Platform)
// ============================================================
// Always use filepath (not path) for OS file paths.
// path is for URL-style forward-slash paths.

func demonstrateFilepath() {
	fmt.Printf("%s▸ 6. filepath — Cross-Platform Paths%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Always use filepath (not path) for OS file paths%s\n", green, reset)

	// Join path segments (uses \ on Windows, / on Unix)
	p := filepath.Join("data", "users", "alice.json")
	fmt.Println("  Joined:", p)

	// Split into directory and file name
	dir, file := filepath.Split(p)
	fmt.Println("  Dir:", dir, "File:", file)

	// Get file extension
	ext := filepath.Ext("report.pdf")
	fmt.Println("  Ext:", ext) // .pdf

	// Base name (filename without directory)
	base := filepath.Base("/home/user/documents/report.pdf")
	fmt.Println("  Base:", base) // report.pdf

	// Absolute path
	abs, _ := filepath.Abs(".")
	fmt.Println("  Abs cwd:", abs)

	// Walk a directory tree
	// filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
	//     fmt.Println(path)
	//     return nil
	// })
	fmt.Println()
}

// ============================================================
// 7. os.Stat — Check if File/Directory Exists
// ============================================================

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func demonstrateFileStat() {
	fmt.Printf("%s▸ 7. os.Stat — File/Directory Existence%s\n", cyan+bold, reset)

	// Create a temp file
	os.WriteFile("check.txt", []byte("hi"), 0644)

	fmt.Println("  check.txt exists:", fileExists("check.txt"))  // true
	fmt.Println("  nope.txt exists:", fileExists("nope.txt"))    // false
	fmt.Println("  check.txt is dir:", isDirectory("check.txt")) // false

	info, _ := os.Stat("check.txt")
	fmt.Printf("  Size: %d bytes, ModTime: %s\n", info.Size(), info.ModTime().Format("2006-01-02"))

	os.Remove("check.txt")
	fmt.Println()
}

// ============================================================
// 8. os.Stdin / os.Stdout / os.Stderr — Fast I/O Concepts
// ============================================================
// These are *os.File values that implement io.Reader/io.Writer.
// Use them for command-line tools and competitive programming.
//
// Fast stdin reading (for competitive programming — avoids slow fmt.Scan):
//
//   reader := bufio.NewReader(os.Stdin)
//   writer := bufio.NewWriter(os.Stdout)
//   defer writer.Flush()
//
//   line, _ := reader.ReadString('\n')
//   fmt.Fprintln(writer, result)

func demonstrateFastIOConcepts() {
	fmt.Printf("%s▸ 8. Fast I/O Pattern (Competitive Programming)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ os.Stdin, os.Stdout, os.Stderr are *os.File implementing io.Reader/Writer%s\n", green, reset)
	fmt.Printf("  %s✔ Wrap in bufio for fast competitive programming I/O%s\n", green, reset)

	// This pattern is used in competitive programming for performance:
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush() // flush at the end

	fmt.Fprintln(writer, "  Fast output via bufio.Writer")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "  Template:")
	fmt.Fprintln(writer, "    reader := bufio.NewReader(os.Stdin)")
	fmt.Fprintln(writer, "    writer := bufio.NewWriter(os.Stdout)")
	fmt.Fprintln(writer, "    defer writer.Flush()")
}
