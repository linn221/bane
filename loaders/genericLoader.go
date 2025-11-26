package middlewares

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

type Identifier interface {
	GetId() int
}

type genericReader[T models.Identifier] struct {
	db         *gorm.DB
	preloads   []string
	getDefault func(id int) T
}

func (r genericReader[T]) BatchFunc(ctx context.Context, ids []int) []*dataloader.Result[*T] {

	var results []T
	dbCtx := r.db.WithContext(ctx).Where("id IN ?", ids)
	if len(r.preloads) > 0 {
		dbCtx = dbCtx.Preload(r.preloads[0])
		if len(r.preloads) > 1 {
			for _, p := range r.preloads[1:] {
				dbCtx.Preload(p)
			}
		}
	}

	err := dbCtx.Find(&results).Error
	if err != nil {
		return handleError[*T](len(ids), err)
	}

	// generate resultMap from results
	resultMap := make(map[int]T, len(results)+1)
	resultMap[0] = r.getDefault(0)
	for _, result := range results {
		resultMap[result.GetId()] = result
	}

	loaderResults := make([]*dataloader.Result[*T], 0, len(ids))
	for _, id := range ids {
		data, ok := resultMap[id]
		if !ok {
			data = r.getDefault(id)
		}
		loaderResults = append(loaderResults, &dataloader.Result[*T]{Data: &data})
	}
	return loaderResults
}

func genericReorderDataloaderResults[T any, K comparable](results []T, keys []K, getKey func(T) K, getDefault func(K) T) []*dataloader.Result[T] {
	m := make(map[K]T, len(results))
	for _, r := range results {
		key := getKey(r)
		m[key] = r
	}

	dataloaderResult := make([]*dataloader.Result[T], 0, len(keys))
	for _, k := range keys {
		result, ok := m[k]
		if !ok {
			result = getDefault(k)
		}
		dataloaderResult = append(dataloaderResult, &dataloader.Result[T]{Data: result})
	}

	return dataloaderResult
}

func GetLeaveType(ctx context.Context, id int) (*models.LeaveType, error) {
	loaders := For(ctx)
	return loaders.leaveTypeLoader.Load(ctx, id)()
}
func GetEmployee(ctx context.Context, id int) (*models.Employee, error) {
	loaders := For(ctx)
	return loaders.employeeLoader.Load(ctx, id)()
}

type HasLoaderKey interface {
	GetLoaderKey() int
}
type genericManyResultLoader[T HasLoaderKey] struct {
	db              *gorm.DB
	foreignIdColumn string
}

func (r *genericManyResultLoader[T]) BatchFunc(ctx context.Context, ids []int) []*dataloader.Result[[]*T] {

	var results []*T
	err := r.db.WithContext(ctx).Where(r.foreignIdColumn+" IN ?", ids).Find(&results).Error
	if err != nil {
		return handleError[[]*T](len(ids), err)
	}

	resultMap := make(map[int][]*T)
	for _, result := range results {
		resultMap[(*result).GetLoaderKey()] = append(resultMap[(*result).GetLoaderKey()], result)
	}
	var loaderResults []*dataloader.Result[[]*T]
	// reordering the results according to ids
	for _, id := range ids {
		results, ok := resultMap[id]
		if !ok {
			var v []*T
			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: v})
		} else {
			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: results})
		}
	}
	return loaderResults
}

type genericArrayResultReader[T any] struct {
	db              *gorm.DB
	foreignIdColumn string
	getLoaderKey    func(*T) int
}

func (r *genericArrayResultReader[T]) BatchFunc(ctx context.Context, ids []int) []*dataloader.Result[[]*T] {
	var results []*T
	err := r.db.WithContext(ctx).Where(r.foreignIdColumn+" IN ?", ids).Find(&results).Error
	if err != nil {
		return handleError[[]*T](len(ids), err)
	}

	resultMap := make(map[int][]*T)
	for _, result := range results {
		resultMap[r.getLoaderKey(result)] = append(resultMap[r.getLoaderKey(result)], result)
	}
	var loaderResults []*dataloader.Result[[]*T]
	// reordering the results according to ids
	for _, id := range ids {
		results, ok := resultMap[id]
		if !ok {
			var v []*T
			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: v})
		} else {
			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: results})
		}
	}
	return loaderResults

}

type genericArrayResultReader2[T any] struct {
	db           *gorm.DB
	fetch        func(db *gorm.DB, ids []int) ([]*T, error)
	getLoaderKey func(*T) int
}

func (r *genericArrayResultReader2[T]) BatchFunc(ctx context.Context, ids []int) []*dataloader.Result[[]*T] {
	// var results []*T
	// err := r.db.WithContext(ctx).Where(r.foreignIdColumn+" IN ?", ids).Find(&results).Error
	// if err != nil {
	// 	return handleError[[]*T](len(ids), err)
	// }

	results, err := r.fetch(r.db.WithContext(ctx), ids)
	if err != nil {
		return handleError[[]*T](len(ids), err)

	}

	resultMap := make(map[int][]*T)
	for _, result := range results {
		resultMap[r.getLoaderKey(result)] = append(resultMap[r.getLoaderKey(result)], result)
	}
	var loaderResults []*dataloader.Result[[]*T]
	// reordering the results according to ids
	for _, id := range ids {
		results, ok := resultMap[id]
		if !ok {
			var v []*T
			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: v})
		} else {
			loaderResults = append(loaderResults, &dataloader.Result[[]*T]{Data: results})
		}
	}
	return loaderResults

}

func fetchExemptWorkingShedules(db *gorm.DB, holidayIds []int) ([]*models.WorkingSchedule, error) {
	raw := `SELECT schedules.id FROM working_schedule_holiday_exemptions exemptions
	INNER JOIN working_schedules schedules ON exemptions.working_schedule_id = schedules.id
	WHERE exemptions.holiday_id IN ?
	`
	var scheduleIds []int
	if err := db.Raw(raw, holidayIds).Scan(&scheduleIds).Error; err != nil {
		return nil, err
	}

	schedules := make([]*models.WorkingSchedule, 0, len(scheduleIds))
	err := db.Where("id IN ?", scheduleIds).Preload("Details").Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func GetCustomersForGroup(ctx context.Context, groupId int) ([]*models.Customer, error) {
	loaders := For(ctx)
	return loaders.customerLoaderForGroup.Load(ctx, groupId)()
}

func GetSuppliersForGroup(ctx context.Context, groupId int) ([]*models.Supplier, error) {
	loaders := For(ctx)
	return loaders.supplierLoaderForGroup.Load(ctx, groupId)()
}

func GetLeaveBalancesForEmployee(ctx context.Context, employeeId int) ([]*models.LeaveBalance, error) {
	loaders := For(ctx)
	return loaders.leaveBalancesLoaderForEmployee.Load(ctx, employeeId)()
}
