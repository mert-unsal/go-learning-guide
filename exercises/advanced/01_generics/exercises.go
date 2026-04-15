package generics

// ============================================================
// EXERCISES -- 01 generics: Type Parameters & Constraints
// ============================================================
// 12 exercises covering Go generics at production depth.
// Focus: type parameters, constraints, when generics help vs hurt.

// ────────────────────────────────────────────────────────────
// Exercise 1: Min -- basic generic function with constraints
// ────────────────────────────────────────────────────────────
// Return the smaller of two values. Works for any ordered type.
// Use the cmp.Ordered constraint.

import "cmp"

func Min[T cmp.Ordered](a, b T) T {
	var zero T
	_ = a
	_ = b
	return zero
}

// ────────────────────────────────────────────────────────────
// Exercise 2: Contains -- generic slice search
// ────────────────────────────────────────────────────────────
// Return true if target exists in the slice.
// Use comparable constraint (supports == operator).

func Contains[T comparable](s []T, target T) bool {
	return false
}

// ────────────────────────────────────────────────────────────
// Exercise 3: Map -- transform a slice with a function
// ────────────────────────────────────────────────────────────
// Apply fn to each element, return new slice of results.
// Two type parameters: input T, output U.

func Map[T any, U any](s []T, fn func(T) U) []U {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 4: Filter -- keep elements matching a predicate
// ────────────────────────────────────────────────────────────

func Filter[T any](s []T, keep func(T) bool) []T {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 5: Reduce -- fold a slice into a single value
// ────────────────────────────────────────────────────────────
// Start with initial value, apply fn(accumulator, element) for each.

func Reduce[T any, U any](s []T, initial U, fn func(U, T) U) U {
	return initial
}

// ────────────────────────────────────────────────────────────
// Exercise 6: Keys and Values -- extract from map
// ────────────────────────────────────────────────────────────

func Keys[K comparable, V any](m map[K]V) []K {
	return nil
}

func Values[K comparable, V any](m map[K]V) []V {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 7: Stack -- generic data structure
// ────────────────────────────────────────────────────────────
// Implement a stack with Push, Pop, Peek, Len.
// Pop and Peek return (T, bool) — false if empty.

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	// TODO: append to items
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	return zero, false
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	return zero, false
}

func (s *Stack[T]) Len() int {
	return 0
}

// ────────────────────────────────────────────────────────────
// Exercise 8: Pair -- generic struct with method
// ────────────────────────────────────────────────────────────
// A typed key-value pair. Swap returns a new Pair with K and V exchanged.

type Pair[K any, V any] struct {
	Key   K
	Value V
}

func (p Pair[K, V]) Swap() Pair[V, K] {
	var zero Pair[V, K]
	return zero
}

// ────────────────────────────────────────────────────────────
// Exercise 9: MaxBy -- generic max with custom comparator
// ────────────────────────────────────────────────────────────
// Return the element for which fn returns the highest value.
// Panics if slice is empty.

func MaxBy[T any](s []T, fn func(T) int) T {
	var zero T
	return zero
}

// ────────────────────────────────────────────────────────────
// Exercise 10: GroupBy -- group elements by a key function
// ────────────────────────────────────────────────────────────
// Returns map[K][]T where K is the result of the key function.

func GroupBy[T any, K comparable](s []T, keyFn func(T) K) map[K][]T {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 11: Result -- generic result type (Ok/Err pattern)
// ────────────────────────────────────────────────────────────
// Implement a Result type that holds either a value or an error.

type Result[T any] struct {
	value T
	err   error
	ok    bool
}

func Ok[T any](v T) Result[T] {
	return Result[T]{}
}

func Err[T any](err error) Result[T] {
	return Result[T]{}
}

func (r Result[T]) Unwrap() (T, error) {
	var zero T
	return zero, r.err
}

func (r Result[T]) IsOk() bool {
	return false
}

// ────────────────────────────────────────────────────────────
// Exercise 12: Set -- generic set using map[T]struct{}
// ────────────────────────────────────────────────────────────

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

func (s *Set[T]) Add(v T) {
	// TODO: s.m[v] = struct{}{}
}

func (s *Set[T]) Has(v T) bool {
	return false
}

func (s *Set[T]) Remove(v T) {
	// TODO: delete(s.m, v)
}

func (s *Set[T]) Len() int {
	return 0
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	return NewSet[T]()
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	return NewSet[T]()
}
