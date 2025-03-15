package gateways

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/clients"
	"github.com/coopersmall/subswag/gateways/groq"
	"github.com/coopersmall/subswag/gateways/openai"
	"github.com/coopersmall/subswag/gateways/redis"
	"github.com/coopersmall/subswag/utils"
	"github.com/tmc/langchaingo/llms"
)

type IGateways interface {
	OpenAIGateway() IAIGateway
	GroqGateway() IAIGateway
	RedisCacheGateway() ICacheGateway
	RedisStreamPublisherGateway() IPublisherGateway
	RedisStreamSubscriberGateway(
		opts *struct {
			BatchSize   int64
			BlockTime   time.Duration
			MaxRetries  int64
			BackoffTime time.Duration
		},
	) ISubscriberGateway
}

type Gateways struct {
	openAIGateway                func() IAIGateway
	groqGateway                  func() IAIGateway
	redisCacheGateway            func() ICacheGateway
	redisStreamPublisherGateway  func() IPublisherGateway
	redisStreamSubscriberGateway func(opts *struct {
		BatchSize   int64
		BlockTime   time.Duration
		MaxRetries  int64
		BackoffTime time.Duration
	}) ISubscriberGateway
}

func GetGateways(
	env iEnv,
	clients clients.IClients,
) IGateways {
	newOpenAIGateway := func() IAIGateway {
		return openai.NewOpenAIGateway(
			env.GetLogger("openai-gateway"),
			env.GetTracer("openai-gateway"),
			clients.OpenAIClient(),
		)
	}
	newGroqGateway := func() IAIGateway {
		return groq.NewGroqGateway(
			env.GetLogger("groq-gateway"),
			env.GetTracer("groq-gateway"),
			clients.GroqClient(),
		)
	}
	newRedisCacheGateway := func() ICacheGateway {
		return redis.NewRedisCacheGateway(
			env.GetLogger("redis-cache-gateway"),
			env.GetTracer("redis-cache-gateway"),
			clients.RedisClient(),
		)
	}
	newRedisStreamPublisherGateway := func() IPublisherGateway {
		return redis.NewRedisStreamPublisherGateway(
			env.GetLogger("redis-stream-publisher-gateway"),
			env.GetTracer("redis-stream-publisher-gateway"),
			clients.RedisClient(),
		)
	}
	newRedisStreamSubscriberGateway := func(
		opts *struct {
			BatchSize   int64
			BlockTime   time.Duration
			MaxRetries  int64
			BackoffTime time.Duration
		},
	) ISubscriberGateway {
		return redis.NewRedisStreamSubscriberGateway(
			env.GetLogger("redis-stream-subscriber-gateway"),
			env.GetTracer("redis-stream-subscriber-gateway"),
			clients.RedisClient(),
			opts,
		)
	}
	return &Gateways{
		openAIGateway:                newOpenAIGateway,
		redisCacheGateway:            newRedisCacheGateway,
		groqGateway:                  newGroqGateway,
		redisStreamPublisherGateway:  newRedisStreamPublisherGateway,
		redisStreamSubscriberGateway: newRedisStreamSubscriberGateway,
	}
}

func (g *Gateways) OpenAIGateway() IAIGateway {
	return g.openAIGateway()
}

func (g *Gateways) GroqGateway() IAIGateway {
	return g.groqGateway()
}

func (g *Gateways) RedisCacheGateway() ICacheGateway {
	return g.redisCacheGateway()
}

func (g *Gateways) RedisStreamPublisherGateway() IPublisherGateway {
	return g.redisStreamPublisherGateway()
}

func (g *Gateways) RedisStreamSubscriberGateway(
	opts *struct {
		BatchSize   int64
		BlockTime   time.Duration
		MaxRetries  int64
		BackoffTime time.Duration
	},
) ISubscriberGateway {
	return g.redisStreamSubscriberGateway(opts)
}

type iEnv interface {
	GetLogger(name string) utils.ILogger
	GetTracer(service string) apm.ITracer
}

type IAIGateway interface {
	GenerateContent(
		ctx context.Context,
		messages []llms.MessageContent,
		opts ...llms.CallOption,
	) (*llms.ContentResponse, error)
}

type ICacheGateway interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteAll(ctx context.Context) error
}

type IPublisherGateway interface {
	Publish(ctx context.Context, stream string, message []byte) error
}

type ISubscriberGateway interface {
	Subscribe(
		ctx context.Context,
		stream, group, name string,
		handler func(context.Context, []byte) error,
	) error
}
