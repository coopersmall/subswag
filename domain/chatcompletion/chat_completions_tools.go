package chatcompletion

import (
	"encoding/json"
)

type CallId string

type ToolCall struct {
	ID        CallId          `json:"id"`
	Name      string          `json:"functionName"`
	Arguments json.RawMessage `json:"arguments"`
}

type Tool struct {
	Name        string            `json:"name"`
	Parameters  map[string]Schema `json:"parameters"`
	Description string            `json:"description,omitempty"`
}

type Schema interface {
	Type() SchemaType
	Title() string
	Description() string
}

type SchemaType string

const (
	SchemaTypeString  SchemaType = "string"
	SchemaTypeNumber  SchemaType = "number"
	SchemaTypeBoolean SchemaType = "boolean"
	SchemaTypeNull    SchemaType = "null"
	SchemaTypeArray   SchemaType = "array"
	SchemaTypeObject  SchemaType = "object"
)

type SchemaDescription struct {
	T           SchemaType
	title       string
	description string
}

func NewSchemaDescription(
	t SchemaType,
	title, description string,
) SchemaDescription {
	return SchemaDescription{
		T:           t,
		title:       title,
		description: description,
	}
}

func (s SchemaDescription) Type() SchemaType {
	return s.T
}

func (s SchemaDescription) Title() string {
	return s.title
}

func (s SchemaDescription) Description() string {
	return s.description
}

type StringSchema struct {
	SchemaDescription
	Enum   []string `json:"enum,omitempty"`
	Format *string  `json:"format,omitempty"`
}

type NumberSchema struct {
	SchemaDescription
	Enum []float64 `json:"enum,omitempty"`
}

type BooleanSchema struct {
	SchemaDescription
}

type NullSchema struct {
	SchemaDescription
}

type ArraySchema struct {
	SchemaDescription
	Items json.RawMessage `json:"items"`
}

type ObjectSchema struct {
	SchemaDescription
	Properties map[string]json.RawMessage `json:"properties"`
	Required   []string                   `json:"required,omitempty"`
}
