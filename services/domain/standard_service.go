package domain

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/utils"
)

type StandardService[ID any, Data any] struct {
	serviceName string
	logger      utils.ILogger
	tracer      apm.ITracer
	repo        struct {
		Get    func(context.Context, ID) (Data, error)
		All    func(context.Context) ([]Data, error)
		Create func(context.Context, Data) error
		Update func(context.Context, Data) error
		Delete func(context.Context, ID) error
	}
	cache *struct {
		Get    func(context.Context, ID, func(context.Context, ID) (Data, error)) (Data, error)
		Set    func(context.Context, ID, Data) error
		Delete func(context.Context, ID) error
	}
	publish *struct {
		Create func(context.Context, Data) error
		Update func(context.Context, Data) error
		Delete func(context.Context, Data) error
	}
}

func NewStandardService[ID any, Data any](
	serviceName string,
	logger utils.ILogger,
	tracer apm.ITracer,
	repo struct {
		Get    func(context.Context, ID) (Data, error)
		All    func(context.Context) ([]Data, error)
		Create func(context.Context, Data) error
		Update func(context.Context, Data) error
		Delete func(context.Context, ID) error
	},
	cache *struct {
		Get    func(context.Context, ID, func(context.Context, ID) (Data, error)) (Data, error)
		Set    func(context.Context, ID, Data) error
		Delete func(context.Context, ID) error
	},
	publish *struct {
		Create func(context.Context, Data) error
		Update func(context.Context, Data) error
		Delete func(context.Context, Data) error
	},
) *StandardService[ID, Data] {
	return &StandardService[ID, Data]{
		serviceName: serviceName,
		logger:      logger,
		tracer:      tracer,
		repo:        repo,
		cache:       cache,
		publish:     publish,
	}
}

func (s *StandardService[ID, Data]) Get(ctx context.Context, id ID) (Data, error) {
	var (
		found Data
		err   error
	)
	s.tracer.Trace(ctx, "get", func(ctx context.Context, span apm.ISpan) error {
		if s.cache == nil {
			span.AddEvent("cache not found")
			found, err = s.repo.Get(ctx, id)
			return utils.ErrorOrNil("failed to get item", utils.NewInternalError, err)
		}
		span.AddEvent("cache found")
		found, err = s.cache.Get(ctx, id, s.repo.Get)
		return err
	})
	return found, utils.ErrorOrNil("failed to get item", utils.NewInternalError, err)

}

func (s *StandardService[ID, Data]) All(ctx context.Context) ([]Data, error) {
	found, err := s.repo.All(ctx)
	return found, utils.ErrorOrNil("failed to get items", utils.NewInternalError, err)
}

func (s *StandardService[ID, Data]) Create(
	ctx context.Context,
	id ID,
	data Data,
) error {
	var err error
	s.tracer.Trace(ctx, "create", func(ctx context.Context, span apm.ISpan) error {
		err = domain.Validate(data)
		if err != nil {
			return err
		}
		err = s.repo.Create(ctx, data)
		if err != nil {
			return utils.NewInternalError("failed to create item", err)
		}
		if s.cache == nil {
			return nil
		}
		err = s.cache.Set(ctx, id, data)
		if err != nil {
			return utils.NewInternalError("failed to set item in cache", err)
		}
		if s.publish == nil || s.publish.Create == nil {
			return nil
		}
		err = s.publish.Create(ctx, data)
		if err != nil {
			s.logger.Error(ctx, "failed to publish create", err, nil)
		}
		return nil
	})
	return err
}

func (s *StandardService[ID, Data]) Update(
	ctx context.Context,
	id ID,
	data Data,
) error {
	var err error
	s.tracer.Trace(ctx, "update", func(ctx context.Context, span apm.ISpan) error {
		err := domain.Validate(data)
		if err != nil {
			return err
		}
		previousState, err := s.Get(ctx, id)
		if err != nil {
			return utils.NewInternalError("failed to get item", err)
		}
		err = s.repo.Update(ctx, data)
		if err != nil {
			return utils.NewInternalError("failed to update item", err)
		}
		if s.cache == nil {
			return nil
		}
		err = s.cache.Set(ctx, id, data)
		if err != nil {
			err = s.repo.Update(ctx, previousState)
			return utils.NewInternalError("failed to set item in cache", err)
		}
		if s.publish == nil || s.publish.Update == nil {
			return nil
		}
		err = s.publish.Update(ctx, data)
		if err != nil {
			s.logger.Error(ctx, "failed to publish update", err, nil)
		}
		return nil
	})
	return err
}

func (s *StandardService[ID, Data]) Delete(
	ctx context.Context,
	id ID,
) error {
	var err error
	s.tracer.Trace(ctx, "delete", func(ctx context.Context, span apm.ISpan) error {
		var data Data
		if s.publish != nil && s.publish.Delete != nil {
			var err error
			data, err = s.Get(ctx, id)
			if err != nil {
				return utils.NewInternalError("failed to get item", err)
			}
		}
		err := s.repo.Delete(ctx, id)
		if err != nil {
			return utils.NewInternalError("failed to delete item", err)
		}
		if s.cache == nil {
			return nil
		}
		err = s.cache.Delete(ctx, id)
		if err != nil {
			return utils.NewInternalError("failed to delete item from cache", err)
		}
		if s.publish == nil || s.publish.Delete == nil {
			return nil
		}
		err = s.publish.Delete(ctx, data)
		if err != nil {
			s.logger.Error(ctx, "failed to publish delete", err, nil)
		}
		return nil
	})
	return err
}
