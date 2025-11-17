package services

import (
	"github.com/linn221/bane/config"
	"gorm.io/gorm"
)

// Deducer interface for reading RId values
type Deducer interface {
	ReadRId(rid int) (int, string)
}

// MyServices contains all service instances
type MyServices struct {
	EndpointService       *endpointService
	TagService            *tagService
	NoteService           *noteService
	ProgramService        *programService
	MyRequestService      *myRequestService
	MemorySheetService    *memorySheetService
	WordService           *wordService
	VulnConnectionService *vulnConnectionService
	MySheetService        *mySheetService
	ProjectService        *projectService
	TodoService           *todoService
	AliasService          *aliasService
}

// NewMyServices creates a new MyServices instance with all services initialized
func NewMyServices(db *gorm.DB, cache config.CacheService, deducer Deducer) *MyServices {
	aliasService := newAliasService(db)
	return &MyServices{
		EndpointService:       newEndpointService(db, deducer, aliasService),
		TagService:            newTagService(db, aliasService),
		NoteService:           newNoteService(db, deducer),
		ProgramService:        newProgramService(db, aliasService),
		MyRequestService:      newMyRequestService(db, deducer),
		MemorySheetService:    newMemorySheetService(db, aliasService),
		WordService:           newWordService(db, aliasService),
		VulnConnectionService: newVulnConnectionService(db),
		MySheetService:        newMySheetService(db, aliasService),
		ProjectService:        newProjectService(db, aliasService),
		TodoService:           newTodoService(db, aliasService),
		AliasService:          aliasService,
	}
}
