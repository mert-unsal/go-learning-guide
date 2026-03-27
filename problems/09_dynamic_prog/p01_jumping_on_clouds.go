package dynamic_prog

// ============================================================
// Jumping on the Clouds — [E]
// ============================================================
// Jump on clouds. 0 = safe cloud, 1 = thunder cloud (skip).
// You can jump 1 or 2 clouds. Find the minimum jumps to reach the last cloud.
//
// Example: c=[0,0,1,0,0,1,0] → 4

// JumpingOnClouds returns the minimum number of jumps to the last cloud.
// Time: O(n)  Space: O(1)
func JumpingOnClouds(c []int) int {
	jumps := 0
	i := 0
	n := len(c)
	for i < n-1 {
		// Prefer jump of 2 if it's safe, else jump 1
		if i+2 < n && c[i+2] == 0 {
			i += 2
		} else {
			i++
		}
		jumps++
	}
	return jumps
}
