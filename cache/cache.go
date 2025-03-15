package cache

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	apitokencache "github.com/coopersmall/subswag/cache/apitoken"
	chatsessionscache "github.com/coopersmall/subswag/cache/chatsessions"
	userscache "github.com/coopersmall/subswag/cache/users"
	"github.com/coopersmall/subswag/domain/apitoken"
	chatsessionsdomain "github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/user"
	usersdomain "github.com/coopersmall/subswag/domain/user"

	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/utils"
)

type ICache interface {
	APITokensCache(userId user.UserID) IAPITokenCache
	ChatSessionsCache() IChatSessionsCache
	UsersCache() IUsersCache
}

type Cache struct {
	apitokensCache    func(userId user.UserID) IAPITokenCache
	chatSessionsCache func() IChatSessionsCache
	usersCache        func() IUsersCache
}

func GetCaches(
	env iEnv,
	gateway gateways.IGateways,
) ICache {
	apitokensCache := func(userId user.UserID) IAPITokenCache {
		return apitokencache.NewAPITokensCache(
			env.GetLogger("api-tokens-cache"),
			env.GetTracer("api-tokens-cache"),
			userId,
			gateway.RedisCacheGateway(),
		)
	}
	chatSessionsCache := func() IChatSessionsCache {
		return chatsessionscache.NewChatSessionsCache(
			env.GetLogger("chat-sessions-cache"),
			env.GetTracer("chat-sessions-cache"),
			gateway.RedisCacheGateway(),
		)
	}
	usersCache := func() IUsersCache {
		return userscache.NewUsersCache(
			env.GetLogger("users-cache"),
			env.GetTracer("users-cache"),
			gateway.RedisCacheGateway(),
		)
	}

	return &Cache{
		apitokensCache:    apitokensCache,
		chatSessionsCache: chatSessionsCache,
		usersCache:        usersCache,
	}
}

func (c *Cache) APITokensCache(userId user.UserID) IAPITokenCache {
	return c.apitokensCache(userId)
}

func (c *Cache) ChatSessionsCache() IChatSessionsCache {
	return c.chatSessionsCache()
}

func (c *Cache) UsersCache() IUsersCache {
	return c.usersCache()
}

type iEnv interface {
	GetLogger(name string) utils.ILogger
	GetTracer(service string) apm.ITracer
}

type IUsersCache interface {
	Get(
		ctx context.Context,
		userId usersdomain.UserID,
		onMiss func(context.Context, usersdomain.UserID) (*usersdomain.User, error),
	) (*usersdomain.User, error)
	Set(
		ctx context.Context,
		userId usersdomain.UserID,
		user *usersdomain.User,
	) error
	Delete(ctx context.Context, userId usersdomain.UserID) error
}

type IChatSessionsCache interface {
	Get(
		ctx context.Context,
		chatSessionId chatsessionsdomain.ChatSessionID,
		onMiss func(context.Context, chatsessionsdomain.ChatSessionID) (*chatsessionsdomain.ChatSession, error),
	) (*chatsessionsdomain.ChatSession, error)
	Set(
		ctx context.Context,
		chatSessionId chatsessionsdomain.ChatSessionID,
		chatSession *chatsessionsdomain.ChatSession,
	) error
	Delete(
		ctx context.Context,
		chatSessionId chatsessionsdomain.ChatSessionID,
	) error
}

type IAPITokenCache interface {
	Get(
		ctx context.Context,
		apiTokenId apitoken.APITokenID,
		onMiss func(context.Context, apitoken.APITokenID) (*apitoken.APIToken, error),
	) (*apitoken.APIToken, error)
	Set(
		ctx context.Context,
		apiTokenId apitoken.APITokenID,
		apiToken *apitoken.APIToken,
	) error
	Delete(
		ctx context.Context,
		apiTokenId apitoken.APITokenID,
	) error
}
