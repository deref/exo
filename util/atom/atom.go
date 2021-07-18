package atom

// An uncoordinated, atomic reference cell in the style of Clojure.
type Atom interface {
	// Sets v to the current value in the atom.
	Deref(v interface{}) error
	// Sets the current value in the atom to v.
	Reset(v interface{}) error
	// Runs f in a compare-and-set loop until a write succeeds uncontended.
	// On each iterate of the loop, v is set to the current value in the atom.
	// If no error occurs, the successfully written value remains in v.
	Swap(v interface{}, f func() error) error
}
