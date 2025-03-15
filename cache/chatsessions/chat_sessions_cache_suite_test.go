package chatsessions_test

import (
	"testing"

	tt "github.com/coopersmall/subswag/testing"
)

type ChatSessionsCacheTestSuite struct {
	*tt.IntegrationTest
}

func TestChatSessionsCacheTestSuite(t *testing.T) {
	config := tt.GetIntegrationSuiteConfig()
	tt.RunIntegrationTest(t, config, func(its *tt.IntegrationTest) *ChatSessionsCacheTestSuite {
		return &ChatSessionsCacheTestSuite{
			IntegrationTest: its,
		}
	})
}
