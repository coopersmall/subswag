package utils

import (
	"encoding/json"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func ValidateJSONSchema(schema string, data any) (bool, error) {
	compiler := jsonschema.NewCompiler()
	err := compiler.AddResource("schema.json", strings.NewReader(schema))
	if err != nil {
		return false, err
	}
	json, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	err = compiler.AddResource("data.json", strings.NewReader(string(json)))
	if err != nil {
		return false, err
	}
	s, err := compiler.Compile("schema.json")
	if err != nil {
		return false, err
	}
	err = s.Validate("data.json")
	if err != nil {
		return false, err
	}
	return true, nil
}

func IsValidJSONSchema(schema string) (bool, error) {
	compiler := jsonschema.NewCompiler()
	err := compiler.AddResource("schema.json", strings.NewReader(schema))
	if err != nil {
		return false, err
	}
	_, err = compiler.Compile("schema.json")
	if err != nil {
		return false, err
	}
	return true, nil
}
