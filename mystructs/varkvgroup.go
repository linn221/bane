package mystructs

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type VarKVGroup struct {
	VarKVs []VarKV
}

type VarKV struct {
	Key   VarString
	Value VarString
}

// Value implements the driver.Valuer interface for GORM
// Stores the VarKVGroup as a string in format "key|{key:default} value|{value:default} ..."
func (vkg VarKVGroup) Value() (driver.Value, error) {
	if len(vkg.VarKVs) == 0 {
		return "", nil
	}

	var parts []string
	for _, kv := range vkg.VarKVs {
		keyStr := kv.Key.OriginalString
		valueStr := kv.Value.OriginalString
		parts = append(parts, fmt.Sprintf("%s|%s", keyStr, valueStr))
	}

	return strings.Join(parts, " "), nil
}

// Scan implements the sql.Scanner interface for GORM
// Parses the stored string format back into VarKVGroup
func (vkg *VarKVGroup) Scan(value interface{}) error {
	if value == nil {
		vkg.VarKVs = []VarKV{}
		return nil
	}

	var input string
	switch v := value.(type) {
	case []byte:
		input = string(v)
	case string:
		input = v
	default:
		return fmt.Errorf("cannot scan %T into VarKVGroup", value)
	}

	if input == "" {
		vkg.VarKVs = []VarKV{}
		return nil
	}

	// Split by spaces to get individual key|value pairs
	parts := strings.Fields(input)
	var varKVs []VarKV

	for _, part := range parts {
		// Split by pipe to separate key and value
		pipeIndex := strings.Index(part, "|")
		if pipeIndex == -1 {
			return fmt.Errorf("invalid format: missing pipe in '%s'", part)
		}

		keyStr := part[:pipeIndex]
		valueStr := part[pipeIndex+1:]

		// Parse key as VarString
		keyVarString, err := NewVarString(keyStr)
		if err != nil {
			return fmt.Errorf("failed to parse key VarString '%s': %v", keyStr, err)
		}

		// Parse value as VarString
		valueVarString, err := NewVarString(valueStr)
		if err != nil {
			return fmt.Errorf("failed to parse value VarString '%s': %v", valueStr, err)
		}

		varKVs = append(varKVs, VarKV{
			Key:   *keyVarString,
			Value: *valueVarString,
		})
	}

	vkg.VarKVs = varKVs
	return nil
}
