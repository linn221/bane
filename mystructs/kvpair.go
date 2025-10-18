package mystructs

import (
	"fmt"
	"strings"
)

type KVPair struct {
	Key   string
	Value string
}

type KVPairGroup struct {
	KVPairs []KVPair
}

type NewKVPairGroup struct {
	KVPairGroup
}

// KVGroupInput represents a GraphQL scalar type for key-value pair groups
// It can be marshaled to/from strings in the format "key1:value1 key2:value2 ..."
type KVGroupInput struct {
	KVPairGroup
}

// MarshalGQL implements the graphql.Marshaler interface for GraphQL serialization
func (kv KVGroupInput) MarshalGQL() ([]byte, error) {
	if len(kv.KVPairs) == 0 {
		return []byte(`""`), nil
	}

	var parts []string
	for _, pair := range kv.KVPairs {
		parts = append(parts, fmt.Sprintf("%s:%s", pair.Key, pair.Value))
	}

	result := strings.Join(parts, " ")
	return []byte(fmt.Sprintf(`"%s"`, result)), nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for GraphQL deserialization
func (kv *KVGroupInput) UnmarshalGQL(v interface{}) error {
	var input string
	switch val := v.(type) {
	case string:
		input = val
	case []byte:
		input = string(val)
	default:
		return fmt.Errorf("KVGroupInput must be a string, got %T", v)
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
func (kv KVGroupInput) String() string {
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
func (kv KVGroupInput) ToKVPairGroup() KVPairGroup {
	return kv.KVPairGroup
}

// NewKVGroupInput creates a new KVGroupInput from a KVPairGroup
func NewKVGroupInput(group KVPairGroup) KVGroupInput {
	return KVGroupInput{KVPairGroup: group}
}

// NewKVGroupInputFromString creates a new KVGroupInput from a string
func NewKVGroupInputFromString(input string) (KVGroupInput, error) {
	kv := &KVGroupInput{}
	err := kv.UnmarshalGQL(input)
	return *kv, err
}
