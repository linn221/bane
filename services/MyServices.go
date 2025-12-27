package services

import (
	"github.com/linn221/bane/config"
	"gorm.io/gorm"
)

// MyServices contains all service instances
type MyServices struct {
	EndpointService *endpointService
	NoteService     *noteService
	MyRequestService *myRequestService
	WordService     *wordService
	ProjectService  *projectService
	AliasService    *aliasService
}

// NewMyServices creates a new MyServices instance with all services initialized
func NewMyServices(db *gorm.DB, cache config.CacheService) *MyServices {
	aliasService := &aliasService{db: db}

	endpointService := &endpointService{
		db:           db,
		aliasService: aliasService,
	}

	noteService := &noteService{
		db: db,
	}

	myRequestService := &myRequestService{
		db: db,
	}

	wordService := &wordService{
		db:           db,
		aliasService: aliasService,
	}

	projectService := &projectService{
		db:           db,
		aliasService: aliasService,
	}

	return &MyServices{
		AliasService:     aliasService,
		EndpointService:  endpointService,
		NoteService:      noteService,
		MyRequestService: myRequestService,
		WordService:      wordService,
		ProjectService:   projectService,
	}
}
