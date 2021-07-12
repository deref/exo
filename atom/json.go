package atom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/natefinch/atomic"
)

func DerefJSON(filename string, v interface{}) error {
	rv := reflect.ValueOf(v)
	rv.Elem().Set(reflect.Zero(rv.Elem().Type()))

	bs, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		bs = []byte("null")
		err = nil
	}
	if err != nil {
		return fmt.Errorf("reading: %w", err)
	}

	if err := json.Unmarshal(bs, v); err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}
	return nil
}

func ResetJSON(filename string, v interface{}) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encoding: %w", err)
	}
	if err := atomic.WriteFile(filename, &buf); err != nil {
		return fmt.Errorf("resetting: %w", err)
	}
	return nil
}

func SwapJSON(filename string, v interface{}, f func() error) error {
	for {
		if err := DerefJSON(filename, v); err != nil {
			return err
		}

		if err := f(); err != nil {
			return err
		}

		// XXX Do a compare-and-set instead of just clobbering.
		return ResetJSON(filename, v)
	}
}
