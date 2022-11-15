package telemetryutils

import (
	"context"

	"github.com/flyteorg/flytestdlib/contextutils"

	//"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

/*func NewSpan(ctx context.Context, tracerName string, spanName string) (context.Context, trace.Span) {
	var attributes []attribute.KeyValue
	for key, value := range contextutils.GetLogFields(ctx) {
		if value, ok := value.(string); ok {
			attributes = append(attributes, attribute.String(key, value))
		}
	}

	return otel.Tracer(tracerName).Start(ctx, spanName, trace.WithAttributes(attributes...))
}*/

func NewSpan(ctx context.Context, serviceName string, spanName string) (context.Context, trace.Span) {
	// TODO @hamersaw - can check ctx.IDK to see if tracing is enabled on

	var attributes []attribute.KeyValue
	for key, value := range contextutils.GetLogFields(ctx) {
		if value, ok := value.(string); ok {
			attributes = append(attributes, attribute.String(key, value))
		}
	}

	tracerProvider := GetTracerProvider(serviceName)
	return tracerProvider.Tracer("default").Start(ctx, spanName, trace.WithAttributes(attributes...))
	// TODO @hamersaw - tracer should be the library name - ex. github.com/flyteorg/flytepropeller
}
