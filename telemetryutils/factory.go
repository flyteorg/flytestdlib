package telemetryutils

import (
	"os"

	"github.com/flyteorg/flytestdlib/version"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	rawtrace "go.opentelemetry.io/otel/trace"
)

const (
	AdminClientTracer       = "admin-client"
	BlobstoreClientTracer   = "blobstore-client"
	DataCatalogClientTracer = "datacatalog-client"
	FlytePropellerTracer    = "flytepropeller"
	K8sClientTracer         = "k8s-client"
)

var tracerProviders = make(map[string]*trace.TracerProvider)
var noopTracerProvider = rawtrace.NewNoopTracerProvider()

func RegisterTracerProvider(serviceName string, config *Config) error {
	if config == nil {
		return nil
	}

	var opts []trace.TracerProviderOption
	if config.FileConfig.Enabled {
		// configure file exporter
		f, err := os.Create(config.FileConfig.Filename)
		if err != nil {
			return err
		}

		exporter, err := stdouttrace.New(
			stdouttrace.WithWriter(f),
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			return err
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
			return err
		}

		opts = append(opts, trace.WithBatcher(exporter))
	}

	// if no exporters are enabled then we can return
	if len(opts) == 0 {
		return nil
	}

	telemetryResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version.Version),
		),
	)
	if err != nil {
		return err
	}

	opts = append(opts, trace.WithResource(telemetryResource))
	tracerProvider := trace.NewTracerProvider(opts...)

	tracerProviders[serviceName] = tracerProvider
	return nil
}

func GetTracerProvider(serviceName string) rawtrace.TracerProvider {
	if t, ok := tracerProviders[serviceName]; ok {
		return t
	}

	return noopTracerProvider
}
