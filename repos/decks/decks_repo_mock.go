package deck

import (
	"context"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/stretchr/testify/mock"
)

type MockDecksRepo struct {
	mock.Mock
}

func (m *MockDecksRepo) Get(ctx context.Context, id card.SerializableDeckID) (*card.SerializableDeck, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*card.SerializableDeck), args.Error(1)
}

func (m *MockDecksRepo) All(ctx context.Context) ([]*card.SerializableDeck, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*card.SerializableDeck), args.Error(1)
}

func (m *MockDecksRepo) Create(ctx context.Context, deck *card.SerializableDeck) error {
	args := m.Called(ctx, deck)
	return args.Error(0)
}

func (m *MockDecksRepo) Update(ctx context.Context, deck *card.SerializableDeck) error {
	args := m.Called(ctx, deck)
	return args.Error(0)
}

func (m *MockDecksRepo) Delete(ctx context.Context, id card.SerializableDeckID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
