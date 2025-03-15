package repos

import (
	"github.com/coopersmall/subswag/domain/user"
	"github.com/stretchr/testify/mock"
)

type MockRepos struct {
	mock.Mock
}

func (m *MockRepos) APITokenRepo(userId user.UserID) IAPITokenRepo {
	args := m.Called(userId)
	return args.Get(0).(IAPITokenRepo)
}

func (m *MockRepos) CardsRepo(userId user.UserID) ICardsRepo {
	args := m.Called(userId)
	return args.Get(0).(ICardsRepo)
}

func (m *MockRepos) ChatSessionsRepo(userId user.UserID) IChatSessionsRepo {
	args := m.Called(userId)
	return args.Get(0).(IChatSessionsRepo)
}

func (m *MockRepos) ChatSessionItemsRepo(userId user.UserID) IChatSessionItemsRepo {
	args := m.Called(userId)
	return args.Get(0).(IChatSessionItemsRepo)
}

func (m *MockRepos) DecksRepo(userId user.UserID) IDecksRepo {
	args := m.Called(userId)
	return args.Get(0).(IDecksRepo)
}

func (m *MockRepos) SecretsRepo(userId user.UserID) ISecretsRepo {
	args := m.Called(userId)
	return args.Get(0).(ISecretsRepo)
}

func (m *MockRepos) RateLimitsRepo(userId user.UserID) IRateLimitsRepo {
	args := m.Called(userId)
	return args.Get(0).(IRateLimitsRepo)
}

func (m *MockRepos) UsersRepo(userId user.UserID) IUsersRepo {
	args := m.Called(userId)
	return args.Get(0).(IUsersRepo)
}
