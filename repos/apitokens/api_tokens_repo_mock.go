package apitokens

import (
	"context"

	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/stretchr/testify/mock"
)

type MockApiTokensRepo struct {
	mock.Mock
}

func (m *MockApiTokensRepo) Get(ctx context.Context, userId user.UserID, tokenId apitoken.APITokenID) (*apitoken.APIToken, error) {
	args := m.Called(ctx, userId, tokenId)
	return args.Get(0).(*apitoken.APIToken), args.Error(1)
}

func (m *MockApiTokensRepo) All(ctx context.Context) ([]*apitoken.APIToken, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*apitoken.APIToken), args.Error(1)
}

func (m *MockApiTokensRepo) Create(ctx context.Context, userId user.UserID, token *apitoken.APIToken) error {
	args := m.Called(ctx, userId, token)
	return args.Error(0)
}

func (m *MockApiTokensRepo) Update(ctx context.Context, userId user.UserID, token *apitoken.APIToken) error {
	args := m.Called(ctx, userId, token)
	return args.Error(0)
}

func (m *MockApiTokensRepo) Delete(ctx context.Context, userId user.UserID, tokenId apitoken.APITokenID) error {
	args := m.Called(ctx, userId, tokenId)
	return args.Error(0)
}
