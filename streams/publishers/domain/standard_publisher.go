package domain

import (
	"context"
	"encoding/json"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/gateways"
	streamsdomain "github.com/coopersmall/subswag/streams/domain"
	"github.com/coopersmall/subswag/utils"
)

type StandardStreamPublisher[Data any] struct {
	stream  func() string
	logger  utils.ILogger
	tracer  apm.ITracer
	gateway gateways.IPublisherGateway
}

func NewStandardStreamPublisher[Data any](
	stream func() string,
	logger utils.ILogger,
	tracer apm.ITracer,
	gateway gateways.IPublisherGateway,
) *StandardStreamPublisher[Data] {
	return &StandardStreamPublisher[Data]{
		stream:  stream,
		logger:  logger,
		tracer:  tracer,
		gateway: gateway,
	}
}

func (s *StandardStreamPublisher[Data]) PublishCreate(ctx context.Context, data Data) error {
	var err error
	s.tracer.Trace(ctx, "publish-create", func(ctx context.Context, span apm.ISpan) error {
		stream := s.stream()
		span.SetAttribute("stream", stream)
		err = s.publish(ctx, data)
		return utils.ErrorOrNil("failed to publish message", utils.NewInternalError, err)
	})
	return err
}

func (s *StandardStreamPublisher[Data]) PublishUpdate(ctx context.Context, data Data) error {
	var err error
	s.tracer.Trace(ctx, "publish-update", func(ctx context.Context, span apm.ISpan) error {
		stream := s.stream()
		span.SetAttribute("stream", stream)
		err = s.publish(ctx, data)
		return utils.ErrorOrNil("failed to publish message", utils.NewInternalError, err)
	})
	return err
}

func (s *StandardStreamPublisher[Data]) PublishDelete(ctx context.Context, data Data) error {
	var err error
	s.tracer.Trace(ctx, "publish-delete", func(ctx context.Context, span apm.ISpan) error {
		stream := s.stream()
		span.SetAttribute("stream", stream)
		err = s.publish(ctx, data)
		return utils.ErrorOrNil("failed to publish message", utils.NewInternalError, err)
	})
	return err
}

func (s *StandardStreamPublisher[Data]) publish(
	ctx context.Context,
	data Data,
) error {
	stream := s.stream()
	event := streamsdomain.NewEvent[Data](stream, data)
	bytes, err := json.Marshal(event)
	if err != nil {
		return utils.NewJSONMarshError("failed to marshal message", err)
	}
	return s.gateway.Publish(ctx, stream, bytes)
}
