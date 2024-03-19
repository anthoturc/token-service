package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	service = "token-service"
)

func traceProvider(envName string) (*tracesdk.TracerProvider, error) {
	exp, err := otlptrace.New(context.Background(), otlptracehttp.NewClient(otlptracehttp.WithInsecure()))
	if err != nil {
		return nil, fmt.Errorf("trace provider: %v", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", envName),
			attribute.Int64("id", 1),
		)),
	)

	return tp, nil
}

// NewSpan is a utility method to make span creation, from a given context, easier
// The span that is created needs to be ended by the caller.
func NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	ctx, span := otel.Tracer("").Start(ctx, fmt.Sprintf("token-service-%s", name))
	return ctx, span
}
