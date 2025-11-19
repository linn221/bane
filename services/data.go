package services

import (
	"fmt"

	"github.com/linn221/bane/models"
)

func toModelStruct(tableName string) any {
	var tableNameToStruct = map[string]any{
		"endpoints": models.Endpoint{},
		// ...
	}
	emptyStruct, ok := tableNameToStruct[tableName]
	if !ok {
		panic(fmt.Sprintf("%s has no associated struct!", tableName))
	}
	return emptyStruct
}
