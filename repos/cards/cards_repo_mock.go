package card

import (
	"context"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/stretchr/testify/mock"
)

type MockCardsRepo struct {
	*mock.Mock
}

func (m *MockCardsRepo) Get(ctx context.Context, cardId card.SerializableCardID) (card.Card, error) {
	args := m.Called(ctx, cardId)
	return args.Get(0).(card.Card), args.Error(1)
}

func (m *MockCardsRepo) All(ctx context.Context) ([]card.Card, error) {
	args := m.Called(ctx)
	return args.Get(0).([]card.Card), args.Error(1)
}

func (m *MockCardsRepo) Create(ctx context.Context, card card.Card) error {
	args := m.Called(ctx, card)
	return args.Error(0)
}

func (m *MockCardsRepo) Update(ctx context.Context, card card.Card) error {
	args := m.Called(ctx, card)
	return args.Error(0)
}

func (m *MockCardsRepo) Delete(ctx context.Context, cardId card.SerializableCardID) error {
	args := m.Called(ctx, cardId)
	return args.Error(0)
}
