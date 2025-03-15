package server

import (
	"context"
	"net/http"

	"github.com/coopersmall/subswag/domain"
)

// Headers
const (
	headerUserAgent     = "User-Agent"
	headerDate          = "Date"
	headerContentLength = "Content-Length"
	headerContentType   = "Content-Type"
	headerAccept        = "Accept"

	CorrelationIDHeader = "X-Correlation-ID"
	userIDHeaderName    = "X-User-ID"
	rateLimiterHeader   = "X-Rate-Limiter"
	authHeaderName      = "Authorization"
)

// Cookies
const (
	sessionCookieName = "session"
)

// Context keys
const CorrelationIDKey = "correlationId"

func setCorrelationID(ctx context.Context, correlationID domain.CorrelationID) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

func getCorrelationID(ctx context.Context) domain.CorrelationID {
	return ctx.Value(CorrelationIDKey).(domain.CorrelationID)
}

type Middleware func(http.HandlerFunc) http.HandlerFunc
