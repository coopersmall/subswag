package users

import (
	"context"

	"github.com/coopersmall/subswag/domain/user"
	"github.com/stretchr/testify/mock"
)

type MockUsersRepo struct {
	mock.Mock
}

func (m *MockUsersRepo) Get(ctx context.Context, userId user.UserID) (*user.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUsersRepo) All(ctx context.Context) ([]*user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUsersRepo) Create(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUsersRepo) Update(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUsersRepo) Delete(ctx context.Context, userId user.UserID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func NewMockUsersRepo() *MockUsersRepo {
	mockRepo := new(MockUsersRepo)
	return mockRepo
}
