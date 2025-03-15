package domain

import (
	"context"
	"strconv"

	"github.com/coopersmall/subswag/utils"
)

type CorrelationID utils.ID

func (c CorrelationID) Validate() error {
	if err := c.Validate(); err != nil {
		return err
	}
	return nil
}

func (c CorrelationID) String() string {
	return utils.ID(c).String()
}

func NewCorrelationID() CorrelationID {
	return CorrelationID(utils.NewID())
}

func NewCorrelationIDFromString(id string) (CorrelationID, error) {
	parsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return CorrelationID(0), err
	}
	return CorrelationID(parsed), nil
}

const correlationIDKey = "correlationID"

func GetCorrelationIDFromContext(ctx context.Context) CorrelationID {
	correlationID, ok := ctx.Value(correlationIDKey).(CorrelationID)
	if ok {
		return correlationID
	}
	return CorrelationID(0)
}

func ContextWithCorrelationID(ctx context.Context, correlationID CorrelationID) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}
