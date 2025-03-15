package booleanexpression_test

import (
	"github.com/coopersmall/subswag/domain/booleanexpression"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
)

func (s *BooleanExpressionTestSuite) TestEvaluateBooleanExpression() {
	s.Run("should evaluate AND conditions correctly", func() {
		data := map[string]interface{}{
			"age":    25,
			"name":   "John",
			"active": true,
			"tags":   []interface{}{"user", "premium"},
		}

		expression := booleanexpression.BooleanExpression{
			Operator: booleanexpression.OperatorAND,
			Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
				booleanexpression.NumericCondition{
					ActualJSONPath: utils.JSONPathQuery("age"),
					Operator:       booleanexpression.OperatorGreaterThan,
					ExpectedValue:  20.0,
				},
				booleanexpression.StringCondition{
					ActualJSONPath: utils.JSONPathQuery("name"),
					Operator:       booleanexpression.OperatorEqual,
					ExpectedValue:  "John",
				},
			},
		}

		result, err := booleanexpression.EvaluateBooleanExpression(expression, data)
		assert.NoError(s.T(), err)
		assert.True(s.T(), result)
	})

	s.Run("should evaluate OR conditions correctly", func() {
		data := map[string]interface{}{
			"score": 75,
			"grade": "B",
		}

		expression := booleanexpression.BooleanExpression{
			Operator: booleanexpression.OperatorOR,
			Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
				booleanexpression.NumericCondition{
					ActualJSONPath: utils.JSONPathQuery("score"),
					Operator:       booleanexpression.OperatorGreaterEqual,
					ExpectedValue:  90.0,
				},
				booleanexpression.StringCondition{
					ActualJSONPath: utils.JSONPathQuery("grade"),
					Operator:       booleanexpression.OperatorEqual,
					ExpectedValue:  "B",
				},
			},
		}

		result, err := booleanexpression.EvaluateBooleanExpression(expression, data)
		assert.NoError(s.T(), err)
		assert.True(s.T(), result)
	})

	s.Run("should evaluate array conditions correctly", func() {
		data := map[string]interface{}{
			"roles": []interface{}{"admin", "user", "editor"},
		}

		expression := booleanexpression.BooleanExpression{
			Operator: booleanexpression.OperatorAND,
			Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
				booleanexpression.ArrayCondition{
					ActualJSONPath: utils.JSONPathQuery("roles"),
					Operator:       booleanexpression.OperatorContains,
					ExpectedValue:  []interface{}{"admin"},
				},
			},
		}

		result, err := booleanexpression.EvaluateBooleanExpression(expression, data)
		assert.NoError(s.T(), err)
		assert.True(s.T(), result)
	})

	s.Run("should handle invalid JSON paths", func() {
		data := map[string]interface{}{
			"value": 42,
		}

		expression := booleanexpression.BooleanExpression{
			Operator: booleanexpression.OperatorAND,
			Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
				booleanexpression.NumericCondition{
					ActualJSONPath: utils.JSONPathQuery("$.nonexistent"),
					Operator:       booleanexpression.OperatorEqual,
					ExpectedValue:  42.0,
				},
			},
		}

		result, err := booleanexpression.EvaluateBooleanExpression(expression, data)
		assert.Error(s.T(), err)
		assert.False(s.T(), result)
	})

	s.Run("should handle nested expressions", func() {
		data := map[string]interface{}{
			"user": map[string]interface{}{
				"age":    30,
				"active": true,
			},
		}

		nestedExpression := booleanexpression.BooleanExpression{
			Operator: booleanexpression.OperatorAND,
			Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
				booleanexpression.NumericCondition{
					ActualJSONPath: utils.JSONPathQuery("user.age"),
					Operator:       booleanexpression.OperatorGreaterThan,
					ExpectedValue:  25.0,
				},
				booleanexpression.BooleanCondition{
					ActualJSONPath: utils.JSONPathQuery("user.active"),
					Operator:       booleanexpression.OperatorEqual,
					ExpectedValue:  true,
				},
			},
		}

		result, err := booleanexpression.EvaluateBooleanExpression(nestedExpression, data)
		assert.NoError(s.T(), err)
		assert.True(s.T(), result)
	})

	s.Run("should handle unsupported operators", func() {
		expression := booleanexpression.BooleanExpression{
			Operator: "INVALID",
			Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
				booleanexpression.NumericCondition{
					ActualJSONPath: utils.JSONPathQuery("value"),
					Operator:       booleanexpression.OperatorEqual,
					ExpectedValue:  42.0,
				},
			},
		}

		result, err := booleanexpression.EvaluateBooleanExpression(expression, nil)
		assert.Error(s.T(), err)
		assert.False(s.T(), result)
		assert.IsType(s.T(), &booleanexpression.EvaluationError{}, err)
	})
}

