package arrays

// ============================================================
// Counting Valleys — [E]
// ============================================================
// A hike is described as a string of U (up) and D (down) steps.
// A valley is a sequence of steps below sea level starting and ending at sea level.
// Count the number of valleys.
//
// Example: "UDDDUDUU" → 1 valley

// CountingValleys counts the number of valleys in the hike.
// Time: O(n)  Space: O(1)
func CountingValleys(steps string) int {
	level := 0
	valleys := 0
	for _, step := range steps {
		if step == 'U' {
			level++
			if level == 0 { // just crossed back to sea level from below
				valleys++
			}
		} else {
			level--
		}
	}
	return valleys
}
