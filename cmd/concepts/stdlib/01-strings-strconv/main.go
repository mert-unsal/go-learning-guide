// Strings & Strconv in Go — demonstrates the strings, strconv, and unicode packages.
//
// These are the most important stdlib packages for string manipulation
// and type conversions in coding problems and production code.
//
// Run: go run cmd/concepts/stdlib/01-strings-strconv/main.go
package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
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
	fmt.Printf("%s%s  Strings, Strconv & Unicode              %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// ============================================================
	// 1. strings PACKAGE — Essential functions
	// ============================================================
	fmt.Printf("%s▸ 1. strings Package — Essential Functions%s\n", cyan+bold, reset)

	s := "Hello, World!"

	// --- Searching ---
	fmt.Printf("\n  %s✔ Searching%s\n", green, reset)
	fmt.Println("  Contains(s, \"World\"):", strings.Contains(s, "World"))   // true
	fmt.Println("  HasPrefix(s, \"Hello\"):", strings.HasPrefix(s, "Hello")) // true
	fmt.Println("  HasSuffix(s, \"!\"):", strings.HasSuffix(s, "!"))         // true
	fmt.Println("  Index(s, \"World\"):", strings.Index(s, "World"))         // 7
	fmt.Println("  LastIndex(s, \"l\"):", strings.LastIndex(s, "l"))         // 10
	fmt.Println("  Count(s, \"l\"):", strings.Count(s, "l"))                 // 3

	// --- Case conversion ---
	fmt.Printf("\n  %s✔ Case Conversion%s\n", green, reset)
	fmt.Println("  ToUpper:", strings.ToUpper(s))           // HELLO, WORLD!
	fmt.Println("  ToLower:", strings.ToLower(s))           // hello, world!

	// --- Trimming ---
	fmt.Printf("\n  %s✔ Trimming%s\n", green, reset)
	fmt.Println("  TrimSpace(\"  hello  \"):", strings.TrimSpace("  hello  "))                 // "hello"
	fmt.Println("  Trim(\"!!!hello!!!\", \"!\"):", strings.Trim("!!!hello!!!", "!"))            // "hello"
	fmt.Println("  TrimLeft(\"!!!hello!!!\", \"!\"):", strings.TrimLeft("!!!hello!!!", "!"))    // "hello!!!"
	fmt.Println("  TrimRight(\"!!!hello!!!\", \"!\"):", strings.TrimRight("!!!hello!!!", "!"))  // "!!!hello"
	fmt.Println("  TrimPrefix(\"hello.go\", \"hello.\"):", strings.TrimPrefix("hello.go", "hello.")) // "go"
	fmt.Println("  TrimSuffix(\"hello.go\", \".go\"):", strings.TrimSuffix("hello.go", ".go"))       // "hello"

	// --- Splitting and joining ---
	fmt.Printf("\n  %s✔ Splitting & Joining%s\n", green, reset)
	parts := strings.Split("a,b,c,d", ",")
	fmt.Println("  Split(\"a,b,c,d\", \",\"):", parts) // [a b c d]

	words := strings.Fields("  hello   world  ") // splits by whitespace
	fmt.Println("  Fields(\"  hello   world  \"):", words) // [hello world]

	joined := strings.Join([]string{"a", "b", "c"}, "-")
	fmt.Println("  Join([a b c], \"-\"):", joined) // a-b-c

	// --- Replacing ---
	fmt.Printf("\n  %s✔ Replacing%s\n", green, reset)
	fmt.Println("  Replace(\"oink oink oink\", \"oink\", \"moo\", 2):", strings.Replace("oink oink oink", "oink", "moo", 2)) // moo moo oink
	fmt.Println("  ReplaceAll(\"oink oink oink\", \"oink\", \"moo\"):", strings.ReplaceAll("oink oink oink", "oink", "moo")) // moo moo moo

	// --- Repeating ---
	fmt.Printf("\n  %s✔ Repeating%s\n", green, reset)
	fmt.Println("  Repeat(\"na\", 4):", strings.Repeat("na", 4)) // nananana

	// ============================================================
	// 2. strings.Builder — efficient string concatenation
	// ============================================================
	// NEVER use + in a loop to concatenate strings — O(n²)!
	// Use strings.Builder instead — O(n)

	fmt.Printf("\n%s▸ 2. strings.Builder — Efficient Concatenation%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ NEVER use + in a loop — O(n²). Use strings.Builder — O(n)%s\n", yellow, reset)

	buildWords := []string{"Hello", ", ", "World", "!"}
	fmt.Println("  BuildString:", BuildString(buildWords))
	fmt.Println("  ReverseString(\"Hello, 世界\"):", ReverseString("Hello, 世界"))

	// ============================================================
	// 3. strconv PACKAGE — String ↔ Number conversions
	// ============================================================
	fmt.Printf("\n%s▸ 3. strconv Package — String ↔ Number Conversions%s\n", cyan+bold, reset)

	// --- int to string ---
	fmt.Printf("\n  %s✔ int ↔ string%s\n", green, reset)
	n := 42
	sv := strconv.Itoa(n) // "42"
	fmt.Printf("  Itoa(42): Type: %T, Value: %q\n", sv, sv)

	// string to int
	n2, err := strconv.Atoi("123")
	if err != nil {
		fmt.Println("  Error:", err)
	} else {
		fmt.Println("  Atoi(\"123\"):", n2)
	}

	// Invalid string
	_, err = strconv.Atoi("abc")
	fmt.Println("  Atoi(\"abc\") Error:", err) // strconv.Atoi: parsing "abc": invalid syntax

	// --- float to string ---
	fmt.Printf("\n  %s✔ float ↔ string%s\n", green, reset)
	f := 3.14159
	fs := strconv.FormatFloat(f, 'f', 2, 64) // "3.14"
	fmt.Println("  FormatFloat(3.14159, 'f', 2, 64):", fs)

	// string to float
	f2, _ := strconv.ParseFloat("3.14", 64)
	fmt.Println("  ParseFloat(\"3.14\", 64) + 1:", f2+1) // 4.140...

	// --- bool ---
	fmt.Printf("\n  %s✔ bool ↔ string%s\n", green, reset)
	b, _ := strconv.ParseBool("true")
	fmt.Println("  ParseBool(\"true\"):", b)                // true
	fmt.Println("  FormatBool(false):", strconv.FormatBool(false)) // "false"

	// --- ParseInt with base ---
	fmt.Printf("\n  %s✔ ParseInt with base%s\n", green, reset)
	hex, _ := strconv.ParseInt("FF", 16, 64)     // parse hex
	fmt.Println("  ParseInt(\"FF\", 16, 64):", hex) // 255
	binary, _ := strconv.ParseInt("1010", 2, 64)  // parse binary
	fmt.Println("  ParseInt(\"1010\", 2, 64):", binary) // 10

	// ============================================================
	// 4. unicode PACKAGE — Character classification
	// ============================================================
	fmt.Printf("\n%s▸ 4. unicode Package — Character Classification%s\n", cyan+bold, reset)

	chars := []rune{'A', 'a', '5', ' ', '!', '世'}
	for _, ch := range chars {
		fmt.Printf("  %c: IsLetter=%s%t%s IsDigit=%s%t%s IsSpace=%s%t%s IsUpper=%s%t%s IsLower=%s%t%s\n",
			ch,
			magenta, unicode.IsLetter(ch), reset,
			magenta, unicode.IsDigit(ch), reset,
			magenta, unicode.IsSpace(ch), reset,
			magenta, unicode.IsUpper(ch), reset,
			magenta, unicode.IsLower(ch), reset,
		)
	}

	// Case conversion
	fmt.Println("  ToUpper('a'):", string(unicode.ToUpper('a'))) // A
	fmt.Println("  ToLower('A'):", string(unicode.ToLower('A'))) // a

	// ============================================================
	// 5. Common String Problems (LeetCode patterns)
	// ============================================================
	fmt.Printf("\n%s▸ 5. Common String Patterns (LeetCode)%s\n", cyan+bold, reset)

	fmt.Println("  IsAnagram(\"anagram\", \"nagaram\"):", IsAnagram("anagram", "nagaram"))                        // true
	fmt.Println("  IsAnagram(\"rat\", \"car\"):", IsAnagram("rat", "car"))                                        // false
	fmt.Println("  IsPalindrome(\"A man, a plan, a canal: Panama\"):", IsPalindrome("A man, a plan, a canal: Panama")) // true
	fmt.Println("  IsPalindrome(\"race a car\"):", IsPalindrome("race a car"))                                    // false
}

// BuildString efficiently concatenates parts using strings.Builder.
func BuildString(parts []string) string {
	var sb strings.Builder
	for _, p := range parts {
		sb.WriteString(p)
	}
	return sb.String()
}

// ReverseString reverses a string handling Unicode correctly via []rune.
func ReverseString(s string) string {
	runes := []rune(s) // handle Unicode correctly
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsAnagram checks if two strings are anagrams using frequency counting.
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

// IsPalindrome checks if a string is a palindrome (ignoring non-alphanumeric).
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
