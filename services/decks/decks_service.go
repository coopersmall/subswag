package decks

import (
	"context"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
)

type DecksService struct {
	logger    utils.ILogger
	decksRepo repos.IDecksRepo
}

func NewDecksService(
	logger utils.ILogger,
	decksRepo repos.IDecksRepo,
) *DecksService {
	return &DecksService{
		logger:    logger,
		decksRepo: decksRepo,
	}
}

func (s *DecksService) GetDeck(ctx context.Context, deckId card.SerializableDeckID) (*card.SerializableDeck, error) {
	return s.decksRepo.Get(ctx, deckId)
}

func (s *DecksService) GetAllDecks(ctx context.Context) ([]*card.SerializableDeck, error) {
	return s.decksRepo.All(ctx)
}

func (s *DecksService) CreateDeck(
	ctx context.Context,
	data card.SerializableDeckData,
) error {
	deck := card.NewDeck(card.NewSerializableDeckID(), data)
	if err := domain.Validate(deck); err != nil {
		return err
	}
	return s.decksRepo.Create(ctx, deck)
}

func (s *DecksService) UpdateDeck(
	ctx context.Context,
	deck *card.SerializableDeck,
) error {
	if err := domain.Validate(deck); err != nil {
		return err
	}
	return s.decksRepo.Update(ctx, deck)
}

func (s *DecksService) DeleteDeck(
	ctx context.Context,
	deckId card.SerializableDeckID,
) error {
	return s.decksRepo.Delete(ctx, deckId)
}
