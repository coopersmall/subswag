package openai_test

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OpenAIGatewayTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func TestOpenAIGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(OpenAIGatewayTestSuite))
}
