package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/deref/util-go/jsonutil"
)

func UnmarshalString(s string, v any) error {
	if s == "" {
		s = "null"
	}
	return json.Unmarshal([]byte(s), v)
}

func MustUnmarshal(bs []byte, v any) {
	if err := json.Unmarshal(bs, v); err != nil {
		panic(err)
	}
}

func MustUnmarshalString(s string, v any) {
	if err := UnmarshalString(s, v); err != nil {
		panic(err)
	}
}

var MustMarshal = jsonutil.MustMarshal
var MarshalString = jsonutil.MarshalString
var MustMarshalString = jsonutil.MustMarshalString
var MarshalIndentString = jsonutil.MarshalIndentString
var MustMarshalIndentString = jsonutil.MustMarshalIndentString

func UnmarshalStringOrEmpty(s string, v any) error {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "{}"
	}
	return UnmarshalString(s, v)
}

func UnmarshalReader(r io.Reader, v any) error {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	bs = bytes.TrimSpace(bs)
	if len(bs) == 0 {
		return nil
	}
	return json.Unmarshal(bs, v)
}

func UnmarshalFile(filePath string, v any) error {
	bs, err := ioutil.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	bs = bytes.TrimSpace(bs)
	if len(bs) == 0 {
		return nil
	}
	return json.Unmarshal(bs, v)
}

func MarshalFile(filePath string, v any) error {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	bs = append(bs, '\n')
	return ioutil.WriteFile(filePath, bs, 0600)
}

func PrettyPrintJSONString(w io.Writer, jsonStr string) error {
	var val any
	if err := UnmarshalString(jsonStr, &val); err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(val)
}

func IsValid(jsonStr string) bool {
	var val any
	err := UnmarshalString(jsonStr, &val)
	return err == nil
}

func Merge(objs ...map[string]any) map[string]any {
	res := map[string]any{}
	for _, o := range objs {
		if o == nil {
			continue
		}
		for k, v := range o {
			res[k] = v
		}
	}
	return res
}

// Returns a new JSON value, simplified by replacing marshalable values with
// the standard JSON-compatible Go types. For example, a (*string)(nil) will be
// replaced by (any)(nil) and json.Marshaler implementations will be
// replaced by the results of round-tripping their MarshalJSON output.
func MustSimplify(v any) any {
	var res any
	MustUnmarshal(MustMarshal(v), &res)
	return res
}
