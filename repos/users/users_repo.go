package users

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	reposdomain "github.com/coopersmall/subswag/repos/domain"
	"github.com/coopersmall/subswag/utils"
)

type UsersRepo struct {
	*reposdomain.SharedRepo[user.UserID, *user.User, db.User]
}

func NewUsersRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
) *UsersRepo {
	return &UsersRepo{
		SharedRepo: reposdomain.NewSharedRepo[user.UserID, *user.User, db.User](
			"user",
			querier,
			tracer,
			convertRowToUser,
			convertUserToRow,
			func(ctx context.Context, iqro db.ISharedQueriesReadOnly, ui user.UserID) (db.User, error) {
				result, err := iqro.GetUser(ctx, int64(ui))
				if isEmptyResult(result) {
					return db.User{}, utils.NewNotFoundError("user not found")
				}
				return result, err
			},
			func(ctx context.Context, iqro db.ISharedQueriesReadOnly) ([]db.User, error) {
				return iqro.GetAllUsers(ctx)
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, u db.User) (sql.Result, error) {
				return iqrw.CreateUser(ctx, db.CreateUserParams{
					ID:        u.ID,
					CreatedAt: u.CreatedAt,
					Data:      u.Data,
				})
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, u db.User) (sql.Result, error) {
				return iqrw.UpdateUser(ctx, db.UpdateUserParams{
					ID: u.ID,
					UpdatedAt: sql.NullTime{
						Time:  u.UpdatedAt.Time,
						Valid: u.UpdatedAt.Valid,
					},
					Data: u.Data,
				})
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, ui user.UserID) (sql.Result, error) {
				return iqrw.DeleteUser(ctx, int64(ui))
			},
		),
	}
}

func isEmptyResult(result db.User) bool {
	return result.ID == 0
}

func convertRowToUser(result db.User) (*user.User, error) {
	var data user.UserData
	err := utils.Unmarshal(result.Data, &data)
	return &user.User{
		ID:       user.UserID(result.ID),
		UserData: data,
		Metadata: &domain.Metadata{
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt.Time,
		},
	}, err
}

func convertUserToRow(user *user.User) (db.User, error) {
	data, err := utils.Marshal(user.UserData)
	return db.User{
		ID:        int64(user.ID),
		CreatedAt: user.Metadata.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  user.Metadata.UpdatedAt,
			Valid: !user.Metadata.UpdatedAt.IsZero(),
		},
		Data: json.RawMessage(data),
	}, err
}
