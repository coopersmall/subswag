package apitoken_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type APITokenCacheTestSuite struct {
	*tt.IntegrationTest
}

func TestAPITokenCacheTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *APITokenCacheTestSuite {
		return &APITokenCacheTestSuite{
			IntegrationTest: its,
		}
	})
}
