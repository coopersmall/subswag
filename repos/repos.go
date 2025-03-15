package repos

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/game"
	"github.com/coopersmall/subswag/domain/integrations"
	"github.com/coopersmall/subswag/domain/ratelimit"
	"github.com/coopersmall/subswag/domain/secret"
	"github.com/coopersmall/subswag/domain/user"
	apitokensrepo "github.com/coopersmall/subswag/repos/apitokens"
	cardsrepo "github.com/coopersmall/subswag/repos/cards"
	chatsessionitemsrepo "github.com/coopersmall/subswag/repos/chatsessionitems"
	chatsessionsrepo "github.com/coopersmall/subswag/repos/chatsessions"
	decksrepo "github.com/coopersmall/subswag/repos/decks"
	gamesstaterepo "github.com/coopersmall/subswag/repos/games"
	gamestateversionsrepo "github.com/coopersmall/subswag/repos/games"
	ratelimitsrepo "github.com/coopersmall/subswag/repos/ratelimits"
	secretsrepo "github.com/coopersmall/subswag/repos/secrets"
	usersrepo "github.com/coopersmall/subswag/repos/users"
)

type IRepos interface {
	APITokenRepo(userId user.UserID) IAPITokenRepo
	CardsRepo() ICardsRepo
	ChatSessionsRepo(userId user.UserID) IChatSessionsRepo
	ChatSessionItemsRepo(userId user.UserID) IChatSessionItemsRepo
	DecksRepo(userId user.UserID) IDecksRepo
	GameStateRepo() IGameStateRepo
	GameStateVersionRepo() IGameStateVersionRepo
	SecretsRepo(userId user.UserID) ISecretsRepo
	RateLimitsRepo(userId user.UserID) IRateLimitsRepo
	UsersRepo() IUsersRepo
}

type Repos struct {
	apiTokensRepo        func(userId user.UserID) *apitokensrepo.APITokenRepo
	cardsRepo            func() *cardsrepo.CardsRepo
	chatSessionsRepo     func(userId user.UserID) *chatsessionsrepo.ChatSessionsRepo
	chatSessionItemsRepo func(userId user.UserID) *chatsessionitemsrepo.ChatSessionItemsRepo
	decksRepo            func(userId user.UserID) *decksrepo.DecksRepo
	gamesStateRepo       func() *gamesstaterepo.GameStateRepo
	gameStateVersionRepo func() *gamestateversionsrepo.GameStateVersionRepo
	secretsRepo          func(userId user.UserID) *secretsrepo.SecretsRepo
	rateLimitsRepo       func(userId user.UserID) *ratelimitsrepo.RateLimitsRepo
	usersRepo            func() *usersrepo.UsersRepo
}

func GetRepos(env iEnv) IRepos {
	apiTokensRepo := func(userId user.UserID) *apitokensrepo.APITokenRepo {
		return NewAPITokenRepo(
			env.GetQuerier(),
			env.GetTracer("api_tokens_repo"),
			userId,
		)
	}

	cardsRepo := func() *cardsrepo.CardsRepo {
		return NewCardsRepo(
			env.GetQuerier(),
			env.GetTracer("cards_repo"),
		)
	}

	chatSessionsRepo := func(userId user.UserID) *chatsessionsrepo.ChatSessionsRepo {
		return NewChatSessionsRepo(
			env.GetQuerier(),
			env.GetTracer("chat_sessions_repo"),
			userId,
		)
	}

	chatSessionItemsRepo := func(userId user.UserID) *chatsessionitemsrepo.ChatSessionItemsRepo {
		return NewChatSessionItemsRepo(
			env.GetQuerier(),
			env.GetTracer("chat_session_items_repo"),
			userId,
		)
	}

	decksRepo := func(userId user.UserID) *decksrepo.DecksRepo {
		return NewDecksRepo(
			env.GetQuerier(),
			env.GetTracer("decks_repo"),
			userId,
		)
	}

	gamesStateRepo := func() *gamesstaterepo.GameStateRepo {
		return gamesstaterepo.NewGameStateRepo(
			env.GetQuerier(),
			env.GetTracer("games_state_repo"),
		)
	}

	gameStateVersionRepo := func() *gamestateversionsrepo.GameStateVersionRepo {
		return gamestateversionsrepo.NewGameStateVersionRepo(
			env.GetQuerier(),
			env.GetTracer("game_state_version_repo"),
		)
	}

	secretsRepo := func(userId user.UserID) *secretsrepo.SecretsRepo {
		return NewSecretsRepo(
			env.GetQuerier(),
			env.GetTracer("secrets_repo"),
			userId,
		)
	}

	rateLimitsRepo := func(userId user.UserID) *ratelimitsrepo.RateLimitsRepo {
		return NewRateLimitRepo(
			env.GetQuerier(),
			env.GetTracer("rate_limits_repo"),
			userId,
		)
	}

	usersRepo := func() *usersrepo.UsersRepo {
		return NewUserRepo(
			env.GetQuerier(),
			env.GetTracer("users_repo"),
		)
	}

	return &Repos{
		apiTokensRepo:        apiTokensRepo,
		cardsRepo:            cardsRepo,
		chatSessionsRepo:     chatSessionsRepo,
		chatSessionItemsRepo: chatSessionItemsRepo,
		decksRepo:            decksRepo,
		gamesStateRepo:       gamesStateRepo,
		gameStateVersionRepo: gameStateVersionRepo,
		secretsRepo:          secretsRepo,
		rateLimitsRepo:       rateLimitsRepo,
		usersRepo:            usersRepo,
	}
}

