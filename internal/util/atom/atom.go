package atom

// Atom are uncoordinated, atomic reference cells.
// Inspired by <https://clojure.org/reference/atoms>.
type Atom interface {
	// Sets v to the current value in the atom.
	Deref(v any) error
	// Sets the current value in the atom to v.
	Reset(v any) error
	// Runs f in a compare-and-set loop until a write succeeds uncontended.
	// On each iterate of the loop, v is set to the current value in the atom.
	// If no error occurs, the successfully written value remains in v.
	Swap(v any, f func() error) error
}
