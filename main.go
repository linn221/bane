package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/linn221/bane/app"
	"github.com/linn221/bane/config"
	"github.com/linn221/bane/middlewares"
	"github.com/linn221/bane/utils"
)

func main() {
	// Initialize dependencies
	db := config.ConnectMySQL()
	cache := config.ConnectRedis(context.Background())

	app := app.NewApp(db, cache)

	port := utils.GetEnv("PORT", "6423")

	mux := SetupRoutes(app)

	secretConfig := middlewares.SecretConfig{
		Host:        "http://localhost:" + port,
		SecretPath:  "start-session",
		RedirectUrl: "http://localhost:" + port + "/graphql",
		SecretFunc: func() string {
			return utils.GenerateRandomString(20)
		},
	}
	secretMiddleware := secretConfig.Middleware()

	// Start server
	srv := http.Server{
		Addr:         ":" + port,
		Handler:      app.WrapMiddlewares(mux, secretMiddleware),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Fatal(srv.ListenAndServe())
}
