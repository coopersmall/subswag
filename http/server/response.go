package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coopersmall/subswag/domain"
)

func writeSuccessfulResponse(
	correlationId domain.CorrelationID,
	method string,
	data []byte,
	w http.ResponseWriter,
) {
	repsonseCode := http.StatusOK
	switch method {
	case http.MethodPost:
		repsonseCode = http.StatusCreated
	case http.MethodPut:
		repsonseCode = http.StatusAccepted
	case http.MethodPatch:
		repsonseCode = http.StatusAccepted
	case http.MethodDelete:
		repsonseCode = http.StatusAccepted
	}

	size := len(data)
	if size > 0 {
		w.Header().Set("Content-Type", "application/json")
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	w.Header().Set("Date", time.Now().Format(time.RFC3339))
	w.Header().Set(CorrelationIDHeader, fmt.Sprintf("%d", correlationId))
	w.WriteHeader(repsonseCode)
	w.Write(data)
}

type errorType string

const (
	invalidRequest  errorType = "Invalid Request"
	notFound        errorType = "Not Found"
	unauthorized    errorType = "Unauthorized"
	tooManyRequests errorType = "Too Many Requests"
)

func writeErrorResponse(
	correlationId domain.CorrelationID,
	errorType errorType,
	body *json.Decoder,
	w http.ResponseWriter,
) {
	data, err := json.Marshal(body)
	if err != nil {
		internalServerError(w)
		return
	}

	responseCode := http.StatusInternalServerError
	switch errorType {
	case invalidRequest:
		responseCode = http.StatusBadRequest
		break
	case notFound:
		responseCode = http.StatusNotFound
		break
	case unauthorized:
		responseCode = http.StatusUnauthorized
		break
	case tooManyRequests:
		responseCode = http.StatusTooManyRequests
		break
	}

	size := len(data)
	if size > 0 {
		w.Header().Set("Content-Type", "application/json")
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	w.Header().Set("Date", time.Now().Format(time.RFC3339))
	w.Header().Set(CorrelationIDHeader, fmt.Sprintf("%d", correlationId))
	w.WriteHeader(responseCode)
	w.Write(data)
}
