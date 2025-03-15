package chatsessionitems_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type ChatSessionItemsRepoTestSuite struct {
	*tt.IntegrationTest
}

func TestChatSessionItemsRepoTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *ChatSessionItemsRepoTestSuite {
		return &ChatSessionItemsRepoTestSuite{
			IntegrationTest: its,
		}
	})
}
