package io_files

import "io"

// ============================================================
// EXERCISES — 04 io & Files
// ============================================================
//
// Go's io package defines two tiny interfaces that power everything:
//
//   type Reader interface { Read(p []byte) (n int, err error) }
//   type Writer interface { Write(p []byte) (n int, err error) }
//
// Every stream in Go — files, HTTP bodies, network connections,
// compression, encryption — implements these interfaces. The power
// comes from COMPOSITION: wrapping readers/writers in decorators.
//
// Exercises 1-4:  Basics — file I/O, line scanning, io.ReadAll
// Exercises 5-8:  Composition — LimitReader, MultiReader, TeeReader, Pipe
// Exercises 9-12: Production patterns — buffered copy, filtering, counting
// ============================================================

// Exercise 1:
// WriteAndReadFile writes content to a file, reads it back, returns the content.
// Use os.WriteFile and os.ReadFile. Delete the file when done (defer os.Remove).
//
// LESSON: os.WriteFile/ReadFile are convenience wrappers for simple cases.
// For large files or streaming, use os.Open + io.Copy instead.
func WriteAndReadFile(filename, content string) (string, error) {
	// TODO: os.WriteFile(filename, []byte(content), 0644)
	//       defer os.Remove(filename)
	//       data, err := os.ReadFile(filename)
	//       return string(data), err
	return "", nil
}

// Exercise 2:
// CountLines counts the number of lines in a string using bufio.Scanner.
//
// LESSON: strings.NewReader wraps a string as an io.Reader — this is
// the adapter pattern. bufio.Scanner reads line-by-line by default.
// An empty string has 0 lines. A string without \n has 1 line.
func CountLines(content string) int {
	// TODO: strings.NewReader → bufio.NewScanner → scan loop → count
	return 0
}

// Exercise 3:
// WordFrequency reads text line-by-line from r (an io.Reader)
// and returns a map of word→frequency. Case-insensitive.
//
// LESSON: Accept io.Reader (not *os.File, not string) — this makes
// your function work with files, HTTP bodies, strings, pipes, anything.
// "Accept interfaces, return concrete types."
func WordFrequency(r io.Reader) map[string]int {
	// TODO: bufio.NewScanner, scanner.Scan loop, strings.Fields, strings.ToLower
	return nil
}

// Exercise 4:
// CopyReader reads all bytes from src and returns them as a string.
//
// LESSON: io.ReadAll reads until EOF. Use it for small payloads.
// For large streams, use io.Copy to avoid loading everything into memory.
func CopyReader(src io.Reader) (string, error) {
	// TODO: io.ReadAll(src), return string(data), err
	return "", nil
}

// ── Composition Patterns — Exercises 5-8 ──

// Exercise 5:
// ReadLimited reads at most maxBytes from the reader and returns the result.
//
// LESSON: io.LimitReader wraps a reader to stop after N bytes.
// This prevents unbounded reads from untrusted sources (HTTP bodies,
// user uploads). Production pattern: limit request body size.
//
// Example: ReadLimited(strings.NewReader("Hello, World!"), 5) → "Hello"
func ReadLimited(r io.Reader, maxBytes int64) (string, error) {
	// TODO: limited := io.LimitReader(r, maxBytes)
	//       data, err := io.ReadAll(limited)
	//       return string(data), err
	return "", nil
}

// Exercise 6:
// ConcatReaders reads from multiple readers in sequence and returns
// the combined content as a single string.
//
// LESSON: io.MultiReader chains readers end-to-end. When the first
// reader hits EOF, it moves to the next. This is the decorator pattern:
// each reader wraps the io.Reader interface.
//
// Example: ConcatReaders(NewReader("Hello"), NewReader(" World")) → "Hello World"
func ConcatReaders(readers ...io.Reader) (string, error) {
	// TODO: combined := io.MultiReader(readers...)
	//       data, err := io.ReadAll(combined)
	return "", nil
}

