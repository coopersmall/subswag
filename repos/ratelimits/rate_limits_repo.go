package ratelimits

import (
	"context"
	"database/sql"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/ratelimit"
	"github.com/coopersmall/subswag/domain/user"
	reposdomain "github.com/coopersmall/subswag/repos/domain"
	"github.com/coopersmall/subswag/utils"
)

type RateLimitsRepo struct {
	*reposdomain.StandardRepo[ratelimit.RateLimitID, *ratelimit.RateLimit, db.RateLimit]
}

func NewRateLimitsRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
	userId user.UserID,
) *RateLimitsRepo {
	return &RateLimitsRepo{
		StandardRepo: reposdomain.NewStandardRepo[ratelimit.RateLimitID, *ratelimit.RateLimit, db.RateLimit](
			"rate_limit",
			querier,
			tracer,
			userId,
			convertRowToRateLimit,
			convertRateLimitToRow,
			isEmptyRateLimit,
			getFunc(userId),
			getAllFunc(),
			createFunc(userId),
			updateFunc(userId),
			deleteFunc(userId),
		),
	}
}

func getFunc(userId user.UserID) func(context.Context, db.IStandardQueriesReadOnly, ratelimit.RateLimitID) (db.RateLimit, error) {
	return func(ctx context.Context, iqro db.IStandardQueriesReadOnly, rli ratelimit.RateLimitID) (db.RateLimit, error) {
		return iqro.GetRateLimit(ctx, db.GetRateLimitParams{
			ID:     int64(rli),
			UserID: int64(userId),
		})
	}
}

func getAllFunc() func(context.Context, db.IStandardQueriesReadOnly) ([]db.RateLimit, error) {
	return func(ctx context.Context, iqro db.IStandardQueriesReadOnly) ([]db.RateLimit, error) {
		return iqro.GetAllRateLimits(ctx)
	}
}

func createFunc(userId user.UserID) func(context.Context, db.IStandardQueriesReadWrite, db.RateLimit) (sql.Result, error) {
	return func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, rl db.RateLimit) (sql.Result, error) {
		return iqrw.CreateRateLimit(ctx, db.CreateRateLimitParams{
			ID:        rl.ID,
			UserID:    int64(userId),
			CreatedAt: rl.CreatedAt,
			Data:      rl.Data,
		})
	}
}

func updateFunc(userId user.UserID) func(context.Context, db.IStandardQueriesReadWrite, db.RateLimit) (sql.Result, error) {
	return func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, rl db.RateLimit) (sql.Result, error) {
		return iqrw.UpdateRateLimit(ctx, db.UpdateRateLimitParams{
			ID:        rl.ID,
			UserID:    int64(userId),
			UpdatedAt: rl.UpdatedAt,
			Data:      rl.Data,
		})
	}
}

func deleteFunc(userId user.UserID) func(context.Context, db.IStandardQueriesReadWrite, ratelimit.RateLimitID) (sql.Result, error) {
	return func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, rli ratelimit.RateLimitID) (sql.Result, error) {
		return iqrw.DeleteRateLimit(ctx, db.DeleteRateLimitParams{
			ID:     int64(rli),
			UserID: int64(userId),
		})
	}
}

func isEmptyRateLimit(rateLimit *ratelimit.RateLimit) bool {
	return rateLimit == nil || rateLimit.ID == 0
}

func convertRowToRateLimit(result db.RateLimit) (*ratelimit.RateLimit, error) {
	var data ratelimit.RateLimitData
	err := utils.Unmarshal([]byte(result.Data), &data)
	return &ratelimit.RateLimit{
		ID:            ratelimit.RateLimitID(result.ID),
		RateLimitData: data,
		Metadata: &domain.Metadata{
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt.Time,
		},
	}, err
}

func convertRateLimitToRow(
	rateLimit *ratelimit.RateLimit,
) (db.RateLimit, error) {
	data, err := utils.Marshal(rateLimit.RateLimitData)
	return db.RateLimit{
		ID:        int64(rateLimit.ID),
		CreatedAt: rateLimit.Metadata.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  rateLimit.Metadata.UpdatedAt,
			Valid: !rateLimit.Metadata.UpdatedAt.IsZero(),
		},
		Data: data,
	}, err
}
