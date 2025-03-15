package card_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type CardsRepoTestSuite struct {
	*tt.IntegrationTest
}

func TestCardsRepoTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *CardsRepoTestSuite {
		return &CardsRepoTestSuite{
			IntegrationTest: its,
		}
	})
}
