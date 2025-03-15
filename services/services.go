package services

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/cache"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/secret"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/repos"
	answerAssistantservice "github.com/coopersmall/subswag/services/answerassistant"
	apitokenservice "github.com/coopersmall/subswag/services/apitoken"
	chatsessionservice "github.com/coopersmall/subswag/services/chatsession"
	chatsessionitemservice "github.com/coopersmall/subswag/services/chatsessionitems"
	deckservice "github.com/coopersmall/subswag/services/decks"
	encryptionservice "github.com/coopersmall/subswag/services/encryption"
	ratelimiterservice "github.com/coopersmall/subswag/services/ratelimiter"
	secretsservice "github.com/coopersmall/subswag/services/secret"
	usersservice "github.com/coopersmall/subswag/services/user"
	userauthenticationservice "github.com/coopersmall/subswag/services/userauthentication"
	"github.com/coopersmall/subswag/streams/publishers"
	"github.com/coopersmall/subswag/utils"
	"github.com/tmc/langchaingo/llms"
)

type IServices interface {
	APITokenService(userId user.UserID) IAPITokenService
	AnswerAssistantService(userId user.UserID) IAnswerAssistantService
	ChatSessionsService(userId user.UserID) IChatSessionsService
	ChatSessionItemsService(userId user.UserID) IChatSessionItemsService
	DecksService(userId user.UserID) IDecksService
	SecretsService(userId user.UserID) ISecretsService
	AuthenticationService() IAuthenticationService
	JWTService() IJWTService
	RSAService() IRSAService
	UsersService() IUsersService
}

type Services struct {
	apiTokenService         func(userId user.UserID) IAPITokenService
	answerAssistantService  func(userId user.UserID) IAnswerAssistantService
	chatSessionsService     func(userId user.UserID) IChatSessionsService
	chatSessionItemsService func(userId user.UserID) IChatSessionItemsService
	decksService            func(userId user.UserID) IDecksService
	secretsService          func(userId user.UserID) ISecretsService
	authenticationService   func() IAuthenticationService
	jwtService              func() IJWTService
	rsaService              func() IRSAService
	usersService            func() IUsersService
}

func GetServices(
	env iEnv,
	vars iEnvVars,
	repos repos.IRepos,
	gateways gateways.IGateways,
	cache cache.ICache,
	publishers publishers.IPublishers,
) IServices {
	newJWTService := func() IJWTService {
		return encryptionservice.NewJWTService(env.GetLogger("jwt-service"))
	}

	newRSAService := func() IRSAService {
		return encryptionservice.NewRSAService(env.GetLogger("rsa-service"))
	}

	newAuthenticationService := func() IAuthenticationService {
		return userauthenticationservice.NewAuthenticationService(
			env.GetLogger("authentication-service"),
			env.GetTracer("authentication-service"),
			newRSAService(),
			newJWTService(),
			vars.GetRequestsRSAPrivateKey,
			vars.GetRequestsRSAPublicKey,
			vars.GetJWTSigningKey,
		)
	}

	newChatSessionItemsService := func(userId user.UserID) IChatSessionItemsService {
		return chatsessionitemservice.NewChatSessionItemsService(
			env.GetLogger("chat-session-items-service"),
			env.GetTracer("chat-session-items-service"),
			repos.ChatSessionItemsRepo(userId),
		)
	}

	newChatSessionsService := func(userId user.UserID) IChatSessionsService {
		return chatsessionservice.NewChatSessionsService(
			env.GetLogger("chat-sessions-service"),
			env.GetTracer("chat-sessions-service"),
			cache.ChatSessionsCache(),
			repos.ChatSessionsRepo(userId),
		)
	}

	newDeckService := func(userId user.UserID) IDecksService {
		return deckservice.NewDecksService(
			env.GetLogger("deck-service"),
			repos.DecksRepo(userId),
		)
	}

	newAnswerAssistantService := func(userId user.UserID) IAnswerAssistantService {
		return answerAssistantservice.NewAnswerAssistantService(
			env.GetLogger("answer-assistant-service"),
			env.GetTracer("answer-assistant-service"),
			gateways.GroqGateway(),
			newChatSessionItemsService(userId),
		)
	}

	newSecretsService := func(userId user.UserID) ISecretsService {
		return secretsservice.NewSecretService(
			env.GetLogger("secrets-service"),
			newRSAService(),
			repos.SecretsRepo(userId),
			vars.GetSecretsRSAPrivateKey,
			vars.GetSecretsRSAPublicKey,
		)
	}

	newAPITokenService := func(userId user.UserID) IAPITokenService {
		return apitokenservice.NewAPITokenService(
			env.GetLogger("api-token-service"),
			env.GetTracer("api-token-service"),
			newJWTService(),
			repos.APITokenRepo(userId),
			cache.APITokensCache(userId),
			vars.GetJWTSigningKey,
		)
	}

	newUsersService := func() IUsersService {
		return usersservice.NewUsersService(
			env.GetLogger("users-service"),
			env.GetTracer("users-service"),
			cache.UsersCache(),
			repos.UsersRepo(),
			publishers.UsersPublisher(),
		)
	}

	return &Services{
		apiTokenService:         newAPITokenService,
		answerAssistantService:  newAnswerAssistantService,
		chatSessionsService:     newChatSessionsService,
		chatSessionItemsService: newChatSessionItemsService,
		decksService:            newDeckService,
		jwtService:              newJWTService,
		rsaService:              newRSAService,
		secretsService:          newSecretsService,
		authenticationService:   newAuthenticationService,
		usersService:            newUsersService,
	}
}

