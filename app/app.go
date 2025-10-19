package app

import (
	"net/http"

	"github.com/linn221/bane/config"
	"gorm.io/gorm"
)

type Middleware func(http.Handler) http.Handler

type App struct {
	Deducer            *Deducer
	DB                 *gorm.DB
	Cache              config.CacheService
	RecoveryMiddleware Middleware
	LoggingMiddleware  Middleware
}

func NewApp(db *gorm.DB, cache config.CacheService) *App {
	return &App{
		DB:                 db,
		Cache:              cache,
		RecoveryMiddleware: recovery,
		LoggingMiddleware:  loggingMiddleware,
		Deducer:            &Deducer{},
	}
}

func (app *App) WrapMiddlewares(mux *http.ServeMux, mdwares ...Middleware) http.Handler {
	handler := app.LoggingMiddleware(app.RecoveryMiddleware(mux))
	for _, middleware := range mdwares {
		handler = middleware(handler)
	}
	return handler
}
