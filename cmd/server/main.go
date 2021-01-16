package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"opentel/clients"
	transportHTTP "opentel/transport/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

const (
	serviceName    = "otel-server"
	serviceVersion = "0.0.1"
	port           = "2021"
)

func main() {
	ctx := context.Background()

	driver := otlpgrpc.NewDriver(otlpgrpc.WithInsecure(), otlpgrpc.WithEndpoint("localhost:55680"))
	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		log.Fatal(err)
	}

	res := resource.NewWithAttributes(
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
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

	otel.SetTracerProvider(tp)

	githubAPI := clients.GithubAPI{
		HTTPClient: &http.Client{
			Transport: clients.NewTracingTransport(http.DefaultTransport),
		},
		URL: "https://api.github.com",
	}

	r := transportHTTP.NewRouter(githubAPI)
	h := otelhttp.NewHandler(r, serviceName, otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))

	fmt.Printf("Running at port %s...", port)
	if err := transportHTTP.ListenAndServe(":"+port, h); err != nil {
		log.Fatal(err)
	}
}