var (
	NewAPITokenService           = apitokenservice.NewAPITokenService
	NewAnswerAssistantService    = answerAssistantservice.NewAnswerAssistantService
	NewChatSessionsService       = chatsessionservice.NewChatSessionsService
	NewChatSessionItemsService   = chatsessionitemservice.NewChatSessionItemsService
	NewDeckService               = deckservice.NewDecksService
	NewJWTService                = encryptionservice.NewJWTService
	NewRateLimiterService        = ratelimiterservice.NewRateLimiterService
	NewRSAService                = encryptionservice.NewRSAService
	NewSecretService             = secretsservice.NewSecretService
	NewUserAuthenticationService = userauthenticationservice.NewAuthenticationService
	NewUserService               = usersservice.NewUsersService
)

func (s *Services) APITokenService(userId user.UserID) IAPITokenService {
	return s.apiTokenService(userId)
}

func (s *Services) AnswerAssistantService(userId user.UserID) IAnswerAssistantService {
	return s.answerAssistantService(userId)
}

func (s *Services) ChatSessionsService(userId user.UserID) IChatSessionsService {
	return s.chatSessionsService(userId)
}

func (s *Services) ChatSessionItemsService(userId user.UserID) IChatSessionItemsService {
	return s.chatSessionItemsService(userId)
}

func (s *Services) DecksService(userId user.UserID) IDecksService {
	return s.decksService(userId)
}

func (s *Services) JWTService() IJWTService {
	return s.jwtService()
}

func (s *Services) RSAService() IRSAService {
	return s.rsaService()
}

func (s *Services) SecretsService(userId user.UserID) ISecretsService {
	return s.secretsService(userId)
}

func (s *Services) AuthenticationService() IAuthenticationService {
	return s.authenticationService()
}

func (s *Services) UsersService() IUsersService {
	return s.usersService()
}

type iEnv interface {
	GetLogger(name string) utils.ILogger
	GetTracer(service string) apm.ITracer
}

type iEnvVars interface {
	GetRequestsRSAPublicKey() (*rsa.PublicKey, error)
	GetRequestsRSAPrivateKey() (*rsa.PrivateKey, error)
	GetSecretsRSAPublicKey() (*rsa.PublicKey, error)
	GetSecretsRSAPrivateKey() (*rsa.PrivateKey, error)
	GetJWTSigningKey() ([]byte, error)
}

type IAPITokenService interface {
	CreateToken(ctx context.Context, data *apitoken.APITokenData) (*apitoken.APITokenWithSecret, error)
	CreateTokenWithId(ctx context.Context, id utils.ID, data *apitoken.APITokenData) (*apitoken.APITokenWithSecret, error)
	GetToken(ctx context.Context, tokenId apitoken.APITokenID) (*apitoken.APIToken, error)
	UpdateToken(ctx context.Context, token *apitoken.APIToken) (*apitoken.APIToken, error)
	DeleteToken(ctx context.Context, tokenID apitoken.APITokenID) error
}

