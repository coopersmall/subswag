package chatsessionitems

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/tmc/langchaingo/llms"
)

type ChatSessionItemsService struct {
	logger utils.ILogger
	tracer apm.ITracer
	repo   repos.IChatSessionItemsRepo
}

func NewChatSessionItemsService(
	logger utils.ILogger,
	tracer apm.ITracer,
	repo repos.IChatSessionItemsRepo,
) *ChatSessionItemsService {
	return &ChatSessionItemsService{
		logger: logger,
		tracer: tracer,
		repo:   repo,
	}
}

func (s *ChatSessionItemsService) GetChatSessionItem(ctx context.Context, itemId chatsession.ChatSessionItemID) (chatsession.ChatSessionItem, error) {
	return s.repo.Get(ctx, itemId)
}

func (s *ChatSessionItemsService) GetAllChatSessionItems(ctx context.Context) ([]chatsession.ChatSessionItem, error) {
	return s.repo.All(ctx)
}

func (s *ChatSessionItemsService) GetChatSessionItemsBySessionId(ctx context.Context, sessionId chatsession.ChatSessionID) ([]chatsession.ChatSessionItem, error) {
	return s.repo.GetBySessionId(ctx, sessionId)
}

func (s *ChatSessionItemsService) CreateChatSessionItem(ctx context.Context, item chatsession.ChatSessionItem) error {
	return s.repo.Create(ctx, item)
}

func (s *ChatSessionItemsService) UpdateChatSessionItem(ctx context.Context, item chatsession.ChatSessionItem) error {
	return s.repo.Update(ctx, item)
}

func (s *ChatSessionItemsService) DeleteChatSessionItem(ctx context.Context, itemId chatsession.ChatSessionItemID) error {
	return s.repo.Delete(ctx, itemId)
}

func (s *ChatSessionItemsService) DeleteChatSessionItemsBySessionId(ctx context.Context, sessionId chatsession.ChatSessionID) error {
	return s.repo.DeleteBySessionId(ctx, sessionId)
}

func (s *ChatSessionItemsService) ConvertChatSessionItemsToLLMMessages(
	ctx context.Context,
	items []chatsession.ChatSessionItem,
) []llms.MessageContent {
	messages := make([]llms.MessageContent, 0, len(items))
	for _, item := range items {
		switch i := item.(type) {
		case *chatsession.UserChatSessionItem:
			messages = append(messages, llms.MessageContent{
				Role: llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{
					llms.TextContent{
						Text: i.Content,
					},
				},
			})
		case *chatsession.AssistantChatSessionItem:
			messages = append(messages, llms.MessageContent{
				Role: llms.ChatMessageTypeAI,
				Parts: []llms.ContentPart{
					llms.TextContent{
						Text: i.Content,
					},
				},
			})
		}
	}
	return messages
}
