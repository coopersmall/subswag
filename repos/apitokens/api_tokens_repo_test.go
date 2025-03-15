package apitokens_test

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ctx           = context.Background()
	userId        = user.UserID(1)
	validUser     = user.NewUser()
	apiTokenId    = apitoken.APITokenID(1)
	validAPIToken = &apitoken.APIToken{
		ID: apiTokenId,
		APITokenData: apitoken.APITokenData{
			UserId:      userId,
			Expiry:      time.Now().Add(12 * time.Hour),
			Permissions: []domain.Permission{domain.APIPermission},
		},
		Metadata: domain.NewMetadata(),
	}
	repo repos.IAPITokenRepo
)

func (s *APITokensRepoTestSuite) SetupSubTest() {
	s.Reset()
	repos, _ := s.GetRepos()
	repo = repos.APITokenRepo(userId)

	validUser.ID = userId
	err := repos.UsersRepo().Create(context.Background(), validUser)
	require.NoError(s.T(), err)
}

func (s *APITokensRepoTestSuite) TestAPITokensRepoSuccess() {
	s.Run("it works", func() {
		err := repo.Create(ctx, validAPIToken)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, apiTokenId)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), apiTokenId, result.ID)
		assert.Len(s.T(), result.Permissions, 1)
		assert.Equal(s.T(), domain.APIPermission, result.Permissions[0])

		apiToken2 := &apitoken.APIToken{}
		utils.DeepClone(validAPIToken, apiToken2)
		apiToken2.ID = 2

		err = repo.Create(ctx, apiToken2)
		assert.NoError(s.T(), err)

		results, err := repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 2)
		assert.Equal(s.T(), apiTokenId, results[0].ID)
		assert.Equal(s.T(), apiToken2.ID, results[1].ID)

		updatedAPIToken := &apitoken.APIToken{}
		utils.DeepClone(validAPIToken, updatedAPIToken)
		updatedAPIToken.Permissions = append(updatedAPIToken.Permissions, domain.AdminPermission)

		err = repo.Update(ctx, updatedAPIToken)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, apiTokenId)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), apiTokenId, result.ID)
		assert.Equal(s.T(), userId, result.APITokenData.UserId)
		assert.Len(s.T(), result.Permissions, 2)
		assert.Equal(s.T(), domain.APIPermission, result.Permissions[0])
		assert.Equal(s.T(), domain.AdminPermission, result.Permissions[1])

		err = repo.Delete(ctx, apiTokenId)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 1)

		err = repo.Delete(ctx, apiToken2.ID)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 0)
	})
}

func (s *APITokensRepoTestSuite) TestAPITokensRepoFailure() {
	s.Run("it throws an error for non-existent user", func() {
		repos, _ := s.GetRepos()
		err := repos.UsersRepo().Delete(ctx, userId)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, apiTokenId)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)

		err = repo.Create(ctx, validAPIToken)
		assert.Error(s.T(), err)
	})

	s.Run("it throws an error for non-existent token", func() {
		result, err := repo.Get(ctx, apiTokenId)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
	})

	s.Run("it returns an error when creating a duplicate token", func() {
		err := repo.Create(ctx, validAPIToken)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validAPIToken)
		assert.Error(s.T(), err)
	})
}
