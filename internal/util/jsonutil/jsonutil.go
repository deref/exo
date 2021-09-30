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

func UnmarshalString(s string, v interface{}) error {
	if s == "" {
		s = "null"
	}
	return json.Unmarshal([]byte(s), v)
}

var MarshalString = jsonutil.MarshalString
var MustMarshalString = jsonutil.MustMarshalString

func UnmarshalStringOrEmpty(s string, v interface{}) error {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "{}"
	}
	return UnmarshalString(s, v)
}

func UnmarshalReader(r io.Reader, v interface{}) error {
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

func UnmarshalFile(filePath string, v interface{}) error {
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

func MarshalFile(filePath string, v interface{}) error {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	bs = append(bs, '\n')
	return ioutil.WriteFile(filePath, bs, 0600)
}

func PrettyPrintJSONString(w io.Writer, jsonStr string) error {
	var val interface{}
	if err := UnmarshalString(jsonStr, &val); err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(val)
}

func IsValid(jsonStr string) bool {
	var val interface{}
	err := UnmarshalString(jsonStr, &val)
	return err == nil
}
