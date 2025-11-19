package loaders

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/linn221/bane/app"
)

type RIdReader struct {
}

func (r *RIdReader) GetRIds(ctx context.Context, keys []app.Reference) []*dataloader.Result[int] {
	// Create dataloader results in the same order as requested IDs
	loaderResults := make([]*dataloader.Result[int], 0, len(keys))
	for idx := range keys {
		loaderResults = append(loaderResults, &dataloader.Result[int]{Data: idx + 1}) // Start RId from 1
	}

	return loaderResults
}

func GetRId(ctx context.Context, refKey app.Reference) (int, error) {
	loaders := For(ctx)
	return loaders.RIdLoader.Load(ctx, refKey)()
}
