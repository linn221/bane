package loaders

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

func newGenericReader[T any, K comparable](db *gorm.DB, getKey func(T) K, getDefault func(K) T, preloads ...string) *genericReader[T, K] {

	return &genericReader[T, K]{
		db:         db,
		getKey:     getKey,
		getDefault: getDefault,
		preloads:   []string(preloads),
	}
}

type genericReader[T any, K comparable] struct {
	db         *gorm.DB
	preloads   []string
	fetch      func([]K) ([]T, error)
	getDefault func(key K) T
	getKey     func(T) K
}

func (r genericReader[T, K]) Loader() *dataloader.Loader[K, T] {
	return dataloader.NewBatchedLoader(r.BatchFunc, dataloader.WithWait[K, T](time.Millisecond))

}

func (r genericReader[T, K]) BatchFunc(ctx context.Context, ids []K) []*dataloader.Result[T] {

	var results []T
	var err error
	if r.fetch == nil {

		dbCtx := r.db.WithContext(ctx).Where("id IN ?", ids)
		if len(r.preloads) > 0 {
			dbCtx = dbCtx.Preload(r.preloads[0])
			if len(r.preloads) > 1 {
				for _, p := range r.preloads[1:] {
					dbCtx.Preload(p)
				}
			}
		}

		err = dbCtx.Find(&results).Error
	} else {
		results, err = r.fetch(ids)
	}
	if err != nil {
		return handleError[T](len(ids), err)
	}

	// generate resultMap from results
	resultMap := make(map[K]T, len(results)+1)
	for _, result := range results {
		resultMap[r.getKey(result)] = result
	}

	loaderResults := make([]*dataloader.Result[T], 0, len(ids))
	for _, id := range ids {
		data, ok := resultMap[id]
		if !ok {
			data = r.getDefault(id)
		}
		loaderResults = append(loaderResults, &dataloader.Result[T]{Data: data})
	}
	return loaderResults
}

func newGenericReaderSlice[T any, K comparable](db *gorm.DB, getKey func(T) K, getDefault func(K) T, foreignIdColumn string, preloads ...string) *genericReaderSlice[T, K] {
	return &genericReaderSlice[T, K]{
		db:              db,
		getKey:          getKey,
		getDefault:      getDefault,
		preloads:        []string(preloads),
		foreignIdColumn: foreignIdColumn,
	}
}

// func genericReorderDataloaderResults[T any, K comparable](results []T, keys []K, getKey func(T) K, getDefault func(K) T) []*dataloader.Result[T] {
// 	m := make(map[K]T, len(results))
// 	for _, r := range results {
// 		key := getKey(r)
// 		m[key] = r
// 	}

// 	dataloaderResult := make([]*dataloader.Result[T], 0, len(keys))
// 	for _, k := range keys {
// 		result, ok := m[k]
// 		if !ok {
// 			result = getDefault(k)
// 		}
// 		dataloaderResult = append(dataloaderResult, &dataloader.Result[T]{Data: result})
// 	}

// 	return dataloaderResult
// }

type genericReaderSlice[T any, K comparable] struct {
	db              *gorm.DB
	preloads        []string
	fetch           func([]K) ([]T, error)
	foreignIdColumn string
	getDefault      func(key K) T
	getKey          func(T) K
}

func (r genericReaderSlice[T, K]) Loader() *dataloader.Loader[K, []T] {
	return dataloader.NewBatchedLoader(r.BatchFunc, dataloader.WithWait[K, []T](time.Millisecond))
}

func (r genericReaderSlice[T, K]) BatchFunc(ctx context.Context, ids []K) []*dataloader.Result[[]T] {
	var results []T
	var err error
	if r.fetch == nil {
		if r.foreignIdColumn == "" {
			r.foreignIdColumn = "id"
		}
		dbCtx := r.db.WithContext(ctx).Where(r.foreignIdColumn+" IN ?", ids)
		if len(r.preloads) > 0 {
			dbCtx = dbCtx.Preload(r.preloads[0])
			if len(r.preloads) > 1 {
				for _, p := range r.preloads[1:] {
					dbCtx.Preload(p)
				}
			}
		}

		err = dbCtx.Find(&results).Error
	} else {
		results, err = r.fetch(ids)
	}
	if err != nil {
		return handleError[[]T](len(ids), err)
	}

	// generate resultMap from results - grouping by key into slices
	resultMap := make(map[K][]T, len(results)+1)
	for _, result := range results {
		key := r.getKey(result)
		resultMap[key] = append(resultMap[key], result)
	}

	loaderResults := make([]*dataloader.Result[[]T], 0, len(ids))
	for _, id := range ids {
		data, ok := resultMap[id]
		if !ok {
			// Return empty slice if not found
			data = []T{}
		}
		loaderResults = append(loaderResults, &dataloader.Result[[]T]{Data: data})
	}
	return loaderResults
}

