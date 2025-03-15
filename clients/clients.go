package clients

import (
	"context"
	"database/sql"
	"time"

	"github.com/coopersmall/subswag/clients/apm"
	"github.com/coopersmall/subswag/clients/groq"
	"github.com/coopersmall/subswag/clients/http"
	"github.com/coopersmall/subswag/clients/openai"
	"github.com/coopersmall/subswag/clients/postgres"
	"github.com/coopersmall/subswag/clients/redis"
	"github.com/tmc/langchaingo/llms"
)

type IClients interface {
	APMClient() IAPMClient
	OpenAIClient() IAIClient
	GroqClient() IAIClient
	HttpClient(baseUrl string, beforeRequest func(req *http.RequestConfig)) IHttpClient
	RedisClient() IRedisClient
	PostgresClient(
		url string,
		opts *struct {
			MaxOpenConns int
			MaxIdleConns int
		},
	) (*sql.DB, error)
}

type Clients struct {
	apmClient      func() IAPMClient
	httpClient     func(baseUrl string, beforeRequest func(req *http.RequestConfig)) IHttpClient
	openaiClient   func() IAIClient
	groqClient     func() IAIClient
	redisClient    func() IRedisClient
	postgresClient func(url string, opts *struct {
		MaxOpenConns int
		MaxIdleConns int
	}) (*sql.DB, error)
}

func GetClients(envVars iEnvVars) IClients {
	newOpenAIClient := func() IAIClient {
		return openai.NewOpenAIClient(envVars.GetOpenAIKey)
	}

	newGroqClient := func() IAIClient {
		return groq.NewGroqClient(envVars.GetGroqKey)
	}

	newHttpClient := func(
		baseUrl string,
		beforeRequest func(req *http.RequestConfig),
	) IHttpClient {
		return http.NewHttpClient(baseUrl, beforeRequest)
	}

	newRedisClient := func() IRedisClient {
		return redis.NewRedisClient(
			envVars.GetRedisURL,
			envVars.GetRedisPassword,
			envVars.GetRedisDB,
		)
	}

	newAPMClient := func() IAPMClient {
		return apm.NewAPMClient(
			envVars.GetAPMURL,
			envVars.GetAPMTimeout,
		)
	}

	newPostgresClient := func(
		url string,
		opts *struct {
			MaxOpenConns int
			MaxIdleConns int
		},
	) (*sql.DB, error) {
		return postgres.NewPostgresClient(url, opts)
	}

	return &Clients{
		apmClient:      newAPMClient,
		openaiClient:   newOpenAIClient,
		httpClient:     newHttpClient,
		redisClient:    newRedisClient,
		groqClient:     newGroqClient,
		postgresClient: newPostgresClient,
	}
}

func (c *Clients) APMClient() IAPMClient {
	return c.apmClient()
}

func (c *Clients) OpenAIClient() IAIClient {
	return c.openaiClient()
}

func (c *Clients) HttpClient(baseUrl string, beforeRequest func(req *http.RequestConfig)) IHttpClient {
	return c.httpClient(baseUrl, beforeRequest)
}

func (c *Clients) GroqClient() IAIClient {
	return c.groqClient()
}

func (c *Clients) RedisClient() IRedisClient {
	return c.redisClient()
}

func (c *Clients) PostgresClient(
	url string,
	opts *struct {
		MaxOpenConns int
		MaxIdleConns int
	},
) (*sql.DB, error) {
	return c.postgresClient(url, opts)
}

type iEnvVars interface {
	GetAPMURL() (string, error)
	GetAPMTimeout() (time.Duration, error)
	GetOpenAIKey() (string, error)
	GetGroqKey() (string, error)
	GetRedisURL() (string, error)
	GetRedisPassword() (string, error)
	GetRedisDB() (int, error)
	GetSharedPostgresURL() (string, error)
	GetStandardPostgresURL() (string, error)
}

type IHttpClient interface {
	Post(ctx context.Context, path string, body any, headers map[string]string) (any, error)
	Patch(ctx context.Context, path string, body any, headers map[string]string) (any, error)
	Put(ctx context.Context, path string, body any, headers map[string]string) (any, error)
	Delete(ctx context.Context, path string, headers map[string]string) (any, error)
	Get(ctx context.Context, path string, headers map[string]string) (any, error)
}

type IRedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
	XAdd(ctx context.Context, stream string, value map[string]any) (string, error)
	XReadGroup(ctx context.Context, args *redis.XReadGroupArgs) ([]redis.XStream, error)
	XGroupCreateConsumer(ctx context.Context, key, group, consumer string) error
	XGroupCreateMkStream(ctx context.Context, key string, consumerGroup string) error
	XAck(ctx context.Context, key string, consumerGroup string, ids ...string) error
	Close() error
}

type IAPMClient apm.IAPMClient

type IAIClient interface {
	GenerateContent(
		ctx context.Context,
		messageContent []llms.MessageContent,
		opts ...llms.CallOption,
	) (*llms.ContentResponse, error)
}
