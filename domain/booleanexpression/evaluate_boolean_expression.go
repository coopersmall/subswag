package booleanexpression

import (
	"encoding/json"
	"fmt"

	"github.com/coopersmall/subswag/utils"
)

// EvaluationError represents an error during evaluation
type EvaluationError struct {
	Message string
	Context map[string]any
}

func (e *EvaluationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Context)
}

// EvaluateBooleanExpression evaluates a boolean expression against provided data
func EvaluateBooleanExpression(expression BooleanExpression, data any) (bool, error) {
	switch expression.Operator {
	case OperatorAND:
		return evaluateAND(expression.Conditions, data)
	case OperatorOR:
		return evaluateOR(expression.Conditions, data)
	default:
		return false, &EvaluationError{
			Message: "unsupported operator",
			Context: map[string]any{
				"operator": expression.Operator,
			},
		}
	}
}

func evaluateAND(conditions []IsBooleanExpressionOrCondition, data any) (bool, error) {
	for _, condition := range conditions {
		result, err := evaluateConditionOrExpression(condition, data)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

func evaluateOR(conditions []IsBooleanExpressionOrCondition, data any) (bool, error) {
	for _, condition := range conditions {
		result, err := evaluateConditionOrExpression(condition, data)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

func evaluateConditionOrExpression(conditionOrExpression IsBooleanExpressionOrCondition, data any) (bool, error) {
	switch c := conditionOrExpression.(type) {
	case NumericCondition:
		return evaluateNumericCondition(c, data)
	case StringCondition:
		return evaluateStringCondition(c, data)
	case BooleanCondition:
		return evaluateBooleanCondition(c, data)
	case ArrayCondition:
		return evaluateArrayCondition(c, data)
	case BooleanExpression:
		return EvaluateBooleanExpression(c, data)
	default:
		return false, &EvaluationError{
			Message: "unsupported condition type",
			Context: map[string]any{
				"type": fmt.Sprintf("%T", c),
			},
		}
	}
}

func evaluateNumericCondition(condition NumericCondition, data any) (bool, error) {
	actual, err := getNumericValue(condition.ActualJSONPath, data)
	if err != nil {
		return false, err
	}
	return evaluateOperator(condition.Operator, actual, condition.ExpectedValue), nil
}

func evaluateStringCondition(condition StringCondition, data any) (bool, error) {
	actual, err := getStringValue(condition.ActualJSONPath, data)
	if err != nil {
		return false, err
	}
	return evaluateOperator(condition.Operator, actual, condition.ExpectedValue), nil
}

func evaluateBooleanCondition(condition BooleanCondition, data any) (bool, error) {
	actual, err := getBooleanValue(condition.ActualJSONPath, data)
	if err != nil {
		return false, err
	}
	return evaluateOperator(condition.Operator, actual, condition.ExpectedValue), nil
}

func evaluateArrayCondition(condition ArrayCondition, data any) (bool, error) {
	actual, err := getArrayValue(condition.ActualJSONPath, data)
	if err != nil {
		return false, err
	}

	switch condition.Operator {
	case OperatorContains:
		return evaluateArrayContains(actual, condition.ExpectedValue), nil
	case OperatorContainsAll:
		return evaluateArrayContainsAll(actual, condition.ExpectedValue), nil
	case OperatorContainsAny:
		return evaluateArrayContainsAny(actual, condition.ExpectedValue), nil
	case OperatorNotContains:
		return !evaluateArrayContains(actual, condition.ExpectedValue), nil
	default:
		return false, &EvaluationError{
			Message: "unsupported array operator",
			Context: map[string]any{
				"operator": condition.Operator,
			},
		}
	}
}

func evaluateOperator(operator Operator, actual, expected any) bool {
	switch operator {
	case OperatorEqual:
		return actual == expected
	case OperatorNotEqual:
		return actual != expected
	case OperatorGreaterThan:
		switch a := actual.(type) {
		case float64:
			if e, ok := expected.(float64); ok {
				return a > e
			}
		}
	case OperatorGreaterEqual:
		switch a := actual.(type) {
		case float64:
			if e, ok := expected.(float64); ok {
				return a >= e
			}
		}
	case OperatorLessThan:
		switch a := actual.(type) {
		case float64:
			if e, ok := expected.(float64); ok {
				return a < e
			}
		}
	case OperatorLessEqual:
		switch a := actual.(type) {
		case float64:
			if e, ok := expected.(float64); ok {
				return a <= e
			}
		}
	}
	return false
}

func evaluateArrayContains(actual []any, expected []any) bool {
	for _, e := range expected {
		found := false
		for _, a := range actual {
			if equalValues(a, e) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func evaluateArrayContainsAll(actual []any, expected []any) bool {
	for _, e := range expected {
		found := false
		for _, a := range actual {
			if equalValues(a, e) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func evaluateArrayContainsAny(actual []any, expected []any) bool {
	for _, e := range expected {
		for _, a := range actual {
			if equalValues(a, e) {
				return true
			}
		}
	}
	return false
}

func equalValues(a, b any) bool {
	aJson, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bJson, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(aJson) == string(bJson)
}

func getStringValue(path utils.JSONPathQuery, data any) (string, error) {
	value, err := utils.EvaluateJSONPathQuery(path, data)
	if err != nil {
		return "", err
	}
	if str, ok := value.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("value is not a string")
}

func getNumericValue(path utils.JSONPathQuery, data any) (float64, error) {
	value, err := utils.EvaluateJSONPathQuery(path, data)
	if err != nil {
		return 0, err
	}
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	}
	return 0, fmt.Errorf("value is not a number")
}

func getBooleanValue(path utils.JSONPathQuery, data any) (bool, error) {
	value, err := utils.EvaluateJSONPathQuery(path, data)
	if err != nil {
		return false, err
	}
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("value is not a boolean")
}

func getArrayValue(path utils.JSONPathQuery, data any) ([]any, error) {
	value, err := utils.EvaluateJSONPathQuery(path, data)
	if err != nil {
		return nil, err
	}
	if arr, ok := value.([]any); ok {
		return arr, nil
	}
	// Handle single value as array
	return []any{value}, nil
}
