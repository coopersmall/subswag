package apitoken

import (
	"fmt"
	"time"

	"github.com/coopersmall/subswag/apm"
	cachedomain "github.com/coopersmall/subswag/cache/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/utils"
)

type APITokensCache struct {
	*cachedomain.StandardCache[apitoken.APITokenID, *apitoken.APIToken]
}

func NewAPITokensCache(
	logger utils.ILogger,
	tracer apm.ITracer,
	userId user.UserID,
	gateway gateways.ICacheGateway,
) *APITokensCache {
	return &APITokensCache{
		StandardCache: cachedomain.NewStandardCache[apitoken.APITokenID, *apitoken.APIToken](
			"api-token-cache",
			logger,
			tracer,
			gateway,
			time.Duration(0),
			func(id apitoken.APITokenID) string {
				return cachedomain.UserIDKey(
					userId,
					"api-token",
					fmt.Sprintf("%d", id),
				)
			},
		),
	}
}
