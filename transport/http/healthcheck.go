package http

import (
	"net/http"

	"github.com/go-chi/render"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status": "working",
	})
}
