package chatsessionitems

import (
	"context"

	"github.com/coopersmall/subswag/domain/chatcompletion"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/stretchr/testify/mock"
)

type MockChatSessionItemsService struct {
	mock.Mock
}

func (m *MockChatSessionItemsService) ConvertChatSessionItemsToChatCompletionMessages(
	ctx context.Context,
	items []chatsession.ChatSessionItem,
) []chatcompletion.Message {
	args := m.Called(ctx, items)
	return args.Get(0).([]chatcompletion.Message)
}
