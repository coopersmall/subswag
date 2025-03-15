package answerassistant_test

import (
	"net/http/httptest"
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type AnswerAssistantServiceTestSuite struct {
	server *httptest.Server
	*tt.IntegrationTest
}

func TestAnswerAssistantServiceSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *AnswerAssistantServiceTestSuite {
		return &AnswerAssistantServiceTestSuite{
			IntegrationTest: its,
		}
	})
}
