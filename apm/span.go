package apm

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type SpanCode codes.Code

const (
	SpanCodeUnset = SpanCode(codes.Unset)
	SpanCodeOk    = SpanCode(codes.Ok)
	SpanCodeError = SpanCode(codes.Error)
)

type ISpan interface {
	RecordError(error)
	SetStatus(SpanCode, string)
	SetAttribute(string, any)
	AddEvent(string)
	End()
}

type span struct {
	span trace.Span
}

func newSpan(s trace.Span) ISpan {
	return &span{
		span: s,
	}
}

func (s *span) RecordError(err error) {
	s.span.RecordError(err)
}

func (s *span) SetStatus(code SpanCode, message string) {
	s.span.SetStatus(codes.Code(code), message)
}

func (s *span) SetAttribute(key string, value any) {
	switch value.(type) {
	case string:
		s.span.SetAttributes(attribute.String(key, value.(string)))
	case int:
		s.span.SetAttributes(attribute.Int(key, value.(int)))
	case int64:
		s.span.SetAttributes(attribute.Int64(key, value.(int64)))
	case float64:
		s.span.SetAttributes(attribute.Float64(key, value.(float64)))
	case bool:
		s.span.SetAttributes(attribute.Bool(key, value.(bool)))
	default:
		s.span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", value)))
	}
}

func (s *span) AddEvent(name string) {
	s.span.AddEvent(name)
}

func (s *span) End() {
	s.span.End()
}

type mockSpan struct{}

func (s *mockSpan) RecordError(error) {}
func (s *mockSpan) SetStatus(SpanCode, string) {
}
func (s *mockSpan) SetAttribute(string, any) {}
func (s *mockSpan) AddEvent(string)          {}
func (s *mockSpan) End()                     {}
