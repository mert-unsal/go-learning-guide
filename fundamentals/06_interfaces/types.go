package interfaces

// Stringer is a shared interface used by exercises.
// Mirrors fmt.Stringer — any type with String() string satisfies it.
type Stringer interface {
	String() string
}
