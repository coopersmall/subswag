package server

import (
	"context"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/services"
	"github.com/stretchr/testify/mock"
)

type MockRequest struct {
	mock.Mock
}

func (m *MockRequest) Ctx() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *MockRequest) CorrelationID() domain.CorrelationID {
	args := m.Called()
	return args.Get(0).(domain.CorrelationID)
}

func (m *MockRequest) UserID() user.UserID {
	args := m.Called()
	return args.Get(0).(user.UserID)
}

func (m *MockRequest) Body() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockRequest) Param(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *MockRequest) SearchParam(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *MockRequest) GetServices() services.IServices {
	args := m.Called()
	return args.Get(0).(services.IServices)
}

func (m *MockRequest) GetRepos() repos.IRepos {
	args := m.Called()
	return args.Get(0).(repos.IRepos)
}

func (m *MockRequest) GetGateways() gateways.IGateways {
	args := m.Called()
	return args.Get(0).(gateways.IGateways)
}
