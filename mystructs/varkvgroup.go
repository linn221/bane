package mystructs

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
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

// MarshalGQL implements the graphql.Marshaler interface for GraphQL serialization
func (vkg VarKVGroup) MarshalGQL(w io.Writer) {
	if len(vkg.VarKVs) == 0 {
		fmt.Fprint(w, `""`)
		return
	}

	// Return the original string format as stored in database
	var parts []string
	for _, kv := range vkg.VarKVs {
		parts = append(parts, fmt.Sprintf("%s|%s", kv.Key.OriginalString, kv.Value.OriginalString))
	}

	result := strings.Join(parts, " ")
	fmt.Fprint(w, strconv.Quote(result))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for GraphQL deserialization
func (vkg *VarKVGroup) UnmarshalGQL(v interface{}) error {
	var input string
	switch val := v.(type) {
	case string:
		input = val
	case []byte:
		input = string(val)
	default:
		return fmt.Errorf("VarKVGroup must be a string, got %T", v)
	}

	if input == "" {
		vkg.VarKVs = []VarKV{}
		return nil
	}

	// Parse the input string in format "key|value key|value ..."
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

// Exec returns the executed string with all variables substituted
func (vkg VarKVGroup) Exec() string {
	if len(vkg.VarKVs) == 0 {
		return ""
	}

	var parts []string
	for _, kv := range vkg.VarKVs {
		parts = append(parts, fmt.Sprintf("%s:%s", kv.Key.Exec(), kv.Value.Exec()))
	}

	return strings.Join(parts, " ")
}
