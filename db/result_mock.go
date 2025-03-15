package db

import (
	"github.com/stretchr/testify/mock"
)

type MockResult struct {
	*mock.Mock
}

func NewMockResult() *MockResult {
	return &MockResult{Mock: new(mock.Mock)}
}

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
