package tracer_test

import (
	"context"
	"testing"

	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/tracer"
)

// GOOGLE_APPLICATION_CREDENTIALS=../../credentials.json go test ./appengine/tracer
func TestTracer(t *testing.T) {
	traceID, spanID := "fe86487e2a6a0b6b202bd69244be420b", "1234567890123456"
	ctx, err := tracer.NewContext(context.Background(), traceID, spanID)
	if err != nil {
		t.Fatalf("new context: %v", err)
	}

	tra, err := tracer.New(dataset.ProjectID, "TestTracer")
	if err != nil {
		t.Fatalf("new tracer: %v", err)
	}
	defer tra.ForceFlush(ctx)

	parent, span := tra.Start(ctx, "parent")
	defer span.End()

	func(parent context.Context) {
		_, span := tra.Start(parent, "hello")
		defer span.End()
	}(parent)

	func(parent context.Context) {
		_, span := tra.Start(parent, "world")
		defer span.End()
	}(parent)
}
