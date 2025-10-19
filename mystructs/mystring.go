package mystructs

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// MyString represents a dynamic string that can display both VarString and VarKVGroup content
// It supports different display formats and separators
type MyString struct {
	Content   string
	Separator string
	Format    string // "string", "kv", "list"
}

// NewMyStringFromVarString creates a MyString from a VarString
func NewMyStringFromVarString(vs VarString) MyString {
	return MyString{
		Content:   vs.Exec(),
		Separator: "",
		Format:    "string",
	}
}

// NewMyStringFromVarKVGroup creates a MyString from a VarKVGroup
func NewMyStringFromVarKVGroup(vkg VarKVGroup, separator string) MyString {
	if len(vkg.VarKVs) == 0 {
		return MyString{
			Content:   "",
			Separator: separator,
			Format:    "kv",
		}
	}

	var parts []string
	for _, kv := range vkg.VarKVs {
		parts = append(parts, fmt.Sprintf("%s:%s", kv.Key.Exec(), kv.Value.Exec()))
	}

	return MyString{
		Content:   strings.Join(parts, separator),
		Separator: separator,
		Format:    "kv",
	}
}

// NewMyStringFromVarKVGroupAsList creates a MyString from a VarKVGroup as a list
func NewMyStringFromVarKVGroupAsList(vkg VarKVGroup, separator string) MyString {
	if len(vkg.VarKVs) == 0 {
		return MyString{
			Content:   "",
			Separator: separator,
			Format:    "list",
		}
	}

	var parts []string
	for _, kv := range vkg.VarKVs {
		parts = append(parts, kv.Key.Exec())
		parts = append(parts, kv.Value.Exec())
	}

	return MyString{
		Content:   strings.Join(parts, separator),
		Separator: separator,
		Format:    "list",
	}
}

// String returns the string representation
func (ms MyString) String() string {
	return ms.Content
}

// MarshalGQL implements the graphql.Marshaler interface for GraphQL serialization
func (ms MyString) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(ms.Content))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for GraphQL deserialization
func (ms *MyString) UnmarshalGQL(v interface{}) error {
	var input string
	switch val := v.(type) {
	case string:
		input = val
	case []byte:
		input = string(val)
	default:
		return fmt.Errorf("MyString must be a string, got %T", v)
	}

	ms.Content = input
	ms.Separator = ""
	ms.Format = "string"
	return nil
}
