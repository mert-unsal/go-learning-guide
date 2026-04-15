package strings_fundamentals

// ============================================================
// EXERCISES — 13 Strings: The Byte/Rune Duality
// ============================================================
//
// Go strings are read-only byte slices with a 2-word header {ptr, len}.
// len() returns BYTES, not characters. Iteration comes in two flavors:
// for-i gives bytes, for-range gives runes (decoded UTF-8).
//
// These exercises test your understanding of:
//   - Bytes vs runes vs characters (§3 of chapter 03)
//   - String immutability and its consequences (§2)
//   - UTF-8 variable-width encoding (§3)
//   - String iteration semantics (§4)
//   - Conversion costs and backing array sharing (§5, §6)
//   - Efficient concatenation (§7)
//
// Exercises 1-5:  Byte/rune duality — the #1 source of Go string bugs
// Exercises 6-8:  Immutability and memory layout
// Exercises 9-12: Production patterns and UTF-8 mastery
// ============================================================

// Exercise 1:
// ByteAndRuneCount returns the byte length and rune count of a string.
//
// KEY INSIGHT: len(s) returns BYTES, not characters. For "世界", len() = 6
// (3 bytes per CJK character in UTF-8), but there are only 2 runes.
//
// Use: len() for bytes, utf8.RuneCountInString() for runes.
func ByteAndRuneCount(s string) (byteLen int, runeCount int) {
	// TODO: return both byte length and rune count
	return 0, 0
}

// Exercise 2:
// ReverseString reverses a string correctly handling multi-byte runes.
//
// THE TRAP: Reversing bytes breaks multi-byte UTF-8 sequences.
// "Hello, 世界" reversed byte-by-byte produces garbage.
// You must convert to []rune, reverse, convert back.
//
// Example: "Hello, 世界" → "界世 ,olleH"
// Example: "café" → "éfac"
func ReverseString(s string) string {
	// TODO: convert to rune slice, reverse, convert back to string
	return ""
}

// Exercise 3:
// NthRune returns the rune at position n (0-indexed) in the string.
// Returns (0, false) if n is out of range.
//
// THE TRAP: s[n] gives you a BYTE, not a rune. For "café", s[3] returns
// 0xC3 (first byte of 'é'), not the rune 'é'.
// You must iterate runes to find the nth one.
func NthRune(s string, n int) (rune, bool) {
	// TODO: iterate runes, return the nth one
	return 0, false
}

// Exercise 4:
// IsASCII returns true if every byte in the string is in the ASCII range (0-127).
//
// INSIGHT: ASCII characters are always single-byte in UTF-8.
// If len(s) == utf8.RuneCountInString(s), the string is pure ASCII.
// Or simply check each byte: s[i] < 128.
func IsASCII(s string) bool {
	// TODO: check if all bytes are < 128
	return false
}

// Exercise 5:
// RuneByteOffsets returns the byte offset where each rune starts in the string.
//
// INSIGHT: for-range over a string gives (byte_index, rune). The byte index
// jumps by 1-4 depending on the UTF-8 encoding width of each rune.
//
// Example: "Go世界" → [0, 1, 2, 5]
//   G starts at byte 0 (1 byte), o at byte 1 (1 byte),
//   世 at byte 2 (3 bytes), 界 at byte 5 (3 bytes)
func RuneByteOffsets(s string) []int {
	// TODO: collect byte indices from for-range iteration
	return nil
}

// Exercise 6:
// ProveImmutability converts s to []byte, sets the first byte to 'X',
// and returns the ORIGINAL string and the modified byte slice.
//
// INSIGHT: []byte(s) copies the backing bytes. The original string is untouched.
// This proves string immutability: the conversion creates a new mutable copy.
//
// Example: ProveImmutability("hello") → ("hello", []byte("Xello"))
func ProveImmutability(s string) (original string, modified []byte) {
	// TODO: convert to []byte, modify first byte, return both
	return "", nil
}

