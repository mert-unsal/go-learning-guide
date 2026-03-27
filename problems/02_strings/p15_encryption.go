package strings_problems

import (
	"math"
	"strings"
)

// ============================================================
// Encryption — [M]
// ============================================================
// Encrypt a string by arranging it in a grid, then reading columns top-to-bottom.
// Grid size: floor(sqrt(len)) rows x ceil(sqrt(len)) cols.
//
// Example: s="haveaniceday"
// Grid (3x4):
//   h a v e
//   a n i c
//   e d a y
// Columns: "hae" "and" "via" "ecy" -> "hae and via ecy"

// Encryption encrypts a string using the grid column method.
// Time: O(n)  Space: O(n)
func Encryption(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	n := len(s)
	rows := int(math.Floor(math.Sqrt(float64(n))))
	cols := int(math.Ceil(math.Sqrt(float64(n))))
	if rows*cols < n {
		rows++
	}

	var result strings.Builder
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			idx := r*cols + c
			if idx < n {
				result.WriteByte(s[idx])
			}
		}
		if c < cols-1 {
			result.WriteByte(' ')
		}
	}
	return result.String()
}
