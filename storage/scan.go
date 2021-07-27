package storage

type Scanner interface {
	// Scan returns an iterator capable of scanning some elements in an index.
	// It is expected to be a fast operation that does not need to perform any
	// I/O. I/O should only be done when calling ScanIter.Next().
	Scan(txn ReadTransaction, args ScanArgs) ScanIter
}

type ScanArgs struct {
	Prefix    []byte
	Direction ScanDirection
	Limit     int

	// KeyOnly is a flag indicating that only the `Key` must be populated in
	// the resulting `ScanEntry` structs. A storage engine can this flag as
	// a hint to perform a key-only scan if there is a performance advantage.
	KeyOnly bool

	// TODO: Add filter and support predicate pushdown.
}

type ScanDirection int

const (
	ScanDirectionASC  ScanDirection = 1
	ScanDirectionDESC ScanDirection = -1
)

// ScanIter is an iterator that should return the successive entries from a scan.
// A ScanIter is acquired by calling the `Scan` method of `Scanner`. If acquiring
// a scanner could fail, that failure should result in Scan returning an error
// rather than a ScanIter in an error state being returned.
//
// Example:
// ```
// it, err := scanner.Scan(txn, args)
// if err != nil {
//	return err
// }
// defer it.Close()
// for it.Next() {
//   entry := it.Item()
//   // ...
// }
// return it.Err()
// ```
type ScanIter interface {
	// Next advances the iterator and returns a boolean indicating whether there are
	// more items available. A call to Next() advances the iterator regardless of whether
	// the item was obtained using Item(). The iterator is not assumed to be in a readible
	// state until the first call to Next(). This behaviour is useful so that `iter.Next()`
	// can be called as the condition of a `for` loop that calls `iter.Item()` in its body.
	// If the iterator is in an error state or has been exhausted, Next() should return false.
	Next() bool

	// Item yields the current item in the iterator. It should return `nil` if the iterator
	// has been exhausted or is in an error state. It may panic if Next() has not yet been
	// called. Calling Item() should not advance the iterator state.
	Item() *ScanEntry

	// Err returns the error if the iterator is in an error state. This method
	// should be called after an iteration loop to check whether the iteration completed
	// normally or was aborted with an error.
	Err() error

	// Close performs any clean-up that may need to be done. The implementation is not expected
	// to perform any clean-up when the iteration is complete or an error is encountered, so this
	// method should always be called by the client code.
	Close()

	// TODO: Reverse()
}

type ScanEntry struct {
	Key, Value []byte
}

func Collect(it ScanIter) ([]ScanEntry, error) {
	defer it.Close()
	var results []ScanEntry
	for it.Next() {
		if entry := it.Item(); entry != nil {
			results = append(results, *entry)
		}
	}
	return results, it.Err()
}
