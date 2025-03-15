package groq

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/clients"
	"github.com/coopersmall/subswag/utils"
	"github.com/tmc/langchaingo/llms"
)

type GroqGateway struct {
	logger     utils.ILogger
	tracer     apm.ITracer
	groqClient clients.IAIClient
}

func NewGroqGateway(
	logger utils.ILogger,
	tracer apm.ITracer,
	clients clients.IAIClient,
) *GroqGateway {
	return &GroqGateway{
		logger:     logger,
		tracer:     tracer,
		groqClient: clients,
	}
}

func (g *GroqGateway) GenerateContent(ctx context.Context, messages []llms.MessageContent, opts ...llms.CallOption) (*llms.ContentResponse, error) {
	var (
		response *llms.ContentResponse
		err      error
	)
	g.tracer.Trace(ctx, "generate-content", func(ctx context.Context, span apm.ISpan) error {
		response, err = g.groqClient.GenerateContent(ctx, messages, opts...)
		return err
	})
	return response, err
}
