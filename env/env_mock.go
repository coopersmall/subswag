package env

import (
	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/cache"
	"github.com/coopersmall/subswag/clients"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/services"
	"github.com/coopersmall/subswag/streams/publishers"
	"github.com/coopersmall/subswag/streams/subscribers"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/mock"
)

type MockEnv struct {
	mock.Mock
}

func (m *MockEnv) GetLogger(name string) utils.ILogger {
	args := m.Called(name)
	return args.Get(0).(utils.ILogger)
}

func (m *MockEnv) GetTracer(service string) apm.ITracer {
	args := m.Called(service)
	return args.Get(0).(apm.ITracer)
}

func (m *MockEnv) GetRepos() (repos.IRepos, func()) {
	args := m.Called()
	return args.Get(0).(repos.IRepos), args.Get(1).(func())
}

func (m *MockEnv) GetPublishers() (publishers.IPublishers, func()) {
	args := m.Called()
	return args.Get(0).(publishers.IPublishers), args.Get(1).(func())
}

func (m *MockEnv) GetSubscribers() (subscribers.ISubscribers, func()) {
	args := m.Called()
	return args.Get(0).(subscribers.ISubscribers), args.Get(1).(func())
}

func (m *MockEnv) GetClients() (clients.IClients, func()) {
	args := m.Called()
	return args.Get(0).(clients.IClients), args.Get(1).(func())
}

func (m *MockEnv) GetGateways() (gateways.IGateways, func()) {
	args := m.Called()
	return args.Get(0).(gateways.IGateways), args.Get(1).(func())
}

func (m *MockEnv) GetCaches() (cache.ICache, func()) {
	args := m.Called()
	return args.Get(0).(cache.ICache), args.Get(1).(func())
}

func (m *MockEnv) GetServices() (services.IServices, func()) {
	args := m.Called()
	return args.Get(0).(services.IServices), args.Get(1).(func())
}

func (m *MockEnv) GetUsersService() (services.IUsersService, func()) {
	args := m.Called()
	return args.Get(0).(services.IUsersService), args.Get(1).(func())
}

func (m *MockEnv) GetAuthenticationService() (services.IAuthenticationService, func()) {
	args := m.Called()
	return args.Get(0).(services.IAuthenticationService), args.Get(1).(func())
}

func (m *MockEnv) GetQuerier() db.IQuerier {
	args := m.Called()
	return args.Get(0).(db.IQuerier)
}

func (m *MockEnv) GetDB() db.IDBManager {
	args := m.Called()
	return args.Get(0).(db.IDBManager)
}

func (m *MockEnv) Shutdown() error {
	args := m.Called()
	return args.Error(0)
}
