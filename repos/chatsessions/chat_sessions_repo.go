package chatsessions

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/user"
	reposdomain "github.com/coopersmall/subswag/repos/domain"
)

type ChatSessionsRepo struct {
	*reposdomain.StandardRepo[chatsession.ChatSessionID, *chatsession.ChatSession, db.ChatSession]
}

func NewChatSessionsRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
	userId user.UserID,
) *ChatSessionsRepo {
	return &ChatSessionsRepo{
		StandardRepo: reposdomain.NewStandardRepo[chatsession.ChatSessionID, *chatsession.ChatSession, db.ChatSession](
			"chat_session",
			querier,
			tracer,
			userId,
			convertRowToChatSession,
			func(a *chatsession.ChatSession) (db.ChatSession, error) {
				return convertChatSessionToRow(userId, a)
			},
			isEmptyChatSession,
			func(ctx context.Context, iqro db.IStandardQueriesReadOnly, csi chatsession.ChatSessionID) (db.ChatSession, error) {
				return iqro.GetChatSession(ctx, db.GetChatSessionParams{
					ID:     int64(csi),
					UserID: int64(userId),
				})
			},
			func(ctx context.Context, iqro db.IStandardQueriesReadOnly) ([]db.ChatSession, error) {
				return iqro.GetAllChatSessions(ctx)
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, cs db.ChatSession) (sql.Result, error) {
				return iqrw.CreateChatSession(ctx, db.CreateChatSessionParams{
					ID:        cs.ID,
					UserID:    int64(userId),
					CreatedAt: cs.CreatedAt,
					Data:      cs.Data,
				})
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, cs db.ChatSession) (sql.Result, error) {
				return iqrw.UpdateChatSession(ctx, db.UpdateChatSessionParams{
					ID:        cs.ID,
					UserID:    int64(userId),
					UpdatedAt: cs.UpdatedAt,
					Data:      cs.Data,
				})
			},
			func(ctx context.Context, iqrw db.IStandardQueriesReadWrite, csi chatsession.ChatSessionID) (sql.Result, error) {
				return iqrw.DeleteChatSession(ctx, db.DeleteChatSessionParams{
					ID:     int64(csi),
					UserID: int64(userId),
				})
			},
		),
	}
}

func isEmptyChatSession(chatSession *chatsession.ChatSession) bool {
	return chatSession == nil || chatSession.ID == 0
}

func convertRowToChatSession(result db.ChatSession) (*chatsession.ChatSession, error) {
	var data chatsession.ChatSessionData
	err := json.Unmarshal(result.Data, &data)
	return &chatsession.ChatSession{
		ID:              chatsession.ChatSessionID(result.ID),
		ChatSessionData: data,
		Metadata: &domain.Metadata{
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt.Time,
		},
	}, err
}

func convertChatSessionToRow(
	userId user.UserID,
	chatSession *chatsession.ChatSession,
) (db.ChatSession, error) {
	data, err := json.Marshal(chatSession.ChatSessionData)
	return db.ChatSession{
		ID:        int64(chatSession.ID),
		UserID:    int64(userId),
		CreatedAt: chatSession.Metadata.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  chatSession.Metadata.UpdatedAt,
			Valid: true,
		},
		Data: data,
	}, err
}
