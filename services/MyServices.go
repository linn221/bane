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
}

// NewMyServices creates a new MyServices instance with all services initialized
func NewMyServices(db *gorm.DB, cache config.CacheService, deducer Deducer) *MyServices {
	return &MyServices{
		EndpointService:       newEndpointService(db, deducer),
		TagService:            &tagService{db: db},
		NoteService:           newNoteService(db, deducer),
		ProgramService:        newProgramService(db),
		MyRequestService:      newMyRequestService(db, deducer),
		MemorySheetService:    newMemorySheetService(db),
		WordService:           newWordService(db),
		VulnConnectionService: newVulnConnectionService(db),
	}
}
