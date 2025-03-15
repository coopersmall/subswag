package apitoken_test

import (
	"context"
	"time"

	c "github.com/coopersmall/subswag/cache"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
)

var (
	ctx           = context.Background()
	userId        = user.UserID(1)
	validUser     = user.NewUser()
	apitokenId    = apitoken.APITokenID(1)
	validAPIToken = &apitoken.APIToken{
		ID: apitokenId,
		APITokenData: apitoken.APITokenData{
			UserId:      userId,
			Expiry:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Permissions: []domain.Permission{domain.APIPermission},
		},
		Metadata: &domain.Metadata{
			CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	cache c.IAPITokenCache
	repo  repos.IAPITokenRepo
)

func (s *APITokenCacheTestSuite) SetupSubTest() {
	s.Reset()
	caches, _ := s.GetCaches()
	cache = caches.APITokensCache(userId)
	repos, _ := s.GetRepos()
	repo = repos.APITokenRepo(userId)
	validUser.ID = userId
	err := repos.UsersRepo().Create(ctx, validUser)
	assert.NoError(s.T(), err)

}

func (s *APITokenCacheTestSuite) TestAPITokenCacheSuccess() {
	s.Run("it works", func() {
		err := cache.Set(
			ctx,
			apitokenId,
			validAPIToken,
		)
		assert.NoError(s.T(), err)

		apitoken2 := &apitoken.APIToken{}
		utils.DeepClone(validAPIToken, apitoken2)
		apitoken2.ID = 2

		err = cache.Set(
			ctx,
			apitoken2.ID,
			apitoken2,
		)
		assert.NoError(s.T(), err)

		result, err := cache.Get(ctx, apitokenId, repo.Get)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), apitokenId, result.ID)

		result, err = cache.Get(ctx, apitoken2.ID, repo.Get)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), apitoken2.ID, result.ID)

		err = repo.Create(ctx, validAPIToken)
		assert.NoError(s.T(), err)
		err = repo.Create(ctx, apitoken2)
		assert.NoError(s.T(), err)

		err = cache.Delete(ctx, apitokenId)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, apitokenId, repo.Get)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), apitokenId, result.ID)

		err = repo.Delete(ctx, apitokenId)
		assert.NoError(s.T(), err)

		err = cache.Delete(ctx, apitokenId)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, apitokenId, repo.Get)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)

		err = cache.Delete(ctx, apitoken2.ID)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, apitoken2.ID, repo.Get)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), apitoken2.ID, result.ID)

		err = repo.Delete(ctx, apitoken2.ID)
		assert.NoError(s.T(), err)

		err = cache.Delete(ctx, apitoken2.ID)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, apitoken2.ID, repo.Get)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
	})
}

func (s *APITokenCacheTestSuite) TestAPITokenCacheFailure() {
	s.Run("it returns error from on miss", func() {
		results, err := cache.Get(ctx, apitokenId, repo.Get)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), results)
	})

	s.Run("it returns error from non-existent token", func() {
		err := cache.Delete(ctx, apitokenId)
		assert.NoError(s.T(), err)
	})
}
