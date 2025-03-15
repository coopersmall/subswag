package openai

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tmc/langchaingo/llms"
)

type MockOpenAIClient struct {
	mock.Mock
}

func (m *MockOpenAIClient) GenerateContent(ctx context.Context, messages []llms.MessageContent, opts ...llms.CallOption) (*llms.ContentResponse, error) {
	args := m.Called(ctx, messages, opts)
	return args.Get(0).(*llms.ContentResponse), args.Error(1)
}
