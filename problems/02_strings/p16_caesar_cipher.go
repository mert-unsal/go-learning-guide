package strings_problems

// ============================================================
// Caesar Cipher — [E]
// ============================================================
// Shift each letter by k positions (wrapping around), preserve case and non-letters.
//
// Example: s="middle-Outz", k=2 → "okffng-Qwvb"

// CaesarCipher applies Caesar cipher with shift k.
// Time: O(n)  Space: O(n)
func CaesarCipher(s string, k int) string {
	k = k % 26 // handle large shifts
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch >= 'a' && ch <= 'z':
			result[i] = byte((int(ch-'a')+k)%26) + 'a'
		case ch >= 'A' && ch <= 'Z':
			result[i] = byte((int(ch-'A')+k)%26) + 'A'
		default:
			result[i] = ch // non-letter: unchanged
		}
	}
	return string(result)
}
