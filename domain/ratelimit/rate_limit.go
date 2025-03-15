package ratelimit

import (
	"github.com/coopersmall/subswag/domain"
)

type RateLimitData struct {
	// Add fields here
}

type RateLimit struct {
	ID               RateLimitID `json:"id" validate:"required,gt=0"`
	RateLimitData    `json:",inline" tstype:",extends"`
	*domain.Metadata `json:"metadata" validate:"required" tstype:"Metadata,required"`
}
