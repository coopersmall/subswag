package booleanexpression

// BooleanOperator represents the type of boolean operation
type BooleanOperator string

const (
	OperatorAND BooleanOperator = "AND"
	OperatorOR  BooleanOperator = "OR"
)

// Returns either a boolean expression or a condition operator
type IsBooleanExpressionOrCondition interface {
	IsBooleanExpressionOrCondition()
}

// BooleanExpression represents a boolean expression that combines multiple conditions
type BooleanExpression struct {
	Operator   BooleanOperator                  `json:"operator"`
	Conditions []IsBooleanExpressionOrCondition `json:"conditions"` // Can contain Condition or BooleanExpression
}

func (b BooleanExpression) IsBooleanExpressionOrCondition() {}

func NewBooleanExpression(operator BooleanOperator, conditions ...IsBooleanExpressionOrCondition) BooleanExpression {
	return BooleanExpression{
		Operator:   operator,
		Conditions: conditions,
	}
}
