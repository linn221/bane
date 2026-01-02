package main

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/linn221/bane/app"
	"github.com/linn221/bane/graph"
	"github.com/linn221/bane/graph/resolvers"
)

func SetupRoutes(app *app.App) *http.ServeMux {

	// Create resolver with dependencies
	resolver := resolvers.NewResolver(app)

	// Create auth handler
	// Create a new mux
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Authentication routes (no middleware needed)
	// mux.HandleFunc("/login", authHandler.Login)

	// Create GraphQL handler
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Serve GraphQL Playground
			playground.Handler("GraphQL Playground", "/graphql").ServeHTTP(w, r)
		case http.MethodPost:
			// Serve GraphQL queries/mutations
			srv.ServeHTTP(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
