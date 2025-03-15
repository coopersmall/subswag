package redis

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/clients"
	"github.com/coopersmall/subswag/clients/redis"
	"github.com/coopersmall/subswag/utils"
)

type RedisStreamSubscriberGateway struct {
	logger      utils.ILogger
	tracer      apm.ITracer
	client      clients.IRedisClient
	consumerId  utils.ID
	batchSize   int64
	blockTime   time.Duration
	maxRetries  int64
	backoffTime time.Duration
}

func NewRedisStreamSubscriberGateway(
	logger utils.ILogger,
	tracer apm.ITracer,
	client clients.IRedisClient,
	opts *struct {
		BatchSize   int64
		BlockTime   time.Duration
		MaxRetries  int64
		BackoffTime time.Duration
	},
) *RedisStreamSubscriberGateway {
	gateway := &RedisStreamSubscriberGateway{
		logger:      logger,
		tracer:      tracer,
		client:      client,
		consumerId:  utils.NewID(),
		batchSize:   defaultBatchSize,
		blockTime:   defaultBlockTime,
		maxRetries:  defaultMaxRetries,
		backoffTime: defaultBackoffTime,
	}
	gateway.setOptions((*options)(opts))
	return gateway
}

func (r *RedisStreamSubscriberGateway) Subscribe(
	ctx context.Context,
	key, group, name string,
	handler func(context.Context, []byte) error,
) error {
	if err := r.ensureConsumerGroup(ctx, key, group, name); err != nil {
		r.logger.Error(ctx, "error ensuring consumer group", err, nil)
	}
	return r.processMessages(ctx, key, group, name, handler)
}

func (r *RedisStreamSubscriberGateway) ensureConsumerGroup(ctx context.Context, key, group, name string) error {
	var err error
	r.tracer.Trace(ctx, "creating-group", func(ctx context.Context, span apm.ISpan) error {
		err := r.client.XGroupCreateMkStream(ctx, key, group)
		if !isConsumerGroupExistsError(err) {
			return utils.NewInternalError("failed to read group", err)
		}

		err = r.client.XGroupCreateConsumer(ctx, key, group, name)
		if !isConsumerExistsError(err) {
			return utils.NewInternalError("failed to create consumer", err)
		}

		return nil
	})
	return err
}

func (r *RedisStreamSubscriberGateway) processMessages(
	ctx context.Context,
	key, group, name string,
	onMessage func(context.Context, []byte) error,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := r.readAndProcessBatch(ctx, key, group, name, onMessage); err != nil {
				r.logger.Error(ctx, "error processing batch", err, nil)
				time.Sleep(r.backoffTime)
			}
		}
	}
}

func (r *RedisStreamSubscriberGateway) readAndProcessBatch(
	ctx context.Context,
	key, group, name string,
	onMessage func(context.Context, []byte) error,
) error {
	messages, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: name,
		Streams:  []string{key, ">"},
		Block:    r.blockTime,
		Count:    r.batchSize,
	})
	if err != nil {
		return utils.ErrorOrNil("failed to read group", utils.NewInternalError, err)
	}
	if len(messages) == 0 {
		return nil
	}
	err = r.processBatch(ctx, messages, key, group, onMessage)
	return err
}

func (r *RedisStreamSubscriberGateway) processBatch(
	ctx context.Context,
	messages []redis.XStream,
	key string,
	consumerGroup string,
	onMessage func(context.Context, []byte) error,
) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(messages))

	for _, message := range messages {
		for _, msg := range message.Messages {
			wg.Add(1)
			go func(messageId string, data any) {
				defer wg.Done()

				err := r.processMessageWithRetry(
					ctx,
					messageId,
					data,
					onMessage,
				)
				if err != nil {
					errChan <- err
					return
				}

				// Acknowledge message after successful processing
				if ackErr := r.client.XAck(ctx, key, consumerGroup, messageId); ackErr != nil {
					r.logger.Error(ctx, "failed to acknowledge message", ackErr, map[string]interface{}{
						"messageID": messageId,
					})
				}
			}(msg.ID, msg.Values["data"])
		}
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return utils.NewInternalError("failed to process some messages", errors...)
	}

	return nil
}

func (r *RedisStreamSubscriberGateway) processMessageWithRetry(
	ctx context.Context,
	messageId string,
	data any,
	onMessage func(context.Context, []byte) error,
) error {
	var err error
	for attempt := 0; int64(attempt) < r.maxRetries; attempt++ {
		r.tracer.Trace(ctx, "processing-message", func(ctx context.Context, span apm.ISpan) error {
			span.SetAttribute("message_id", messageId)
			span.SetAttribute("attempt", attempt)
			str := fmt.Sprintf("%v", data)
			return onMessage(ctx, []byte(str))
		})
		if err == nil {
			return nil
		}
		time.Sleep(r.backoffTime * time.Duration(attempt+1))
	}
	return err
}

type options struct {
	BatchSize   int64
	BlockTime   time.Duration
	MaxRetries  int64
	BackoffTime time.Duration
}

func (r *RedisStreamSubscriberGateway) setOptions(opts *options) {
	if opts == nil {
		return
	}
	r.batchSize = defaultInt(opts.BatchSize, r.batchSize)
	r.blockTime = defaultDuration(opts.BlockTime, r.blockTime)
	r.maxRetries = defaultInt(opts.MaxRetries, r.maxRetries)
	r.backoffTime = defaultDuration(opts.BackoffTime, r.backoffTime)
}

const (
	defaultBatchSize   = 10
	defaultBlockTime   = 1 * time.Second
	defaultMaxRetries  = 3
	defaultBackoffTime = 100 * time.Millisecond
)

func defaultInt(value, defaultValue int64) int64 {
	if value == 0 {
		return defaultValue
	}
	return value
}

func defaultDuration(value, defaultValue time.Duration) time.Duration {
	if value == 0 {
		return defaultValue
	}
	return value
}

func isConsumerGroupExistsError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "BUSYGROUP")
}

func isConsumerExistsError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Consumer Group name already exists")
}
