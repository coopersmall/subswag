package redis

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/utils"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(
	getRedisAddress func() (string, error),
	getRedisPassword func() (string, error),
	getRedisDB func() (int, error),
) *RedisClient {
	address, err := getRedisAddress()
	if err != nil {
		panic(utils.NewInternalError("failed to get Redis address for client", err))
	}
	password, err := getRedisPassword()
	if err != nil {
		panic(utils.NewInternalError("failed to get Redis password for client", err))
	}
	db, err := getRedisDB()
	if err != nil {
		panic(utils.NewInternalError("failed to get Redis DB for client", err))
	}
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return &RedisClient{
		client: client,
	}
}

// key-value
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) FlushAll(ctx context.Context) error {
	err := r.client.FlushAll(ctx).Err()
	if err != nil {
		return err
	}
	return nil
}

// stream
func (r *RedisClient) XGroupCreateConsumer(ctx context.Context, stream, group, consumer string) error {
	err := r.client.XGroupCreateConsumer(ctx, stream, group, consumer).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) XAdd(ctx context.Context, stream string, value map[string]any) (string, error) {
	id, err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: value,
	}).Result()
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *RedisClient) XReadGroup(
	ctx context.Context,
	args *XReadGroupArgs,
) ([]XStream, error) {
	streamsData, err := r.client.XReadGroup(ctx, (*redis.XReadGroupArgs)(args)).Result()
	if err == redis.Nil {
		// No messages available
		return []XStream{}, nil
	}
	if err != nil {
		return nil, err
	}
	result := make([]XStream, len(streamsData))
	for i, streamData := range streamsData {
		result[i] = XStream(streamData)
	}
	return result, nil
}

func (r *RedisClient) XGroupCreateMkStream(
	ctx context.Context,
	stream, group string,
) error {
	err := r.client.XGroupCreateMkStream(ctx, stream, group, "$").Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) XAck(
	ctx context.Context,
	stream, group string,
	ids ...string,
) error {
	err := r.client.XAck(ctx, stream, group, ids...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

type XReadGroupArgs redis.XReadGroupArgs
type XStream redis.XStream

func IsRedisNil(err error) bool {
	return err == redis.Nil
}
