package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/utils"
)

type StandardCache[ID any, Data any] struct {
	name          string
	logger        utils.ILogger
	tracer        apm.ITracer
	gateway       gateways.ICacheGateway
	itemTTL       time.Duration
	formatItemKey func(ID) string
}

func NewStandardCache[ID any, Data any](
	name string,
	logger utils.ILogger,
	tracer apm.ITracer,
	gateway gateways.ICacheGateway,
	itemTTL time.Duration,
	formatItemKey func(ID) string,
) *StandardCache[ID, Data] {
	return &StandardCache[ID, Data]{
		name:          name,
		logger:        logger,
		tracer:        tracer,
		gateway:       gateway,
		itemTTL:       itemTTL,
		formatItemKey: formatItemKey,
	}
}

func (s *StandardCache[ID, Data]) Get(ctx context.Context, id ID, onMiss func(context.Context, ID) (Data, error)) (Data, error) {
	var (
		data Data
		err  error
	)
	s.tracer.Trace(ctx, s.name+".get", func(ctx context.Context, span apm.ISpan) error {
		key := s.formatItemKey(id)
		var found []byte
		found, err = s.gateway.Get(ctx, key)
		if err != nil {
			err = utils.NewInternalError("failed to get item", err)
			return err
		}
		if found != nil && len(found) > 0 {
			span.AddEvent("cache hit")
			err = json.Unmarshal(found, &data)
			if err != nil {
				err = utils.NewJSONMarshError("failed to unmarshal item", err)
				return err
			}
			return nil
		}
		span.AddEvent("cache miss")
		data, err = onMiss(ctx, id)
		if err != nil {
			err = utils.NewInternalError("failed to get item", err)
			return err
		}
		err = s.Set(ctx, id, data)
		return err
	})
	return data, err
}

func (s *StandardCache[ID, Data]) Set(ctx context.Context, id ID, data Data) error {
	var err error
	s.tracer.Trace(ctx, s.name+".set", func(ctx context.Context, span apm.ISpan) error {
		key := s.formatItemKey(id)
		var bytes []byte
		bytes, err = json.Marshal(data)
		if err != nil {
			err = utils.NewJSONMarshError("failed to marshal item", err)
			return err
		}
		err = s.gateway.Set(ctx, key, bytes, s.itemTTL)
		return err
	})
	return utils.ErrorOrNil("failed to set item", utils.NewInternalError, err)
}

func (s *StandardCache[ID, Data]) Delete(ctx context.Context, id ID) error {
	var err error
	s.tracer.Trace(ctx, s.name+".delete", func(ctx context.Context, span apm.ISpan) error {
		key := s.formatItemKey(id)
		err = s.gateway.Delete(ctx, key)
		return err
	})
	return utils.ErrorOrNil("failed to delete item", utils.NewInternalError, err)
}
