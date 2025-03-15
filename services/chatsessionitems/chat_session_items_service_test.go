package chatsessionitems_test

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/services"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ctx = context.Background()

	userId    = user.UserID(1)
	validUser = user.NewUser()

	session1Id = chatsession.ChatSessionID(1)
	session1   = chatsession.NewChatSession(nil, nil)
	session2Id = chatsession.ChatSessionID(2)
	session2   = chatsession.NewChatSession(nil, nil)

	item1Id = chatsession.ChatSessionItemID(1)
	item2Id = chatsession.ChatSessionItemID(2)
	item3Id = chatsession.ChatSessionItemID(3)

	items = []chatsession.ChatSessionItem{
		&chatsession.UserChatSessionItem{
			ChatSessionItemBase: chatsession.ChatSessionItemBase{
				ID:        item1Id,
				SessionID: session1Id,
				Metadata: &domain.Metadata{
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			Type: chatsession.ChatSessionItemTypeUser,
			UserChatSessionItemData: chatsession.UserChatSessionItemData{
				Content: "Hi there!",
			},
		},
		&chatsession.AssistantChatSessionItem{
			ChatSessionItemBase: chatsession.ChatSessionItemBase{
				ID:        item2Id,
				SessionID: session1Id,
				Metadata: &domain.Metadata{
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 1, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 1, time.UTC),
				},
			},
			Type: chatsession.ChatSessionItemTypeAssistant,
			AssistantChatSessionItemData: chatsession.AssistantChatSessionItemData{
				Content: "Hi there!",
			},
		},
		&chatsession.UserChatSessionItem{
			ChatSessionItemBase: chatsession.ChatSessionItemBase{
				ID:        item3Id,
				SessionID: session1Id,
				Metadata: &domain.Metadata{
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 2, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 2, time.UTC),
				},
			},
			Type: chatsession.ChatSessionItemTypeUser,
			UserChatSessionItemData: chatsession.UserChatSessionItemData{
				Content: "Hi there!",
			},
		},
	}

	repo    repos.IChatSessionItemsRepo
	service services.IChatSessionItemsService
)

func (s *ChatSessionItemsServiceTestSuite) SetupSubTest() {
	s.Reset()
	repos, _ := s.GetRepos()
	repo = repos.ChatSessionItemsRepo(userId)
	services, _ := s.GetServices()
	service = services.ChatSessionItemsService(userId)

	validUser.ID = userId
	err := repos.UsersRepo().Create(ctx, validUser)
	assert.NoError(s.T(), err)
	session1.ID = session1Id
	err = repos.ChatSessionsRepo(userId).Create(ctx, session1)
	assert.NoError(s.T(), err)
	session2.ID = session2Id
	err = repos.ChatSessionsRepo(userId).Create(ctx, session2)
	assert.NoError(s.T(), err)
}

func (s *ChatSessionItemsServiceTestSuite) TestChatSessionItemsServiceSuccess() {
	s.Run("it works", func() {
		for _, item := range items {
			err := service.CreateChatSessionItem(ctx, item)
			assert.NoError(s.T(), err)

			result, err := service.GetChatSessionItem(ctx, item.GetID())
			require.NoError(s.T(), err)
			assert.NotNil(s.T(), result)
			assert.Equal(s.T(), item.GetID(), result.GetID())
		}

		results, err := service.GetAllChatSessionItems(ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), results, len(items))
		assert.Equal(s.T(), items[0].GetID(), results[0].GetID())
		assert.Equal(s.T(), items[1].GetID(), results[1].GetID())
		assert.Equal(s.T(), items[2].GetID(), results[2].GetID())

		session2ItemId := chatsession.ChatSessionItemID(4)
		session2Item := &chatsession.UserChatSessionItem{}
		utils.DeepClone(items[0], session2Item)
		session2Item.ID = session2ItemId
		session2Item.SessionID = session2Id

		err = service.CreateChatSessionItem(ctx, session2Item)
		assert.NoError(s.T(), err)

		results, err = service.GetAllChatSessionItems(ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), results, len(items)+1)
		assert.Equal(s.T(), session2Item.GetID(), results[3].GetID())

		results, err = service.GetChatSessionItemsBySessionId(ctx, session1Id)
		require.NoError(s.T(), err)
		assert.Len(s.T(), results, len(items))
		assert.Equal(s.T(), items[0].GetID(), results[0].GetID())
		assert.Equal(s.T(), items[1].GetID(), results[1].GetID())
		assert.Equal(s.T(), items[2].GetID(), results[2].GetID())

		err = service.DeleteChatSessionItem(ctx, items[0].GetID())
		assert.NoError(s.T(), err)

		results, err = service.GetChatSessionItemsBySessionId(ctx, session1Id)
		require.NoError(s.T(), err)
		assert.Len(s.T(), results, len(items)-1)
		assert.Equal(s.T(), items[1].GetID(), results[0].GetID())
		assert.Equal(s.T(), items[2].GetID(), results[1].GetID())

		err = service.DeleteChatSessionItemsBySessionId(ctx, session1Id)
		assert.NoError(s.T(), err)

		results, err = service.GetAllChatSessionItems(ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), results, 1)

		results, err = service.GetChatSessionItemsBySessionId(ctx, session2Id)
		require.NoError(s.T(), err)
		assert.Len(s.T(), results, 1)

		err = service.DeleteChatSessionItemsBySessionId(ctx, session2Id)
		assert.NoError(s.T(), err)

		results, err = service.GetAllChatSessionItems(ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), results, 0)
	})

}
