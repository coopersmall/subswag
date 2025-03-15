package chatsessionitems

import (
	"context"

	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/stretchr/testify/mock"
)

type MockChatSessionItemsRepo struct {
	*mock.Mock
}

func (m *MockChatSessionItemsRepo) Get(ctx context.Context, chatSessionItemId chatsession.ChatSessionItemID) (chatsession.ChatSessionItem, error) {
	args := m.Called(ctx, chatSessionItemId)
	return args.Get(0).(chatsession.ChatSessionItem), args.Error(1)
}

func (m *MockChatSessionItemsRepo) GetBySessionId(ctx context.Context, chatSessionId chatsession.ChatSessionID) ([]chatsession.ChatSessionItem, error) {
	args := m.Called(ctx, chatSessionId)
	return args.Get(0).([]chatsession.ChatSessionItem), args.Error(1)
}

func (m *MockChatSessionItemsRepo) All(ctx context.Context) ([]chatsession.ChatSessionItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]chatsession.ChatSessionItem), args.Error(1)
}

func (m *MockChatSessionItemsRepo) Create(ctx context.Context, chatSessionItem chatsession.ChatSessionItem) error {
	args := m.Called(ctx, chatSessionItem)
	return args.Error(0)
}

func (m *MockChatSessionItemsRepo) Update(ctx context.Context, chatSessionItem chatsession.ChatSessionItem) error {
	args := m.Called(ctx, chatSessionItem)
	return args.Error(0)
}

func (m *MockChatSessionItemsRepo) Delete(ctx context.Context, chatSessionItemId chatsession.ChatSessionItemID) error {
	args := m.Called(ctx, chatSessionItemId)
	return args.Error(0)
}

func (m *MockChatSessionItemsRepo) DeleteBySessionId(ctx context.Context, chatSessionId chatsession.ChatSessionID) error {
	args := m.Called(ctx, chatSessionId)
	return args.Error(0)
}
