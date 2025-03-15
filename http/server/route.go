package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/env"
	"github.com/coopersmall/subswag/services"
	"github.com/coopersmall/subswag/utils"
	"github.com/gorilla/mux"
)

type Route interface {
	Method() string
	Path() string
	HandlerFunc() http.HandlerFunc
}

type route struct {
	method      string
	path        string
	handlerFunc http.HandlerFunc
}

func (r route) Method() string {
	return r.method
}

func (r route) Path() string {
	return r.path
}

func (r route) HandlerFunc() http.HandlerFunc {
	return r.handlerFunc
}

func NewAPIRoute(
	resource string,
	extension string,
	method string,
	routeFunc func(IRequest) (any, error),
	env env.IEnv,
	permissions ...domain.Permission,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := env.GetLogger("api-route")
		ctx := r.Context()
		correlationId := domain.GetCorrelationIDFromContext(ctx)
		logger.Info(ctx, "running api route", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		})

		tracer := env.GetTracer(resource + extension)
		tracer.Trace(ctx, method, func(ctx context.Context, span apm.ISpan) error {
			authService, close := env.GetAuthenticationService()
			defer close()

			userId, apiTokenId, err := authService.AuthenticateToken(
				ctx,
				r.Header.Get("Authorization"),
				time.Time(utils.Now()),
				permissions...,
			)
			if err != nil {
				logger.Error(ctx, "failed to authenticate token", err, nil)
				writeErrorResponse(correlationId, unauthorized, nil, w)
				return err
			}

			srvs, close := env.GetServices()
			defer close()

			fmt.Println("userId: ", userId)
			fmt.Println("apiTokenId: ", apiTokenId)
			_, err = srvs.APITokenService(userId).GetToken(ctx, apiTokenId)
			if err != nil {
				logger.Error(ctx, "failed to get token", err, nil)
				internalServerError(w)
				return err
			}

			// ratelimited, err := ss.RateLimiterService(userId).IsRateLimited(r.Context(), userId)
			// if err != nil {
			// 	internalServerError(w)
			// 	return
			// }
			// if ratelimited {
			// 	writeErrorResponse(correlationId, tooManyRequests, nil, w)
			// 	return
			// }

			var reader io.Reader
			if r.Body != nil {
				reader = io.LimitReader(r.Body, 1048576)
				defer r.Body.Close()
			}

			getBody := func() ([]byte, error) {
				if reader == nil {
					return nil, nil
				}
				bytes, err := io.ReadAll(reader)
				if err != nil {
					return nil, utils.NewJSONMarshError("failed to read request body", err)
				}
				return bytes, nil
			}

			getParam := func(key string) (string, error) {
				value := mux.Vars(r)[key]
				if value == "" {
					return "", utils.NewNotFoundError("param not found")
				}
				return value, nil
			}

			getSearchParam := func(key string) (string, error) {
				value := r.URL.Query().Get(key)
				if value == "" {
					return "", utils.NewNotFoundError("search param not found")
				}
				return value, nil
			}

			response, err := routeFunc(&Request{
				ctx:            ctx,
				correlationID:  correlationId,
				userID:         userId,
				getBody:        getBody,
				getParam:       getParam,
				getSearchParam: getSearchParam,
				getServices: func() services.IServices {
					return srvs
				},
			})
			if err != nil {
				logger.Error(ctx, "failed to run route", err, nil)
				internalServerError(w)
				return err
			}

			data, err := json.Marshal(response)
			if err != nil {
				internalServerError(w)
				return err
			}

			span.SetAttribute("response", string(data))

			writeSuccessfulResponse(correlationId, r.Method, data, w)
			return nil
		})
	}
}
