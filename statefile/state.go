package statefile

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/natefinch/atomic"
)

type State struct {
	mx       sync.Mutex
	filename string
}

func New(filename string) *State {
	return &State{
		filename: filename,
	}
}

type root struct {
	Components map[string]component `json:"components"`
}

type component struct {
	IRI string `json:"iri"`
}

func (state *State) update(ctx context.Context, f func(root *root) error) error {
	state.mx.Lock()
	defer state.mx.Unlock()

	in, err := ioutil.ReadFile(state.filename)
	if os.IsNotExist(err) {
		in = []byte("{}")
		err = nil
	}
	if err != nil {
		return fmt.Errorf("reading: %w", err)
	}

	var root root
	if err := json.Unmarshal(in, &root); err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}
	if root.Components == nil {
		root.Components = make(map[string]component)
	}

	if err := f(&root); err != nil {
		return err
	}

	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&root); err != nil {
		panic(err)
	}

	if err := atomic.WriteFile(state.filename, &out); err != nil {
		return fmt.Errorf("writing: %w", err)
	}
	return nil
}

func (state *State) Remember(ctx context.Context, name string, iri string) error {
	return state.update(ctx, func(root *root) error {
		root.Components[name] = component{
			IRI: iri,
		}
		return nil
	})
}
