package tracer

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gcptrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func Setup() (func(), error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	serviceName := os.Getenv("GAE_SERVICE")
	version := os.Getenv("GAE_VERSION")

	exporter, err := gcptrace.New(gcptrace.WithProjectID(projectID))
	if err != nil {
		return nil, fmt.Errorf("new exporter: %v", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String(version),
			),
		),
	)

	otel.SetTracerProvider(provider)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := provider.ForceFlush(ctx); err != nil {
			log.Printf("provider trace flush: %v", err)
		}

		if err := provider.Shutdown(ctx); err != nil {
			log.Printf("provider shutdown: %v", err)
		}
	}, nil
}

func Span(t trace.Tracer, parent context.Context, spanName string, f func(child context.Context, span trace.Span) error) error {
	child, span := t.Start(parent, spanName)
	defer span.End()

	return f(child, span)
}

func NewContext(ctx context.Context, traceID, spanID string, isSampled bool) (context.Context, error) {
	tID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return nil, fmt.Errorf("traceID from hex(%v): %v", traceID, err)
	}

	// hex encoded span-id must have length equals to 16
	sID, err := trace.SpanIDFromHex(spanID[:16])
	if err != nil {
		return nil, fmt.Errorf("spanID from hex(%v): %v", spanID, err)
	}

	flags := trace.TraceFlags(00)
	if isSampled {
		flags = 01
	}

	return trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tID,
		SpanID:     sID,
		TraceFlags: flags,
		Remote:     false,
	})), nil
}
