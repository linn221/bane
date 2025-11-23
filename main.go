package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/linn221/bane/app"
	"github.com/linn221/bane/config"
	"github.com/linn221/bane/loaders"
	"github.com/linn221/bane/middlewares"
	"github.com/linn221/bane/utils"
)

func main() {
	// Initialize dependencies
	db := config.ConnectSQLite()
	cache := config.NewInMemoryCache()

	app := app.NewApp(db, cache)

	port := utils.GetEnv("PORT", "6423")

	mux := SetupRoutes(app)

	secretConfig := middlewares.SecretConfig{
		Host:        "http://localhost:" + port,
		SecretPath:  "start-session",
		RedirectUrl: "http://localhost:" + port + "/graphql",
		SecretFunc: func() string {
			// return utils.GenerateRandomString(20)
			secretFilename := "secret.txt"
			bs, err := os.ReadFile(secretFilename)
			if err != nil {
				secret := utils.GenerateRandomString(20)
				err := os.WriteFile(secretFilename, []byte(secret), 0666)
				if err != nil {
					panic(err)
				}
				return secret
			}
			return string(bs)
		},
	}
	secretMiddleware := secretConfig.Middleware()
	loaderMiddleware := loaders.LoaderMiddleware(app.DB)

	// Start server
	srv := http.Server{
		Addr:         ":" + port,
		Handler:      app.WrapMiddlewares(mux, secretMiddleware, loaderMiddleware),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Fatal(srv.ListenAndServe())
}
