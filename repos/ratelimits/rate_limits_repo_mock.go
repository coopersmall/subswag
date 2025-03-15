package ratelimits

import (
	"context"

	"github.com/coopersmall/subswag/domain/ratelimit"
	"github.com/stretchr/testify/mock"
)

type MockRateLimitsRepo struct {
	mock.Mock
}

func (m *MockRateLimitsRepo) Get(ctx context.Context, rateLimitId ratelimit.RateLimitID) (*ratelimit.RateLimit, error) {
	args := m.Called(ctx, rateLimitId)
	return args.Get(0).(*ratelimit.RateLimit), args.Error(1)
}

func (m *MockRateLimitsRepo) All(ctx context.Context) ([]*ratelimit.RateLimit, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*ratelimit.RateLimit), args.Error(1)
}

func (m *MockRateLimitsRepo) Create(ctx context.Context, rateLimit *ratelimit.RateLimit) error {
	args := m.Called(ctx, rateLimit)
	return args.Error(0)
}

func (m *MockRateLimitsRepo) Update(ctx context.Context, rateLimit *ratelimit.RateLimit) error {
	args := m.Called(ctx, rateLimit)
	return args.Error(0)
}

func (m *MockRateLimitsRepo) Delete(ctx context.Context, rateLimitId ratelimit.RateLimitID) error {
	args := m.Called(ctx, rateLimitId)
	return args.Error(0)
}
