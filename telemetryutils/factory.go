package telemetryutils

import (
	"os"

	"github.com/flyteorg/flytestdlib/version"

	//"go.opentelemetry.io/otel"
	//"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func NewTracerProvider(appName string, config *Config) (*trace.TracerProvider, error) {
	if config == nil {
		return nil, nil
	}

	var opts []trace.TracerProviderOption
	if config.FileConfig.Enabled {
		// configure file exporter
		f, err := os.Create(config.FileConfig.Filename)
		if err != nil {
			return nil, err
		}

		exporter, err := stdouttrace.New(
			stdouttrace.WithWriter(f),
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			return nil, err
		}

		opts = append(opts, trace.WithBatcher(exporter))
	}

	if config.JaegerConfig.Enabled {
		// configure jaeger exporter
		exporter, err := jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(config.JaegerConfig.Endpoint),
			),
		)
		if err != nil {
			return nil, err
		}

		opts = append(opts, trace.WithBatcher(exporter))
	}

	// if no exporters are enabled then we return a nil TracerProvider
	if len(opts) == 0 {
		return nil, nil
	}

	telemetryResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
			semconv.ServiceVersionKey.String(version.Version),
		),
	)
	if err != nil {
		return nil, err
	}

	opts = append(opts, trace.WithResource(telemetryResource))
	return trace.NewTracerProvider(opts...), nil
}
