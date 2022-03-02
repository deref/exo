package scalars

import (
	"reflect"

	graphql "github.com/graph-gophers/graphql-go/decode"
	"github.com/mitchellh/mapstructure"
)

func DecodeStruct(input, output interface{}) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     output,
		Squash:     true,
		TagName:    "json",
		DecodeHook: decodeStructHook,
	})
	if err != nil {
		panic(err)
	}
	err = dec.Decode(input)
	return err
}

func decodeStructHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	switch to.Kind() {
	case reflect.Int, reflect.String, reflect.Bool:
		return data, nil
	}
	result := reflect.New(to).Interface()
	unmarshaller, ok := result.(graphql.Unmarshaler)
	if !ok {
		return data, nil
	}
	if err := unmarshaller.UnmarshalGraphQL(data); err != nil {
		return nil, err
	}
	return result, nil
}
