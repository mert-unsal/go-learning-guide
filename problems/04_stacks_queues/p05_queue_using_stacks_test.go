package stacks_queues

import "testing"

func TestMyQueue(t *testing.T) {
	q := &MyQueue{}

	q.Push(1)
	q.Push(2)

	if got := q.Peek(); got != 1 {
		t.Errorf("Peek() = %v, want 1", got)
	}
	if got := q.Pop(); got != 1 {
		t.Errorf("Pop() = %v, want 1", got)
	}
	if got := q.Empty(); got != false {
		t.Errorf("Empty() = %v, want false", got)
	}
	if got := q.Pop(); got != 2 {
		t.Errorf("Pop() = %v, want 2", got)
	}
	if got := q.Empty(); got != true {
		t.Errorf("Empty() = %v, want true", got)
	}
}
