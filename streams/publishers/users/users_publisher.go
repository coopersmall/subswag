package users

import (
	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/gateways"
	publisherdomain "github.com/coopersmall/subswag/streams/publishers/domain"
	"github.com/coopersmall/subswag/utils"
)

type UsersPublisher struct {
	logger utils.ILogger
	tracer apm.ITracer
	*publisherdomain.StandardStreamPublisher[*user.User]
}

func NewUsersPublisher(
	logger utils.ILogger,
	tracer apm.ITracer,
	gateway gateways.IPublisherGateway,
) *UsersPublisher {
	return &UsersPublisher{
		logger: logger,
		tracer: tracer,
		StandardStreamPublisher: publisherdomain.NewStandardStreamPublisher[*user.User](
			stream,
			logger,
			tracer,
			gateway,
		),
	}
}

func stream() string {
	return "users"
}
