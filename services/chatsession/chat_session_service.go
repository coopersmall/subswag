package chatsession

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/cache"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/repos"
	servicesdomain "github.com/coopersmall/subswag/services/domain"
	"github.com/coopersmall/subswag/utils"
)

type ChatSessionsService struct {
	standardService  *servicesdomain.StandardService[chatsession.ChatSessionID, *chatsession.ChatSession]
	chatSessionsRepo repos.IChatSessionsRepo
}

func NewChatSessionsService(
	logger utils.ILogger,
	tracer apm.ITracer,
	chatSessionsCache cache.IChatSessionsCache,
	chatSessionsRepo repos.IChatSessionsRepo,
) *ChatSessionsService {
	standardService := servicesdomain.NewStandardService[chatsession.ChatSessionID, *chatsession.ChatSession](
		"chat-sessions",
		logger,
		tracer,
		struct {
			Get    func(context.Context, chatsession.ChatSessionID) (*chatsession.ChatSession, error)
			All    func(context.Context) ([]*chatsession.ChatSession, error)
			Create func(context.Context, *chatsession.ChatSession) error
			Update func(context.Context, *chatsession.ChatSession) error
			Delete func(context.Context, chatsession.ChatSessionID) error
		}{
			Get:    chatSessionsRepo.Get,
			All:    chatSessionsRepo.All,
			Create: chatSessionsRepo.Create,
			Update: chatSessionsRepo.Update,
			Delete: chatSessionsRepo.Delete,
		},
		&struct {
			Get    func(context.Context, chatsession.ChatSessionID, func(context.Context, chatsession.ChatSessionID) (*chatsession.ChatSession, error)) (*chatsession.ChatSession, error)
			Set    func(context.Context, chatsession.ChatSessionID, *chatsession.ChatSession) error
			Delete func(context.Context, chatsession.ChatSessionID) error
		}{
			Get:    chatSessionsCache.Get,
			Set:    chatSessionsCache.Set,
			Delete: chatSessionsCache.Delete,
		},
		nil,
	)
	return &ChatSessionsService{
		standardService:  standardService,
		chatSessionsRepo: chatSessionsRepo,
	}
}

func (s *ChatSessionsService) GetChatSession(ctx context.Context, sessionId chatsession.ChatSessionID) (*chatsession.ChatSession, error) {
	return s.standardService.Get(ctx, sessionId)
}

func (s *ChatSessionsService) GetAllChatSessions(ctx context.Context) ([]*chatsession.ChatSession, error) {
	return s.standardService.All(ctx)
}

func (s *ChatSessionsService) CreateChatSession(ctx context.Context, data chatsession.ChatSessionData) error {
	chatSession := chatsession.NewChatSession(data.UserIDs, data.ChatSessionItemIDs)
	return s.standardService.Create(ctx, chatSession.ID, chatSession)
}

func (s *ChatSessionsService) UpdateChatSession(ctx context.Context, session *chatsession.ChatSession) error {
	return s.standardService.Update(ctx, session.ID, session)
}

func (s *ChatSessionsService) DeleteChatSession(ctx context.Context, sessionId chatsession.ChatSessionID) error {
	return s.standardService.Delete(ctx, sessionId)
}
