package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"opentel/clients"
	"opentel/telemetry"
	transportHTTP "opentel/transport/http"

	"go.opentelemetry.io/otel/exporters/metric/prometheus"
)

const (
	serviceName       = "otel-server"
	serviceVersion    = "0.0.1"
	collectorEndpoint = "localhost:55680"
	port              = "2021"
)

func main() {
	prometheusExporter, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		log.Fatal(err)
	}

	telemetry := telemetry.New(serviceName, serviceVersion, collectorEndpoint)
	if err := telemetry.Init(context.Background()); err != nil {
		log.Fatal(err)
	}
	defer telemetry.Shutdown(context.Background())

	githubAPI := clients.GithubAPI{
		HTTPClient: &http.Client{
			Transport: clients.NewTracingTransport(http.DefaultTransport),
		},
		URL: "https://api.github.com",
	}

	r := transportHTTP.NewRouter(
		serviceName,
		prometheusExporter,
		githubAPI,
	)

	fmt.Printf("Running at port %s...", port)
	if err := transportHTTP.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
