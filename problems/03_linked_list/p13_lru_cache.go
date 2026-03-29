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
	// TODO: implement — create head/tail sentinels, link them
	var head = &lruNode{}
	var tail = &lruNode{}
	head.next = tail
	tail.prev = head
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*lruNode),
		head:     head, tail: tail,
	}
}

func (c *LRUCache) Get(key int) int {
	// TODO: implement
	if localLruNode, ok := c.cache[key]; ok {
		c.remove(localLruNode)
		c.addToFront(localLruNode)
		return localLruNode.val
	}
	return -1
}

func (c *LRUCache) Put(key, value int) {
	// TODO: implement
	if localLruNode, ok := c.cache[key]; ok {
		c.cache[key].val = value
		c.remove(localLruNode)
		c.addToFront(localLruNode)
	} else {
		newLruNode := &lruNode{key: key, val: value}
		c.addToFront(newLruNode)
		c.cache[key] = newLruNode
	}

	if len(c.cache) > c.capacity {
		delete(c.cache, c.tail.prev.key)
		c.remove(c.tail.prev)
	}
}

func (c *LRUCache) remove(node *lruNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *LRUCache) addToFront(node *lruNode) {
	node.prev = c.head
	node.next = c.head.next
	c.head.next.prev = node
	c.head.next = node
}
