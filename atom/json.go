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
	rv.Set(reflect.Zero(rv.Type()))

	bs, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
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

func SwapJSON(filename string, v interface{}, f func() error) error {
	for {
		if err := DerefJSON(filename, v); err != nil {
			return err
		}

		if err := f(); err != nil {
			return err
		}

		// XXX Do a compare-and-set instead of just clobbering.
		return atomic.WriteFile(filename, bytes.NewBuffer(bs))
	}
}
