package groq

import (
	"github.com/coopersmall/subswag/utils"
	"github.com/tmc/langchaingo/llms/openai"
)

type GroqClient struct {
	*openai.LLM
}

func NewGroqClient(
	getGroqKey func() (string, error),
) *GroqClient {
	apiKey, err := getGroqKey()
	if err != nil {
		panic(utils.NewInternalError("failed to get Groq key for client", err))
	}
	llm, _ := openai.New(
		openai.WithBaseURL("https://api.groq.com/openai/v1"),
		openai.WithToken(apiKey),
	)
	return &GroqClient{
		llm,
	}
}
