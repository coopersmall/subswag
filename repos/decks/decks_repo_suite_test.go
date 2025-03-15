package deck_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type DecksRepoTestSuite struct {
	*tt.IntegrationTest
}

func TestDecksRepoTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *DecksRepoTestSuite {
		return &DecksRepoTestSuite{
			IntegrationTest: its,
		}
	})
}
