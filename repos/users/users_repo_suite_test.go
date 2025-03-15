package users_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type UsersRepoTestSuite struct {
	*tt.IntegrationTest
}

func TestUsersRepoTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *UsersRepoTestSuite {
		return &UsersRepoTestSuite{
			IntegrationTest: its,
		}
	})
}
