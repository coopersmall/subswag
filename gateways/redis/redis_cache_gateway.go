package redis

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/clients"
	"github.com/coopersmall/subswag/clients/redis"
	"github.com/coopersmall/subswag/utils"
)

type RedisCacheGateway struct {
	logger      utils.ILogger
	tracer      apm.ITracer
	redisClient clients.IRedisClient
}

func NewRedisCacheGateway(
	logger utils.ILogger,
	tracer apm.ITracer,
	redisClient clients.IRedisClient,
) *RedisCacheGateway {
	return &RedisCacheGateway{
		logger:      logger,
		tracer:      tracer,
		redisClient: redisClient,
	}
}

func (r *RedisCacheGateway) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	var err error
	r.tracer.Trace(ctx, "set-cache", func(ctx context.Context, span apm.ISpan) error {
		err := r.redisClient.Set(ctx, key, string(value), ttl)
		return utils.ErrorOrNil("failed to set cache", utils.NewInternalError, err)
	})
	return err
}

func (r *RedisCacheGateway) Get(ctx context.Context, key string) ([]byte, error) {
	var (
		value string
		err   error
	)
	r.tracer.Trace(ctx, "get-cache", func(ctx context.Context, span apm.ISpan) error {
		value, err = r.redisClient.Get(ctx, key)
		if redis.IsRedisNil(err) {
			err = nil
			return err
		}
		return utils.ErrorOrNil("failed to get cache", utils.NewInternalError, err)
	})
	if value == "" {
		return nil, err
	}
	return []byte(value), err
}

func (r *RedisCacheGateway) Delete(ctx context.Context, key string) error {
	var err error
	r.tracer.Trace(ctx, "delete-cache", func(ctx context.Context, span apm.ISpan) error {
		err := r.redisClient.Del(ctx, key)
		return utils.ErrorOrNil("failed to delete cache", utils.NewInternalError, err)
	})
	return err
}

func (r *RedisCacheGateway) DeleteAll(ctx context.Context) error {
	var err error
	r.tracer.Trace(ctx, "flush-all-cache", func(ctx context.Context, span apm.ISpan) error {
		err := r.redisClient.FlushAll(ctx)
		return utils.ErrorOrNil("failed to flush all cache", utils.NewInternalError, err)
	})
	return err
}
