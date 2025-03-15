package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/services"
	subscriberdomain "github.com/coopersmall/subswag/streams/subscribers/domain"
	"github.com/coopersmall/subswag/utils"
)

type UserUpdateSubscriber struct {
	logger utils.ILogger
	tracer apm.ITracer
	*subscriberdomain.StandardStreamSubscriber[*user.User]
	services services.IServices
}

func NewUserUpdateSubscriber(
	logger utils.ILogger,
	tracer apm.ITracer,
	gateway gateways.ISubscriberGateway,
	services services.IServices,
) *UserUpdateSubscriber {
	subscriber := &UserUpdateSubscriber{
		logger:   logger,
		tracer:   tracer,
		services: services,
	}
	subscriber.StandardStreamSubscriber = subscriberdomain.NewStandardStreamSubscriber[*user.User](
		"users",
		"user-update",
		stream,
		logger,
		tracer,
		gateway,
		subscriber.Handle,
	)
	return subscriber
}

func (s *UserUpdateSubscriber) Handle(
	ctx context.Context,
	data *user.User,
) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println("UserUpdateSubscriber")
	fmt.Println(string(bytes))
	return nil
}

func stream() string {
	return "users"
}
