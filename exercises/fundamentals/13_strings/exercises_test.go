package strings_fundamentals

import (
	"fmt"
	"reflect"
	"testing"
)

// ────────────────────────────────────────────────────────────
// Test helpers
// ────────────────────────────────────────────────────────────

func assertEq[T comparable](t *testing.T, name string, got, want T, hint string) {
	t.Helper()
	if got != want {
		t.Errorf("❌ %s = %v, want %v\n\t\tHint: %s", name, got, want, hint)
	} else {
		t.Logf("✅ %s = %v", name, got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 1: ByteAndRuneCount
// ────────────────────────────────────────────────────────────

func TestByteAndRuneCount(t *testing.T) {
	tests := []struct {
		input     string
		wantBytes int
		wantRunes int
	}{
		{"Hello", 5, 5},                 // ASCII: 1 byte per rune
		{"世界", 6, 2},                    // CJK: 3 bytes per rune
		{"café", 5, 4},                  // 'é' = 2 bytes in UTF-8
		{"😀🎉", 8, 2},                    // emoji: 4 bytes per rune
		{"", 0, 0},                      // empty string
		{"Go世界!", 10, 5},                // mixed: 2×1 + 2×3 + 1×1 = 9... wait
	}

	// Fix the last test case: "Go世界!" = G(1) + o(1) + 世(3) + 界(3) + !(1) = 9 bytes, 5 runes
	tests[5] = struct {
		input     string
		wantBytes int
		wantRunes int
	}{"Go世界!", 9, 5}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.input), func(t *testing.T) {
			gotBytes, gotRunes := ByteAndRuneCount(tt.input)
			assertEq(t, fmt.Sprintf("ByteAndRuneCount(%q) bytes", tt.input), gotBytes, tt.wantBytes,
				"len(s) returns BYTES. Use utf8.RuneCountInString(s) for rune count")
			assertEq(t, fmt.Sprintf("ByteAndRuneCount(%q) runes", tt.input), gotRunes, tt.wantRunes,
				"'世' is 3 bytes in UTF-8 but 1 rune. See learnings/03 §3")
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: ReverseString
// ────────────────────────────────────────────────────────────

func TestReverseString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "olleh"},
		{"Hello, 世界", "界世 ,olleH"},   // multi-byte runes must stay intact
		{"café", "éfac"},              // 2-byte rune 'é' must not break
		{"😀🎉", "🎉😀"},                  // 4-byte emojis
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.input), func(t *testing.T) {
			got := ReverseString(tt.input)
			assertEq(t, fmt.Sprintf("ReverseString(%q)", tt.input), got, tt.want,
				"Convert to []rune first, reverse the rune slice, then string(). "+
					"Reversing bytes breaks multi-byte UTF-8. See learnings/03 §3")
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: NthRune
// ────────────────────────────────────────────────────────────

func TestNthRune(t *testing.T) {
	tests := []struct {
		input   string
		n       int
		want    rune
		wantOK  bool
	}{
		{"hello", 0, 'h', true},
		{"hello", 4, 'o', true},
		{"café", 3, 'é', true},       // s[3] would give 0xC3, NOT 'é'
		{"Go世界", 2, '世', true},       // rune index 2, byte index 2
		{"Go世界", 3, '界', true},       // rune index 3, byte index 5
		{"hello", 5, 0, false},        // out of range
		{"hello", -1, 0, false},       // negative index
		{"", 0, 0, false},            // empty string
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q[%d]", tt.input, tt.n), func(t *testing.T) {
			got, ok := NthRune(tt.input, tt.n)
			assertEq(t, fmt.Sprintf("NthRune(%q, %d) ok", tt.input, tt.n), ok, tt.wantOK,
				"s[n] gives a BYTE, not a rune! Iterate with for-range to find the nth rune")
			if ok {
				assertEq(t, fmt.Sprintf("NthRune(%q, %d) rune", tt.input, tt.n), got, tt.want,
					"For 'café', s[3]=0xC3 (byte), but rune index 3 is 'é'. See learnings/03 §3")
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: IsASCII
// ────────────────────────────────────────────────────────────

func TestIsASCII(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"hello", true},
		{"Hello, World!", true},
		{"café", false},               // 'é' is 2 bytes
		{"世界", false},
		{"", true},                    // empty is vacuously ASCII
		{"abc123!@#", true},
		{"naïve", false},              // 'ï' is non-ASCII
		{"\x00\x7f", true},           // boundary: 0x00 and 0x7F are valid ASCII
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.input), func(t *testing.T) {
			got := IsASCII(tt.input)
			assertEq(t, fmt.Sprintf("IsASCII(%q)", tt.input), got, tt.want,
				"ASCII = bytes 0-127. Check each byte: s[i] < 128. "+
					"Or compare len(s) == utf8.RuneCountInString(s)")
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: RuneByteOffsets
// ────────────────────────────────────────────────────────────

func TestRuneByteOffsets(t *testing.T) {
	tests := []struct {
		input string
		want  []int
	}{
		{"Go世界", []int{0, 1, 2, 5}},      // G(1) o(1) 世(3) 界(3)
		{"hello", []int{0, 1, 2, 3, 4}},  // all single-byte
		{"café", []int{0, 1, 2, 3}},      // c(1) a(1) f(1) é(2) — é starts at byte 3
		{"😀", []int{0}},                    // single 4-byte rune
		{"", []int{}},                     // empty → empty offsets (not nil)
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.input), func(t *testing.T) {
			got := RuneByteOffsets(tt.input)
			if got == nil && len(tt.want) == 0 {
				got = []int{} // normalize nil to empty for comparison
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("❌ RuneByteOffsets(%q) = %v, want %v\n\t\tHint: "+
					"for i, _ := range s gives byte offsets. 世 is 3 bytes, so indices jump. See learnings/03 §4",
					tt.input, got, tt.want)
			} else {
				t.Logf("✅ RuneByteOffsets(%q) = %v", tt.input, got)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: ProveImmutability
// ────────────────────────────────────────────────────────────

func TestProveImmutability(t *testing.T) {
	tests := []struct {
		input        string
		wantOriginal string
		wantFirst    byte
	}{
		{"hello", "hello", 'X'},
		{"world", "world", 'X'},
		{"Go", "Go", 'X'},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.input), func(t *testing.T) {
			original, modified := ProveImmutability(tt.input)
			assertEq(t, "original unchanged", original, tt.wantOriginal,
				"[]byte(s) COPIES the backing bytes. The original string is untouched. "+
					"See learnings/03 §2 and §5")
			if len(modified) == 0 {
				t.Errorf("❌ modified slice is empty — did you forget to convert and modify?")
				return
			}
			assertEq(t, "modified[0]", modified[0], tt.wantFirst,
				"Set modified[0] = 'X'. This changes the []byte copy, not the string")
			// Verify the rest of the bytes are preserved
			if len(modified) != len(tt.input) {
				t.Errorf("❌ modified length = %d, want %d", len(modified), len(tt.input))
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: ReplaceAtRuneIndex
// ────────────────────────────────────────────────────────────

func TestReplaceAtRuneIndex(t *testing.T) {
	tests := []struct {
		s       string
		idx     int
		newRune rune
		want    string
	}{
		{"hello", 0, 'H', "Hello"},
		{"café", 3, 'o', "cafo"},
		{"Go世界", 2, '地', "Go地界"},
		{"hello", 4, '!', "hell!"},
		{"hello", 5, '!', "hello"},    // out of range → unchanged
		{"hello", -1, '!', "hello"},   // negative → unchanged
		{"", 0, 'x', ""},             // empty → unchanged
		{"a", 0, '😀', "😀"},           // replace ASCII with emoji (grows in bytes)
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q[%d]='%c'", tt.s, tt.idx, tt.newRune), func(t *testing.T) {
			got := ReplaceAtRuneIndex(tt.s, tt.idx, tt.newRune)
			assertEq(t, fmt.Sprintf("ReplaceAtRuneIndex(%q, %d, '%c')", tt.s, tt.idx, tt.newRune),
				got, tt.want,
				"Strings are immutable — you can't s[i]='X'. Convert to []rune, replace, convert back. "+
					"See learnings/03 §2")
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: SafeTruncate
// ────────────────────────────────────────────────────────────

func TestSafeTruncate(t *testing.T) {
	tests := []struct {
		s        string
		maxRunes int
		want     string
	}{
		{"Hello, 世界!", 7, "Hello, …"},    // 9 runes → truncate at 7 + "…"
		{"Hello", 10, "Hello"},            // fits → unchanged
		{"Hello", 5, "Hello"},             // exact fit → unchanged
		{"café latte", 4, "café…"},        // truncate mid-word
		{"", 5, ""},                       // empty → empty
		{"Hello", 0, "…"},                 // maxRunes=0 but has content → just ellipsis
		{"Hi", 1, "H…"},                   // truncate after 1 rune
		{"世界", 1, "世…"},                   // truncate multi-byte
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q/%d", tt.s, tt.maxRunes), func(t *testing.T) {
			got := SafeTruncate(tt.s, tt.maxRunes)
			assertEq(t, fmt.Sprintf("SafeTruncate(%q, %d)", tt.s, tt.maxRunes), got, tt.want,
				"s[:n] truncates by BYTES and breaks multi-byte runes! "+
					"Iterate runes to find the byte position after maxRunes runes. "+
					"Append '…' (U+2026). See learnings/03 §3")
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: CountByteClasses
// ────────────────────────────────────────────────────────────

func TestCountByteClasses(t *testing.T) {
	tests := []struct {
		input string
		want  [4]int
	}{
		{"Aé世😀", [4]int{1, 1, 1, 1}},       // one of each class
		{"hello", [4]int{5, 0, 0, 0}},       // all ASCII
		{"世界你好", [4]int{0, 0, 4, 0}},         // all 3-byte
		{"😀🎉🚀💯", [4]int{0, 0, 0, 4}},          // all 4-byte
		{"", [4]int{0, 0, 0, 0}},            // empty
		{"café", [4]int{3, 1, 0, 0}},        // c,a,f = 1-byte; é = 2-byte
		{"naïve", [4]int{4, 1, 0, 0}},       // n,a,v,e = 1-byte; ï = 2-byte
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.input), func(t *testing.T) {
			got := CountByteClasses(tt.input)
			if got != tt.want {
				t.Errorf("❌ CountByteClasses(%q) = %v, want %v\n\t\tHint: "+
					"Use utf8.RuneLen(r) after decoding each rune with for-range. "+
					"ASCII=1, accented=2, CJK=3, emoji=4. See learnings/03 §3",
					tt.input, got, tt.want)
			} else {
				t.Logf("✅ CountByteClasses(%q) = %v", tt.input, got)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 10: ConcatRepeat
// ────────────────────────────────────────────────────────────

func TestConcatRepeat(t *testing.T) {
	tests := []struct {
		s    string
		n    int
		want string
	}{
		{"Go", 3, "GoGoGo"},
		{"ab", 0, ""},
		{"", 5, ""},
		{"x", 1, "x"},
		{"世", 3, "世世世"},
		{"Hi", 5, "HiHiHiHiHi"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q×%d", tt.s, tt.n), func(t *testing.T) {
			got := ConcatRepeat(tt.s, tt.n)
			assertEq(t, fmt.Sprintf("ConcatRepeat(%q, %d)", tt.s, tt.n), got, tt.want,
				"Use strings.Builder with Grow(len(s)*n) to pre-allocate. "+
					"Never use += in a loop — it's O(n²). See learnings/03 §7")
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: DetachSubstring
// ────────────────────────────────────────────────────────────

func TestDetachSubstring(t *testing.T) {
	tests := []struct {
		s         string
		runeStart int
		runeEnd   int
		want      string
	}{
		{"Hello, 世界!", 7, 9, "世界"},
		{"Hello", 0, 5, "Hello"},
		{"Hello", 1, 3, "el"},
		{"café", 2, 4, "fé"},
		{"Hello", 3, 3, ""},          // start == end → empty
		{"Hello", 4, 2, ""},          // start > end → empty
		{"Hello", 0, 10, ""},         // end out of range → empty
		{"", 0, 0, ""},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q[%d:%d]", tt.s, tt.runeStart, tt.runeEnd), func(t *testing.T) {
			got := DetachSubstring(tt.s, tt.runeStart, tt.runeEnd)
			assertEq(t, fmt.Sprintf("DetachSubstring(%q, %d, %d)", tt.s, tt.runeStart, tt.runeEnd),
				got, tt.want,
				"Convert rune indices to byte indices, then use strings.Clone(s[byteStart:byteEnd]) "+
					"to detach from the backing array. See learnings/03 §6")
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: EqualFoldASCII
// ────────────────────────────────────────────────────────────

func TestEqualFoldASCII(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"Hello", "hELLO", true},
		{"hello", "hello", true},
		{"hello", "world", false},
		{"Go", "go", true},
		{"Go", "GO", true},
		{"abc", "ab", false},          // different lengths
		{"", "", true},
		{"ABC123", "abc123", true},    // digits are case-insensitive (they have no case)
		{"café", "CAFÉ", false},       // non-ASCII: é (0xC3A9) ≠ É (0xC389) at byte level
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q==%q", tt.a, tt.b), func(t *testing.T) {
			got := EqualFoldASCII(tt.a, tt.b)
			assertEq(t, fmt.Sprintf("EqualFoldASCII(%q, %q)", tt.a, tt.b), got, tt.want,
				"Compare byte-by-byte. To fold A-Z to a-z: if b >= 'A' && b <= 'Z' { b |= 0x20 }. "+
					"This only works for ASCII. Non-ASCII bytes are compared as-is. "+
					"For full Unicode folding, use strings.EqualFold. See learnings/03 §8")
		})
	}
}
