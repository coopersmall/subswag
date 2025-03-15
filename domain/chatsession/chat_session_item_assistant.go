package chatsession

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/utils"
)

type AssistantChatSessionItemData struct {
	Content string `json:"content" validate:"required,min=1"`
}

func (AssistantChatSessionItemData) IsChatSessionItemDataUnion() {}

type AssistantChatSessionItem struct {
	ChatSessionItemBase          `json:",inline" validate:"required" tstype:",extends"`
	Type                         ChatSessionItemType `json:"type" validate:"required,eq=assistant" tstype:"'assistant'"`
	AssistantChatSessionItemData `json:",inline" validate:"required" tstype:",extends"`
}

func (t *AssistantChatSessionItem) GetType() ChatSessionItemType {
	return t.Type
}

func (t *AssistantChatSessionItem) GetData() ChatSessionItemDataUnion {
	return t.AssistantChatSessionItemData
}

func NewAssistantChatSessionItem(
	sessionID ChatSessionID,
	content string,
) *AssistantChatSessionItem {
	return &AssistantChatSessionItem{
		ChatSessionItemBase: ChatSessionItemBase{
			ID:        ChatSessionItemID(utils.NewID()),
			SessionID: sessionID,
			Metadata:  domain.NewMetadata(),
		},
		Type: ChatSessionItemTypeAssistant,
		AssistantChatSessionItemData: AssistantChatSessionItemData{
			Content: content,
		},
	}
}
