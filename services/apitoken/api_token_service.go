package apitoken

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/cache"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	servicesdomain "github.com/coopersmall/subswag/services/domain"
	"github.com/coopersmall/subswag/utils"
)

type APITokenService struct {
	logger          utils.ILogger
	jwtService      jwtService
	tokenSigningKey []byte
	standardService *servicesdomain.StandardService[apitoken.APITokenID, *apitoken.APIToken]
}

func NewAPITokenService(
	logger utils.ILogger,
	tracer apm.ITracer,
	jwtService jwtService,
	apiTokenRepo repos.IAPITokenRepo,
	apiTokenCache cache.IAPITokenCache,
	getTokenSigningKey func() ([]byte, error),
) *APITokenService {
	tokenSigningKey, err := getTokenSigningKey()
	if err != nil {
		panic(utils.NewInternalError("failed to get token signing key", err))
	}
	return &APITokenService{
		logger:          logger,
		jwtService:      jwtService,
		tokenSigningKey: tokenSigningKey,
		standardService: servicesdomain.NewStandardService[apitoken.APITokenID, *apitoken.APIToken](
			"apitoken",
			logger,
			tracer,
			struct {
				Get    func(context.Context, apitoken.APITokenID) (*apitoken.APIToken, error)
				All    func(context.Context) ([]*apitoken.APIToken, error)
				Create func(context.Context, *apitoken.APIToken) error
				Update func(context.Context, *apitoken.APIToken) error
				Delete func(context.Context, apitoken.APITokenID) error
			}{
				Get:    apiTokenRepo.Get,
				All:    apiTokenRepo.All,
				Create: apiTokenRepo.Create,
				Update: apiTokenRepo.Update,
				Delete: apiTokenRepo.Delete,
			},
			&struct {
				Get    func(context.Context, apitoken.APITokenID, func(context.Context, apitoken.APITokenID) (*apitoken.APIToken, error)) (*apitoken.APIToken, error)
				Set    func(context.Context, apitoken.APITokenID, *apitoken.APIToken) error
				Delete func(context.Context, apitoken.APITokenID) error
			}{
				Get:    apiTokenCache.Get,
				Set:    apiTokenCache.Set,
				Delete: apiTokenCache.Delete,
			},
			nil,
		),
	}
}

func (a *APITokenService) CreateToken(
	ctx context.Context,
	data *apitoken.APITokenData,
) (*apitoken.APITokenWithSecret, error) {
	return a.createToken(ctx, utils.NewID(), data)
}

func (a *APITokenService) CreateTokenWithId(
	ctx context.Context,
	id utils.ID,
	data *apitoken.APITokenData,
) (*apitoken.APITokenWithSecret, error) {
	return a.createToken(ctx, utils.ID(id), data)
}

func (a *APITokenService) createToken(
	ctx context.Context,
	tokenId utils.ID,
	data *apitoken.APITokenData,
) (*apitoken.APITokenWithSecret, error) {
	userId := data.UserId
	expiry := data.Expiry
	permissions := data.Permissions

	token, err := a.jwtService.CreateTokenWithID(
		ctx,
		userId,
		tokenId,
		a.tokenSigningKey,
	)
	if err != nil {
		a.logger.Error(ctx, "failed to create token", err, nil)
		return nil, err
	}
	apiToken := &apitoken.APIToken{
		ID: apitoken.APITokenID(tokenId),
		APITokenData: apitoken.APITokenData{
			UserId:      userId,
			Expiry:      expiry,
			Permissions: permissions,
		},
		Metadata: domain.NewMetadata(),
	}
	if err != nil {
		return nil, err
	}
	err = a.standardService.Create(ctx, apiToken.ID, apiToken)
	if err != nil {
		return nil, err
	}
	return &apitoken.APITokenWithSecret{
		APIToken: *apiToken,
		Secret:   token,
	}, nil
}

func (a *APITokenService) GetToken(
	ctx context.Context,
	tokenId apitoken.APITokenID,
) (*apitoken.APIToken, error) {
	return a.standardService.Get(
		ctx,
		tokenId,
	)
}

func (a *APITokenService) UpdateToken(
	ctx context.Context,
	token *apitoken.APIToken,
) (*apitoken.APIToken, error) {
	err := a.standardService.Update(ctx, token.ID, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *APITokenService) DeleteToken(
	ctx context.Context,
	tokenID apitoken.APITokenID,
) error {
	return a.standardService.Delete(ctx, tokenID)
}

type jwtService interface {
	CreateTokenWithID(ctx context.Context, userId user.UserID, tokenId utils.ID, signingKey []byte) (string, error)
}