func (r *Repos) APITokenRepo(userId user.UserID) IAPITokenRepo {
	return r.apiTokensRepo(userId)
}

func (r *Repos) CardsRepo() ICardsRepo {
	return r.cardsRepo()
}

func (r *Repos) ChatSessionsRepo(userId user.UserID) IChatSessionsRepo {
	return r.chatSessionsRepo(userId)
}

func (r *Repos) ChatSessionItemsRepo(userId user.UserID) IChatSessionItemsRepo {
	return r.chatSessionItemsRepo(userId)
}

func (r *Repos) DecksRepo(userId user.UserID) IDecksRepo {
	return r.decksRepo(userId)
}

func (r *Repos) GameStateRepo() IGameStateRepo {
	return r.gamesStateRepo()
}

func (r *Repos) GameStateVersionRepo() IGameStateVersionRepo {
	return r.gameStateVersionRepo()
}

func (r *Repos) SecretsRepo(userId user.UserID) ISecretsRepo {
	return r.secretsRepo(userId)
}

func (r *Repos) RateLimitsRepo(userId user.UserID) IRateLimitsRepo {
	return r.rateLimitsRepo(userId)
}

func (r *Repos) UsersRepo() IUsersRepo {
	return r.usersRepo()
}

var (
	NewAPITokenRepo         = apitokensrepo.NewAPITokenRepo
	NewCardsRepo            = cardsrepo.NewCardsRepo
	NewChatSessionsRepo     = chatsessionsrepo.NewChatSessionsRepo
	NewChatSessionItemsRepo = chatsessionitemsrepo.NewChatSessionItemsRepo
	NewDecksRepo            = decksrepo.NewDecksRepo
	NewGameStateRepo        = gamesstaterepo.NewGameStateRepo
	NewGameStateVersionRepo = gamestateversionsrepo.NewGameStateVersionRepo
	NewSecretsRepo          = secretsrepo.NewSecretsRepo
	NewRateLimitRepo        = ratelimitsrepo.NewRateLimitsRepo
	NewUserRepo             = usersrepo.NewUsersRepo
)

type iEnv interface {
	GetQuerier() db.IQuerier
	GetTracer(string) apm.ITracer
}

type IAPITokenRepo interface {
	Get(ctx context.Context, tokenId apitoken.APITokenID) (*apitoken.APIToken, error)
	All(ctx context.Context) ([]*apitoken.APIToken, error)
	Create(ctx context.Context, token *apitoken.APIToken) error
	Update(ctx context.Context, token *apitoken.APIToken) error
	Delete(ctx context.Context, tokenId apitoken.APITokenID) error
}

