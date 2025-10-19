package loaders

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/linn221/bane/app"
)

type RIdReader struct {
	deducer *app.Deducer
}

func (r *RIdReader) GetRIds(ctx context.Context, keys []app.Reference) []*dataloader.Result[int] {
	// Create dataloader results in the same order as requested IDs
	unlock := r.deducer.Lock()
	defer unlock()

	loaderResults := make([]*dataloader.Result[int], 0, len(keys))
	references := make([]app.Reference, 0, len(keys))
	for idx, ref := range keys {
		loaderResults = append(loaderResults, &dataloader.Result[int]{Data: idx + 1}) // Start RId from 1
		references = append(references, ref)
	}
	// re assign references
	r.deducer.References = references

	return loaderResults
}

func GetRId(ctx context.Context, refKey app.Reference) (int, error) {
	loaders := For(ctx)
	return loaders.RIdLoader.Load(ctx, refKey)()
}
