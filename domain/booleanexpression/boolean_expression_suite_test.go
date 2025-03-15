package booleanexpression_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type BooleanExpressionTestSuite struct {
	suite.Suite
}

func TestBooleanExpressionSuite(t *testing.T) {
	suite.Run(t, new(BooleanExpressionTestSuite))
}
