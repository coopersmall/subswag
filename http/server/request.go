package server

import (
	"context"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/services"
)

type IRequest interface {
	Ctx() context.Context
	CorrelationID() domain.CorrelationID
	UserID() user.UserID
	Body() ([]byte, error)
	Param(key string) (string, error)
	SearchParam(key string) (string, error)
	GetServices() services.IServices
}

type Request struct {
	ctx            context.Context
	correlationID  domain.CorrelationID
	userID         user.UserID
	getBody        func() ([]byte, error)
	getParam       func(key string) (string, error)
	getSearchParam func(key string) (string, error)
	getServices    func() services.IServices
}

func (r *Request) Ctx() context.Context {
	return r.ctx
}

func (r *Request) CorrelationID() domain.CorrelationID {
	return r.correlationID
}

func (r *Request) UserID() user.UserID {
	return r.userID
}

func (r *Request) Body() ([]byte, error) {
	return r.getBody()
}

func (r *Request) Param(key string) (string, error) {
	return r.getParam(key)
}

func (r *Request) SearchParam(key string) (string, error) {
	return r.getSearchParam(key)
}

func (r *Request) GetServices() services.IServices {
	return r.getServices()
}