type IChatSessionsRepo interface {
	Get(ctx context.Context, chatSessionId chatsession.ChatSessionID) (*chatsession.ChatSession, error)
	All(ctx context.Context) ([]*chatsession.ChatSession, error)
	Create(ctx context.Context, chatSession *chatsession.ChatSession) error
	Update(ctx context.Context, chatSession *chatsession.ChatSession) error
	Delete(ctx context.Context, chatSessionId chatsession.ChatSessionID) error
}

type IChatSessionItemsRepo interface {
	Get(ctx context.Context, itemId chatsession.ChatSessionItemID) (chatsession.ChatSessionItem, error)
	GetBySessionId(ctx context.Context, sessionId chatsession.ChatSessionID) ([]chatsession.ChatSessionItem, error)
	All(ctx context.Context) ([]chatsession.ChatSessionItem, error)
	Create(ctx context.Context, item chatsession.ChatSessionItem) error
	Update(ctx context.Context, item chatsession.ChatSessionItem) error
	Delete(ctx context.Context, itemId chatsession.ChatSessionItemID) error
	DeleteBySessionId(ctx context.Context, sessionId chatsession.ChatSessionID) error
}

type IDecksRepo interface {
	Get(ctx context.Context, deckId card.SerializableDeckID) (*card.SerializableDeck, error)
	All(ctx context.Context) ([]*card.SerializableDeck, error)
	Create(ctx context.Context, deck *card.SerializableDeck) error
	Update(ctx context.Context, deck *card.SerializableDeck) error
	Delete(ctx context.Context, deckId card.SerializableDeckID) error
}

type ISecretsRepo interface {
	Get(ctx context.Context, secretId secret.SecretID) (*secret.StoredSecret, error)
	All(ctx context.Context) ([]*secret.StoredSecret, error)
	Create(ctx context.Context, secret *secret.StoredSecret) error
	Update(ctx context.Context, secret *secret.StoredSecret) error
	Delete(ctx context.Context, secretId secret.SecretID) error
}

type IRateLimitsRepo interface {
	Get(ctx context.Context, rateLimitId ratelimit.RateLimitID) (*ratelimit.RateLimit, error)
	All(ctx context.Context) ([]*ratelimit.RateLimit, error)
	Create(ctx context.Context, rateLimit *ratelimit.RateLimit) error
	Update(ctx context.Context, rateLimit *ratelimit.RateLimit) error
	Delete(ctx context.Context, rateLimitId ratelimit.RateLimitID) error
}

type ICardsRepo interface {
	Get(ctx context.Context, cardId card.SerializableCardID) (card.Card, error)
	All(ctx context.Context) ([]card.Card, error)
	Create(ctx context.Context, card card.Card) error
	Update(ctx context.Context, card card.Card) error
	Delete(ctx context.Context, cardId card.SerializableCardID) error
}

type IGameStateRepo interface {
	Get(ctx context.Context, gameStateId game.GameStateID) (*game.GameState, error)
	All(ctx context.Context) ([]*game.GameState, error)
	Create(ctx context.Context, gameState *game.GameState) error
	Update(ctx context.Context, gameState *game.GameState) error
	Delete(ctx context.Context, gameStateId game.GameStateID) error
}

type IGameStateVersionRepo interface {
	Get(ctx context.Context, versionId game.GameStateVersionID) (*game.GameStateVersion, error)
	Create(ctx context.Context, version *game.GameStateVersion) error
	GetVersionsForGameState(ctx context.Context, gameStateId game.GameStateID) ([]*game.GameStateVersion, error)
	GetLatestVersion(ctx context.Context, gameStateId game.GameStateID) (*game.GameStateVersion, error)
}

type IIntegrationsRepo interface {
	Get(ctx context.Context, integrationId integrations.IntegrationID) (integrations.Integration, error)
	All(ctx context.Context) ([]integrations.Integration, error)
	Create(ctx context.Context, integration integrations.Integration) error
	Update(ctx context.Context, integration integrations.Integration) error
	Delete(ctx context.Context, integrationId integrations.IntegrationID) error
}

type IUsersRepo interface {
	Get(ctx context.Context, userId user.UserID) (*user.User, error)
	All(ctx context.Context) ([]*user.User, error)
	Create(ctx context.Context, user *user.User) error
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, userId user.UserID) error
}
