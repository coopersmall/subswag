package secrets_test

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/secret"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
)

var (
	ctx         = context.Background()
	userId      = user.UserID(1)
	secretId    = secret.SecretID(1)
	validUser   = user.NewUser()
	validSecret = &secret.StoredSecret{
		SecretBase: secret.SecretBase{
			ID: secretId,
			Metadata: &domain.Metadata{
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		Type: secret.SecretTypeStored,
		StoredSecretData: secret.StoredSecretData{
			Salt:  []byte("test-salt"),
			Value: []byte("test-value"),
		},
	}
	repo repos.ISecretsRepo
)

func (s *SecretsRepoTestSuite) SetupSubTest() {
	s.Reset()
	repos, _ := s.GetRepos()
	repo = repos.SecretsRepo(userId)
	validUser.ID = userId
	err := repos.UsersRepo().Create(ctx, validUser)
	assert.NoError(s.T(), err)
}

func (s *SecretsRepoTestSuite) TestSecretsRepo() {
	s.Run("it works", func() {
		err := repo.Create(ctx, validSecret)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, secretId)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), secretId, result.ID)
		assert.Equal(s.T(), validSecret.Type, result.Type)
		assert.Equal(s.T(), validSecret.Salt, result.Salt)
		assert.Equal(s.T(), validSecret.Value, result.Value)

		secret2 := &secret.StoredSecret{}
		utils.DeepClone(validSecret, secret2)
		secret2.ID = 2

		err = repo.Create(ctx, secret2)
		assert.NoError(s.T(), err)

		results, err := repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 2)

		updatedSecret := &secret.StoredSecret{}
		utils.DeepClone(validSecret, updatedSecret)
		updatedSecret.Value = []byte("updated-value")

		err = repo.Update(ctx, updatedSecret)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, secretId)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), secretId, result.ID)
		assert.Equal(s.T(), updatedSecret.Value, result.Value)

		err = repo.Delete(ctx, secretId)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, secretId)
		assert.Error(s.T(), err)
		assert.True(s.T(), errors.Is(err, sql.ErrNoRows))
		assert.Nil(s.T(), result)
	})

	s.Run("it deletes secret when user is deleted", func() {
		err := repo.Create(ctx, validSecret)
		assert.NoError(s.T(), err)

		repos, _ := s.GetRepos()
		err = repos.UsersRepo().Delete(ctx, userId)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, secretId)
		assert.Error(s.T(), err)
		assert.True(s.T(), errors.Is(err, sql.ErrNoRows))
		assert.Nil(s.T(), result)
	})
}

func (s *SecretsRepoTestSuite) TestSecretsRepoFailure() {
	s.Run("it fails with non-existant user", func() {
		repos, _ := s.GetRepos()
		err := repos.UsersRepo().Delete(ctx, userId)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validSecret)
		assert.Error(s.T(), err)
	})

	s.Run("it fails with non-existant secret", func() {
		err := repo.Create(ctx, validSecret)
		assert.NoError(s.T(), err)

		err = repo.Delete(ctx, secretId)
		assert.NoError(s.T(), err)

		err = repo.Delete(ctx, secretId)
		assert.Error(s.T(), err)
	})

	s.Run("it fails with duplicate secret", func() {
		err := repo.Create(ctx, validSecret)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validSecret)
		assert.Error(s.T(), err)
	})
}
