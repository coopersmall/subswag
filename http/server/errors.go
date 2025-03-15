package server

import (
	"errors"
	"net/http"
	"time"
)

var (
	errEmptyUserIdHeader    = errors.New("user id header is empty")
	errEmptyRequest         = errors.New("request is empty")
	errEmptyRequestHeader   = errors.New("request header is empty")
	errEmptyRequestMethod   = errors.New("request method is empty")
	errEmptyRequestURL      = errors.New("request url is empty")
	errInvalidRequestMethod = errors.New("request method is invalid")
	errInvalidSchema        = errors.New("invalid schema")
	errTooManyRequests      = errors.New("rate limit exceeded")
	errUnauthorized         = errors.New("unauthorized")
)

func internalServerError(w http.ResponseWriter) {
	w.Header().Set(headerContentType, "application/text")
	w.Header().Set(headerDate, time.Now().Format(time.RFC3339))
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	w.Write([]byte("Internal Server Error"))
}
