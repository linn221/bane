package mystructs

import (
	"database/sql/driver"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// VarString represents a string with variable placeholders that can be injected with values
// It supports GORM serialization and implements the Stringer interface
type VarString struct {
	OriginalString string            `json:"original_string"` // The original string with placeholders
	Variables      map[string]string `json:"variables"`       // Injected variable values
	ParsedTemplate string            `json:"parsed_template"` // The template after parsing placeholders
	Placeholders   []string          `json:"placeholders"`    // List of placeholder names found in the string
}

// NewVarString parses a string with variable placeholders in the format {name:default}
// Returns a VarString with parsed placeholders and default values
func NewVarString(s string) (*VarString, error) {
	vs := &VarString{
		OriginalString: s,
		Variables:      make(map[string]string),
		Placeholders:   make([]string, 0),
	}

	// Parse the string to extract placeholders in format {name:default}
	// This regex matches {name:default} where name is alphanumeric and default can contain any character except }
	re := regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*):([^}]*)\}`)
	matches := re.FindAllStringSubmatch(s, -1)

	// Build the parsed template by replacing placeholders with variable references
	parsedTemplate := s
	for _, match := range matches {
		placeholder := match[0]  // Full match like "{id:1}"
		varName := match[1]      // Variable name like "id"
		defaultValue := match[2] // Default value like "1"

		// Store the default value
		vs.Variables[varName] = defaultValue
		vs.Placeholders = append(vs.Placeholders, varName)

		// Replace the placeholder with a simple variable reference for later substitution
		parsedTemplate = strings.ReplaceAll(parsedTemplate, placeholder, "{"+varName+"}")
	}

	vs.ParsedTemplate = parsedTemplate
	return vs, nil
}

// Inject sets variable values for substitution
// Variables not provided will use their default values from the original string
func (vs *VarString) Inject(vars map[string]string) {
	for key, value := range vars {
		vs.Variables[key] = value
	}
}

// Exec returns the final string with all variables substituted
// Uses the most efficient method: direct string replacement
func (vs *VarString) Exec() string {
	result := vs.ParsedTemplate

	// Replace each variable with its value
	for varName, value := range vs.Variables {
		placeholder := "{" + varName + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// String implements the Stringer interface
func (vs *VarString) String() string {
	return vs.Exec()
}

// Value implements the driver.Valuer interface for GORM
// Stores only the original string to database
func (vs VarString) Value() (driver.Value, error) {
	return vs.OriginalString, nil
}

// Scan implements the sql.Scanner interface for GORM
// Reads the original string from database and parses it
func (vs *VarString) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var originalString string
	switch v := value.(type) {
	case []byte:
		originalString = string(v)
	case string:
		originalString = v
	default:
		return fmt.Errorf("cannot scan %T into VarString", value)
	}

	// Parse the original string to recreate the VarString
	parsedVs, err := NewVarString(originalString)
	if err != nil {
		return err
	}

	// Copy the parsed values to this instance
	*vs = *parsedVs
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for GraphQL serialization
func (vs VarString) MarshalGQL(w io.Writer) {
	// Return the original string (as stored in database)
	fmt.Fprint(w, strconv.Quote(vs.OriginalString))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for GraphQL deserialization
func (vs *VarString) UnmarshalGQL(v interface{}) error {
	var input string
	switch val := v.(type) {
	case string:
		input = val
	case []byte:
		input = string(val)
	default:
		return fmt.Errorf("VarString must be a string, got %T", v)
	}

	// Parse the input string as a VarString
	parsedVs, err := NewVarString(input)
	if err != nil {
		return err
	}

	// Copy the parsed values to this instance
	*vs = *parsedVs
	return nil
}
