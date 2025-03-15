package booleanexpression

import "github.com/coopersmall/subswag/utils"

// Operator represents comparison operators
type Operator string

const (
	OperatorEqual        Operator = "=="
	OperatorNotEqual     Operator = "!="
	OperatorGreaterThan  Operator = ">"
	OperatorGreaterEqual Operator = ">="
	OperatorLessThan     Operator = "<"
	OperatorLessEqual    Operator = "<="
)

// ConditionType represents the type of condition
type ConditionType string

const (
	ConditionTypeNumeric ConditionType = "numeric"
	ConditionTypeString  ConditionType = "string"
	ConditionTypeBoolean ConditionType = "boolean"
	ConditionTypeArray   ConditionType = "array"
)

// Condition represents a base condition interface
type Condition interface {
	GetType() ConditionType
	GetActualJSONPath() utils.JSONPathQuery
}

// NumericCondition represents a condition for numeric values
type NumericCondition struct {
	Type           ConditionType       `json:"type"`
	ActualJSONPath utils.JSONPathQuery `json:"actual_json_path"`
	Operator       Operator            `json:"operator"`
	ExpectedValue  float64             `json:"expected_value"`
}

func (n NumericCondition) IsBooleanExpressionOrCondition() {}

func (n NumericCondition) GetType() ConditionType {
	return n.Type
}

func (n NumericCondition) GetActualJSONPath() utils.JSONPathQuery {
	return n.ActualJSONPath
}

// StringCondition represents a condition for string values
type StringCondition struct {
	Type           ConditionType       `json:"type"`
	ActualJSONPath utils.JSONPathQuery `json:"actual_json_path"`
	Operator       Operator            `json:"operator"`
	ExpectedValue  string              `json:"expected_value"`
}

func (s StringCondition) IsBooleanExpressionOrCondition() {}

func (s StringCondition) GetType() ConditionType {
	return s.Type
}

func (s StringCondition) GetActualJSONPath() utils.JSONPathQuery {
	return s.ActualJSONPath
}

type BooleanCondition struct {
	Type           ConditionType       `json:"type"`
	ActualJSONPath utils.JSONPathQuery `json:"actual_json_path"`
	Operator       Operator            `json:"operator"`
	ExpectedValue  bool                `json:"expected_value"`
}

func (b BooleanCondition) IsBooleanExpressionOrCondition() {}

func (b BooleanCondition) GetType() ConditionType {
	return b.Type
}

func (b BooleanCondition) GetActualJSONPath() utils.JSONPathQuery {
	return b.ActualJSONPath
}

type ArryConditionOperator string

const (
	OperatorContains    ArryConditionOperator = "contains"
	OperatorContainsAll ArryConditionOperator = "contains_all"
	OperatorContainsAny ArryConditionOperator = "contains_any"
	OperatorNotContains ArryConditionOperator = "not_contains"
)

type ArrayCondition struct {
	Type           ConditionType         `json:"type"`
	ActualJSONPath utils.JSONPathQuery   `json:"actual_json_path"`
	Operator       ArryConditionOperator `json:"operator"`
	ExpectedValue  []any                 `json:"expected_value"`
}

func (a ArrayCondition) IsBooleanExpressionOrCondition() {}

func (a ArrayCondition) GetType() ConditionType {
	return a.Type
}

func (a ArrayCondition) GetActualJSONPath() utils.JSONPathQuery {
	return a.ActualJSONPath
}

// NewNumericCondition creates a new NumericCondition
func NewNumericCondition(jsonPath utils.JSONPathQuery, operator Operator, expectedValue float64) NumericCondition {
	return NumericCondition{
		Type:           ConditionTypeNumeric,
		ActualJSONPath: jsonPath,
		Operator:       operator,
		ExpectedValue:  expectedValue,
	}
}

// NewStringCondition creates a new StringCondition
func NewStringCondition(jsonPath utils.JSONPathQuery, operator Operator, expectedValue string) StringCondition {
	return StringCondition{
		Type:           ConditionTypeString,
		ActualJSONPath: jsonPath,
		Operator:       operator,
		ExpectedValue:  expectedValue,
	}
}

func NewBooleanCondition(jsonPath utils.JSONPathQuery, operator Operator, expectedValue bool) BooleanCondition {
	return BooleanCondition{
		Type:           ConditionTypeBoolean,
		ActualJSONPath: jsonPath,
		Operator:       operator,
		ExpectedValue:  expectedValue,
	}
}

func NewArrayCondition(jsonPath utils.JSONPathQuery, operator ArryConditionOperator, expectedValue []any) ArrayCondition {
	return ArrayCondition{
		Type:           ConditionTypeArray,
		ActualJSONPath: jsonPath,
		Operator:       operator,
		ExpectedValue:  expectedValue,
	}
}
