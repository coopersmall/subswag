package user

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/utils"
)

type User struct {
	ID       UserID `json:"id" validate:"required,gt=0"`
	UserData `json:",inline" validate:"required" tstype:",extends"`
	Metadata *domain.Metadata `json:"metadata" validate:"required" tstype:"Metadata"`
}

type UserData struct {
	/*
	   @format email
	*/
	Email string `json:"email,omitempty" validate:"omitempty,email" tstype:",optional"`
	/*
	   @minLength 1
	*/
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=1" tstype:",optional"`
	/*
	   @minLength 1
	*/
	LastName string `json:"last_name,omitempty" validate:"omitempty,min=1" tstype:",optional"`
}

func NewUser() *User {
	return &User{
		ID:       UserID(utils.NewID()),
		UserData: UserData{},
		Metadata: domain.NewMetadata(),
	}
}
