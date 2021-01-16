module opentel

go 1.15

require (
	github.com/go-chi/chi v1.5.1
	github.com/go-chi/render v1.0.1
	go.opencensus.io v0.22.5 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http v0.11.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.15.1
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.15.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp v0.16.0
	go.opentelemetry.io/otel/exporters/stdout v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
)
