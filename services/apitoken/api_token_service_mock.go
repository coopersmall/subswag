package apitoken

import (
	"context"

	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/mock"
)

type MockAPITokenService struct {
	mock.Mock
}

func (m *MockAPITokenService) CreateToken(
	ctx context.Context,
	data *apitoken.APITokenData,
) (*apitoken.APITokenWithSecret, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apitoken.APITokenWithSecret), args.Error(1)
}

func (m *MockAPITokenService) CreateTokenWithId(
	ctx context.Context,
	id utils.ID,
	data *apitoken.APITokenData,
) (*apitoken.APITokenWithSecret, error) {
	args := m.Called(ctx, id, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apitoken.APITokenWithSecret), args.Error(1)
}

func (m *MockAPITokenService) GetToken(
	ctx context.Context,
	tokenId apitoken.APITokenID,
) (*apitoken.APIToken, error) {
	args := m.Called(ctx, tokenId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apitoken.APIToken), args.Error(1)
}

func (m *MockAPITokenService) UpdateToken(
	ctx context.Context,
	token *apitoken.APIToken,
) (*apitoken.APIToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apitoken.APIToken), args.Error(1)
}

func (m *MockAPITokenService) DeleteToken(
	ctx context.Context,
	tokenID apitoken.APITokenID,
) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}
