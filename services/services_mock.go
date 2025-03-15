package services

import (
	"github.com/coopersmall/subswag/domain/user"
	"github.com/stretchr/testify/mock"
)

type MockServices struct {
	mock.Mock
}

func (m *MockServices) APITokenService(userId user.UserID) IAPITokenService {
	args := m.Called(userId)
	return args.Get(0).(IAPITokenService)
}

func (m *MockServices) AnswerAssistantService(userId user.UserID) IAnswerAssistantService {
	args := m.Called(userId)
	return args.Get(0).(IAnswerAssistantService)
}

func (m *MockServices) ChatSessionsService(userId user.UserID) IChatSessionsService {
	args := m.Called(userId)
	return args.Get(0).(IChatSessionsService)
}

func (m *MockServices) ChatSessionItemsService(userId user.UserID) IChatSessionItemsService {
	args := m.Called(userId)
	return args.Get(0).(IChatSessionItemsService)
}

func (m *MockServices) DecksService(userId user.UserID) IDecksService {
	args := m.Called(userId)
	return args.Get(0).(IDecksService)
}

func (m *MockServices) SecretsService(userId user.UserID) ISecretsService {
	args := m.Called(userId)
	return args.Get(0).(ISecretsService)
}

func (m *MockServices) AuthenticationService() IAuthenticationService {
	args := m.Called()
	return args.Get(0).(IAuthenticationService)
}

func (m *MockServices) JWTService() IJWTService {
	args := m.Called()
	return args.Get(0).(IJWTService)
}

func (m *MockServices) RSAService() IRSAService {
	args := m.Called()
	return args.Get(0).(IRSAService)
}

func (m *MockServices) UsersService() IUsersService {
	args := m.Called()
	return args.Get(0).(IUsersService)
}
