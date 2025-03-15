package chatsessions

import (
	"fmt"
	"time"

	"github.com/coopersmall/subswag/apm"
	cachedomain "github.com/coopersmall/subswag/cache/domain"
	chatsessiondomain "github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/utils"
)

const (
	chatSessionCacheKeyPrefix = "chat-session"
)

type ChatSessionsCache struct {
	*cachedomain.StandardCache[chatsessiondomain.ChatSessionID, *chatsessiondomain.ChatSession]
}

func NewChatSessionsCache(
	logger utils.ILogger,
	tracer apm.ITracer,
	gateway gateways.ICacheGateway,
) *ChatSessionsCache {
	return &ChatSessionsCache{
		StandardCache: cachedomain.NewStandardCache[chatsessiondomain.ChatSessionID, *chatsessiondomain.ChatSession](
			"chat-sessions-cache",
			logger,
			tracer,
			gateway,
			time.Duration(0),
			func(id chatsessiondomain.ChatSessionID) string {
				return fmt.Sprintf("%s:%d", chatSessionCacheKeyPrefix, id)
			},
		),
	}
}
