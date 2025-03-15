package chatsessions

import (
	"context"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/stretchr/testify/mock"
)

type MockChatSessionsRepo struct {
	mock.Mock
}

func (m *MockChatSessionsRepo) Get(ctx context.Context, id chatsession.ChatSessionID) (*chatsession.ChatSession, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*chatsession.ChatSession), args.Error(1)
}

func (m *MockChatSessionsRepo) All(ctx context.Context) ([]*chatsession.ChatSession, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*chatsession.ChatSession), args.Error(1)
}

func (m *MockChatSessionsRepo) Create(ctx context.Context, chatSession *chatsession.ChatSession) error {
	args := m.Called(ctx, chatSession)
	return args.Error(0)
}

func (m *MockChatSessionsRepo) Update(ctx context.Context, chatSession *chatsession.ChatSession) error {
	args := m.Called(ctx, chatSession)
	return args.Error(0)
}

func (m *MockChatSessionsRepo) Delete(ctx context.Context, id chatsession.ChatSessionID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func NewMockChatSessionsRepo() *MockChatSessionsRepo {
	mockRepo := new(MockChatSessionsRepo)
	return mockRepo
}
