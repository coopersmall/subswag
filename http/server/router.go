package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type IRouter interface {
	MustStart(vars iEnvVars) func() error
	Start(vars iEnvVars) (func() error, error)
}

type Router struct {
	logger utils.ILogger
	tracer apm.ITracer
	mux    http.Handler
}

func NewRouter(
	path string,
	logger utils.ILogger,
	tracer apm.ITracer,
	middlewares []Middleware,
	handlers ...IHandler,
) *Router {
	mux := mux.NewRouter()
	mux.StrictSlash(true)
	mux.UseEncodedPath()

	subRouter := mux.PathPrefix(path).Subrouter()
	for _, handler := range handlers {
		routes := handler.Routes()
		for _, route := range routes {
			handlerFunc := route.HandlerFunc()
			for _, middleware := range middlewares {
				handlerFunc = middleware(handlerFunc)
			}
			logger.Info(context.Background(), "Adding route", map[string]any{
				"path":   route.Path(),
				"method": route.Method(),
			})
			subRouter.HandleFunc(route.Path(), handlerFunc).Methods(route.Method())
		}
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Your frontend origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(mux)

	return &Router{
		logger: logger,
		tracer: tracer,
		mux:    handler,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	correlationId := GetCorrelationIdFromRequest(req)
	ctx := domain.ContextWithCorrelationID(req.Context(), correlationId)
	r.tracer.Trace(ctx, "handle-request", func(ctx context.Context, span apm.ISpan) error {
		req = req.WithContext(ctx)
		span.SetAttribute("host", req.Host)
		span.SetAttribute("path", req.URL.Path)
		span.SetAttribute("method", req.Method)
		r.mux.ServeHTTP(w, req)
		return nil
	})
}

func (r *Router) MustStart(
	vars iEnvVars,
) func() error {
	closer, err := r.Start(vars)
	if err != nil {
		closer()
		panic(err)
	}
	return closer
}

func (r *Router) Start(
	vars iEnvVars,
) (func() error, error) {
	url, err := vars.GetAPIURL()
	if err != nil {
		return nil, err
	}
	timeout, err := vars.GetAPITimeout()
	if err != nil {
		return nil, err
	}
	server := newServer(r, timeout, timeout)
	listener, err := net.Listen("tcp", url)
	if err != nil {
		return nil, err
	}
	r.logger.Info(context.Background(), "Starting server", map[string]any{
		"url": url,
	})
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				r.logger.Error(context.Background(), "failed to accept connection", err, nil)
			}
			go server.Serve(conn)
		}
	}()
	return listener.Close, nil
}

type iEnvVars interface {
	GetAPIURL() (string, error)
	GetAPITimeout() (time.Duration, error)
}
