package deck

import (
	"context"
	"database/sql"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/coopersmall/subswag/domain/user"
	reposdomain "github.com/coopersmall/subswag/repos/domain"
	"github.com/coopersmall/subswag/utils"
)

type DecksRepo struct {
	*reposdomain.StandardRepo[card.SerializableDeckID, *card.SerializableDeck, db.Deck]
}

func NewDecksRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
	userId user.UserID,
) *DecksRepo {
	return &DecksRepo{
		StandardRepo: reposdomain.NewStandardRepo[card.SerializableDeckID, *card.SerializableDeck, db.Deck](
			"deck",
			querier,
			tracer,
			userId,
			convertRowToDeck,
			convertDeckToRow,
			isEmptyDeck,
			func(ctx context.Context, iqro db.IStandardQueriesReadOnly, sdi card.SerializableDeckID) (db.Deck, error) {
				return iqro.GetDeck(ctx, db.GetDeckParams{
					ID:     int64(sdi),
					UserID: int64(userId),
				})
			},
			func(ctx context.Context, iqro db.IStandardQueriesReadOnly) ([]db.Deck, error) {
				return iqro.GetAllDecks(ctx)
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, d db.Deck) (sql.Result, error) {
				return iqrw.CreateDeck(ctx, db.CreateDeckParams{
					ID:        d.ID,
					UserID:    int64(userId),
					CreatedAt: d.CreatedAt,
					Data:      d.Data,
				})
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, d db.Deck) (sql.Result, error) {
				return iqrw.UpdateDeck(ctx, db.UpdateDeckParams{
					ID:        d.ID,
					UserID:    int64(userId),
					UpdatedAt: d.UpdatedAt,
					Data:      d.Data,
				})
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, sdi card.SerializableDeckID) (sql.Result, error) {
				return iqrw.DeleteDeck(ctx, db.DeleteDeckParams{
					ID:     int64(sdi),
					UserID: int64(userId),
				})
			},
		),
	}
}

func isEmptyDeck(deck *card.SerializableDeck) bool {
	return deck == nil || deck.ID == 0
}

func convertRowToDeck(result db.Deck) (*card.SerializableDeck, error) {
	var deck card.SerializableDeck
	err := utils.Unmarshal(result.Data, &deck)
	return &deck, err
}

func convertDeckToRow(
	deck *card.SerializableDeck,
) (db.Deck, error) {
	data, err := utils.Marshal(deck)
	return db.Deck{
		ID:        int64(deck.ID),
		UserID:    int64(deck.UserID),
		CreatedAt: deck.Metadata.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  deck.Metadata.UpdatedAt,
			Valid: deck.Metadata.UpdatedAt != time.Time{},
		},
		Data: data,
	}, err
}
