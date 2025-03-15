package chatsessionitems_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type ChatSessionItemsServiceTestSuite struct {
	*tt.IntegrationTest
}

func TestChatSessionItemsServiceSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *ChatSessionItemsServiceTestSuite {
		return &ChatSessionItemsServiceTestSuite{
			IntegrationTest: its,
		}
	})
}