// type HasLoaderKey interface {
// 	GetLoaderKey() int
// }
// type genericManyResultLoader[T HasLoaderKey] struct {
// 	db              *gorm.DB
// 	foreignIdColumn string
// }

// func (r *genericManyResultLoader[T]) BatchFunc(ctx context.Context, ids []int) []*dataloader.Result[[]*T] {

// 	var results []*T
// 	err := r.db.WithContext(ctx).Where(r.foreignIdColumn+" IN ?", ids).Find(&results).Error
// 	if err != nil {
// 		return handleError[[]*T](len(ids), err)
// 	}

// 	resultMap := make(map[int][]*T)
// 	for _, result := range results {
// 		resultMap[(*result).GetLoaderKey()] = append(resultMap[(*result).GetLoaderKey()], result)
// 	}
// 	var loaderResults []*dataloader.Result[[]*T]
// 	// reordering the results according to ids
// 	for _, id := range ids {
// 		results, ok := resultMap[id]
// 		if !ok {
// 			var v []*T
// 			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: v})
// 		} else {
// 			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: results})
// 		}
// 	}
// 	return loaderResults
// }

// type genericArrayResultReader[T any] struct {
// 	db              *gorm.DB
// 	foreignIdColumn string
// 	getLoaderKey    func(*T) int
// }

// func (r *genericArrayResultReader[T]) BatchFunc(ctx context.Context, ids []int) []*dataloader.Result[[]*T] {
// 	var results []*T
// 	err := r.db.WithContext(ctx).Where(r.foreignIdColumn+" IN ?", ids).Find(&results).Error
// 	if err != nil {
// 		return handleError[[]*T](len(ids), err)
// 	}

// 	resultMap := make(map[int][]*T)
// 	for _, result := range results {
// 		resultMap[r.getLoaderKey(result)] = append(resultMap[r.getLoaderKey(result)], result)
// 	}
// 	var loaderResults []*dataloader.Result[[]*T]
// 	// reordering the results according to ids
// 	for _, id := range ids {
// 		results, ok := resultMap[id]
// 		if !ok {
// 			var v []*T
// 			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: v})
// 		} else {
// 			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: results})
// 		}
// 	}
// 	return loaderResults

// }

// type genericArrayResultReader2[T any] struct {
// 	db           *gorm.DB
// 	fetch        func(db *gorm.DB, ids []int) ([]*T, error)
// 	getLoaderKey func(*T) int
// }

// func (r *genericArrayResultReader2[T]) BatchFunc(ctx context.Context, ids []int) []*dataloader.Result[[]*T] {
// 	// var results []*T
// 	// err := r.db.WithContext(ctx).Where(r.foreignIdColumn+" IN ?", ids).Find(&results).Error
// 	// if err != nil {
// 	// 	return handleError[[]*T](len(ids), err)
// 	// }

// 	results, err := r.fetch(r.db.WithContext(ctx), ids)
// 	if err != nil {
// 		return handleError[[]*T](len(ids), err)

// 	}

// 	resultMap := make(map[int][]*T)
// 	for _, result := range results {
// 		resultMap[r.getLoaderKey(result)] = append(resultMap[r.getLoaderKey(result)], result)
// 	}
// 	var loaderResults []*dataloader.Result[[]*T]
// 	// reordering the results according to ids
// 	for _, id := range ids {
// 		results, ok := resultMap[id]
// 		if !ok {
// 			var v []*T
// 			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: v})
// 		} else {
// 			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: results})
// 		}
// 	}
// 	return loaderResults

// }
