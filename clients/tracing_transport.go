package clients

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewTracingTransport returns a new openteletry transport with provisioned metrics set.
func NewTracingTransport(transport http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(transport,
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
	)
}

func spanNameFormatter(operation string, r *http.Request) string {
	return fmt.Sprintf("%s %s", r.Method, r.URL.String())
}
