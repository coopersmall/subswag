package userauthentication

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/mock"
)

type MockAuthenticationService struct {
	mock.Mock
}

func (m *MockAuthenticationService) AuthenticateMAC(
	ctx context.Context,
	timestamp time.Time,
	signature string,
) error {
	args := m.Called(ctx, timestamp, signature)
	return args.Error(0)
}

func (m *MockAuthenticationService) AuthenticateToken(
	ctx context.Context,
	token string,
	now time.Time,
	permissions ...domain.Permission,
) (user.UserID, apitoken.APITokenID, error) {
	args := m.Called(ctx, token, now, permissions)
	return args.Get(0).(user.UserID), args.Get(1).(apitoken.APITokenID), args.Error(2)
}

func (m *MockAuthenticationService) AuthenticateSession(
	ctx context.Context,
	session string,
) (user.UserID, error) {
	args := m.Called(ctx, session)
	return args.Get(0).(user.UserID), args.Error(1)
}

// MockRSAService is a mock of the rsaService interface
type MockRSAService struct {
	mock.Mock
}

func (m *MockRSAService) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	args := m.Called(ctx, data)
	return args.Get(0).([]byte), args.Error(1)
}

// MockJWTService is a mock of the jwtService interface
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) ValidateToken(ctx context.Context, token string) (user.UserID, utils.ID, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(user.UserID), args.Get(1).(utils.ID), args.Error(2)
}
