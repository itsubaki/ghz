package tracer

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	otltrace "go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	t otltrace.Tracer
	p *sdktrace.TracerProvider
}

func New(projectID, path string) (*Tracer, error) {
	exporter, err := trace.New(trace.WithProjectID(projectID))
	if err != nil {
		return nil, fmt.Errorf("new exporter: %v", err)
	}
	provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))

	return &Tracer{
		t: provider.Tracer(path),
		p: provider,
	}, nil
}

func (t *Tracer) ForceFlush(ctx context.Context) {
	t.p.ForceFlush(ctx)
}

func (t *Tracer) Start(ctx context.Context, spanName string, opts ...otltrace.SpanStartOption) (context.Context, otltrace.Span) {
	return t.t.Start(ctx, spanName, opts...)
}
