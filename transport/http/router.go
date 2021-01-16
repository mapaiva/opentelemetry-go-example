package http

import (
	"net/http"
	"opentel/clients"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// NewRouter returns an HTTP router.
func NewRouter(githubAPI clients.GithubAPI) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/healthcheck", healthcheck)
	r.Method(http.MethodGet, "/users/{username}", &retrieveUserHandler{githubAPI})

	return r
}

// ListenAndServe serves requests to an address routed by a given http.Handler.
func ListenAndServe(addr string, handler http.Handler) error {
	s := http.Server{
		Addr:        addr,
		Handler:     handler,
		ReadTimeout: 2 * time.Second,
	}
	return s.ListenAndServe()
}
