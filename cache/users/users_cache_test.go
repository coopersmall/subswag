package users_test

import (
	"context"
	"fmt"

	c "github.com/coopersmall/subswag/cache"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/stretchr/testify/assert"
)

var (
	ctx       = context.Background()
	userId    = user.UserID(1)
	validUser = user.NewUser()
	cache     c.IUsersCache
	repo      repos.IUsersRepo
)

func (s *UsersCacheTestSuite) SetupSubTest() {
	s.Reset()
	caches, _ := s.GetCaches()
	cache = caches.UsersCache()
	validUser.ID = userId
	repos, _ := s.GetRepos()
	repo = repos.UsersRepo()
}

func (s *UsersCacheTestSuite) TestUsersCacheSuccess() {
	s.Run("it works", func() {
		err := cache.Set(
			ctx,
			userId,
			validUser,
		)
		assert.NoError(s.T(), err)

		user2 := user.NewUser()
		user2.ID = 2

		err = cache.Set(
			ctx,
			user2.ID,
			user2,
		)
		assert.NoError(s.T(), err)

		result, err := cache.Get(ctx, userId, repo.Get)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), userId, result.ID)

		result, err = cache.Get(ctx, user2.ID, repo.Get)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), user2.ID, result.ID)

		err = repo.Create(ctx, validUser)
		assert.NoError(s.T(), err)
		err = repo.Create(ctx, user2)
		assert.NoError(s.T(), err)

		err = cache.Delete(ctx, userId)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, userId, repo.Get)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), userId, result.ID)

		err = repo.Delete(ctx, userId)
		assert.NoError(s.T(), err)

		err = cache.Delete(ctx, userId)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, userId, repo.Get)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)

		err = cache.Delete(ctx, user2.ID)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, user2.ID, repo.Get)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), user2.ID, result.ID)

		err = repo.Delete(ctx, user2.ID)
		assert.NoError(s.T(), err)

		err = cache.Delete(ctx, user2.ID)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, user2.ID, repo.Get)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
	})
}

func (s *UsersCacheTestSuite) TestUsersCacheFailures() {
	s.Run("it returns error from on miss", func() {
		onMiss := func(ctx context.Context, userId user.UserID) (*user.User, error) {
			return nil, fmt.Errorf("on miss error")
		}
		results, err := cache.Get(ctx, userId, onMiss)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), results)
	})

	s.Run("it returns error from non-existent user", func() {
		err := cache.Delete(ctx, userId)
		assert.NoError(s.T(), err)
	})
}
