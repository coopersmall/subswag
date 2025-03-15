package publishers

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	usersdomain "github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/gateways"
	userspublisher "github.com/coopersmall/subswag/streams/publishers/users"
	"github.com/coopersmall/subswag/utils"
)

type IPublishers interface {
	UsersPublisher() IUsersPublisher
}

type Publishers struct {
	usersPublisher func() IUsersPublisher
}

func GetPublishers(
	env iEnv,
	gateways gateways.IGateways,
) IPublishers {
	newUserUpdateSubscriber := func() IUsersPublisher {
		return userspublisher.NewUsersPublisher(
			env.GetLogger("users-publisher"),
			env.GetTracer("users-publisher"),
			gateways.RedisStreamPublisherGateway(),
		)
	}
	return &Publishers{
		usersPublisher: newUserUpdateSubscriber,
	}
}

func (p *Publishers) UsersPublisher() IUsersPublisher {
	return p.usersPublisher()
}

type iEnv interface {
	GetLogger(name string) utils.ILogger
	GetTracer(service string) apm.ITracer
}

type IUsersPublisher interface {
	PublishUpdate(ctx context.Context, user *usersdomain.User) error
}
