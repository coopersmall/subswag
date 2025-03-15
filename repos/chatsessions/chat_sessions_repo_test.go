package chatsessions_test

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
)

var (
	ctx                  = context.Background()
	userId               = user.UserID(1)
	chatSessionId        = chatsession.ChatSessionID(1)
	validUser            = user.NewUser()
	validChatSessionData = chatsession.ChatSessionData{
		ChatSessionItemIDs: []chatsession.ChatSessionItemID{1, 2, 3},
		UserIDs:            []user.UserID{userId},
	}
	validChatSession = &chatsession.ChatSession{
		ID:              chatSessionId,
		ChatSessionData: validChatSessionData,
		Metadata: &domain.Metadata{
			CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	repo repos.IChatSessionsRepo
)

func (s *ChatSessionsRepoTestSuite) SetupSubTest() {
	s.Reset()
	repos, _ := s.GetRepos()
	repo = repos.ChatSessionsRepo(userId)
	validUser.ID = userId
	err := repos.UsersRepo().Create(context.Background(), validUser)
	assert.NoError(s.T(), err)
}

func (s *ChatSessionsRepoTestSuite) TestChatSessionsRepoSuccess() {
	s.Run("it works", func() {
		err := repo.Create(ctx, validChatSession)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, chatSessionId)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), chatSessionId, result.ID)
		assert.Equal(s.T(), validChatSessionData, result.ChatSessionData)

		chatSession2 := &chatsession.ChatSession{}
		utils.DeepClone(validChatSession, chatSession2)
		chatSession2.ID = 2

		err = repo.Create(ctx, chatSession2)
		assert.NoError(s.T(), err)

		results, err := repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 2)

		updatedChatSession := &chatsession.ChatSession{}
		utils.DeepClone(validChatSession, updatedChatSession)
		updatedChatSession.ChatSessionData.ChatSessionItemIDs = []chatsession.ChatSessionItemID{4, 5, 6}

		err = repo.Update(ctx, updatedChatSession)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, chatSessionId)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), chatSessionId, result.ID)
		assert.Equal(s.T(), updatedChatSession.ChatSessionData, result.ChatSessionData)
	})

	s.Run("it deletes chat session when user is deleted", func() {
		err := repo.Create(ctx, validChatSession)
		assert.NoError(s.T(), err)

		repos, _ := s.GetRepos()
		err = repos.UsersRepo().Delete(ctx, userId)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, chatSessionId)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
	})
}

func (s *ChatSessionsRepoTestSuite) TestChatSessionsRepoFailure() {
	s.Run("it fails for non-existent user", func() {
		repos, _ := s.GetRepos()
		err := repos.UsersRepo().Delete(ctx, userId)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validChatSession)
		assert.Error(s.T(), err)
	})

	s.Run("it fails for non-existent chat session", func() {
		result, err := repo.Get(ctx, chatSessionId)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)

		err = repo.Update(ctx, validChatSession)
		assert.Error(s.T(), err)

		err = repo.Delete(ctx, chatSessionId)
		assert.Error(s.T(), err)
	})
}
