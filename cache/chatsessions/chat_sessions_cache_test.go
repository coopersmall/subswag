package chatsessions_test

import (
	"context"
	"time"

	c "github.com/coopersmall/subswag/cache"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
)

var (
	ctx              = context.Background()
	userId           = user.UserID(1)
	validUser        = user.NewUser()
	chatSessionId    = chatsession.ChatSessionID(1)
	validChatSession = &chatsession.ChatSession{
		ID: chatSessionId,
		ChatSessionData: chatsession.ChatSessionData{
			ChatSessionItemIDs: []chatsession.ChatSessionItemID{1, 2, 3},
			UserIDs:            []user.UserID{1, 2},
		},
		Metadata: &domain.Metadata{
			CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	cache c.IChatSessionsCache
	repo  repos.IChatSessionsRepo
)

func (s *ChatSessionsCacheTestSuite) SetupSubTest() {
	s.Reset()
	caches, _ := s.GetCaches()
	cache = caches.ChatSessionsCache()
	repos, _ := s.GetRepos()
	repo = repos.ChatSessionsRepo(userId)
	validUser.ID = userId
	err := repos.UsersRepo().Create(ctx, validUser)
	assert.NoError(s.T(), err)
}

func (s *ChatSessionsCacheTestSuite) TestChatSessionsCacheSuccess() {
	s.Run("it works", func() {
		err := cache.Set(
			ctx,
			chatSessionId,
			validChatSession,
		)
		assert.NoError(s.T(), err)

		chatSession2 := &chatsession.ChatSession{}
		utils.DeepClone(validChatSession, chatSession2)
		chatSession2.ID = 2

		err = cache.Set(
			ctx,
			chatSession2.ID,
			chatSession2,
		)
		assert.NoError(s.T(), err)

		result, err := cache.Get(ctx, chatSessionId, repo.Get)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), chatSessionId, result.ID)

		result, err = cache.Get(ctx, chatSession2.ID, repo.Get)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), chatSession2.ID, result.ID)

		err = repo.Create(ctx, validChatSession)
		assert.NoError(s.T(), err)
		err = repo.Create(ctx, chatSession2)
		assert.NoError(s.T(), err)

		err = cache.Delete(ctx, chatSessionId)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, chatSessionId, repo.Get)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), chatSessionId, result.ID)

		err = repo.Delete(ctx, chatSessionId)
		assert.NoError(s.T(), err)

		err = cache.Delete(ctx, chatSessionId)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, chatSessionId, repo.Get)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)

		err = cache.Delete(ctx, chatSession2.ID)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, chatSession2.ID, repo.Get)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), chatSession2.ID, result.ID)

		err = repo.Delete(ctx, chatSession2.ID)
		assert.NoError(s.T(), err)

		err = cache.Delete(ctx, chatSession2.ID)
		assert.NoError(s.T(), err)

		result, err = cache.Get(ctx, chatSession2.ID, repo.Get)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
	})
}

func (s *ChatSessionsCacheTestSuite) TestChatSessionsCacheFailure() {
	s.Run("it fails with on miss error", func() {
		results, err := cache.Get(ctx, chatSessionId, repo.Get)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), results)
	})
}
