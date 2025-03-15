package users_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type UsersCacheTestSuite struct {
	*tt.IntegrationTest
}

func TestUsersCacheTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *UsersCacheTestSuite {
		return &UsersCacheTestSuite{
			IntegrationTest: its,
		}
	})
}
