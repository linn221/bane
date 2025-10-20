package mystructs

import (
	"fmt"
	"io"
	"strings"
)

// MyString represents a dynamic string that can display both VarString and VarKVGroup content
// It supports different display formats and separators
type MyString struct {
	Content   string
	Separator string
	Format    string // "string", "kv", "list"
	Sep       string // separator for VarKVGroup (default: "\n")
	Limit     int    // limit for VarKVGroup items (default: 20)
}

// NewMyStringFromVarString creates a MyString from a VarString
func NewMyStringFromVarString(vs VarString) MyString {
	return MyString{
		Content:   vs.Exec(),
		Separator: "",
		Format:    "string",
		Sep:       "\n",
		Limit:     20,
	}
}

// NewMyStringFromVarKVGroup creates a MyString from a VarKVGroup
func NewMyStringFromVarKVGroup(vkg VarKVGroup, sep *string, limit *int) MyString {
	// Set defaults
	separator := "\n"
	if sep != nil {
		separator = *sep
	}

	lim := 20
	if limit != nil {
		lim = *limit
	}

	if len(vkg.VarKVs) == 0 {
		return MyString{
			Content:   "",
			Separator: "",
			Format:    "kv",
			Sep:       separator,
			Limit:     lim,
		}
	}

	// Apply limit
	varKVs := vkg.VarKVs
	if lim > 0 && len(varKVs) > lim {
		varKVs = varKVs[:lim]
	}

	var parts []string
	for _, kv := range varKVs {
		parts = append(parts, fmt.Sprintf("%s: %s", kv.Key.Exec(), kv.Value.Exec()))
	}

	// Join with separator (which should be actual newline, not escaped)
	content := strings.Join(parts, separator)

	return MyString{
		Content:   content,
		Separator: "",
		Format:    "kv",
		Sep:       separator,
		Limit:     lim,
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

// NewMyStringFromVarKV creates a MyString from a single VarKV
func NewMyStringFromVarKV(kv VarKV) MyString {
	content := fmt.Sprintf("%s: %s", kv.Key.Exec(), kv.Value.Exec())
	// Convert \n to actual newlines
	content = strings.ReplaceAll(content, "\\n", "\n")
	return MyString{
		Content:   content,
		Separator: "",
		Format:    "kv",
		Sep:       "\n",
		Limit:     20,
	}
}

// String returns the string representation
func (ms MyString) String() string {
	return ms.Content
}

// MarshalGQL implements the graphql.Marshaler interface for GraphQL serialization
func (ms MyString) MarshalGQL(w io.Writer) {
	// Use block string syntax to preserve newlines in GraphQL Playground
	fmt.Fprint(w, `"""`)
	fmt.Fprint(w, ms.Content)
	fmt.Fprint(w, `"""`)
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
