package loaders

import (
	"context"
	"net/http"

	"github.com/linn221/bane/app"
	"gorm.io/gorm"
)

// LoaderMiddleware creates a middleware that injects dataloaders into the request context
// This middleware should be used in your HTTP server setup to enable dataloader functionality
// The middleware creates new loader instances for each request to ensure data consistency
func LoaderMiddleware(db *gorm.DB) app.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create new loaders for this request
			loaders := NewLoaders(db)

			// Inject loaders into the request context
			ctx := r.Context()
			ctx = contextWithLoaders(ctx, loaders)

			// Update the request with the new context
			r = r.WithContext(ctx)

			// Continue to the next handler
			h.ServeHTTP(w, r)
		})
	}
}

// contextWithLoaders adds the loaders to the context
func contextWithLoaders(ctx context.Context, loaders *Loaders) context.Context {
	return context.WithValue(ctx, loadersKey, loaders)
}
