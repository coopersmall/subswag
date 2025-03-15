package openai

import (
	"github.com/coopersmall/subswag/utils"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIClient struct {
	*openai.LLM
}

func NewOpenAIClient(
	getOpenAIKey func() (string, error),
	opts ...openai.Option,
) *OpenAIClient {
	apiKey, err := getOpenAIKey()
	if err != nil {
		panic(utils.NewInternalError("failed to get OpenAI key for client", err))
	}
	options := []openai.Option{
		openai.WithToken(apiKey),
	}
	for _, opt := range opts {
		options = append(options, opt)
	}
	llm, _ := openai.New(
		options...,
	)
	return &OpenAIClient{
		llm,
	}
}
