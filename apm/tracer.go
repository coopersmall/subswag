package apm

import (
	"context"

	"github.com/coopersmall/subswag/domain"
	"go.opentelemetry.io/otel/trace"
)

type ITracer interface {
	Trace(context.Context, string, func(context.Context, ISpan) error)
}

type tracer struct {
	t trace.Tracer
}

func (t *tracer) Trace(ctx context.Context, name string, fn func(context.Context, ISpan) error) {
	ctx, span := t.t.Start(ctx, name)
	defer span.End()
	iSpan := newSpan(span)

	iSpan.SetAttribute("correlationId", domain.GetCorrelationIDFromContext(ctx))
	err := fn(ctx, iSpan)
	if err != nil {
		span.RecordError(err)
	}
	iSpan.SetStatus(SpanCodeOk, "")
}

type MockTracer struct {
	*mockTracer
}

type mockTracer struct{}

func (m *mockTracer) Trace(ctx context.Context, name string, fn func(context.Context, ISpan) error) {
	fn(ctx, &mockSpan{})
}
