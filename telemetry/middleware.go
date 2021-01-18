package telemetry

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

// Telemetry ...
type Telemetry interface {
	Init(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type telemetry struct {
	serviceName       string
	serviceVersion    string
	collectorEndpoint string
	exporter          *otlp.Exporter
	traceProvider     *sdktrace.TracerProvider
}

// New creates a new telemetry adapter.
func New(serviceName, serviceVersion, collectorEndpoint string) Telemetry {
	return &telemetry{
		serviceName:       serviceName,
		serviceVersion:    serviceVersion,
		collectorEndpoint: collectorEndpoint,
	}
}

// Init inits open telemetry components.
func (t *telemetry) Init(ctx context.Context) error {
	driver := otlpgrpc.NewDriver(otlpgrpc.WithInsecure(), otlpgrpc.WithEndpoint(t.collectorEndpoint))
	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		return err
	}

	res := resource.NewWithAttributes(
		semconv.ServiceNameKey.String(t.serviceName),
		semconv.ServiceVersionKey.String(t.serviceVersion),
		semconv.TelemetrySDKNameKey.String("opentelemetry"),
		semconv.TelemetrySDKLanguageKey.String("go"),
		semconv.TelemetrySDKVersionKey.String("0.16.0"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5),
			sdktrace.WithMaxExportBatchSize(10),
		),
		sdktrace.WithResource(res),
	)

	t.traceProvider = tp
	t.exporter = exporter

	otel.SetTracerProvider(tp)

	return nil
}

// Shutdown shutsdown open telemetry components.
func (t *telemetry) Shutdown(ctx context.Context) error {
	if err := t.traceProvider.Shutdown(ctx); err != nil {
		return err
	}

	if err := t.exporter.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// Midlleware represents a tracing HTTP middleware.
func Midlleware(serviceName string) func(next http.Handler) http.Handler {
	// FIXME: Implement our own metrics set.
	fn := func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, serviceName, otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))
	}

	return fn
}
