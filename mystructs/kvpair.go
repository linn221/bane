package mystructs

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type KVPair struct {
	Key   string
	Value string
}

type KVGroup struct {
	KVPairs []KVPair
}

type NewKVPairGroup struct {
	KVGroup
}

// KVGroupInput represents a GraphQL scalar type for key-value pair groups
// It can be marshaled to/from strings in the format "key1:value1 key2:value2 ..."
type KVGroupInput struct {
	KVGroup
}

// MarshalGQL implements the graphql.Marshaler interface for GraphQL serialization
func (kv KVGroup) MarshalGQL(w io.Writer) {
	if len(kv.KVPairs) == 0 {
		fmt.Fprint(w, `""`)
		return
	}

	var parts []string
	for _, pair := range kv.KVPairs {
		parts = append(parts, fmt.Sprintf("%s:%s", pair.Key, pair.Value))
	}

	result := strings.Join(parts, " ")
	fmt.Fprint(w, strconv.Quote(result))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for GraphQL deserialization
func (kv *KVGroup) UnmarshalGQL(v interface{}) error {
	var input string
	switch val := v.(type) {
	case string:
		input = val
	case []byte:
		input = string(val)
	default:
		return fmt.Errorf("KVGroup must be a string, got %T", v)
	}

	// Parse the input string in format "key1:value1 key2:value2 ..."
	if input == "" {
		kv.KVPairs = []KVPair{}
		return nil
	}

	var pairs []KVPair
	// Parse the input string by finding key:value pairs
	// Each pair is separated by spaces, but values can contain spaces
	// We need to be smart about splitting: look for the pattern "key:value " (space after value)

	// First, let's try a simpler approach: split by spaces and handle each part
	parts := strings.Fields(input) // Split by whitespace

	for _, part := range parts {
		// Split each part by the first colon
		colonIndex := strings.Index(part, ":")
		if colonIndex == -1 {
			return fmt.Errorf("invalid format: missing colon in '%s'", part)
		}

		key := part[:colonIndex]
		value := part[colonIndex+1:]

		pairs = append(pairs, KVPair{
			Key:   key,
			Value: value,
		})
	}

	kv.KVPairs = pairs
	return nil
}

// String returns the string representation of the KVGroupInput
func (kv KVGroup) String() string {
	if len(kv.KVPairs) == 0 {
		return ""
	}

	var parts []string
	for _, pair := range kv.KVPairs {
		parts = append(parts, fmt.Sprintf("%s:%s", pair.Key, pair.Value))
	}

	return strings.Join(parts, " ")
}

// ToKVPairGroup converts KVGroupInput to KVPairGroup
func (kv KVGroup) ToKVGroup() KVGroup {
	return kv
}

// NewKVGroupInput creates a new KVGroupInput from a KVPairGroup
func NewKVGroup(group KVGroup) KVGroup {
	return group
}

// NewKVGroupInputFromString creates a new KVGroupInput from a string
func NewKVGroupFromString(input string) (KVGroup, error) {
	kv := &KVGroup{}
	err := kv.UnmarshalGQL(input)
	return *kv, err
}

// Value implements the driver.Valuer interface for GORM
// Stores the KVPairGroup as a string in format "key:value key:value ..."
func (kg KVGroup) Value() (driver.Value, error) {
	if len(kg.KVPairs) == 0 {
		return "", nil
	}

	var parts []string
	for _, kv := range kg.KVPairs {
		parts = append(parts, fmt.Sprintf("%s:%s", kv.Key, kv.Value))
	}

	return strings.Join(parts, " "), nil
}

// Scan implements the sql.Scanner interface for GORM
// Parses the stored string format back into KVPairGroup
func (kg *KVGroup) Scan(value interface{}) error {
	if value == nil {
		kg.KVPairs = []KVPair{}
		return nil
	}

	var input string
	switch v := value.(type) {
	case []byte:
		input = string(v)
	case string:
		input = v
	default:
		return fmt.Errorf("cannot scan %T into KVPairGroup", value)
	}

	if input == "" {
		kg.KVPairs = []KVPair{}
		return nil
	}

	// Split by spaces to get individual key:value pairs
	parts := strings.Fields(input)
	var pairs []KVPair

	for _, part := range parts {
		// Split by colon to separate key and value
		colonIndex := strings.Index(part, ":")
		if colonIndex == -1 {
			return fmt.Errorf("invalid format: missing colon in '%s'", part)
		}

		key := part[:colonIndex]
		value := part[colonIndex+1:]

		pairs = append(pairs, KVPair{
			Key:   key,
			Value: value,
		})
	}

	kg.KVPairs = pairs
	return nil
}
