package chatsessionitems_test

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ctx               = context.Background()
	userId            = user.UserID(1)
	chatSessionId     = chatsession.ChatSessionID(1)
	chatSessionItemId = chatsession.ChatSessionItemID(1)
	validUser         = user.NewUser()
	validChatSession  = &chatsession.ChatSession{
		ID: chatSessionId,
		ChatSessionData: chatsession.ChatSessionData{
			ChatSessionItemIDs: []chatsession.ChatSessionItemID{chatSessionItemId},
			UserIDs:            []user.UserID{userId},
		},
		Metadata: domain.NewMetadata(),
	}
	validUserMsgData = chatsession.UserChatSessionItemData{
		Content: "Hello I'm a user message",
	}
	validAssistantMsgData = chatsession.AssistantChatSessionItemData{
		Content: "Hello I'm an assistant message",
	}
	validUserChatItem = &chatsession.UserChatSessionItem{
		ChatSessionItemBase: chatsession.ChatSessionItemBase{
			ID:        chatSessionItemId,
			SessionID: chatSessionId,
			Metadata: &domain.Metadata{
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		Type:                    chatsession.ChatSessionItemTypeUser,
		UserChatSessionItemData: validUserMsgData,
	}
	validAssistantChatItem = &chatsession.AssistantChatSessionItem{
		ChatSessionItemBase: chatsession.ChatSessionItemBase{
			ID:        chatSessionItemId + 1,
			SessionID: chatSessionId,
			Metadata: &domain.Metadata{
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 1, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 1, time.UTC),
			},
		},
		Type:                         chatsession.ChatSessionItemTypeAssistant,
		AssistantChatSessionItemData: validAssistantMsgData,
	}
	repo repos.IChatSessionItemsRepo
)

func (s *ChatSessionItemsRepoTestSuite) SetupSubTest() {
	s.Reset()
	repos, _ := s.GetRepos()
	repo = repos.ChatSessionItemsRepo(userId)
	validUser.ID = userId
	err := repos.UsersRepo().Create(ctx, validUser)
	require.NoError(s.T(), err)
	err = repos.ChatSessionsRepo(userId).Create(ctx, validChatSession)
	require.NoError(s.T(), err)
}

func (s *ChatSessionItemsRepoTestSuite) TestChatSessionItemsRepo() {
	s.Run("it works", func() {
		err := repo.Create(ctx, validUserChatItem)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, chatSessionItemId)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), validUserChatItem.ID, result.GetID())
		assert.Equal(s.T(), validUserChatItem.Type, result.GetType())

		err = repo.Create(ctx, validAssistantChatItem)
		assert.NoError(s.T(), err)

		results, err := repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 2)
		assert.Equal(s.T(), validUserChatItem.ID, results[0].GetID())
		assert.Equal(s.T(), validAssistantChatItem.ID, results[1].GetID())

		updatedUserChatItem := &chatsession.UserChatSessionItem{}
		utils.DeepClone(validUserChatItem, updatedUserChatItem)
		updatedUserChatItem.Content = "Hello I'm an updated user message"

		err = repo.Update(ctx, updatedUserChatItem)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, chatSessionItemId)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), updatedUserChatItem.ID, result.GetID())
		assert.Equal(s.T(), updatedUserChatItem.Content, result.(*chatsession.UserChatSessionItem).Content)

		err = repo.Delete(ctx, chatSessionItemId)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 1)

		err = repo.Delete(ctx, validAssistantChatItem.ID)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 0)
	})
}

func (s *ChatSessionItemsRepoTestSuite) TestChatSessionItemsRepoFailure() {
	s.Run("it fails with non-existent user", func() {
		repos, _ := s.GetRepos()
		err := repos.UsersRepo().Delete(ctx, userId)
		require.NoError(s.T(), err)

		err = repo.Create(ctx, validUserChatItem)
		assert.Error(s.T(), err)

		_, err = repo.Get(ctx, chatSessionItemId)
		assert.Error(s.T(), err)

		err = repo.Delete(ctx, chatSessionItemId)
		assert.Error(s.T(), err)
	})

	s.Run("it fails with non-existent chat session", func() {
		repos, _ := s.GetRepos()
		err := repos.ChatSessionsRepo(userId).Delete(ctx, chatSessionId)
		require.NoError(s.T(), err)

		err = repo.Create(ctx, validUserChatItem)
		assert.Error(s.T(), err)

		_, err = repo.Get(ctx, chatSessionItemId)
		assert.Error(s.T(), err)

		err = repo.Delete(ctx, chatSessionItemId)
		assert.Error(s.T(), err)
	})

	s.Run("it fails with non-existent chat session item", func() {
		_, err := repo.Get(ctx, chatSessionItemId)
		assert.Error(s.T(), err)

		err = repo.Delete(ctx, chatSessionItemId)
		assert.Error(s.T(), err)
	})

	s.Run("it fails with duplicate chat session item", func() {
		err := repo.Create(ctx, validUserChatItem)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validUserChatItem)
		assert.Error(s.T(), err)
	})
}
