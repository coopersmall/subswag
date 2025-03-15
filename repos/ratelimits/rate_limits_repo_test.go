package ratelimits_test

import (
	"context"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/ratelimit"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ctx                = context.Background()
	userId             = user.UserID(1)
	rateLimitId        = ratelimit.RateLimitID(1)
	validUser          = user.NewUser()
	validRateLimitData = ratelimit.RateLimitData{}
	validRateLimit     = &ratelimit.RateLimit{
		ID:            rateLimitId,
		RateLimitData: validRateLimitData,
		Metadata:      domain.NewMetadata(),
	}
	repo repos.IRateLimitsRepo
)

func (s *RateLimitsRepoTestSuite) SetupSubTest() {
	s.Reset()
	repos, _ := s.GetRepos()
	repo = repos.RateLimitsRepo(userId)
	validUser.ID = userId
	err := repos.UsersRepo().Create(ctx, validUser)
	require.NoError(s.T(), err)
}

func (s *RateLimitsRepoTestSuite) TestRateLimitsRepo() {
	s.Run("it works", func() {
		err := repo.Create(ctx, validRateLimit)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, rateLimitId)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), rateLimitId, result.ID)
		assert.Equal(s.T(), validRateLimitData, result.RateLimitData)

		rateLimit2 := &ratelimit.RateLimit{}
		utils.DeepClone(validRateLimit, rateLimit2)
		rateLimit2.ID = 2

		err = repo.Create(ctx, rateLimit2)
		assert.NoError(s.T(), err)

		results, err := repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 2)

		updatedRateLimit := &ratelimit.RateLimit{}
		utils.DeepClone(validRateLimit, updatedRateLimit)
		updatedRateLimit.RateLimitData = ratelimit.RateLimitData{}

		err = repo.Update(ctx, updatedRateLimit)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), updatedRateLimit.RateLimitData, result.RateLimitData)

		err = repo.Delete(ctx, rateLimitId)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 1)
		assert.Equal(s.T(), rateLimit2.ID, results[0].ID)

		err = repo.Delete(ctx, rateLimit2.ID)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 0)
	})

	s.Run("it deletes rate limit when user is deleted", func() {
		err := repo.Create(ctx, validRateLimit)
		assert.NoError(s.T(), err)

		repos, _ := s.GetRepos()
		err = repos.UsersRepo().Delete(ctx, userId)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, rateLimitId)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
	})
}

func (s *RateLimitsRepoTestSuite) TestRateLimitsRepoFailure() {
	s.Run("it fails for non-existent user", func() {
		repos, _ := s.GetRepos()
		err := repos.UsersRepo().Delete(ctx, userId)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validRateLimit)
		assert.Error(s.T(), err)
	})

	s.Run("it fails for non-existent rate limit", func() {
		result, err := repo.Get(ctx, rateLimitId)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)

		err = repo.Update(ctx, validRateLimit)
		assert.Error(s.T(), err)

		err = repo.Delete(ctx, rateLimitId)
		assert.Error(s.T(), err)
	})

	s.Run("it fails for duplicate rate limit", func() {
		err := repo.Create(ctx, validRateLimit)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validRateLimit)
		assert.Error(s.T(), err)
	})
}
