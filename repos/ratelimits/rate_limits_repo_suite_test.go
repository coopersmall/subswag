package ratelimits_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type RateLimitsRepoTestSuite struct {
	*tt.IntegrationTest
}

func TestRateLimitsRepoTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *RateLimitsRepoTestSuite {
		return &RateLimitsRepoTestSuite{
			IntegrationTest: its,
		}
	})
}
