package chatsessions_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type ChatSessionsRepoTestSuite struct {
	*tt.IntegrationTest
}

func TestChatSessionsRepo(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *ChatSessionsRepoTestSuite {
		return &ChatSessionsRepoTestSuite{
			IntegrationTest: its,
		}
	})
}
