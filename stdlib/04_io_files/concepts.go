// Package io_files covers reading and writing files in Go using
// the os, bufio, and io packages — essential for real-world programs.
//
// Topics:
//   - os.ReadFile / os.WriteFile  (simple, whole-file operations)
//   - os.Open / os.Create         (fine-grained file control)
//   - bufio.Scanner               (line-by-line reading)
//   - bufio.Writer                (buffered writing)
//   - io.Reader / io.Writer       (interfaces behind everything)
//   - Working with paths          (filepath package)
package io_files

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ============================================================
// 1. SIMPLE READ/WRITE — os.ReadFile & os.WriteFile
// ============================================================
// These are the easiest APIs — load or save the ENTIRE file into memory.
// Perfect for configuration files, small data files, JSON payloads.

func DemonstrateSimpleReadWrite() {
	// --- Write a file ---
	content := []byte("Hello, Go!\nLine two.\nLine three.\n")
	err := os.WriteFile("example.txt", content, 0644)
	// 0644 = owner can read+write, group/others can read
	if err != nil {
		fmt.Println("write error:", err)
		return
	}
	fmt.Println("Written example.txt")

	// --- Read entire file into memory ---
	data, err := os.ReadFile("example.txt")
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	fmt.Println("File contents:\n" + string(data))

	// Cleanup
	os.Remove("example.txt")
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

func DemonstrateOpenCreate() {
	// --- Create and write ---
	f, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("create error:", err)
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
		fmt.Println("open error:", err)
		return
	}
	defer rf.Close()

	buf := make([]byte, 1024)
	n, _ := rf.Read(buf) // reads up to 1024 bytes
	fmt.Println("Read bytes:", n)
	fmt.Print("Content: ", string(buf[:n]))

	rf.Close()

	// --- Append to existing file ---
	af, _ := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, 0644)
	defer af.Close()
	af.WriteString("Appended line\n")
	af.Close()

	// Cleanup
	os.Remove("output.txt")
}

// ============================================================
// 3. bufio.Scanner — Line-by-Line Reading (MOST COMMON PATTERN)
// ============================================================
// bufio.Scanner reads a file line by line without loading the whole file.
// Essential for large files. Also works with stdin (os.Stdin).
//
// Default split function: ScanLines (splits on \n)
// Other options: bufio.ScanWords, bufio.ScanRunes, bufio.ScanBytes

func DemonstrateScanner() {
	// Create a test file with multiple lines
	os.WriteFile("data.txt", []byte("alice 30\nbob 25\ncharlie 35\n"), 0644)

	f, err := os.Open("data.txt")
	if err != nil {
		fmt.Println("open error:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// --- Read line by line ---
	lineNum := 1
	for scanner.Scan() { // Scan() returns false at EOF or error
		line := scanner.Text() // current line (without newline)
		fmt.Printf("Line %d: %q\n", lineNum, line)
		lineNum++
	}

	// Always check for scanner errors after the loop
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
	}

	// --- Scan words instead of lines ---
	f.Seek(0, io.SeekStart) // rewind to beginning
	wordScanner := bufio.NewScanner(f)
	wordScanner.Split(bufio.ScanWords) // split on whitespace
	var words []string
	for wordScanner.Scan() {
		words = append(words, wordScanner.Text())
	}
	fmt.Println("Words:", words)

	os.Remove("data.txt")
}

// ============================================================
// 4. bufio.Writer — Buffered Writing (FASTER for many small writes)
// ============================================================
// Buffered writing accumulates writes in memory and flushes them
// in larger batches — much faster than writing one line at a time.
//
// CRITICAL: always call Flush() before closing, or data may be lost!

func DemonstrateBufferedWrite() {
	f, _ := os.Create("buffered.txt")
	defer f.Close()

	w := bufio.NewWriter(f) // wrap file in buffered writer
	// or: bufio.NewWriterSize(f, 65536) for custom buffer size

	for i := 1; i <= 5; i++ {
		fmt.Fprintf(w, "Line %d: hello\n", i)
		// Data is in the buffer, NOT written to disk yet
	}

	err := w.Flush() // FLUSH: write buffered data to the underlying file
	if err != nil {
		fmt.Println("flush error:", err)
	}
	fmt.Printf("Bytes in buffer before flush: %d\n", w.Buffered())

	os.Remove("buffered.txt")
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

func DemonstrateIOInterfaces() {
	// io.Copy — copy from any Reader to any Writer
	src := strings.NewReader("Hello from strings.Reader!")
	var dst strings.Builder // strings.Builder implements io.Writer
	n, err := io.Copy(&dst, src)
	fmt.Printf("Copied %d bytes: %q, err=%v\n", n, dst.String(), err)

	// io.ReadAll — read everything from a Reader
	r := strings.NewReader("Read all of this")
	data, _ := io.ReadAll(r)
	fmt.Println("ReadAll:", string(data))

	// Write to multiple destinations simultaneously
	f, _ := os.Create("tee.txt")
	defer f.Close()
	mw := io.MultiWriter(f, os.Stdout) // write to file AND stdout
	fmt.Fprintln(mw, "Written to both!")

	os.Remove("tee.txt")
}

// ============================================================
// 6. filepath — Working with File Paths (Cross-Platform)
// ============================================================
// Always use filepath (not path) for OS file paths.
// path is for URL-style forward-slash paths.

func DemonstrateFilepath() {
	// Join path segments (uses \ on Windows, / on Unix)
	p := filepath.Join("data", "users", "alice.json")
	fmt.Println("Joined:", p) // data\users\alice.json (on Windows)

	// Split into directory and file name
	dir, file := filepath.Split(p)
	fmt.Println("Dir:", dir, "File:", file)

	// Get file extension
	ext := filepath.Ext("report.pdf")
	fmt.Println("Ext:", ext) // .pdf

	// Base name (filename without directory)
	base := filepath.Base("/home/user/documents/report.pdf")
	fmt.Println("Base:", base) // report.pdf

	// Absolute path
	abs, _ := filepath.Abs(".")
	fmt.Println("Abs cwd:", abs)

	// Walk a directory tree
	// filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
	//     fmt.Println(path)
	//     return nil
	// })
}

// ============================================================
// 7. os.Stat — Check if File/Directory Exists
// ============================================================

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func DemonstrateFileStat() {
	// Create a temp file
	os.WriteFile("check.txt", []byte("hi"), 0644)

	fmt.Println("check.txt exists:", FileExists("check.txt"))  // true
	fmt.Println("nope.txt exists:", FileExists("nope.txt"))    // false
	fmt.Println("check.txt is dir:", IsDirectory("check.txt")) // false

	info, _ := os.Stat("check.txt")
	fmt.Printf("Size: %d bytes, ModTime: %s\n", info.Size(), info.ModTime().Format("2006-01-02"))

	os.Remove("check.txt")
}

// ============================================================
// 8. os.Stdin / os.Stdout / os.Stderr
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

func FastIOTemplate() {
	// This pattern is used in competitive programming for performance:
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush() // flush at the end

	// Read a line (use in a loop for multiple test cases)
	_ = reader // suppress unused warning in non-interactive demo
	fmt.Fprintln(writer, "Fast output via bufio.Writer")
}