// Exercise 7:
// ReadAndCapture reads all content from r and returns both the content
// and a copy captured in a separate buffer.
//
// LESSON: io.TeeReader creates a reader that writes to w everything it reads.
// Like the Unix `tee` command: data flows through to the reader AND is
// captured in the writer simultaneously. Use for logging/auditing streams.
//
// Returns (content read, captured copy, error)
func ReadAndCapture(r io.Reader) (content string, captured string, err error) {
	// TODO: var buf bytes.Buffer
	//       tee := io.TeeReader(r, &buf)
	//       data, err := io.ReadAll(tee)
	//       return string(data), buf.String(), err
	return "", "", nil
}

// Exercise 8:
// PipeTransfer writes the given message through an io.Pipe and returns
// what was read from the other end.
//
// LESSON: io.Pipe creates a synchronous, in-memory pipe connecting a
// Writer to a Reader. Writes block until the data is read. This is how
// you connect a producer goroutine to a consumer without buffering.
//
// You MUST write in a separate goroutine (or the write blocks forever).
func PipeTransfer(message string) (string, error) {
	// TODO: pr, pw := io.Pipe()
	//       go func() { pw.Write([]byte(message)); pw.Close() }()
	//       data, err := io.ReadAll(pr)
	//       return string(data), err
	return "", nil
}

// ── Production Patterns — Exercises 9-12 ──

// Exercise 9:
// CopyWithBuffer copies from src to a bytes.Buffer using io.Copy
// and returns the number of bytes copied and the result.
//
// LESSON: io.Copy is the fundamental streaming primitive in Go.
// It reads from src in 32KB chunks and writes to dst — never loads
// the entire content in memory. Use it for large files, HTTP proxies,
// and any stream-to-stream transfer.
func CopyWithBuffer(src io.Reader) (int64, string, error) {
	// TODO: var dst bytes.Buffer
	//       n, err := io.Copy(&dst, src)
	//       return n, dst.String(), err
	return 0, "", nil
}

// Exercise 10:
// FilterLines reads lines from r, keeps only lines where the predicate
// returns true, and writes them to w (one per line, with \n).
// Returns the number of lines written.
//
// PRODUCTION PATTERN: This is the reader→filter→writer pipeline.
// Accept io.Reader + io.Writer, use bufio.Scanner for line-by-line.
// The same pattern works for log filtering, CSV processing, etc.
func FilterLines(r io.Reader, w io.Writer, keep func(string) bool) (int, error) {
	// TODO: bufio.NewScanner(r), scan lines, if keep(line) → fmt.Fprintln(w, line)
	return 0, nil
}

// Exercise 11:
// CountingWriter wraps an io.Writer and counts total bytes written.
// Implement the WriterCounter type and its Write method.
//
// LESSON: Implementing io.Writer is trivial — just one method.
// Wrapping adds cross-cutting concerns (counting, logging, rate limiting)
// without modifying the original writer. This is the decorator pattern.
//
// After writing through a CountingWriter, call BytesWritten() to get the count.

// WriterCounter wraps an io.Writer and counts bytes written through it.
type WriterCounter struct {
	W     io.Writer
	Count int64
}

// Write implements io.Writer. Write to the underlying writer and track bytes.
func (wc *WriterCounter) Write(p []byte) (int, error) {
	// TODO: n, err := wc.W.Write(p)
	//       wc.Count += int64(n)
	//       return n, err
	return 0, nil
}

// BytesWritten returns the total number of bytes written.
func (wc *WriterCounter) BytesWritten() int64 {
	return wc.Count
}

// Exercise 12:
// SectionRead reads length bytes starting at offset from the reader.
// The reader must implement io.ReaderAt (random access).
//
// LESSON: io.SectionReader provides a windowed view into a ReaderAt.
// Useful for reading specific portions of large files without reading
// the whole thing. *os.File implements io.ReaderAt via pread(2) syscall.
//
// Example: SectionRead(strings.NewReader("Hello, World!"), 7, 5) → "World"
func SectionRead(r io.ReaderAt, offset int64, length int64) (string, error) {
	// TODO: section := io.NewSectionReader(r, offset, length)
	//       data, err := io.ReadAll(section)
	//       return string(data), err
	return "", nil
}
