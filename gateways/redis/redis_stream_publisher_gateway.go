package redis

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/clients"
	"github.com/coopersmall/subswag/utils"
)

type RedisStreamPublisherGateway struct {
	logger utils.ILogger
	tracer apm.ITracer
	client clients.IRedisClient
}

func NewRedisStreamPublisherGateway(
	logger utils.ILogger,
	tracer apm.ITracer,
	client clients.IRedisClient,
) *RedisStreamPublisherGateway {
	return &RedisStreamPublisherGateway{
		logger: logger,
		tracer: tracer,
		client: client,
	}
}

func (r *RedisStreamPublisherGateway) Publish(
	ctx context.Context,
	stream string,
	message []byte,
) error {
	var err error
	r.tracer.Trace(ctx, "publish-message", func(ctx context.Context, span apm.ISpan) error {
		span.SetAttribute("stream", stream)
		var id string
		value := map[string]interface{}{
			"data": message,
		}
		id, err = r.client.XAdd(ctx, stream, value)
		if err != nil {
			return utils.NewInternalError("failed to publish message", err)
		}
		span.SetAttribute("message_id", id)
		return nil
	})
	return err
}
