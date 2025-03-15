package server

import (
	"net/http"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/env"
)

type IHandler interface {
	Routes() []Route
}

type Handler struct {
	path        string
	env         env.IEnv
	routes      []route
	middlewares []Middleware
	permissions []domain.Permission
}

func NewHandler(
	path string,
	env env.IEnv,
	permissions []domain.Permission,
	middlewares []Middleware,
	opts ...HandlerOption,
) *Handler {
	handler := &Handler{
		path:        path,
		env:         env,
		permissions: permissions,
	}

	for _, opt := range opts {
		opt(handler)
	}

	for i, route := range handler.routes {
		handlerFunc := route.handlerFunc
		for _, middleware := range handler.middlewares {
			handler.routes[i].handlerFunc = func(w http.ResponseWriter, r *http.Request) {
				handlerFunc = middleware(handlerFunc)
			}
		}
		handler.routes[i].handlerFunc = handlerFunc
	}

	return handler
}

func (h *Handler) Routes() []Route {
	routes := make([]Route, len(h.routes))
	for i, route := range h.routes {
		routes[i] = route
	}
	return routes
}

type HandlerOption func(*Handler)

func APIGetRoute(
	suffix string,
	routeFunc func(IRequest) (any, error),
	permissions ...domain.Permission,
) HandlerOption {
	return func(h *Handler) {
		perms := append(h.permissions, permissions...)
		route := route{
			method: http.MethodGet,
			path:   h.path + suffix,
			handlerFunc: NewAPIRoute(
				h.path,
				suffix,
				http.MethodGet,
				routeFunc,
				h.env,
				perms...,
			),
		}

		h.routes = append(h.routes, route)
	}
}

func APIPostRoute(
	suffix string,
	routeFunc func(IRequest) (any, error),
	permissions ...domain.Permission,
) HandlerOption {
	return func(h *Handler) {
		perms := append(h.permissions, permissions...)
		route := route{
			method: http.MethodPost,
			path:   h.path + suffix,
			handlerFunc: NewAPIRoute(
				h.path,
				suffix,
				http.MethodPost,
				routeFunc,
				h.env,
				perms...,
			),
		}
		h.routes = append(h.routes, route)
	}
}

func APIPutRoute(
	suffix string,
	routeFunc func(IRequest) (any, error),
	permissions ...domain.Permission,
) HandlerOption {
	return func(h *Handler) {
		perms := append(h.permissions, permissions...)
		route := route{
			method: http.MethodPut,
			path:   h.path + suffix,
			handlerFunc: NewAPIRoute(
				h.path,
				suffix,
				http.MethodPut,
				routeFunc,
				h.env,
				perms...,
			),
		}
		h.routes = append(h.routes, route)
	}
}

func APIDeleteRoute(
	suffix string,
	routeFunc func(IRequest) (any, error),
	permissions ...domain.Permission,
) HandlerOption {
	return func(h *Handler) {
		perms := append(h.permissions, permissions...)
		route := route{
			method: http.MethodDelete,
			path:   h.path + suffix,
			handlerFunc: NewAPIRoute(
				h.path,
				suffix,
				http.MethodDelete,
				routeFunc,
				h.env,
				perms...,
			),
		}
		h.routes = append(h.routes, route)
	}
}

func WithMiddleware(middlewares ...Middleware) HandlerOption {
	return func(h *Handler) {
		h.middlewares = append(h.middlewares, middlewares...)
	}
}
