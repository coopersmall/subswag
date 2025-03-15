package gateways

import (
	"context"

	"github.com/coopersmall/subswag/domain/chatcompletion"
	"github.com/stretchr/testify/mock"
)

type MockGateways struct {
	mock.Mock
}

func (m *MockGateways) OpenAI() IAIGateway {
	args := m.Called()
	return args.Get(0).(IAIGateway)
}

type MockAIGateway struct {
	mock.Mock
}

func (m *MockAIGateway) ChatCompletion(
	ctx context.Context,
	params chatcompletion.ChatCompletionRequest,
) (chatcompletion.ChatCompletionResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(chatcompletion.ChatCompletionResponse), args.Error(1)
}
