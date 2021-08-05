package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

func UnmarshalString(s string, v interface{}) error {
	if s == "" {
		s = "null"
	}
	return json.Unmarshal([]byte(s), v)
}

func MarshalString(v interface{}) (string, error) {
	bs, err := json.Marshal(v)
	return string(bs), err
}

func MustMarshalString(v interface{}) string {
	s, err := MarshalString(v)
	if err != nil {
		panic(err)
	}
	return s
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
