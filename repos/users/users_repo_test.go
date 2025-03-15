package users_test

import (
	"context"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"

	"github.com/stretchr/testify/assert"
)

var (
	ctx       = context.Background()
	userId    = user.UserID(1)
	validUser = &user.User{
		ID: userId,
		UserData: user.UserData{
			Email:     "test@example.com",
			FirstName: "Jon",
			LastName:  "Doe",
		},
		Metadata: domain.NewMetadata(),
	}
	repo repos.IUsersRepo
)

func (s *UsersRepoTestSuite) SetupSubTest() {
	s.Reset()
	repos, _ := s.GetRepos()
	repo = repos.UsersRepo()
}

func (s *UsersRepoTestSuite) TestUsersRepoSuccess() {
	s.Run("it works", func() {
		result, err := repo.Get(ctx, userId)
		assert.Error(s.T(), err)

		err = repo.Create(ctx, validUser)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, userId)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), userId, result.ID)
		assert.Equal(s.T(), validUser.Email, result.Email)

		user2 := &user.User{}
		utils.DeepClone(validUser, user2)
		user2.ID = 2

		err = repo.Create(ctx, user2)
		assert.NoError(s.T(), err)

		results, err := repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 2)

		updatedUser := &user.User{}
		utils.DeepClone(validUser, updatedUser)
		updatedUser.Email = "another@example.com"

		err = repo.Update(ctx, updatedUser)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, userId)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), userId, result.ID)
		assert.Equal(s.T(), updatedUser.Email, result.Email)

		err = repo.Delete(ctx, userId)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 1)

		err = repo.Delete(ctx, user2.ID)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 0)
	})
}

func (s *UsersRepoTestSuite) TestUsersRepoErrors() {
	s.Run("it returns an error when creating a duplicate user", func() {
		err := repo.Create(ctx, validUser)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validUser)
		assert.Error(s.T(), err)
	})

	s.Run("it throws an error for non-existent user", func() {
		user, err := repo.Get(ctx, userId)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), user)

		err = repo.Delete(ctx, userId)
		assert.Error(s.T(), err)
	})
}
