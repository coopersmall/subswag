package chatsession

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/utils"
)

type UserChatSessionItemData struct {
	/*
	   @minLength 1
	*/
	Content string      `json:"content" validate:"required,min=1"`
	UserID  user.UserID `json:"userId" validate:"required" tstype:"UserID"`
}

func (UserChatSessionItemData) IsChatSessionItemDataUnion() {}

type UserChatSessionItem struct {
	ChatSessionItemBase     `json:",inline" validate:"required" tstype:",extends"`
	Type                    ChatSessionItemType `json:"type" validate:"required,eq=user" tstype:"'user'"`
	UserChatSessionItemData `json:",inline" validate:"required" tstype:",extends"`
}

func (u UserChatSessionItem) GetType() ChatSessionItemType {
	return u.Type
}

func (u UserChatSessionItem) GetData() ChatSessionItemDataUnion {
	return u.UserChatSessionItemData
}

func NewUserChatSessionItem(
	sessionID ChatSessionID,
	content string,
) *UserChatSessionItem {
	return &UserChatSessionItem{
		ChatSessionItemBase: ChatSessionItemBase{
			ID:        ChatSessionItemID(utils.NewID()),
			SessionID: sessionID,
			Metadata:  domain.NewMetadata(),
		},
		Type: ChatSessionItemTypeUser,
		UserChatSessionItemData: UserChatSessionItemData{
			Content: content,
		},
	}
}
