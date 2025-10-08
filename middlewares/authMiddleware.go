package middlewares

import (
	"net/http"

	"github.com/linn221/bane/app"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, ok := app.GetContextKeyAuth(r.Context())
		if !ok || !auth {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
