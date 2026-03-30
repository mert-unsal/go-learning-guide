package trees

import "testing"

func TestTrie(t *testing.T) {
	trie := NewTrie()
	if trie == nil {
		t.Skip("NewTrie not implemented yet")
	}
	trie.Insert("apple")
	if got := trie.Search("apple"); !got {
		t.Errorf("Search(apple) = %v, want true", got)
	}
	if got := trie.Search("app"); got {
		t.Errorf("Search(app) = %v, want false", got)
	}
	if got := trie.StartsWith("app"); !got {
		t.Errorf("StartsWith(app) = %v, want true", got)
	}
}
