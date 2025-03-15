package user

import (
	"context"

	"github.com/coopersmall/subswag/domain/user"
	"github.com/stretchr/testify/mock"
)

type MockUsersService struct {
	mock.Mock
}

func (m *MockUsersService) CreateUser(ctx context.Context, data user.UserData) (*user.User, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUsersService) CreateUserWithId(ctx context.Context, id user.UserID, data user.UserData) (*user.User, error) {
	args := m.Called(ctx, id, data)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUsersService) UpdateUser(ctx context.Context, existing *user.User) (*user.User, error) {
	args := m.Called(ctx, existing)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUsersService) GetUser(ctx context.Context, userId user.UserID) (*user.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUsersService) DeleteUser(ctx context.Context, userId user.UserID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
