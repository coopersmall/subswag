package secrets

import (
	"context"

	"github.com/coopersmall/subswag/domain/secret"
	"github.com/stretchr/testify/mock"
)

type MockSecretsRepo struct {
	mock.Mock
}

func (m *MockSecretsRepo) Get(ctx context.Context, secretId secret.SecretID) (*secret.StoredSecret, error) {
	args := m.Called(ctx, secretId)
	return args.Get(0).(*secret.StoredSecret), args.Error(1)
}

func (m *MockSecretsRepo) All(ctx context.Context) ([]*secret.StoredSecret, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*secret.StoredSecret), args.Error(1)
}

func (m *MockSecretsRepo) Create(ctx context.Context, secret *secret.StoredSecret) error {
	args := m.Called(ctx, secret)
	return args.Error(0)
}

func (m *MockSecretsRepo) Update(ctx context.Context, secret *secret.StoredSecret) error {
	args := m.Called(ctx, secret)
	return args.Error(0)
}

func (m *MockSecretsRepo) Delete(ctx context.Context, userId secret.SecretID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func NewMockSecretsRepo() *MockSecretsRepo {
	mockRepo := new(MockSecretsRepo)
	return mockRepo
}
