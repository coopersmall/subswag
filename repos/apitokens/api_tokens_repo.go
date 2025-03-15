package apitokens

import (
	"context"
	"database/sql"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	reposdomain "github.com/coopersmall/subswag/repos/domain"
	"github.com/coopersmall/subswag/utils"
)

type APITokenRepo struct {
	*reposdomain.StandardRepo[apitoken.APITokenID, *apitoken.APIToken, db.ApiToken]
}

func NewAPITokenRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
	userId user.UserID,
) *APITokenRepo {
	return &APITokenRepo{
		StandardRepo: reposdomain.NewStandardRepo[apitoken.APITokenID, *apitoken.APIToken, db.ApiToken](
			"apitoken",
			querier,
			tracer,
			userId,
			convertRowToAPIToken,
			convertAPITokenToRow,
			isEmptyAPIToken,
			func(ctx context.Context, iqro db.IStandardQueriesReadOnly, ai apitoken.APITokenID) (db.ApiToken, error) {
				return iqro.GetAPIToken(ctx, db.GetAPITokenParams{
					ID:     int64(ai),
					UserID: int64(userId),
				})
			},
			func(ctx context.Context, iqro db.IStandardQueriesReadOnly) ([]db.ApiToken, error) {
				return iqro.GetAllAPITokens(ctx)
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, at db.ApiToken) (sql.Result, error) {
				return iqrw.CreateAPIToken(ctx, db.CreateAPITokenParams{
					ID:        at.ID,
					UserID:    int64(userId),
					CreatedAt: at.CreatedAt,
					Data:      at.Data,
				})
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, at db.ApiToken) (sql.Result, error) {
				return iqrw.UpdateAPIToken(ctx, db.UpdateAPITokenParams{
					ID:        at.ID,
					UserID:    int64(userId),
					UpdatedAt: at.UpdatedAt,
					Data:      at.Data,
				})
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, ai apitoken.APITokenID) (sql.Result, error) {
				return iqrw.DeleteAPIToken(ctx, db.DeleteAPITokenParams{
					ID:     int64(ai),
					UserID: int64(userId),
				})
			},
		),
	}
}

func isEmptyAPIToken(token *apitoken.APIToken) bool {
	return token == nil || token.ID == 0
}

func convertRowToAPIToken(result db.ApiToken) (*apitoken.APIToken, error) {
	var data apitoken.APITokenData
	err := utils.Unmarshal([]byte(result.Data), &data)
	return &apitoken.APIToken{
		ID:           apitoken.APITokenID(result.ID),
		APITokenData: data,
		Metadata: &domain.Metadata{
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt.Time,
		},
	}, err
}

func convertAPITokenToRow(
	token *apitoken.APIToken,
) (db.ApiToken, error) {
	data, err := utils.Marshal(token.APITokenData)
	return db.ApiToken{
		ID:        int64(token.ID),
		CreatedAt: token.Metadata.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  token.Metadata.UpdatedAt,
			Valid: !token.Metadata.UpdatedAt.IsZero(),
		},
		Data: data,
	}, err
}
