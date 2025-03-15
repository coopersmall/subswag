package chatsession

import "github.com/coopersmall/subswag/domain"

type ChatSessionItemBase struct {
	ID ChatSessionItemID `json:"id" validate:"required,gt=0"`
	//tygo:emit type: ChatSessionItemType
	SessionID ChatSessionID    `json:"sessionId" validate:"required,gt=0"`
	Metadata  *domain.Metadata `json:"metadata" validate:"required" tstype:"Metadata"`
}

func (b ChatSessionItemBase) GetID() ChatSessionItemID {
	return b.ID
}

func (b ChatSessionItemBase) GetSessionID() ChatSessionID {
	return b.SessionID
}

func (b ChatSessionItemBase) GetMetadata() *domain.Metadata {
	return b.Metadata
}
