package chatsession

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
)

type ChatSession struct {
	ID              ChatSessionID `validate:"required,gt=0"`
	ChatSessionData `json:",inline" validate:"required" tstype:",extends"`
	Metadata        *domain.Metadata `validate:"required" tstype:"Metadata"`
}

type ChatSessionData struct {
	ChatSessionItemIDs []ChatSessionItemID `validate:"optional"`
	UserIDs            []user.UserID       `validate:"optional"`
}

func NewChatSession(
	userIds []user.UserID,
	itemIds []ChatSessionItemID,
) *ChatSession {
	return &ChatSession{
		ID: NewChatSessionID(),
		ChatSessionData: ChatSessionData{
			ChatSessionItemIDs: itemIds,
			UserIDs:            userIds,
		},
		Metadata: domain.NewMetadata(),
	}
}
