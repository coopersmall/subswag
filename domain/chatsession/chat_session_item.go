//ts:ignore
package chatsession

import (
	"github.com/coopersmall/subswag/domain"
)

type ChatSessionItemType string

const (
	ChatSessionItemTypeUser      ChatSessionItemType = "user"
	ChatSessionItemTypeAssistant ChatSessionItemType = "assistant"
)

type ChatSessionItem interface {
	GetID() ChatSessionItemID
	GetSessionID() ChatSessionID
	GetData() ChatSessionItemDataUnion
	GetType() ChatSessionItemType
	GetMetadata() *domain.Metadata
}

type ChatSessionItemDataUnion interface {
	IsChatSessionItemDataUnion()
}

func NewChatSessionItem(
	sessionID ChatSessionID,
	data ChatSessionItemDataUnion,
) ChatSessionItem {
	switch d := data.(type) {
	case UserChatSessionItemData:
		return NewUserChatSessionItem(
			sessionID,
			d.Content,
		)
	case AssistantChatSessionItemData:
		return NewAssistantChatSessionItem(
			sessionID,
			d.Content,
		)
	default:
		return nil
	}
}
