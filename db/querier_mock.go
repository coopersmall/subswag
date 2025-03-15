package db

import (
	"context"
	"database/sql"

	"github.com/coopersmall/subswag/domain/user"
	"github.com/stretchr/testify/mock"
)

type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) Shared(ctx context.Context, fn func(ISharedQueriesReadOnly) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockQuerier) SharedWrite(ctx context.Context, fn func(ISharedQueriesReadWrite) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockQuerier) Standard(ctx context.Context, userId user.UserID, fn func(IStandardQueriesReadOnly) error) error {
	args := m.Called(ctx, userId, fn)
	return args.Error(0)
}

func (m *MockQuerier) StandardWrite(ctx context.Context, userId user.UserID, fn func(IStandardQueriesReadWrite) error) error {
	args := m.Called(ctx, userId, fn)
	return args.Error(0)
}

type MockSharedQueriesReadOnly struct {
	mock.Mock
}

func (m *MockSharedQueriesReadOnly) GetCard(ctx context.Context, id int64) (Card, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Card), args.Error(1)
}

func (m *MockSharedQueriesReadOnly) GetAllCards(ctx context.Context) ([]Card, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Card), args.Error(1)
}

func (m *MockSharedQueriesReadOnly) GetUser(ctx context.Context, id int64) (User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockSharedQueriesReadOnly) GetAllUsers(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockSharedQueriesReadOnly) GetGameState(ctx context.Context, id int64) (GameState, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(GameState), args.Error(1)
}

func (m *MockSharedQueriesReadOnly) GetAllGameStates(ctx context.Context) ([]GameState, error) {
	args := m.Called(ctx)
	return args.Get(0).([]GameState), args.Error(1)
}

func (m *MockSharedQueriesReadOnly) GetGameStateVersion(ctx context.Context, id int64) (GameStateVersion, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(GameStateVersion), args.Error(1)
}

func (m *MockSharedQueriesReadOnly) GetAllGameStateVersions(ctx context.Context) ([]GameStateVersion, error) {
	args := m.Called(ctx)
	return args.Get(0).([]GameStateVersion), args.Error(1)
}

func (m *MockSharedQueriesReadOnly) GetGameStateVersionsByGameStateID(ctx context.Context, id int64) ([]GameStateVersion, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]GameStateVersion), args.Error(1)
}

func (m *MockSharedQueriesReadOnly) GetLatestGameStateVersionByGameStateID(ctx context.Context, id int64) (GameStateVersion, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(GameStateVersion), args.Error(1)
}

type MockSharedQueriesReadWrite struct {
	MockSharedQueriesReadOnly
}

func (m *MockSharedQueriesReadWrite) CreateCard(ctx context.Context, arg CreateCardParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) UpdateCard(ctx context.Context, arg UpdateCardParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) DeleteCard(ctx context.Context, id int64) (sql.Result, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) CreateGameState(ctx context.Context, arg CreateGameStateParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) UpdateGameState(ctx context.Context, arg UpdateGameStateParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) DeleteGameState(ctx context.Context, id int64) (sql.Result, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) CreateGameStateVersion(ctx context.Context, arg CreateGameStateVersionParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) DeleteUser(ctx context.Context, id int64) (sql.Result, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSharedQueriesReadWrite) WithTx(tx *sql.Tx) *Queries {
	args := m.Called(tx)
	return args.Get(0).(*Queries)
}

type MockStandardQueriesReadOnly struct {
	mock.Mock
}

func (m *MockStandardQueriesReadOnly) GetAPIToken(ctx context.Context, arg GetAPITokenParams) (ApiToken, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(ApiToken), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetAllAPITokens(ctx context.Context) ([]ApiToken, error) {
	args := m.Called(ctx)
	return args.Get(0).([]ApiToken), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetChatSession(ctx context.Context, arg GetChatSessionParams) (ChatSession, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(ChatSession), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetAllChatSessions(ctx context.Context) ([]ChatSession, error) {
	args := m.Called(ctx)
	return args.Get(0).([]ChatSession), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetChatSessionItem(ctx context.Context, arg GetChatSessionItemParams) (ChatSessionItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(ChatSessionItem), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetChatSessionItemsBySessionID(ctx context.Context, arg GetChatSessionItemsBySessionIDParams) ([]ChatSessionItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]ChatSessionItem), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetAllChatSessionItems(ctx context.Context) ([]ChatSessionItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]ChatSessionItem), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetDeck(ctx context.Context, arg GetDeckParams) (Deck, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(Deck), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetAllDecks(ctx context.Context) ([]Deck, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Deck), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetSecret(ctx context.Context, arg GetSecretParams) (Secret, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(Secret), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetAllSecrets(ctx context.Context) ([]Secret, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Secret), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetRateLimit(ctx context.Context, arg GetRateLimitParams) (RateLimit, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(RateLimit), args.Error(1)
}

func (m *MockStandardQueriesReadOnly) GetAllRateLimits(ctx context.Context) ([]RateLimit, error) {
	args := m.Called(ctx)
	return args.Get(0).([]RateLimit), args.Error(1)
}

type MockStandardQueriesReadWrite struct {
	MockStandardQueriesReadOnly
}

func (m *MockStandardQueriesReadWrite) CreateAPIToken(ctx context.Context, arg CreateAPITokenParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) UpdateAPIToken(ctx context.Context, arg UpdateAPITokenParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) DeleteAPIToken(ctx context.Context, arg DeleteAPITokenParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) CreateChatSession(ctx context.Context, arg CreateChatSessionParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) UpdateChatSession(ctx context.Context, arg UpdateChatSessionParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) DeleteChatSession(ctx context.Context, arg DeleteChatSessionParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) CreateChatSessionItem(ctx context.Context, arg CreateChatSessionItemParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) UpdateChatSessionItem(ctx context.Context, arg UpdateChatSessionItemParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) DeleteChatSessionItem(ctx context.Context, arg DeleteChatSessionItemParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) DeleteChatSessionItemsBySessionID(ctx context.Context, arg DeleteChatSessionItemsBySessionIDParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockStandardQueriesReadWrite) CreateDeck(ctx context.Context, arg CreateDeckParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) UpdateDeck(ctx context.Context, arg UpdateDeckParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) DeleteDeck(ctx context.Context, arg DeleteDeckParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) CreateSecret(ctx context.Context, arg CreateSecretParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) UpdateSecret(ctx context.Context, arg UpdateSecretParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) DeleteSecret(ctx context.Context, arg DeleteSecretParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) CreateRateLimit(ctx context.Context, arg CreateRateLimitParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) UpdateRateLimit(ctx context.Context, arg UpdateRateLimitParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) DeleteRateLimit(ctx context.Context, arg DeleteRateLimitParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockStandardQueriesReadWrite) WithTx(tx *sql.Tx) *Queries {
	args := m.Called(tx)
	return args.Get(0).(*Queries)
}
