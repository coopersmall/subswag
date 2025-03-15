package apitoken

import (
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/utils"
)

type APITokenData struct {
	UserId      user.UserID         `json:"userId" validate:"required" tstype:"UserID"`
	Expiry      time.Time           `json:"expiry" validate:"required,nonZeroTime"`
	Permissions []domain.Permission `json:"permissions" validate:"required,min=1,dive,permission" tstype:"Array<Permission>"`
}

type APIToken struct {
	ID           APITokenID `validate:"required,gt=0"`
	APITokenData `json:",inline" validate:"required" tstype:",extends"`
	Metadata     *domain.Metadata `validate:"required" tstype:"Metadata"`
}

type APITokenWithSecret struct {
	APIToken `json:",inline" validate:"required" tstype:",extends"`
	Secret   string `json:"secret" validate:"required"`
}

func NewAPIToken(data APITokenData) *APIToken {
	return &APIToken{
		ID:           APITokenID(utils.NewID()),
		APITokenData: data,
		Metadata:     domain.NewMetadata(),
	}
}

func (a *APIToken) HasPermission(permission domain.Permission) bool {
	for _, p := range a.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}
