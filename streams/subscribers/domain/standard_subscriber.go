package domain

import (
	"context"
	"encoding/json"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/gateways"
	streamsdomain "github.com/coopersmall/subswag/streams/domain"
	"github.com/coopersmall/subswag/utils"
)

type StandardStreamSubscriber[Data any] struct {
	group   string
	name    string
	stream  func() string
	logger  utils.ILogger
	tracer  apm.ITracer
	gateway gateways.ISubscriberGateway
	handler func(
		context.Context,
		Data,
	) error
}

func NewStandardStreamSubscriber[Data any](
	group string,
	name string,
	stream func() string,
	logger utils.ILogger,
	tracer apm.ITracer,
	gateway gateways.ISubscriberGateway,
	handler func(
		context.Context,
		Data,
	) error,
) *StandardStreamSubscriber[Data] {
	subscriber := &StandardStreamSubscriber[Data]{
		group:   group,
		stream:  stream,
		logger:  logger,
		tracer:  tracer,
		gateway: gateway,
		handler: handler,
	}
	return subscriber
}

func (s *StandardStreamSubscriber[Data]) Subscribe(ctx context.Context) error {
	return s.gateway.Subscribe(
		ctx,
		s.stream(),
		s.group,
		s.name,
		func(ctx context.Context, data []byte) error {
			var err error
			s.tracer.Trace(ctx, "process", func(ctx context.Context, span apm.ISpan) error {
				var event streamsdomain.Event[Data]
				if err = json.Unmarshal(data, &event); err != nil {
					return utils.NewJSONMarshError("failed to unmarshal message", err)
				}
				span.SetAttribute("event_id", event.ID)
				err = s.handler(ctx, event.Data)
				return utils.ErrorOrNil("failed to handle message", utils.NewInternalError, err)
			})
			return err
		},
	)
}