// Exercise 7:
// ReplaceAtRuneIndex returns a new string with the rune at position idx
// replaced by newRune. Returns the original string if idx is out of range.
//
// THE CHALLENGE: You can't do s[idx] = newRune — strings are immutable,
// and indexing gives bytes anyway. You must:
// 1. Convert to []rune (or iterate to find the byte range)
// 2. Replace the target rune
// 3. Convert back to string
//
// Example: ReplaceAtRuneIndex("café", 3, 'o') → "cafo"
// Example: ReplaceAtRuneIndex("Go世界", 2, '地') → "Go地界"
func ReplaceAtRuneIndex(s string, idx int, newRune rune) string {
	// TODO: convert to rune slice, replace at idx, convert back
	return ""
}

// Exercise 8:
// SafeTruncate truncates a string to at most maxRunes runes.
// If truncated, appends "…" (U+2026 HORIZONTAL ELLIPSIS).
// If the string already fits, returns it unchanged.
//
// PRODUCTION PATTERN: Log messages, UI labels, API responses — anywhere
// you need to limit display length without breaking UTF-8.
//
// THE TRAP: s[:n] truncates by BYTES, which splits multi-byte runes.
// "café"[:4] = "caf\xc3" — broken UTF-8!
//
// Example: SafeTruncate("Hello, 世界!", 7) → "Hello, …"
// Example: SafeTruncate("Hi", 10) → "Hi"
func SafeTruncate(s string, maxRunes int) string {
	// TODO: count runes, if over maxRunes truncate at rune boundary + "…"
	return ""
}

// Exercise 9:
// CountByteClasses counts how many 1-byte, 2-byte, 3-byte, and 4-byte
// UTF-8 sequences are in the string. Returns [4]int where index 0 = 1-byte count, etc.
//
// INSIGHT: UTF-8 is variable-width. ASCII = 1 byte, accented = 2, CJK = 3, emoji = 4.
// You can determine the width from utf8.RuneLen(r) or utf8.DecodeRuneInString.
//
// Example: "Aé世😀" → [4]int{1, 1, 1, 1}
//   A = 1 byte, é = 2 bytes, 世 = 3 bytes, 😀 = 4 bytes
func CountByteClasses(s string) [4]int {
	// TODO: iterate runes, check byte width of each
	return [4]int{}
}

// Exercise 10:
// ConcatRepeat returns the string s repeated n times.
// Must use strings.Builder for efficiency — not the + operator.
//
// WHY: "x" + "x" + "x" ... is O(n²) — each + allocates and copies everything so far.
// strings.Builder uses an internal []byte with amortized O(1) appends.
// Pre-allocate with Grow(len(s)*n) for zero intermediate allocations.
//
// Example: ConcatRepeat("Go", 3) → "GoGoGo"
func ConcatRepeat(s string, n int) string {
	// TODO: use strings.Builder with Grow for efficient concatenation
	return ""
}

// Exercise 11:
// DetachSubstring returns the substring from rune index start to end (exclusive),
// ensuring the result does NOT share its backing array with the original string.
//
// THE MEMORY LEAK: s[start:end] shares the backing array. If s is 10MB
// and you keep a 5-rune substring, the ENTIRE 10MB stays in memory.
// Use strings.Clone() (Go 1.20+) to detach.
//
// Returns "" if start >= end or indices are out of range.
//
// Example: DetachSubstring("Hello, 世界!", 7, 9) → "世界"
func DetachSubstring(s string, runeStart, runeEnd int) string {
	// TODO: find byte range for rune indices, extract, Clone to detach
	return ""
}

// Exercise 12:
// EqualFoldASCII performs case-insensitive comparison of two strings.
// Only handles ASCII letters (A-Z, a-z). Do NOT use strings.EqualFold.
//
// INSIGHT: In Go, strings are comparable with == (byte-by-byte, O(n)).
// Case-insensitive comparison requires normalizing case manually.
// For ASCII: if byte is in 'A'-'Z', OR with 0x20 to make lowercase.
//
// Example: EqualFoldASCII("Hello", "hELLO") → true
// Example: EqualFoldASCII("café", "CAFÉ") → false (non-ASCII é ≠ É with this method)
func EqualFoldASCII(a, b string) bool {
	// TODO: compare byte-by-byte, folding A-Z to a-z
	return false
}
