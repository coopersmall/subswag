package utils

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

type JSONPathQuery string

func EvaluateJSONPathQuery(query JSONPathQuery, data any) (any, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return "", NewJSONMarshError("failed to marshal json while evaluating json path", err)
	}
	result := gjson.Get(string(json), string(query))
	if !result.Exists() {
		return "", nil
	}
	if result.IsArray() {
		var results []any
		for _, r := range result.Array() {
			results = append(results, r.Value())
		}
		return results, nil
	}
	return result.Value(), nil
}
