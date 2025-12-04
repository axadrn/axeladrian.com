package ui

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	err := c.Render(r.Context(), w)
	if err != nil {
		slog.Error("render failed", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
