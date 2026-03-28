// Standalone demo: Constructor Functions & the Stringer Interface
//
// Go has no constructors. Convention: use a New* function that returns
// a pointer to a fully initialized struct. This is the composition root
// pattern — all setup lives in one explicit place.
//
// Implementing fmt.Stringer (the String() method) gives your type custom
// output when passed to fmt.Println, fmt.Sprintf, etc. The fmt package
// checks for Stringer via a type assertion at runtime.
//
// Run:  go run .
package main

import "fmt"

// Stack is a LIFO data structure backed by a slice.
type Stack struct {
	items []int
}

// NewStack is the constructor — returns a ready-to-use *Stack.
func NewStack() *Stack {
	return &Stack{items: make([]int, 0)}
}

func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	n := len(s.items)
	item := s.items[n-1]
	s.items = s.items[:n-1]
	return item, true
}

func (s *Stack) Peek() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack) Len() int {
	return len(s.items)
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// String implements fmt.Stringer for custom printing.
func (s *Stack) String() string {
	return fmt.Sprintf("Stack%v", s.items)
}

func main() {
	s := NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	fmt.Println(s) // Stack[1 2 3]

	if top, ok := s.Pop(); ok {
		fmt.Println("Popped:", top) // 3
	}
	fmt.Println("Length:", s.Len()) // 2
}
