package apm

import (
	"strings"
	"time"

	"github.com/coopersmall/subswag/utils"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

type IAPMClient otlptrace.Client

type APMClient struct {
	otlptrace.Client
}

func NewAPMClient(
	getAPMURL func() (string, error),
	getAPMTimeout func() (time.Duration, error),
) *APMClient {
	url, err := getAPMURL()
	if err != nil {
		panic(utils.NewInternalError("failed to get APM URL for client", err))
	}
	timeout, err := getAPMTimeout()
	if err != nil {
		panic(utils.NewInternalError("failed to get APM timeout for client", err))
	}
	opts := make([]otlptracehttp.Option, 0)
	opts = append(opts, otlptracehttp.WithTimeout(timeout))
	opts = append(opts, otlptracehttp.WithEndpoint(url))

	if strings.Contains(url, "localhost") {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	return &APMClient{
		otlptracehttp.NewClient(opts...),
	}
}
