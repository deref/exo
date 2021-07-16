package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

func UnmarshalString(s string, v interface{}) error {
	if s == "" {
		s = "null"
	}
	return json.Unmarshal([]byte(s), v)
}

func MustMarshalString(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bs)
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
	if err != nil {
		return err
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
