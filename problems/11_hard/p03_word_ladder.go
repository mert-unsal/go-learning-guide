package hard

// ============================================================
// PROBLEM 3: Word Ladder (LeetCode #127) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   Given two words, beginWord and endWord, and a dictionary
//   wordList, return the number of words in the shortest
//   transformation sequence from beginWord to endWord, such that
//   only one letter can be changed at a time and each transformed
//   word must exist in the word list. Return 0 if no such sequence.
//
// PARAMETERS:
//   beginWord string   — the starting word
//   endWord   string   — the target word
//   wordList  []string — dictionary of allowed intermediate words
//
// RETURN:
//   int — length of the shortest transformation sequence, or 0 if none exists
//
// CONSTRAINTS:
//   • 1 <= beginWord.length <= 10
//   • endWord.length == beginWord.length
//   • 1 <= wordList.length <= 5000
//   • wordList[i].length == beginWord.length
//   • beginWord, endWord, and wordList[i] consist of lowercase English letters
//   • beginWord != endWord
//   • All words in wordList are unique
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log","cog"]
//   Output: 5
//   Why:    "hit" -> "hot" -> "dot" -> "dog" -> "cog" (5 words)
//
// Example 2:
//   Input:  beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log"]
//   Output: 0
//   Why:    endWord "cog" is not in wordList, no valid transformation.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • BFS from beginWord; each level is one transformation step
// • Use wildcard patterns (e.g., "h*t") to find neighbors efficiently
// • Bidirectional BFS can reduce search space dramatically
// • Target: O(M^2 * N) time, O(M^2 * N) space where M=word length, N=wordList size
func LadderLength(beginWord string, endWord string, wordList []string) int {
	return 0
}
