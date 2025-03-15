package secrets

import (
	"context"
	"database/sql"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/secret"
	"github.com/coopersmall/subswag/domain/user"
	reposdomain "github.com/coopersmall/subswag/repos/domain"
	"github.com/coopersmall/subswag/utils"
)

type SecretsRepo struct {
	*reposdomain.StandardRepo[secret.SecretID, *secret.StoredSecret, db.Secret]
}

func NewSecretsRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
	userId user.UserID,
) *SecretsRepo {
	return &SecretsRepo{
		StandardRepo: reposdomain.NewStandardRepo[secret.SecretID, *secret.StoredSecret, db.Secret](
			"secret",
			querier,
			tracer,
			userId,
			convertRowToSecret,
			convertSecretToRow,
			isEmptySecret,
			func(ctx context.Context, iqro db.IStandardQueriesReadOnly, si secret.SecretID) (db.Secret, error) {
				return iqro.GetSecret(ctx, db.GetSecretParams{
					ID:     int64(si),
					UserID: int64(userId),
				})
			},
			func(ctx context.Context, iqro db.IStandardQueriesReadOnly) ([]db.Secret, error) {
				return iqro.GetAllSecrets(ctx)
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, s db.Secret) (sql.Result, error) {
				return iqrw.CreateSecret(ctx, db.CreateSecretParams{
					ID:        s.ID,
					UserID:    int64(userId),
					CreatedAt: s.CreatedAt,
					Data:      s.Data,
				})
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, s db.Secret) (sql.Result, error) {
				return iqrw.UpdateSecret(ctx, db.UpdateSecretParams{
					ID:        s.ID,
					UserID:    int64(userId),
					UpdatedAt: s.UpdatedAt,
					Data:      s.Data,
				})
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, si secret.SecretID) (sql.Result, error) {
				return iqrw.DeleteSecret(ctx, db.DeleteSecretParams{
					ID:     int64(si),
					UserID: int64(userId),
				})
			},
		),
	}
}

func isEmptySecret(secret *secret.StoredSecret) bool {
	return secret == nil || secret.ID == 0
}

func convertRowToSecret(result db.Secret) (*secret.StoredSecret, error) {
	var data secret.StoredSecretData
	err := utils.Unmarshal([]byte(result.Data), &data)
	return &secret.StoredSecret{
		SecretBase: secret.SecretBase{
			ID: secret.SecretID(result.ID),
			Metadata: &domain.Metadata{
				CreatedAt: result.CreatedAt,
				UpdatedAt: result.UpdatedAt.Time,
			},
		},
		Type:             secret.SecretTypeStored,
		StoredSecretData: data,
	}, err
}

func convertSecretToRow(secret *secret.StoredSecret) (db.Secret, error) {
	data, err := utils.Marshal(secret.StoredSecretData)
	return db.Secret{
		ID:        int64(secret.ID),
		CreatedAt: secret.Metadata.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  secret.Metadata.UpdatedAt,
			Valid: !secret.Metadata.UpdatedAt.IsZero(),
		},
		Data: data,
	}, err
}
