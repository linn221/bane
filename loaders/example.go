package loaders

import (
	"net/http"

	"github.com/linn221/bane/app"
	"gorm.io/gorm"
)

// Example of how to use the LoaderMiddleware in your HTTP server setup
// This shows how to integrate the dataloader middleware with your existing app structure

func ExampleUsage(db *gorm.DB) {
	// Create your app instance
	app := app.NewApp(db, nil) // Replace nil with your cache service if you have one

	// Create HTTP mux
	mux := http.NewServeMux()

	// Add your routes here
	// mux.HandleFunc("/graphql", yourGraphQLHandler)

	// Wrap middlewares including the dataloader middleware
	handler := app.WrapMiddlewares(mux, LoaderMiddleware(db))

	// Start your server
	// http.ListenAndServe(":8080", handler)
	_ = handler // Prevent unused variable warning
}

// The middleware will automatically inject loaders into the request context
// Your GraphQL resolvers can then use:
// - loaders.GetProgram(ctx, id) to load a single Program
// - loaders.GetNotesForProgram(ctx, programId) to load Notes for a Program
// - loaders.GetPrograms(ctx, ids) to load multiple Programs
// - loaders.GetNotesForPrograms(ctx, programIds) to load Notes for multiple Programs
