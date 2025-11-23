package services

import (
	"fmt"

	"github.com/linn221/bane/models"
)

func toModelStruct(tableName string) any {
	var tableNameToStruct = map[string]any{
		"endpoints":        models.Endpoint{},
		"tasks":            models.Task{},
		"projects":         models.Project{},
		"notes":            models.Note{},
		"my_sheets":        models.MySheet{},
		"words":            models.Word{},
		"wordlists":        models.WordList{},
		"my_requests":      models.MyRequest{},
		"vulns":            models.Vuln{},
		"vuln_connections": models.VulnConnection{},
		"aliases":          models.Alias{},
	}
	emptyStruct, ok := tableNameToStruct[tableName]
	if !ok {
		panic(fmt.Sprintf("%s has no associated struct!", tableName))
	}
	return emptyStruct
}
