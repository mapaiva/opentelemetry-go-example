package http

import (
	"net/http"
	"opentel/clients"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

type retrieveUserHandler struct {
	client clients.GithubAPI
}

func (u *retrieveUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	labeler, _ := otelhttp.LabelerFromContext(ctx)
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "retrieveUser")
	defer span.End()

	username := chi.URLParam(r, "username")

	user, err := u.client.UserByUsername(ctx, username)
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
