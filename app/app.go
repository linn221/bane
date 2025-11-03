package app

import (
	"net/http"

	"github.com/linn221/bane/config"
	"github.com/linn221/bane/services"
	"gorm.io/gorm"
)

type Middleware func(http.Handler) http.Handler

type App struct {
	Deducer            *Deducer
	DB                 *gorm.DB
	Cache              config.CacheService
	Services           *services.MyServices
	RecoveryMiddleware Middleware
	LoggingMiddleware  Middleware
}

func NewApp(db *gorm.DB, cache config.CacheService) *App {
	app := &App{
		DB:                 db,
		Cache:              cache,
		RecoveryMiddleware: recovery,
		LoggingMiddleware:  loggingMiddleware,
		Deducer:            &Deducer{},
	}
	app.Services = services.NewMyServices(db, cache, app.Deducer)
	return app
}

func (app *App) WrapMiddlewares(mux *http.ServeMux, mdwares ...Middleware) http.Handler {
	handler := app.LoggingMiddleware(app.RecoveryMiddleware(mux))
	for _, middleware := range mdwares {
		handler = middleware(handler)
	}
	return handler
}