func (s *BooleanExpressionTestSuite) TestNumericConditions() {
	s.Run("should evaluate all numeric operators correctly", func() {
		data := map[string]interface{}{
			"value": 50,
		}

		// Test each numeric operator
		testCases := []struct {
			operator      booleanexpression.Operator
			expectedValue float64
			expected      bool
		}{
			{booleanexpression.OperatorEqual, 50.0, true},
			{booleanexpression.OperatorNotEqual, 51.0, true},
			{booleanexpression.OperatorGreaterThan, 49.0, true},
			{booleanexpression.OperatorGreaterEqual, 50.0, true},
			{booleanexpression.OperatorLessThan, 51.0, true},
			{booleanexpression.OperatorLessEqual, 50.0, true},
		}

		for _, tc := range testCases {
			expression := booleanexpression.BooleanExpression{
				Operator: booleanexpression.OperatorAND,
				Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
					booleanexpression.NumericCondition{
						ActualJSONPath: utils.JSONPathQuery("value"),
						Operator:       tc.operator,
						ExpectedValue:  tc.expectedValue,
					},
				},
			}

			result, err := booleanexpression.EvaluateBooleanExpression(expression, data)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tc.expected, result, "Operator: %s", tc.operator)
		}
	})
}

func (s *BooleanExpressionTestSuite) TestStringConditions() {
	s.Run("should evaluate string conditions correctly", func() {
		data := map[string]interface{}{
			"text": "hello",
		}

		testCases := []struct {
			operator      booleanexpression.Operator
			expectedValue string
			expected      bool
		}{
			{booleanexpression.OperatorEqual, "hello", true},
			{booleanexpression.OperatorEqual, "world", false},
			{booleanexpression.OperatorNotEqual, "world", true},
			{booleanexpression.OperatorNotEqual, "hello", false},
		}

		for _, tc := range testCases {
			expression := booleanexpression.BooleanExpression{
				Operator: booleanexpression.OperatorAND,
				Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
					booleanexpression.StringCondition{
						ActualJSONPath: utils.JSONPathQuery("text"),
						Operator:       tc.operator,
						ExpectedValue:  tc.expectedValue,
					},
				},
			}

			result, err := booleanexpression.EvaluateBooleanExpression(expression, data)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tc.expected, result, "Operator: %s, Expected: %s", tc.operator, tc.expectedValue)
		}
	})
}

func (s *BooleanExpressionTestSuite) TestArrayConditions() {
	s.Run("should evaluate all array operators correctly", func() {
		data := map[string]interface{}{
			"tags": []interface{}{"tag1", "tag2", "tag3"},
		}

		testCases := []struct {
			operator      booleanexpression.ArryConditionOperator
			expectedValue []interface{}
			expected      bool
		}{
			{booleanexpression.OperatorContains, []interface{}{"tag1"}, true},
			{booleanexpression.OperatorContains, []interface{}{"tag4"}, false},
			{booleanexpression.OperatorContainsAll, []interface{}{"tag1", "tag2"}, true},
			{booleanexpression.OperatorContainsAll, []interface{}{"tag1", "tag4"}, false},
			{booleanexpression.OperatorContainsAny, []interface{}{"tag1", "tag4"}, true},
			{booleanexpression.OperatorContainsAny, []interface{}{"tag4", "tag5"}, false},
			{booleanexpression.OperatorNotContains, []interface{}{"tag4"}, true},
			{booleanexpression.OperatorNotContains, []interface{}{"tag1"}, false},
		}

		for _, tc := range testCases {
			expression := booleanexpression.BooleanExpression{
				Operator: booleanexpression.OperatorAND,
				Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
					booleanexpression.ArrayCondition{
						ActualJSONPath: utils.JSONPathQuery("tags"),
						Operator:       tc.operator,
						ExpectedValue:  tc.expectedValue,
					},
				},
			}

			result, err := booleanexpression.EvaluateBooleanExpression(expression, data)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tc.expected, result, "Operator: %s", tc.operator)
		}
	})
}

func (s *BooleanExpressionTestSuite) TestBooleanConditions() {
	s.Run("should evaluate boolean conditions correctly", func() {
		data := map[string]interface{}{
			"flag": true,
		}

		testCases := []struct {
			operator      booleanexpression.Operator
			expectedValue bool
			expected      bool
		}{
			{booleanexpression.OperatorEqual, true, true},
			{booleanexpression.OperatorEqual, false, false},
			{booleanexpression.OperatorNotEqual, false, true},
			{booleanexpression.OperatorNotEqual, true, false},
		}

		for _, tc := range testCases {
			expression := booleanexpression.BooleanExpression{
				Operator: booleanexpression.OperatorAND,
				Conditions: []booleanexpression.IsBooleanExpressionOrCondition{
					booleanexpression.BooleanCondition{
						ActualJSONPath: utils.JSONPathQuery("flag"),
						Operator:       tc.operator,
						ExpectedValue:  tc.expectedValue,
					},
				},
			}

			result, err := booleanexpression.EvaluateBooleanExpression(expression, data)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tc.expected, result, "Operator: %s, Expected: %v", tc.operator, tc.expectedValue)
		}
	})
}
