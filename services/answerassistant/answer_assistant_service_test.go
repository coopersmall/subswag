package answerassistant_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	openaiclient "github.com/coopersmall/subswag/clients/openai"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/gateways"
	openaigateway "github.com/coopersmall/subswag/gateways/openai"
	"github.com/coopersmall/subswag/services"
	"github.com/coopersmall/subswag/services/answerassistant"
	"github.com/stretchr/testify/assert"
	llmso "github.com/tmc/langchaingo/llms/openai"
)

var (
	ctx = context.Background()

	userId    = user.UserID(1)
	validUser = user.NewUser()

	sessionId = chatsession.ChatSessionID(1)
	session   = chatsession.NewChatSession(nil, nil)

	sessionItem1Id = chatsession.ChatSessionItemID(1)
	sessionItem2Id = chatsession.ChatSessionItemID(2)
	sessionItem3Id = chatsession.ChatSessionItemID(3)

	sessionItems = []chatsession.ChatSessionItem{
		&chatsession.UserChatSessionItem{
			ChatSessionItemBase: chatsession.ChatSessionItemBase{
				ID:        sessionItem1Id,
				SessionID: sessionId,
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
				ID:        sessionItem2Id,
				SessionID: sessionId,
				Metadata: &domain.Metadata{
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 1, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 1, time.UTC),
				},
			},
			Type: chatsession.ChatSessionItemTypeAssistant,
			AssistantChatSessionItemData: chatsession.AssistantChatSessionItemData{
				Content: "What's Up?",
			},
		},
	}
	gateway                 gateways.IAIGateway
	chatSessionItemsService services.IChatSessionItemsService
	service                 services.IAnswerAssistantService
)

func (s *AnswerAssistantServiceTestSuite) SetupSubTest() {
	s.Reset()
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	getAPIToken := func() (string, error) {
		return "test-token", nil
	}
	client := openaiclient.NewOpenAIClient(getAPIToken, llmso.WithBaseURL(s.server.URL))
	gateway = openaigateway.NewOpenAIGateway(s.GetLogger("ai-gateway"), s.GetTracer("ai-gateway"), client)

	services, _ := s.GetServices()
	chatSessionItemsService = services.ChatSessionItemsService(userId)
	service = answerassistant.NewAnswerAssistantService(
		s.GetLogger("answer-assistant"),
		s.GetTracer("answer-assistant"),
		gateway,
		chatSessionItemsService,
	)

	repos, _ := s.GetRepos()
	validUser.ID = userId
	err := repos.UsersRepo().Create(ctx, validUser)
	assert.NoError(s.T(), err)
	err = repos.ChatSessionsRepo(userId).Create(ctx, session)
	assert.NoError(s.T(), err)
}

func (s *AnswerAssistantServiceTestSuite) AfterTest(suiteName, testName string) {
	s.server.Close()
}

func (s *AnswerAssistantServiceTestSuite) TestAnswerQuestionSuccess() {
	s.Run("successful answer generation", func() {
		s.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(s.T(), "/chat/completions", r.URL.Path)
			assert.Equal(s.T(), "Bearer test-token", r.Header.Get("Authorization"))
			assert.Equal(s.T(), "POST", r.Method)

			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)

			messages := body["messages"].([]interface{})
			assert.Len(s.T(), messages, 3)

			message1 := messages[0].(map[string]interface{})
			assert.Equal(s.T(), "user", message1["role"])
			assert.Equal(s.T(), "Hi there!", message1["content"])

			message2 := messages[1].(map[string]interface{})
			assert.Equal(s.T(), "assistant", message2["role"])
			assert.Equal(s.T(), "What's Up?", message2["content"])

			message3 := messages[2].(map[string]interface{})
			assert.Equal(s.T(), "system", message3["role"])
			assert.Contains(s.T(), message3["content"], "You are helping answer a user's question. Please provide a response.")
			assert.Contains(s.T(), message3["content"], `{
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
}`)

			response := map[string]interface{}{
				"id":     "test-id",
				"object": "chat.completion",
				"choices": []map[string]interface{}{
					{
						"message": map[string]interface{}{
							"role": "assistant",
							"content": `{
								"answer": "This is a test answer",
								"confidence": 0.95,
								"sources": ["source1", "source2"],
								"tags": ["tag1", "tag2"]
							}`,
						},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		})

		result, err := service.AnswerQuestion(ctx, userId, sessionItems)

		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.NotZero(s.T(), result.ID)
		assert.Equal(s.T(), chatsession.ChatSessionItemTypeAssistant, result.Type)
		assert.Equal(s.T(), sessionId, result.SessionID)
		assert.Contains(s.T(), result.Content, "This is a test answer")
	})
}
