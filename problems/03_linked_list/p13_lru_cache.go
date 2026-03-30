package linked_list

// PROBLEM 13: LRU Cache (LeetCode #146) — MEDIUM
// Design a data structure: get(key) and put(key, value) both O(1).
// Approach: doubly-linked list + hash map.
// Hash map: key → node pointer. DLL: most recently used at head, LRU at tail.
// On get: move node to head. On put: add to head, evict tail if over capacity.

type lruNode struct {
	key, val   int
	prev, next *lruNode
}

// LRUCache is an LRU cache with O(1) get and put.
type LRUCache struct {
	capacity   int
	cache      map[int]*lruNode
	head, tail *lruNode // sentinel nodes
}

func NewLRUCache(capacity int) *LRUCache {
	return nil
}

func (c *LRUCache) Get(key int) int {
	return 0
}

func (c *LRUCache) Put(key, value int) {
}

func (c *LRUCache) remove(node *lruNode) {
}

func (c *LRUCache) addToFront(node *lruNode) {
}
