package apitokens_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type APITokensRepoTestSuite struct {
	*tt.IntegrationTest
}

func TestAPITokensRepoTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *APITokensRepoTestSuite {
		return &APITokensRepoTestSuite{
			IntegrationTest: its,
		}
	})
}
