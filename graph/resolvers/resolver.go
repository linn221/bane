package resolvers

import (
	"github.com/linn221/bane/app"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB  *gorm.DB
	app *app.App
	// TagService *services.TagService
	// ProductCategoryService *services.ProductCategoryService
	// ProductTagService      *services.ProductTagService
	// ProductUnitService     *services.ProductUnitService
}

// NewResolver creates a new resolver with dependencies
func NewResolver(app *app.App) *Resolver {
	return &Resolver{
		app: app,
		DB:  app.DB,
		// TagService: &services.TagService{
		// 	DB: app.DB,
		// },
	}
}
