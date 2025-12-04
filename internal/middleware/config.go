package middleware

import (
	"net/http"

	"github.com/axadrn/axeladrian/internal/config"
	"github.com/axadrn/axeladrian/internal/ctxkeys"
)

func WithConfig(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ctxkeys.WithConfig(r.Context(), cfg)
			ctx = ctxkeys.WithURLPath(ctx, r.URL.Path)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