type IAnswerAssistantService interface {
	AnswerQuestion(ctx context.Context, userId user.UserID, items []chatsession.ChatSessionItem) (*chatsession.AssistantChatSessionItem, error)
}

type IChatSessionsService interface {
	GetChatSession(ctx context.Context, sessionId chatsession.ChatSessionID) (*chatsession.ChatSession, error)
	GetAllChatSessions(ctx context.Context) ([]*chatsession.ChatSession, error)
	CreateChatSession(ctx context.Context, data chatsession.ChatSessionData) error
	UpdateChatSession(ctx context.Context, session *chatsession.ChatSession) error
	DeleteChatSession(ctx context.Context, sessionId chatsession.ChatSessionID) error
}

type IChatSessionItemsService interface {
	GetChatSessionItem(ctx context.Context, itemId chatsession.ChatSessionItemID) (chatsession.ChatSessionItem, error)
	GetAllChatSessionItems(ctx context.Context) ([]chatsession.ChatSessionItem, error)
	GetChatSessionItemsBySessionId(ctx context.Context, sessionId chatsession.ChatSessionID) ([]chatsession.ChatSessionItem, error)
	CreateChatSessionItem(ctx context.Context, item chatsession.ChatSessionItem) error
	UpdateChatSessionItem(ctx context.Context, item chatsession.ChatSessionItem) error
	DeleteChatSessionItem(ctx context.Context, itemId chatsession.ChatSessionItemID) error
	DeleteChatSessionItemsBySessionId(ctx context.Context, sessionId chatsession.ChatSessionID) error
	ConvertChatSessionItemsToLLMMessages(ctx context.Context, items []chatsession.ChatSessionItem) []llms.MessageContent
}

type IDecksService interface {
	GetDeck(ctx context.Context, deckId card.SerializableDeckID) (*card.SerializableDeck, error)
	GetAllDecks(ctx context.Context) ([]*card.SerializableDeck, error)
	CreateDeck(ctx context.Context, data card.SerializableDeckData) error
	UpdateDeck(ctx context.Context, deck *card.SerializableDeck) error
	DeleteDeck(ctx context.Context, deckId card.SerializableDeckID) error
}

type IJWTService interface {
	CreateToken(ctx context.Context, userId user.UserID, signingKey []byte) (string, error)
	CreateTokenWithID(ctx context.Context, userId user.UserID, tokenId utils.ID, signingKey []byte) (string, error)
	ValidateToken(ctx context.Context, token string, signingKey []byte) (user.UserID, utils.ID, error)
}

type IRSAService interface {
	Encrypt(ctx context.Context, data []byte, publicKey *rsa.PublicKey) ([]byte, error)
	Decrypt(ctx context.Context, data []byte, privateKey *rsa.PrivateKey) ([]byte, error)
}

type ISecretsService interface {
	GetSecret(ctx context.Context, secretId secret.SecretID) (secret.Secret, error)
	CreateSecret(ctx context.Context, value string) (*secret.StoredSecret, error)
	CreateSecretWithId(ctx context.Context, secretId secret.SecretID, value string) (*secret.StoredSecret, error)
	UpdateSecret(ctx context.Context, secretId secret.SecretID, value string) (*secret.StoredSecret, error)
	DeleteSecret(ctx context.Context, secretId secret.SecretID) error
}

type IRateLimitService interface {
	IsRateLimited(ctx context.Context, userID user.UserID) (bool, error)
}

type IAuthenticationService interface {
	AuthenticateMAC(ctx context.Context, timestamp time.Time, signature string) error
	AuthenticateToken(ctx context.Context, token string, now time.Time, permissions ...domain.Permission) (user.UserID, apitoken.APITokenID, error)
	AuthenticateSession(ctx context.Context, session string) (user.UserID, error)
}

type IUsersService interface {
	CreateUser(ctx context.Context, data user.UserData) (*user.User, error)
	CreateUserWithId(ctx context.Context, id user.UserID, data user.UserData) (*user.User, error)
	UpdateUser(ctx context.Context, user *user.User) (*user.User, error)
	GetUser(ctx context.Context, userId user.UserID) (*user.User, error)
	GetAllUsers(ctx context.Context) ([]*user.User, error)
	DeleteUser(ctx context.Context, userId user.UserID) error
}
