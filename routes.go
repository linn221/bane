package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/linn221/bane/app"
	"github.com/linn221/bane/graph"
	"github.com/linn221/bane/graph/resolvers"
	"github.com/linn221/bane/services"
	"github.com/linn221/bane/utils"
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

	mux.HandleFunc("/importWordlist/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if r.Method == http.MethodGet {
			t := template.Must(template.ParseFiles("views/importWordlist.html"))
			err := t.Execute(w, map[string]interface{}{
				"WordListId": id,
			})
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			// Parse the wordlist ID from the URL
			wordListId, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, "Invalid wordlist ID", http.StatusBadRequest)
				return
			}

			// Read words from the uploaded file using the utility function
			words, err := utils.ReadFileFromForm(r, "file")
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusBadRequest)
				return
			}

			// Add words to wordlist
			err = services.AddWordsToWordList(app.DB, wordListId, words)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to add words: %v", err), http.StatusInternalServerError)
				return
			}

			// Return success response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Successfully added %d words to wordlist", len(words))))
		}
	})

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h2>hello world!</h2>"))
	})

	return mux
}
