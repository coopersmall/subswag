package openai_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/coopersmall/subswag/apm"
	c "github.com/coopersmall/subswag/clients/openai"
	"github.com/coopersmall/subswag/gateways/openai"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	llmso "github.com/tmc/langchaingo/llms/openai"
)

var (
	ctx      = context.Background()
	messages = []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: "Test message",
				},
			},
		},
	}
	getToken = func() (string, error) {
		return "test-token", nil
	}
)

func (s *OpenAIGatewayTestSuite) SetupTest() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
}

func (s *OpenAIGatewayTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *OpenAIGatewayTestSuite) TestGenerateContent() {
	s.Run("successful content generation", func() {
		s.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(s.T(), "/chat/completions", r.URL.Path)
			assert.Equal(s.T(), "POST", r.Method)
			assert.Equal(s.T(), "Bearer test-token", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"id": "test-id",
				"object": "chat.completion",
				"choices": [{
					"message": {
						"role": "assistant",
						"content": "Test response"
					}
				}]
			}`))
		})

		client := c.NewOpenAIClient(getToken, llmso.WithBaseURL(s.server.URL))

		mockLogger := new(utils.MockLogger)
		mockTracer := new(apm.MockTracer)
		gateway := openai.NewOpenAIGateway(mockLogger, mockTracer, client)

		response, err := gateway.GenerateContent(ctx, messages)

		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), response)
		assert.Equal(s.T(), "Test response", response.Choices[0].Content)
	})

	s.Run("handles error response", func() {
		s.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": {"message": "Test error"}}`))
		})

		client := c.NewOpenAIClient(getToken, llmso.WithBaseURL(s.server.URL))

		mockLogger := new(utils.MockLogger)
		mockTracer := new(apm.MockTracer)
		gateway := openai.NewOpenAIGateway(mockLogger, mockTracer, client)

		response, err := gateway.GenerateContent(ctx, messages)

		assert.Error(s.T(), err)
		assert.Nil(s.T(), response)
	})

	s.Run("handles context cancellation", func() {
		s.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate delay
			select {
			case <-r.Context().Done():
				return
			}
		})

		client := c.NewOpenAIClient(getToken, llmso.WithBaseURL(s.server.URL))

		mockLogger := new(utils.MockLogger)
		mockTracer := new(apm.MockTracer)
		gateway := openai.NewOpenAIGateway(mockLogger, mockTracer, client)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		response, err := gateway.GenerateContent(ctx, messages)

		assert.Error(s.T(), err)
		assert.Nil(s.T(), response)
	})
}
