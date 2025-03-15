package apm

import (
	"context"

	"github.com/coopersmall/subswag/clients"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type ITraceProvider interface {
	NewTracer(service string) ITracer
	Stop(context.Context) error
}

func GetTraceProvider(client clients.IAPMClient) (ITraceProvider, error) {
	return NewTraceProvider(context.Background(), client)
}

type TraceProvider struct {
	exporter *otlptrace.Exporter
}

func NewTraceProvider(ctx context.Context, client clients.IAPMClient) (ITraceProvider, error) {
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}
	return &TraceProvider{
		exporter: exporter,
	}, nil
}

func (t *TraceProvider) NewTracer(service string) ITracer {
	opt := sdktrace.WithResource(resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(service),
	))

	provider := sdktrace.NewTracerProvider(opt)
	provider.RegisterSpanProcessor(sdktrace.NewBatchSpanProcessor(t.exporter))
	return &tracer{
		t: provider.Tracer(service),
	}
}

func (t *TraceProvider) Stop(ctx context.Context) error {
	return t.exporter.Shutdown(ctx)
}
