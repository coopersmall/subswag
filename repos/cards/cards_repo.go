package card

import (
	"context"
	"database/sql"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/card"
	reposdomain "github.com/coopersmall/subswag/repos/domain"
	"github.com/coopersmall/subswag/utils"
)

type CardsRepo struct {
	*reposdomain.SharedRepo[card.SerializableCardID, card.Card, db.Card]
}

func NewCardsRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
) *CardsRepo {
	return &CardsRepo{
		SharedRepo: reposdomain.NewSharedRepo[card.SerializableCardID, card.Card, db.Card](
			"card",
			querier,
			tracer,
			convertRowToCard,
			convertCardToRow,
			func(ctx context.Context, iqro db.ISharedQueriesReadOnly, sci card.SerializableCardID) (db.Card, error) {
				return iqro.GetCard(ctx, int64(sci))
			},
			func(ctx context.Context, iqro db.ISharedQueriesReadOnly) ([]db.Card, error) {
				return iqro.GetAllCards(ctx)
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, c db.Card) (sql.Result, error) {
				return iqrw.CreateCard(ctx, db.CreateCardParams{
					ID:        c.ID,
					CreatedAt: c.CreatedAt,
					Type:      c.Type,
					Data:      c.Data,
				})
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, c db.Card) (sql.Result, error) {
				return iqrw.UpdateCard(ctx, db.UpdateCardParams{
					ID:        c.ID,
					UpdatedAt: c.UpdatedAt,
					Type:      c.Type,
					Data:      c.Data,
				})
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, sci card.SerializableCardID) (sql.Result, error) {
				return iqrw.DeleteCard(ctx, int64(sci))
			},
		),
	}
}

func convertRowToCard(result db.Card) (card.Card, error) {
	cardType := card.SerializableCardType(result.Type)
	metadata := &domain.Metadata{
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt.Time,
	}

	switch cardType {
	case card.SerializableCardTypeFace:
		var c card.SerializableFaceCard
		err := utils.Unmarshal(result.Data, &c)
		if err != nil {
			return nil, err
		}
		c.Metadata = metadata
		return &c, nil
	case card.SerializableCardTypeNumber:
		var c card.SerializableNumberCard
		err := utils.Unmarshal(result.Data, &c)
		if err != nil {
			return nil, err
		}
		c.Metadata = metadata
		return &c, nil
	default:
		return nil, utils.NewUnableToHandleError("unknown card type")
	}
}

func convertCardToRow(card card.Card) (db.Card, error) {
	metadata := card.GetMetadata()
	data, err := utils.Marshal(card)
	if err != nil {
		return db.Card{}, err
	}

	return db.Card{
		ID:        int64(card.GetID()),
		CreatedAt: metadata.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  metadata.UpdatedAt,
			Valid: !metadata.UpdatedAt.IsZero(),
		},
		Type: string(card.GetType()),
		Data: data,
	}, nil
}
