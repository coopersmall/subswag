package card

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/utils"
)

type UserSerializableCardID utils.ID

func NewUserSerializableCardID() UserSerializableCardID {
	return UserSerializableCardID(utils.NewID())
}

type UserSerializableCard struct {
	ID                       UserSerializableCardID `json:"id" validate:"required,gt=0" tstype:"string"`
	UserSerializableCardData `json:",inline" validate:"required" tstype:",extends"`
	*domain.Metadata         `json:"metadata" validate:"required" tstype:"Metadata"`
}

type UserSerializableCardData struct {
	UserID user.UserID        `json:"user_id" validate:"required" tstype:"string"`
	CardID SerializableCardID `json:"card_id" validate:"required" tstype:"string"`
}
