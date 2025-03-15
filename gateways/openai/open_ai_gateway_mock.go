package openai

import (
	"context"

	"github.com/coopersmall/subswag/domain/chatcompletion"
	"github.com/stretchr/testify/mock"
)

type MockOpenAIGateway struct {
	mock.Mock
}

func (m *MockOpenAIGateway) ChatCompletion(ctx context.Context, params chatcompletion.ChatCompletionRequest) (chatcompletion.ChatCompletionResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return chatcompletion.ChatCompletionResponse{}, args.Error(1)
	}
	return args.Get(0).(chatcompletion.ChatCompletionResponse), args.Error(1)
}
