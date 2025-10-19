package loaders

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

// NotesReader handles efficient loading of Note slices for Programs using dataloader
// This reader implements the BatchFunc interface for dataloader to batch multiple
// Note requests into a single database query, reducing N+1 query problems
type NotesReader struct {
	db            *gorm.DB
	referenceType string
}

// GetNotes is the batch function that loads Notes for multiple Program IDs
// It receives a slice of Program IDs and returns a slice of dataloader.Result containing Note slices
// This function is called by dataloader when multiple Note requests are batched
func (r *NotesReader) GetNotes(ctx context.Context, referenceIds []int) []*dataloader.Result[[]*models.Note] {
	// Query the database for all Notes that belong to the given Program IDs
	// Notes are linked to Programs via polymorphic relationship (ReferenceType="programs", ReferenceID=Program.Id)
	var results []models.Note
	err := r.db.WithContext(ctx).
		Where("reference_type = ? AND reference_id IN ?", r.referenceType, referenceIds).
		Order("note_date DESC"). // Order by most recent notes first
		Find(&results).Error

	if err != nil {
		// If there's an error, return error results for all requested Program IDs
		return handleError[[]*models.Note](len(referenceIds), err)
	}

	// Group Notes by Program ID
	notesByProgramId := make(map[int][]*models.Note)
	for _, note := range results {
		// Create a copy to avoid pointer issues
		noteCopy := note
		notesByProgramId[note.ReferenceID] = append(notesByProgramId[note.ReferenceID], &noteCopy)
	}

	// Create dataloader results in the same order as requested Program IDs
	loaderResults := make([]*dataloader.Result[[]*models.Note], 0, len(referenceIds))
	for _, programId := range referenceIds {
		if notes, exists := notesByProgramId[programId]; exists {
			loaderResults = append(loaderResults, &dataloader.Result[[]*models.Note]{Data: notes})
		} else {
			// If no Notes found for this Program, return empty slice (not an error)
			loaderResults = append(loaderResults, &dataloader.Result[[]*models.Note]{Data: []*models.Note{}})
		}
	}

	return loaderResults
}
