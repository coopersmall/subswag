//ts:ignore
package chatsession

import (
	"github.com/coopersmall/subswag/utils"
)

type ChatSessionID utils.ID

func NewChatSessionID() ChatSessionID {
	return ChatSessionID(utils.NewID())
}
