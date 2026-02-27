// Package strings_strconv covers the strings and strconv packages —
// the most important stdlib packages for coding exam problems.
package strings_strconv

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ============================================================
// 1. strings PACKAGE — Essential functions
// ============================================================

func DemonstrateStrings() {
	s := "Hello, World!"

	// --- Searching ---
	fmt.Println(strings.Contains(s, "World"))  // true
	fmt.Println(strings.HasPrefix(s, "Hello")) // true
	fmt.Println(strings.HasSuffix(s, "!"))     // true
	fmt.Println(strings.Index(s, "World"))     // 7
	fmt.Println(strings.LastIndex(s, "l"))     // 10
	fmt.Println(strings.Count(s, "l"))         // 3

	// --- Case conversion ---
	fmt.Println(strings.ToUpper(s))           // HELLO, WORLD!
	fmt.Println(strings.ToLower(s))           // hello, world!
	fmt.Println(strings.Title("hello world")) // Hello World

	// --- Trimming ---
	fmt.Println(strings.TrimSpace("  hello  "))           // "hello"
	fmt.Println(strings.Trim("!!!hello!!!", "!"))         // "hello"
	fmt.Println(strings.TrimLeft("!!!hello!!!", "!"))     // "hello!!!"
	fmt.Println(strings.TrimRight("!!!hello!!!", "!"))    // "!!!hello"
	fmt.Println(strings.TrimPrefix("hello.go", "hello.")) // "go"
	fmt.Println(strings.TrimSuffix("hello.go", ".go"))    // "hello"

	// --- Splitting and joining ---
	parts := strings.Split("a,b,c,d", ",")
	fmt.Println(parts) // [a b c d]

	words := strings.Fields("  hello   world  ") // splits by whitespace
	fmt.Println(words)                           // [hello world]

	joined := strings.Join([]string{"a", "b", "c"}, "-")
	fmt.Println(joined) // a-b-c

	// --- Replacing ---
	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", 2)) // moo moo oink
	fmt.Println(strings.ReplaceAll("oink oink oink", "oink", "moo")) // moo moo moo

	// --- Repeating ---
	fmt.Println(strings.Repeat("na", 4)) // nananana
}

// ============================================================
// 2. strings.Builder — efficient string concatenation
// ============================================================
// NEVER use + in a loop to concatenate strings — O(n²)!
// Use strings.Builder instead — O(n)

func BuildString(parts []string) string {
	var sb strings.Builder
	for _, p := range parts {
		sb.WriteString(p)
	}
	return sb.String()
}

func ReverseString(s string) string {
	runes := []rune(s) // handle Unicode correctly
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func DemonstrateBuilder() {
	words := []string{"Hello", ", ", "World", "!"}
	fmt.Println(BuildString(words))

	fmt.Println(ReverseString("Hello, 世界"))
}

// ============================================================
// 3. strconv PACKAGE — String ↔ Number conversions
// ============================================================

func DemonstrateStrconv() {
	// --- int to string ---
	n := 42
	s := strconv.Itoa(n) // "42"
	fmt.Printf("Type: %T, Value: %q\n", s, s)

	// string to int
	n2, err := strconv.Atoi("123")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Parsed:", n2)
	}

	// Invalid string
	_, err = strconv.Atoi("abc")
	fmt.Println("Error:", err) // strconv.Atoi: parsing "abc": invalid syntax

	// --- float to string ---
	f := 3.14159
	fs := strconv.FormatFloat(f, 'f', 2, 64) // "3.14"
	fmt.Println(fs)

	// string to float
	f2, _ := strconv.ParseFloat("3.14", 64)
	fmt.Println(f2 + 1) // 4.140...

	// --- bool ---
	b, _ := strconv.ParseBool("true")
	fmt.Println(b)                         // true
	fmt.Println(strconv.FormatBool(false)) // "false"

	// --- ParseInt with base ---
	hex, _ := strconv.ParseInt("FF", 16, 64)     // parse hex
	fmt.Println(hex)                             // 255
	binary, _ := strconv.ParseInt("1010", 2, 64) // parse binary
	fmt.Println(binary)                          // 10
}

// ============================================================
// 4. unicode PACKAGE — Character classification
// ============================================================

func DemonstrateUnicode() {
	chars := []rune{'A', 'a', '5', ' ', '!', '世'}
	for _, ch := range chars {
		fmt.Printf("%c: IsLetter=%t IsDigit=%t IsSpace=%t IsUpper=%t IsLower=%t\n",
			ch,
			unicode.IsLetter(ch),
			unicode.IsDigit(ch),
			unicode.IsSpace(ch),
			unicode.IsUpper(ch),
			unicode.IsLower(ch),
		)
	}

	// Case conversion
	fmt.Println(string(unicode.ToUpper('a'))) // A
	fmt.Println(string(unicode.ToLower('A'))) // a
}

// ============================================================
// 5. Common String Problems (LeetCode patterns)
// ============================================================

// Check if two strings are anagrams
func IsAnagram(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	count := [26]int{} // fixed-size array for 'a'-'z'
	for i := 0; i < len(s); i++ {
		count[s[i]-'a']++
		count[t[i]-'a']--
	}
	for _, c := range count {
		if c != 0 {
			return false
		}
	}
	return true
}

// Check if a string is a palindrome (ignoring non-alphanumeric)
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; {
		for i < j && !unicode.IsLetter(runes[i]) && !unicode.IsDigit(runes[i]) {
			i++
		}
		for i < j && !unicode.IsLetter(runes[j]) && !unicode.IsDigit(runes[j]) {
			j--
		}
		if runes[i] != runes[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func DemonstratePatterns() {
	fmt.Println("anagram:", IsAnagram("anagram", "nagaram"))                   // true
	fmt.Println("anagram:", IsAnagram("rat", "car"))                           // false
	fmt.Println("palindrome:", IsPalindrome("A man, a plan, a canal: Panama")) // true
	fmt.Println("palindrome:", IsPalindrome("race a car"))                     // false
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== strings Package ===")
	DemonstrateStrings()
	fmt.Println("\n=== strings.Builder ===")
	DemonstrateBuilder()
	fmt.Println("\n=== strconv ===")
	DemonstrateStrconv()
	fmt.Println("\n=== unicode ===")
	DemonstrateUnicode()
	fmt.Println("\n=== String Patterns ===")
	DemonstratePatterns()
}
