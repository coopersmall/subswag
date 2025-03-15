package ratelimiter

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/domain/ratelimit"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/utils"
)

const (
	DefaultMaxRequests = 10
	DefaultDuration    = 5 * time.Second // Use time.Duration for clearer representation
)

type RateLimiterService struct {
	logger        utils.ILogger
	rateLimitRepo rateLimitRepo
}

func NewRateLimiterService(
	logger utils.ILogger,
	rateLimitRepo rateLimitRepo,
) *RateLimiterService {
	return &RateLimiterService{
		logger:        logger,
		rateLimitRepo: rateLimitRepo,
	}
}

func (b *RateLimiterService) IsRateLimited(ctx context.Context, userID user.UserID) (bool, error) {
	// rateLimit, err := b.rateLimitRepo.GetRateLimit(userID)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// now := time.Now()
	//
	// if rateLimit == nil {
	// 	// Initialize a new rate limit if none exists
	// 	rateLimit = &domain.RateLimit{
	// 		MaxRequests:         DefaultMaxRequests,
	// 		ResetDuration:       DefaultDuration,
	// 		LatestRequest:       now,
	// 		CurrentRequestCount: 0,
	// 	}
	// }
	//
	// // Calculate time since the last request
	// timeSinceLastRequest := now.Sub(rateLimit.LatestRequest)
	//
	// // If the time since the last request is greater than the duration, reset the rate limit
	// if timeSinceLastRequest > rateLimit.ResetDuration {
	// 	rateLimit.CurrentRequestCount = 0
	// 	rateLimit.LatestRequest = now
	// }
	//
	// // Calculate the reset time
	// resetTime := rateLimit.ResetDuration - timeSinceLastRequest
	//
	// // Check if the current request count exceeds the max requests allowed
	// if rateLimit.CurrentRequestCount >= rateLimit.MaxRequests {
	// 	return &domain.RateLimitResponse{
	// 		RemainingRequests: 0,
	// 		ResetDuration:     resetTime,
	// 	}, nil
	// }
	//
	// // Increment the request count and update the latest request time
	// rateLimit.CurrentRequestCount++
	// rateLimit.LatestRequest = now
	//
	// // Save the updated rate limit back to the repository
	// if err := b.rateLimitRepo.SetRateLimit(userID, rateLimit); err != nil {
	// 	return nil, err
	// }
	//
	// return &domain.RateLimitResponse{
	// 	RemainingRequests: rateLimit.MaxRequests - rateLimit.CurrentRequestCount,
	// 	ResetDuration:     resetTime,
	// }, nil
	return false, nil
}

type rateLimitRepo interface {
	Get(context.Context, ratelimit.RateLimitID) (*ratelimit.RateLimit, error)
}
