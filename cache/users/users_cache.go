package users

import (
	"time"

	"github.com/coopersmall/subswag/apm"
	cachedomain "github.com/coopersmall/subswag/cache/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/utils"
)

type UsersCache struct {
	*cachedomain.StandardCache[user.UserID, *user.User]
}

func NewUsersCache(
	logger utils.ILogger,
	tracer apm.ITracer,
	gateway gateways.ICacheGateway,
) *UsersCache {
	return &UsersCache{
		StandardCache: cachedomain.NewStandardCache[user.UserID, *user.User](
			"users-cache",
			logger,
			tracer,
			gateway,
			time.Duration(0),
			func(id user.UserID) string {
				return cachedomain.UserIDKey(id)
			},
		),
	}
}
