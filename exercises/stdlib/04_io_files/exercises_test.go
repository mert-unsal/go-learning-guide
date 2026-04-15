package io_files

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: WriteAndReadFile
// ────────────────────────────────────────────────────────────

func TestWriteAndReadFile(t *testing.T) {
	content := "hello, world!\nline two"
	got, err := WriteAndReadFile("testfile_tmp.txt", content)
	if err != nil {
		t.Fatalf("❌ WriteAndReadFile error: %v\n\t\tHint: os.WriteFile(name, []byte(content), 0644), "+
			"defer os.Remove(name), os.ReadFile(name)", err)
	}
	if got != content {
		t.Errorf("❌ WriteAndReadFile = %q, want %q", got, content)
	} else {
		t.Logf("✅ WriteAndReadFile round-trip OK")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: CountLines
// ────────────────────────────────────────────────────────────

func TestCountLines(t *testing.T) {
	tests := []struct {
		content string
		want    int
	}{
		{"line1\nline2\nline3", 3},
		{"single", 1},
		{"", 0},
		{"a\nb\nc\nd", 4},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("CountLines(%q)", tt.content), func(t *testing.T) {
			got := CountLines(tt.content)
			if got != tt.want {
				t.Errorf("❌ CountLines(%q) = %d, want %d\n\t\t"+
					"Hint: strings.NewReader → bufio.NewScanner → for scanner.Scan() { count++ }",
					tt.content, got, tt.want)
			} else {
				t.Logf("✅ CountLines(%q) = %d", tt.content, got)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: WordFrequency
// ────────────────────────────────────────────────────────────

func TestWordFrequency(t *testing.T) {
	r := strings.NewReader("Go is great\ngo is fun\nGO")
	freq := WordFrequency(r)
	if freq == nil {
		t.Fatal("❌ WordFrequency returned nil\n\t\t" +
			"Hint: Accept io.Reader, not string. bufio.NewScanner(r), " +
			"strings.Fields, strings.ToLower for case-insensitive counting")
	}
	if freq["go"] != 3 {
		t.Errorf("❌ freq[\"go\"] = %d, want 3\n\t\t"+
			"Hint: 'Go', 'go', 'GO' should all map to 'go' (strings.ToLower)",
			freq["go"])
	} else {
		t.Logf("✅ freq[\"go\"] = 3")
	}
	if freq["is"] != 2 {
		t.Errorf("❌ freq[\"is\"] = %d, want 2", freq["is"])
	} else {
		t.Logf("✅ freq[\"is\"] = 2")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: CopyReader
// ────────────────────────────────────────────────────────────

func TestCopyReader(t *testing.T) {
	r := strings.NewReader("hello world")
	got, err := CopyReader(r)
	if err != nil {
		t.Fatalf("❌ CopyReader error: %v\n\t\tHint: io.ReadAll(src)", err)
	}
	if got != "hello world" {
		t.Errorf("❌ CopyReader = %q, want %q\n\t\t"+
			"Hint: io.ReadAll reads until EOF, returns []byte",
			got, "hello world")
	} else {
		t.Logf("✅ CopyReader = %q", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: ReadLimited
// ────────────────────────────────────────────────────────────

func TestReadLimited(t *testing.T) {
	tests := []struct {
		input    string
		max      int64
		want     string
	}{
		{"Hello, World!", 5, "Hello"},
		{"short", 100, "short"},
		{"abc", 0, ""},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("limit=%d", tt.max), func(t *testing.T) {
			got, err := ReadLimited(strings.NewReader(tt.input), tt.max)
			if err != nil {
				t.Fatalf("❌ error: %v", err)
			}
			if got != tt.want {
				t.Errorf("❌ ReadLimited(%q, %d) = %q, want %q\n\t\t"+
					"Hint: io.LimitReader(r, n) wraps r to stop after n bytes. "+
					"Use it to cap untrusted input (HTTP bodies, uploads)",
					tt.input, tt.max, got, tt.want)
			} else {
				t.Logf("✅ ReadLimited(%q, %d) = %q", tt.input, tt.max, got)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: ConcatReaders
// ────────────────────────────────────────────────────────────

func TestConcatReaders(t *testing.T) {
	t.Run("two_readers", func(t *testing.T) {
		got, err := ConcatReaders(strings.NewReader("Hello"), strings.NewReader(" World"))
		if err != nil {
			t.Fatalf("❌ error: %v", err)
		}
		if got != "Hello World" {
			t.Errorf("❌ ConcatReaders = %q, want %q\n\t\t"+
				"Hint: io.MultiReader(readers...) chains them. "+
				"When first hits EOF, moves to next. Decorator pattern",
				got, "Hello World")
		} else {
			t.Logf("✅ ConcatReaders = %q", got)
		}
	})

	t.Run("three_readers", func(t *testing.T) {
		got, _ := ConcatReaders(
			strings.NewReader("a"),
			strings.NewReader("b"),
			strings.NewReader("c"),
		)
		if got != "abc" {
			t.Errorf("❌ ConcatReaders(a, b, c) = %q, want \"abc\"", got)
		} else {
			t.Logf("✅ ConcatReaders(a, b, c) = \"abc\"")
		}
	})

	t.Run("empty", func(t *testing.T) {
		got, _ := ConcatReaders()
		if got != "" {
			t.Errorf("❌ ConcatReaders() = %q, want \"\"", got)
		} else {
			t.Logf("✅ ConcatReaders() = \"\"")
		}
	})
}

// ────────────────────────────────────────────────────────────
// Exercise 7: ReadAndCapture
// ────────────────────────────────────────────────────────────

func TestReadAndCapture(t *testing.T) {
	r := strings.NewReader("stream data here")
	content, captured, err := ReadAndCapture(r)
	if err != nil {
		t.Fatalf("❌ error: %v", err)
	}
	if content != "stream data here" {
		t.Errorf("❌ content = %q, want %q\n\t\t"+
			"Hint: io.TeeReader(r, &buf) creates a reader that writes "+
			"everything it reads to buf. Like Unix tee command",
			content, "stream data here")
	} else {
		t.Logf("✅ content = %q", content)
	}
	if captured != "stream data here" {
		t.Errorf("❌ captured = %q, want %q (should match content)\n\t\t"+
			"Hint: The buffer receives an exact copy of what was read",
			captured, "stream data here")
	} else {
		t.Logf("✅ captured = %q (matches content)", captured)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: PipeTransfer
// ────────────────────────────────────────────────────────────

func TestPipeTransfer(t *testing.T) {
	tests := []struct {
		msg string
	}{
		{"hello via pipe"},
		{""},
		{"multi\nline\nmessage"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.msg), func(t *testing.T) {
			got, err := PipeTransfer(tt.msg)
			if err != nil {
				t.Fatalf("❌ error: %v", err)
			}
			if got != tt.msg {
				t.Errorf("❌ PipeTransfer(%q) = %q\n\t\t"+
					"Hint: io.Pipe() returns (reader, writer). "+
					"Write in a goroutine (or it blocks forever), "+
					"then io.ReadAll from the reader. Close the writer when done",
					tt.msg, got)
			} else {
				t.Logf("✅ PipeTransfer(%q) round-trip OK", tt.msg)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: CopyWithBuffer
// ────────────────────────────────────────────────────────────

func TestCopyWithBuffer(t *testing.T) {
	src := strings.NewReader("streaming copy test")
	n, got, err := CopyWithBuffer(src)
	if err != nil {
		t.Fatalf("❌ error: %v", err)
	}
	if got != "streaming copy test" {
		t.Errorf("❌ CopyWithBuffer content = %q, want %q\n\t\t"+
			"Hint: var dst bytes.Buffer; io.Copy(&dst, src) "+
			"— streams in 32KB chunks, never loads all into memory",
			got, "streaming copy test")
	} else {
		t.Logf("✅ CopyWithBuffer content = %q", got)
	}
	if n != 19 {
		t.Errorf("❌ bytes copied = %d, want 19", n)
	} else {
		t.Logf("✅ bytes copied = %d", n)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 10: FilterLines
// ────────────────────────────────────────────────────────────

func TestFilterLines(t *testing.T) {
	input := "error: disk full\ninfo: started\nerror: timeout\ninfo: done"
	r := strings.NewReader(input)
	var w bytes.Buffer

	keepErrors := func(line string) bool {
		return strings.HasPrefix(line, "error:")
	}

	count, err := FilterLines(r, &w, keepErrors)
	if err != nil {
		t.Fatalf("❌ error: %v", err)
	}
	if count != 2 {
		t.Errorf("❌ FilterLines count = %d, want 2\n\t\t"+
			"Hint: bufio.NewScanner(r), for scanner.Scan() { if keep(line) { fmt.Fprintln(w, line) } }. "+
			"This is the Reader→Filter→Writer pipeline pattern",
			count)
	} else {
		t.Logf("✅ FilterLines count = 2")
	}

	got := w.String()
	want := "error: disk full\nerror: timeout\n"
	if got != want {
		t.Errorf("❌ filtered output = %q, want %q", got, want)
	} else {
		t.Logf("✅ filtered output correct")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: WriterCounter
// ────────────────────────────────────────────────────────────

func TestWriterCounter(t *testing.T) {
	var buf bytes.Buffer
	wc := &WriterCounter{W: &buf}

	n1, err := fmt.Fprint(wc, "Hello")
	if err != nil {
		t.Fatalf("❌ first write error: %v", err)
	}
	n2, err := fmt.Fprint(wc, ", World!")
	if err != nil {
		t.Fatalf("❌ second write error: %v", err)
	}

	if buf.String() != "Hello, World!" {
		t.Errorf("❌ underlying buffer = %q, want %q\n\t\t"+
			"Hint: Delegate to wc.W.Write(p) — your Write wraps the original",
			buf.String(), "Hello, World!")
	} else {
		t.Logf("✅ underlying buffer = %q", buf.String())
	}

	total := int64(n1 + n2)
	if wc.BytesWritten() != total {
		t.Errorf("❌ BytesWritten() = %d, want %d\n\t\t"+
			"Hint: Track cumulative bytes: wc.Count += int64(n) after each Write. "+
			"This is the decorator pattern — add counting without changing the writer",
			wc.BytesWritten(), total)
	} else {
		t.Logf("✅ BytesWritten() = %d", total)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: SectionRead
// ────────────────────────────────────────────────────────────

func TestSectionRead(t *testing.T) {
	full := "Hello, World! This is Go."
	r := strings.NewReader(full)

	tests := []struct {
		offset int64
		length int64
		want   string
	}{
		{7, 5, "World"},
		{0, 5, "Hello"},
		{14, 11, "This is Go."},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("offset=%d/len=%d", tt.offset, tt.length), func(t *testing.T) {
			got, err := SectionRead(r, tt.offset, tt.length)
			if err != nil {
				t.Fatalf("❌ error: %v", err)
			}
			if got != tt.want {
				t.Errorf("❌ SectionRead(offset=%d, len=%d) = %q, want %q\n\t\t"+
					"Hint: io.NewSectionReader(r, offset, length) creates a windowed view. "+
					"*strings.Reader implements io.ReaderAt. *os.File does too (via pread syscall)",
					tt.offset, tt.length, got, tt.want)
			} else {
				t.Logf("✅ SectionRead(offset=%d, len=%d) = %q", tt.offset, tt.length, got)
			}
		})
	}
}

