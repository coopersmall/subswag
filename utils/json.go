package utils

import (
	"encoding/json"
	// "github.com/ohler55/ojg/oj"
)

func Marshal(v any) ([]byte, error) {
	// data, err := oj.Marshal(v, &oj.Options{
	// 	UseTags:    true,
	// 	TimeFormat: "RFC3339",
	// 	NoReflect:  true,
	// })
	data, err := json.Marshal(v)
	return data, ErrorOrNil("unable to marshal JSON", NewJSONMarshError, err)
}

func Unmarshal(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	return ErrorOrNil("unable to unmarshal JSON", NewJSONMarshError, err)
}
