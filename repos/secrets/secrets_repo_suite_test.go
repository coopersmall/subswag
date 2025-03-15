package secrets_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type SecretsRepoTestSuite struct {
	*tt.IntegrationTest
}

func TestSecretRepoTestSuite(t *testing.T) {
	tt.RunIntegrationTest(t, tt.GetIntegrationSuiteConfig(), func(its *tt.IntegrationTest) *SecretsRepoTestSuite {
		return &SecretsRepoTestSuite{
			IntegrationTest: its,
		}
	})
}
