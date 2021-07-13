package jsonutil

import "encoding/json"

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
