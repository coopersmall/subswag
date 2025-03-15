package answerassistant

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/ai"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/utils"
	"github.com/tmc/langchaingo/llms"
)

type AnswerAssistantService struct {
	logger                  utils.ILogger
	tracer                  apm.ITracer
	aiGateway               gateways.IAIGateway
	chatSessionItemsService iChatSessionItemsService
}

func NewAnswerAssistantService(
	logger utils.ILogger,
	tracer apm.ITracer,
	aiGateway gateways.IAIGateway,
	chatSessionItemsService iChatSessionItemsService,
) *AnswerAssistantService {
	return &AnswerAssistantService{
		logger:                  logger,
		tracer:                  tracer,
		aiGateway:               aiGateway,
		chatSessionItemsService: chatSessionItemsService,
	}
}

func (s *AnswerAssistantService) AnswerQuestion(
	ctx context.Context,
	userId user.UserID,
	items []chatsession.ChatSessionItem,
) (*chatsession.AssistantChatSessionItem, error) {
	var (
		assistantChatSessionItem *chatsession.AssistantChatSessionItem
		err                      error
	)
	s.tracer.Trace(ctx, "answer-question", func(ctx context.Context, span apm.ISpan) error {
		if len(items) == 0 {
			return utils.NewInvalidArgumentError("items cannot be empty")
		}
		sessionId := items[0].GetSessionID()

		messages := s.chatSessionItemsService.ConvertChatSessionItemsToLLMMessages(ctx, items)
		messages = append(messages, llms.MessageContent{
			Role: llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: ai.RenderPrompt(
						s.prompt(),
						ai.PromptWithJSONSchema(responseSchema()),
					),
				},
			},
		})
		var response *llms.ContentResponse
		response, err = s.aiGateway.GenerateContent(
			ctx,
			messages,
			llms.WithModel(ai.MODEL_LLAMA_3_3_8b),
			llms.WithTemperature(0.5),
			llms.WithJSONMode(),
		)
		if err != nil {
			return utils.NewInternalError("chat completion failed", err)
		}
		content := response.Choices[0].Content
		if content == "" {
			return utils.NewInternalError("chat completion response content is nil")
		}
		assistantChatSessionItem = &chatsession.AssistantChatSessionItem{
			ChatSessionItemBase: chatsession.ChatSessionItemBase{
				ID:        chatsession.ChatSessionItemID(utils.NewID()),
				SessionID: sessionId,
				Metadata:  domain.NewMetadata(),
			},
			Type: chatsession.ChatSessionItemTypeAssistant,
			AssistantChatSessionItemData: chatsession.AssistantChatSessionItemData{
				Content: string(content),
			},
		}
		return nil
	})
	return assistantChatSessionItem, err
}

type iChatSessionItemsService interface {
	ConvertChatSessionItemsToLLMMessages(
		ctx context.Context,
		items []chatsession.ChatSessionItem,
	) []llms.MessageContent
}

func (s *AnswerAssistantService) prompt() string {
	return "You are helping answer a user's question. Please provide a response."
}

func responseSchema() string {
	return `{
    "$schema": {
        "type": "object",
        "properties": {
            "answer": {
                "type": "string",
                "description": "The detailed answer to the user's question"
            },
            "confidence": {
                "type": "number",
                "minimum": 0,
                "maximum": 1,
                "description": "Confidence score of the answer (0-1)"
            },
            "sources": {
                "type": "array",
                "items": {
                    "type": "string"
                },
                "description": "Optional list of sources or references"
            },
            "tags": {
                "type": "array",
                "items": {
                    "type": "string"
                },
                "description": "Optional relevant tags or categories"
            }
        },
        "required": ["answer", "confidence"]
    }
}`
}
