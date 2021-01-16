package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const (
	serviceName    = "test-otel"
	serviceVersion = "0.0.1"
)

// Handler is an interface that wraps the standard library's http.Handler.
type Handler interface {
	http.Handler
}

func main() {
	r := NewRouter()

	fmt.Printf("Running at port %d...", 2021)
	if err := http.ListenAndServe(":2021", otelhttp.NewHandler(r, serviceName, otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))); err != nil {
		log.Fatal(err)
	}
}

// NewRouter ...
func NewRouter() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	ctx := context.Background()

	driver := otlpgrpc.NewDriver(otlpgrpc.WithInsecure(), otlpgrpc.WithEndpoint("localhost:55680"))
	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		log.Fatal(err)
	}

	res := resource.NewWithAttributes(
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceName),
		semconv.TelemetrySDKNameKey.String("opentelemetry"),
		semconv.TelemetrySDKLanguageKey.String("go"),
		semconv.TelemetrySDKVersionKey.String("0.16.0"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithBatchTimeout(5),
			sdktrace.WithMaxExportBatchSize(10),
		),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	// exporter, err := prometheus.InstallNewPipeline(prometheus.Config{})
	// if err != nil {
	// 	log.Fatal(err) }
	// metrics := func(w http.ResponseWriter, req *http.Request) {
	// 	exporter.ServeHTTP(w, req)
	// }

	githubAPI := Github{
		HTTPClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport,
				otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
					return fmt.Sprintf("%s %s", r.Method, r.URL.String())
				}),
			),
		},
		URL: "https://api.github.com",
	}

	r.Get("/healthcheck", healthcheck)
	// r.Get("/metrics", metrics)
	r.Method(http.MethodGet, "/users/{username}", &retrieveUserHandler{githubAPI})

	return r
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status": "working",
	})
}

type retrieveUserHandler struct {
	client Github
}

func (u *retrieveUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	labeler, _ := otelhttp.LabelerFromContext(ctx)
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "retrieveUser")
	defer span.End()

	username := chi.URLParam(r, "username")

	user, err := u.client.User(ctx, username)
	if err != nil {
		render.JSON(w, r, map[string]string{
			"mesage": err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		labeler.Add(label.Bool("error", true))
		return
	}

	render.JSON(w, r, user)
}

type Github struct {
	HTTPClient *http.Client
	URL        string
}

func (g Github) User(ctx context.Context, username string) (map[string]interface{}, error) {
	// ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "github.user")
	// defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%s", g.URL, username), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := g.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var userResponse map[string]interface{}
	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(&userResponse)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return userResponse, nil
}
