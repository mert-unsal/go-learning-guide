package trees

// ============================================================
// PROBLEM 16: Implement Trie / Prefix Tree (LeetCode #208) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Implement a trie (prefix tree) with Insert, Search, and
//   StartsWith operations. A trie efficiently stores and retrieves
//   keys in a dataset of strings.
//
// OPERATIONS:
//   NewTrie()                      — Initialize the trie object
//   Insert(word string)            — Insert word into the trie
//   Search(word string) bool       — Return true if word is in the trie
//   StartsWith(prefix string) bool — Return true if any word starts with prefix
//
// CONSTRAINTS:
//   • 1 <= word.length, prefix.length <= 2000
//   • word and prefix consist only of lowercase English letters
//   • At most 3 * 10^4 calls to Insert, Search, and StartsWith
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   trie := NewTrie()
//   trie.Insert("apple")
//   trie.Search("apple")   → true
//   trie.Search("app")     → false
//   trie.StartsWith("app") → true
//   trie.Insert("app")
//   trie.Search("app")     → true
//
// Example 2:
//   trie := NewTrie()
//   trie.Insert("hello")
//   trie.Search("hell")     → false
//   trie.StartsWith("hell") → true
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Each node has up to 26 children (one per letter) and an isEnd flag
// • Insert: walk/create nodes for each character, mark last as end
// • Search: walk nodes for each character, return false if missing; check isEnd at last
// • StartsWith: same as Search but don't check isEnd — just verify the path exists
// • Target: O(m) time per operation where m = word/prefix length, O(n*m) space total
func NewTrie() *Trie {
	return nil
}
func (t *Trie) Insert(word string) {
}
func (t *Trie) Search(word string) bool {
	return false
}
func (t *Trie) StartsWith(prefix string) bool {
	return false
}
