package db

import (
	"context"
	"database/sql"

	"github.com/coopersmall/subswag/domain/user"
)

// Shared Queries - Read Only
type ISharedQueriesReadOnly interface {
	GetCard(ctx context.Context, id int64) (Card, error)
	GetAllCards(ctx context.Context) ([]Card, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetGameState(ctx context.Context, id int64) (GameState, error)
	GetAllGameStates(ctx context.Context) ([]GameState, error)
	GetGameStateVersion(ctx context.Context, id int64) (GameStateVersion, error)
	GetAllGameStateVersions(ctx context.Context) ([]GameStateVersion, error)
	GetGameStateVersionsByGameStateID(ctx context.Context, id int64) ([]GameStateVersion, error)
	GetLatestGameStateVersionByGameStateID(ctx context.Context, id int64) (GameStateVersion, error)
	GetIntegration(ctx context.Context, id int64) (Integration, error)
	GetAllIntegrations(ctx context.Context) ([]Integration, error)
}

// Shared Queries - Read Write
type ISharedQueriesReadWrite interface {
	ISharedQueriesReadOnly
	CreateIntegration(ctx context.Context, arg CreateIntegrationParams) (sql.Result, error)
	UpdateIntegration(ctx context.Context, arg UpdateIntegrationParams) (sql.Result, error)
	DeleteIntegration(ctx context.Context, id int64) (sql.Result, error)
	CreateCard(ctx context.Context, arg CreateCardParams) (sql.Result, error)
	UpdateCard(ctx context.Context, arg UpdateCardParams) (sql.Result, error)
	DeleteCard(ctx context.Context, id int64) (sql.Result, error)
	CreateGameState(ctx context.Context, arg CreateGameStateParams) (sql.Result, error)
	UpdateGameState(ctx context.Context, arg UpdateGameStateParams) (sql.Result, error)
	DeleteGameState(ctx context.Context, id int64) (sql.Result, error)
	CreateGameStateVersion(ctx context.Context, arg CreateGameStateVersionParams) (sql.Result, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error)
	DeleteUser(ctx context.Context, id int64) (sql.Result, error)
	WithTx(tx *sql.Tx) *Queries
}

// Standard Queries - Read Only
type IStandardQueriesReadOnly interface {
	GetAPIToken(ctx context.Context, arg GetAPITokenParams) (ApiToken, error)
	GetAllAPITokens(ctx context.Context) ([]ApiToken, error)
	GetChatSession(ctx context.Context, arg GetChatSessionParams) (ChatSession, error)
	GetAllChatSessions(ctx context.Context) ([]ChatSession, error)
	GetChatSessionItem(ctx context.Context, arg GetChatSessionItemParams) (ChatSessionItem, error)
	GetChatSessionItemsBySessionID(ctx context.Context, arg GetChatSessionItemsBySessionIDParams) ([]ChatSessionItem, error)
	GetAllChatSessionItems(ctx context.Context) ([]ChatSessionItem, error)
	GetDeck(ctx context.Context, arg GetDeckParams) (Deck, error)
	GetAllDecks(ctx context.Context) ([]Deck, error)
	GetSecret(ctx context.Context, arg GetSecretParams) (Secret, error)
	GetAllSecrets(ctx context.Context) ([]Secret, error)
	GetRateLimit(ctx context.Context, arg GetRateLimitParams) (RateLimit, error)
	GetAllRateLimits(ctx context.Context) ([]RateLimit, error)
}

// Standard Queries - Read Write
type IStandardQueriesReadWrite interface {
	IStandardQueriesReadOnly
	CreateAPIToken(ctx context.Context, arg CreateAPITokenParams) (sql.Result, error)
	UpdateAPIToken(ctx context.Context, arg UpdateAPITokenParams) (sql.Result, error)
	DeleteAPIToken(ctx context.Context, arg DeleteAPITokenParams) (sql.Result, error)
	CreateChatSession(ctx context.Context, arg CreateChatSessionParams) (sql.Result, error)
	UpdateChatSession(ctx context.Context, arg UpdateChatSessionParams) (sql.Result, error)
	DeleteChatSession(ctx context.Context, arg DeleteChatSessionParams) (sql.Result, error)
	CreateChatSessionItem(ctx context.Context, arg CreateChatSessionItemParams) (sql.Result, error)
	UpdateChatSessionItem(ctx context.Context, arg UpdateChatSessionItemParams) (sql.Result, error)
	DeleteChatSessionItem(ctx context.Context, arg DeleteChatSessionItemParams) (sql.Result, error)
	DeleteChatSessionItemsBySessionID(ctx context.Context, arg DeleteChatSessionItemsBySessionIDParams) error
	CreateDeck(ctx context.Context, arg CreateDeckParams) (sql.Result, error)
	UpdateDeck(ctx context.Context, arg UpdateDeckParams) (sql.Result, error)
	DeleteDeck(ctx context.Context, arg DeleteDeckParams) (sql.Result, error)
	CreateSecret(ctx context.Context, arg CreateSecretParams) (sql.Result, error)
	UpdateSecret(ctx context.Context, arg UpdateSecretParams) (sql.Result, error)
	DeleteSecret(ctx context.Context, arg DeleteSecretParams) (sql.Result, error)
	CreateRateLimit(ctx context.Context, arg CreateRateLimitParams) (sql.Result, error)
	UpdateRateLimit(ctx context.Context, arg UpdateRateLimitParams) (sql.Result, error)
	DeleteRateLimit(ctx context.Context, arg DeleteRateLimitParams) (sql.Result, error)
	WithTx(tx *sql.Tx) *Queries
}

type IQuerier interface {
	Shared(ctx context.Context, fn func(ISharedQueriesReadOnly) error) error
	SharedWrite(ctx context.Context, fn func(ISharedQueriesReadWrite) error) error
	Standard(ctx context.Context, userId user.UserID, fn func(IStandardQueriesReadOnly) error) error
	StandardWrite(ctx context.Context, userId user.UserID, fn func(IStandardQueriesReadWrite) error) error
}

type Querier struct {
	manager IDBManager
}

func NewQuerier(manager IDBManager) IQuerier {
	return &Querier{
		manager: manager,
	}
}

func (q *Querier) Shared(ctx context.Context, fn func(ISharedQueriesReadOnly) error) error {
	db := q.manager.ReadOnly()
	return fn(New(db))
}

func (q *Querier) SharedWrite(ctx context.Context, fn func(ISharedQueriesReadWrite) error) error {
	db := q.manager.ReadWrite()
	return fn(New(db))
}

func (q *Querier) Standard(ctx context.Context, userId user.UserID, fn func(IStandardQueriesReadOnly) error) error {
	db := q.manager.ReadOnly()
	return fn(New(db))
}

func (q *Querier) StandardWrite(ctx context.Context, userId user.UserID, fn func(IStandardQueriesReadWrite) error) error {
	db := q.manager.ReadWrite()
	return fn(New(db))
}
