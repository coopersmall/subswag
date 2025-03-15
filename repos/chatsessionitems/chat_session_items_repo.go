package chatsessionitems

import (
	"context"
	"database/sql"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/user"
	reposdomain "github.com/coopersmall/subswag/repos/domain"
	"github.com/coopersmall/subswag/utils"
)

type ChatSessionItemsRepo struct {
	*reposdomain.StandardRepo[chatsession.ChatSessionItemID, chatsession.ChatSessionItem, db.ChatSessionItem]
}

func NewChatSessionItemsRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
	userId user.UserID,
) *ChatSessionItemsRepo {
	return &ChatSessionItemsRepo{
		StandardRepo: reposdomain.NewStandardRepo[chatsession.ChatSessionItemID, chatsession.ChatSessionItem, db.ChatSessionItem](
			"chat_session_item",
			querier,
			tracer,
			userId,
			convertRowToChatSessionItem,
			convertChatSessionItemToRow,
			isEmptyChatSessionItem,
			getFunc(userId),
			getAllFunc(userId),
			createFunc(userId),
			updateFunc(userId),
			deleteFunc(userId),
		),
	}
}

func getFunc(userId user.UserID) func(ctx context.Context, q db.IStandardQueriesReadOnly, id chatsession.ChatSessionItemID) (db.ChatSessionItem, error) {
	return func(ctx context.Context, q db.IStandardQueriesReadOnly, id chatsession.ChatSessionItemID) (db.ChatSessionItem, error) {
		return q.GetChatSessionItem(ctx, db.GetChatSessionItemParams{
			ID:     int64(id),
			UserID: int64(userId),
		})
	}
}

func getAllFunc(userId user.UserID) func(ctx context.Context, q db.IStandardQueriesReadOnly) ([]db.ChatSessionItem, error) {
	return func(ctx context.Context, q db.IStandardQueriesReadOnly) ([]db.ChatSessionItem, error) {
		return q.GetAllChatSessionItems(ctx)
	}
}

func createFunc(userId user.UserID) func(ctx context.Context, q db.IStandardQueriesReadWrite, item db.ChatSessionItem) (sql.Result, error) {
	return func(ctx context.Context, q db.IStandardQueriesReadWrite, item db.ChatSessionItem) (sql.Result, error) {
		return q.CreateChatSessionItem(ctx, db.CreateChatSessionItemParams{
			ID:        item.ID,
			UserID:    int64(userId),
			SessionID: item.SessionID,
			CreatedAt: item.CreatedAt,
			Type:      item.Type,
			Data:      item.Data,
		})
	}
}

func updateFunc(userId user.UserID) func(ctx context.Context, q db.IStandardQueriesReadWrite, item db.ChatSessionItem) (sql.Result, error) {
	return func(ctx context.Context, q db.IStandardQueriesReadWrite, item db.ChatSessionItem) (sql.Result, error) {
		return q.UpdateChatSessionItem(ctx, db.UpdateChatSessionItemParams{
			ID:        item.ID,
			UserID:    int64(userId),
			SessionID: item.SessionID,
			UpdatedAt: item.UpdatedAt,
			Type:      item.Type,
			Data:      item.Data,
		})
	}
}

func deleteFunc(userId user.UserID) func(ctx context.Context, q db.IStandardQueriesReadWrite, id chatsession.ChatSessionItemID) (sql.Result, error) {
	return func(ctx context.Context, q db.IStandardQueriesReadWrite, id chatsession.ChatSessionItemID) (sql.Result, error) {
		return q.DeleteChatSessionItem(ctx, db.DeleteChatSessionItemParams{
			ID:     int64(id),
			UserID: int64(userId),
		})
	}
}

func (r *ChatSessionItemsRepo) GetBySessionId(ctx context.Context, sessionId chatsession.ChatSessionID) ([]chatsession.ChatSessionItem, error) {
	return r.StandardRepo.Query(ctx, func(ctx context.Context, q db.IStandardQueriesReadOnly, userId user.UserID) ([]db.ChatSessionItem, error) {
		return q.GetChatSessionItemsBySessionID(ctx, db.GetChatSessionItemsBySessionIDParams{
			SessionID: int64(sessionId),
			UserID:    int64(userId),
		})
	})
}

func (r *ChatSessionItemsRepo) DeleteBySessionId(ctx context.Context, sessionId chatsession.ChatSessionID) error {
	return r.StandardRepo.Execute(ctx, func(ctx context.Context, q db.IStandardQueriesReadWrite, userId user.UserID) error {
		return q.DeleteChatSessionItemsBySessionID(ctx, db.DeleteChatSessionItemsBySessionIDParams{
			SessionID: int64(sessionId),
			UserID:    int64(userId),
		})
	})
}

func isEmptyChatSessionItem(item chatsession.ChatSessionItem) bool {
	return item == nil || item.GetID() == 0
}

func convertRowToChatSessionItem(result db.ChatSessionItem) (chatsession.ChatSessionItem, error) {
	chatType := chatsession.ChatSessionItemType(result.Type)
	switch chatType {
	case chatsession.ChatSessionItemTypeUser:
		var data chatsession.UserChatSessionItemData
		err := utils.Unmarshal(result.Data, &data)
		return &chatsession.UserChatSessionItem{
			ChatSessionItemBase: chatsession.ChatSessionItemBase{
				ID:        chatsession.ChatSessionItemID(result.ID),
				SessionID: chatsession.ChatSessionID(result.SessionID),
				Metadata: &domain.Metadata{
					CreatedAt: result.CreatedAt,
					UpdatedAt: result.UpdatedAt.Time,
				},
			},
			Type:                    chatType,
			UserChatSessionItemData: data,
		}, err
	case chatsession.ChatSessionItemTypeAssistant:
		var data chatsession.AssistantChatSessionItemData
		err := utils.Unmarshal(result.Data, &data)
		return &chatsession.AssistantChatSessionItem{
			ChatSessionItemBase: chatsession.ChatSessionItemBase{
				ID:        chatsession.ChatSessionItemID(result.ID),
				SessionID: chatsession.ChatSessionID(result.SessionID),
				Metadata: &domain.Metadata{
					CreatedAt: result.CreatedAt,
					UpdatedAt: result.UpdatedAt.Time,
				},
			},
			Type:                         chatType,
			AssistantChatSessionItemData: data,
		}, err
	default:
		return nil, utils.NewUnableToHandleError("unknown chat session item type")
	}
}

func convertChatSessionItemToRow(
	item chatsession.ChatSessionItem,
) (db.ChatSessionItem, error) {
	data, err := utils.Marshal(item)
	return db.ChatSessionItem{
		ID:        int64(item.GetID()),
		SessionID: int64(item.GetSessionID()),
		CreatedAt: item.GetMetadata().CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  item.GetMetadata().UpdatedAt,
			Valid: item.GetMetadata().UpdatedAt != time.Time{},
		},
		Type: string(item.GetType()),
		Data: data,
	}, err
}
