package utils

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(ctx context.Context, msg string, metadata map[string]any) {
	m.Called(ctx, msg, metadata)
}

func (m *MockLogger) Debug(ctx context.Context, msg string, metadata map[string]any) {
	m.Called(ctx, msg, metadata)
}

func (m *MockLogger) Warn(ctx context.Context, msg string, metadata map[string]any) {
	m.Called(ctx, msg, metadata)
}

func (m *MockLogger) Error(ctx context.Context, msg string, err error, metadata map[string]any) {
	m.Called(ctx, msg, metadata)
}
