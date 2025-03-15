package subscribers

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/services"
	userupdate "github.com/coopersmall/subswag/streams/subscribers/users/update"
	"github.com/coopersmall/subswag/utils"
)

type ISubscribers interface {
	UserUpdateSubscriber() ISubscriber
}

type ISubscriber interface {
	Subscribe(context.Context) error
}

type Subscribers struct {
	userUpdateSubscriber func() ISubscriber
}

func GetSubscribers(
	env iEnv,
	gateways gateways.IGateways,
	services services.IServices,
) ISubscribers {
	newUserUpdateSubscriber := func() ISubscriber {
		return userupdate.NewUserUpdateSubscriber(
			env.GetLogger("user-update"),
			env.GetTracer("user-update"),
			gateways.RedisStreamSubscriberGateway(nil),
			services,
		)
	}
	return &Subscribers{
		userUpdateSubscriber: newUserUpdateSubscriber,
	}
}

func (s *Subscribers) UserUpdateSubscriber() ISubscriber {
	return s.userUpdateSubscriber()
}

type iEnv interface {
	GetLogger(name string) utils.ILogger
	GetTracer(service string) apm.ITracer
}
