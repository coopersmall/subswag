package subscribers

import (
	"context"
	"sync"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/utils"
)

func StartSubscribers(
	env iEnv,
	subscribers ISubscribers,
) error {
	correlationID := domain.NewCorrelationID()
	ctx := domain.ContextWithCorrelationID(context.Background(), correlationID)
	logger := env.GetLogger("subscribers")

	var errChan = make(chan error)
	defer close(errChan)

	wg := sync.WaitGroup{}

	s := []ISubscriber{
		subscribers.UserUpdateSubscriber(),
	}

	for _, subscriber := range s {
		go func(logger utils.ILogger) {
			logger.Info(ctx, "Starting subscriber", nil)
			wg.Add(1)
			defer wg.Done()
			err := subscriber.Subscribe(ctx)
			if err != nil {
				errChan <- err
			}
			logger.Info(ctx, "Started subscriber", nil)
		}(logger)
	}

	wg.Wait()
	err := <-errChan
	if err != nil {
		return err
	}
	return nil
}
