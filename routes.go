package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/linn221/bane/app"
	"github.com/linn221/bane/graph"
	"github.com/linn221/bane/graph/resolvers"
	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"github.com/linn221/bane/views"
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
			// For now, return a simple HTML response for the import wordlist page
			// TODO: Create a templ template for this page
			html := fmt.Sprintf(`
				<!DOCTYPE html>
				<html>
				<head>
					<title>Import Wordlist</title>
					<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
				</head>
				<body>
					<div class="container mt-4">
						<h1>Import Wordlist #%s</h1>
						<form method="POST" enctype="multipart/form-data">
							<div class="mb-3">
								<label for="file" class="form-label">Select file to import:</label>
								<input type="file" class="form-control" id="file" name="file" required>
							</div>
							<button type="submit" class="btn btn-primary">Import</button>
						</form>
					</div>
				</body>
				</html>
			`, id)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
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
			err = app.Services.WordService.AddWordsToWordList(wordListId, words)
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

	// Memory Sheet routes
	mux.HandleFunc("GET /memory-sheets", func(w http.ResponseWriter, r *http.Request) {
		// Get all memory sheets for today
		today := utils.Today()
		memorySheets, err := app.Services.MemorySheetService.GetTodayNotes(today)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get memory sheets: %v", err), http.StatusInternalServerError)
			return
		}

		// Render the memory sheet list page
		component := views.MemorySheetList(memorySheets)
		err = component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /memory-sheets/create", func(w http.ResponseWriter, r *http.Request) {
		// Render the memory sheet creation form
		component := views.MemorySheetForm()
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /memory-sheets", func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		value := r.FormValue("value")
		if value == "" {
			http.Error(w, "Value is required", http.StatusBadRequest)
			return
		}

		alias := r.FormValue("alias")
		dateStr := r.FormValue("date")

		// Create new memory sheet
		newMemorySheet := &models.MemorySheetInput{
			Value: value,
			Alias: alias,
		}

		// Handle custom date if provided
		if dateStr != "" {
			customDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				http.Error(w, "Invalid date format", http.StatusBadRequest)
				return
			}
			newMemorySheet.Date = &models.MyDate{Time: customDate}
		}

		// Save to database
		_, err = app.Services.MemorySheetService.Create(newMemorySheet)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create memory sheet: %v", err), http.StatusInternalServerError)
			return
		}

		// Redirect to memory sheets list
		http.Redirect(w, r, "/memory-sheets", http.StatusSeeOther)
	})

	return mux
}
