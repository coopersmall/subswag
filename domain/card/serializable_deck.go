package card

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
)

type SerializableDeck struct {
	ID                   SerializableDeckID `json:"id" validate:"required,gt=0"`
	SerializableDeckData `json:",inline" validate:"required" tstype:",extends"`
	Metadata             *domain.Metadata `json:"metadata" validate:"required" tstype:"Metadata"`
}

type SerializableDeckData struct {
	user.UserID `json:"user_id" validate:"required" tstype:"string"`
	CardIDs     []SerializableCardID `json:"cards" validate:"required" tstype:"Array<SerializableCardID>"`
	Name        string               `json:"name" validate:"required" tstype:"string"`
	Favorited   bool                 `json:"favorited" validate:"required" tstype:"boolean"`
	GamesPlayed int                  `json:"games_played" validate:"required" tstype:"number"`
	GamesWon    int                  `json:"games_won" validate:"required" tstype:"number"`
}

func NewDeck(
	id SerializableDeckID,
	data SerializableDeckData,
) *SerializableDeck {
	if data.Name == "" {
		data.Name = "New Deck"
	}
	return &SerializableDeck{
		ID:                   id,
		SerializableDeckData: data,
		Metadata:             domain.NewMetadata(),
	}
}
