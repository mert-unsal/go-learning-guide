package strings_problems

// ============================================================
// Sherlock and the Valid String — [M]
// ============================================================
// A string is "valid" if all characters have the same frequency,
// OR if removing exactly ONE character makes all frequencies equal.
//
// Example: "aabbcc" -> "YES" (all freq 2)
// Example: "aabbccc" -> "YES" (remove one 'c')

// IsValid returns "YES" if the string can be made valid by removing at most 1 char.
// Time: O(n)  Space: O(n)
func IsValid(s string) string {
	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++
	}
	freqOfFreq := make(map[int]int)
	for _, f := range freq {
		freqOfFreq[f]++
	}

	if len(freqOfFreq) == 1 {
		return "YES"
	}
	if len(freqOfFreq) == 2 {
		freqs := make([]int, 0, 2)
		counts := make([]int, 0, 2)
		for f, c := range freqOfFreq {
			freqs = append(freqs, f)
			counts = append(counts, c)
		}
		f1, f2 := freqs[0], freqs[1]
		c1, c2 := counts[0], counts[1]
		if (f1 == f2+1 && c1 == 1) || (f2 == f1+1 && c2 == 1) {
			return "YES"
		}
		if (f1 == 1 && c1 == 1) || (f2 == 1 && c2 == 1) {
			return "YES"
		}
		if (f1 == f2+1 && c2 == len(freq)-1) || (f2 == f1+1 && c1 == len(freq)-1) {
			return "YES"
		}
	}
	return "NO"
}
