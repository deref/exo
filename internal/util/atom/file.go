package atom

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	"github.com/gofrs/flock"
	"github.com/natefinch/atomic"
)

// A atom durable on the local file system.
type FileAtom struct {
	filename string
	codec    Codec
}

func NewFileAtom(filename string, codec Codec) *FileAtom {
	return &FileAtom{
		filename: filename,
		codec:    codec,
	}
}

func (a *FileAtom) Deref(v interface{}) error {
	rv := reflect.ValueOf(v)
	rv.Elem().Set(reflect.Zero(rv.Elem().Type()))

	bs, err := ioutil.ReadFile(a.filename)
	if os.IsNotExist(err) {
		bs = []byte("null")
		err = nil
	}
	if err != nil {
		return fmt.Errorf("reading: %w", err)
	}

	if err := a.codec.Unmarshal(bs, v); err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}
	return nil
}

func (a *FileAtom) Reset(v interface{}) error {
	bs, err := a.codec.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}
	if err := atomic.WriteFile(a.filename, bytes.NewBuffer(bs)); err != nil {
		return fmt.Errorf("resetting: %w", err)
	}
	return nil
}

func (a *FileAtom) Swap(v interface{}, f func() error) error {
	// Establish timeout.
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Acquire lock.
	// TODO: Actual compare-and-set semantics, rather than mutual exclusion.
	lockPath := a.filename + ".lock"
	lock := flock.New(lockPath)
	retryDelay := 5 * time.Millisecond
	locked, err := lock.TryLockContext(ctx, retryDelay)
	if err != nil {
		return fmt.Errorf("locking %q: %w", lockPath, err)
	}
	if !locked {
		return fmt.Errorf("locking %q timed out", lockPath)
	}
	defer lock.Unlock()

	// Read, update-in-memory, Write.
	if err := a.Deref(v); err != nil {
		return err
	}
	if err := f(); err != nil {
		return err
	}
	return a.Reset(v)
}
