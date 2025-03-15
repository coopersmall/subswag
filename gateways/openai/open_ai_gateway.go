package openai

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/clients"
	"github.com/coopersmall/subswag/utils"
	"github.com/tmc/langchaingo/llms"
)

type OpenAIGateway struct {
	logger utils.ILogger
	tracer apm.ITracer
	client clients.IAIClient
}

func NewOpenAIGateway(
	logger utils.ILogger,
	tracer apm.ITracer,
	client clients.IAIClient,
) *OpenAIGateway {
	return &OpenAIGateway{
		logger: logger,
		tracer: tracer,
		client: client,
	}
}

func (g *OpenAIGateway) GenerateContent(ctx context.Context, messages []llms.MessageContent, opts ...llms.CallOption) (*llms.ContentResponse, error) {
	var (
		response *llms.ContentResponse
		err      error
	)
	g.tracer.Trace(ctx, "generate-content", func(ctx context.Context, span apm.ISpan) error {
		response, err = g.client.GenerateContent(ctx, messages, opts...)
		return err
	})
	return response, err
}
