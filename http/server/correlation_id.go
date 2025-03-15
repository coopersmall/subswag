package server

import (
	"net/http"

	"github.com/coopersmall/subswag/domain"
)

func GetCorrelationIdFromRequest(r *http.Request) domain.CorrelationID {
	correlationID, ok := r.Context().Value(CorrelationIDKey).(domain.CorrelationID)
	if ok {
		return correlationID
	}
	if parsed := r.Header.Get(CorrelationIDHeader); parsed != "" {
		var err error
		correlationID, err = domain.NewCorrelationIDFromString(parsed)
		if err == nil {
			return correlationID
		}
	}
	return domain.NewCorrelationID()
}
