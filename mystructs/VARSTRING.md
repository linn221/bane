# VarString Documentation

## Overview

`VarString` represents a string with variable placeholders that can be injected with values at runtime. It supports a template-like syntax where placeholders are defined in the format `{name=default}`.

## Syntax

Placeholders in VarString use the format: `{variableName=defaultValue}`

- `variableName`: Must start with a letter or underscore, followed by letters, numbers, or underscores
- `defaultValue`: Can contain any characters except `}` (the closing brace)

## How It Works

1. **Parsing**: When you create a VarString with `NewVarString()`, it parses the original string to extract all placeholders
2. **Default Values**: Each placeholder's default value is stored in the `Variables` map
3. **Injection**: You can override default values using the `Inject()` method
4. **Execution**: The `Exec()` method substitutes all placeholders with their current values (injected or default)

## Usage Examples

### Basic Usage

```go
// Create a VarString with placeholders
vs, err := NewVarString("Hello {name=World}!")
if err != nil {
    log.Fatal(err)
}

// Use default value
fmt.Println(vs.Exec()) // Output: "Hello World!"

// Inject a new value
vs.Inject(map[string]string{"name": "Go"})
fmt.Println(vs.Exec()) // Output: "Hello Go!"
```

### Multiple Placeholders

```go
vs, _ := NewVarString("{id=1}-{name=Unknown}")
fmt.Println(vs.Exec()) // Output: "1-Unknown"

vs.Inject(map[string]string{"id": "42", "name": "Alice"})
fmt.Println(vs.Exec()) // Output: "42-Alice"
```

### Empty Default Values

```go
vs, _ := NewVarString("prefix-{var=}-suffix")
fmt.Println(vs.Exec()) // Output: "prefix--suffix"

vs.Inject(map[string]string{"var": "value"})
fmt.Println(vs.Exec()) // Output: "prefix-value-suffix"
```

### Database Storage

VarString implements GORM's `Valuer` and `Scanner` interfaces, so it can be used directly in GORM models:

```go
type Endpoint struct {
    HttpPath mystructs.VarString `gorm:"not null"`
}

// When saving, only the original string is stored
endpoint := Endpoint{
    HttpPath: mystructs.VarString{OriginalString: "/api/{id=1}/data"},
}
db.Create(&endpoint)

// When loading, the VarString is automatically parsed
var loaded Endpoint
db.First(&loaded, endpoint.ID)
fmt.Println(loaded.HttpPath.Exec()) // Uses default: "/api/1/data"
```

### GraphQL Integration

VarString implements GraphQL's `Marshaler` and `Unmarshaler` interfaces:

```go
// In GraphQL schema:
// scalar VarString

// The original string (with placeholders) is serialized to GraphQL
// Clients receive: "{id=1}-{name=Unknown}"
// Clients can inject values and execute on their side if needed
```

## Methods

### NewVarString(s string) (*VarString, error)

Creates a new VarString by parsing the input string. Returns an error if the string contains invalid placeholder syntax.

### Inject(vars map[string]string)

Sets variable values for substitution. Variables not provided will use their default values from the original string.

```go
vs.Inject(map[string]string{
    "id": "123",
    "name": "Test",
})
```

### Exec() string

Returns the final string with all variables substituted. Uses injected values if available, otherwise uses defaults.

### String() string

Implements the `fmt.Stringer` interface. Calls `Exec()` internally.

### Value() (driver.Value, error)

Implements GORM's `Valuer` interface. Returns only the `OriginalString` for database storage.

### Scan(value interface{}) error

Implements GORM's `Scanner` interface. Parses the stored string and recreates the VarString structure.

### MarshalGQL(w io.Writer)

Implements GraphQL's `Marshaler` interface. Serializes the original string (with placeholders) to GraphQL.

### UnmarshalGQL(v interface{}) error

Implements GraphQL's `Unmarshaler` interface. Parses a GraphQL input string into a VarString.

### IsZero() bool

Returns true if the VarString is in its zero state (empty original string, no variables, etc.).

## Use Cases

- **API Endpoints**: Parameterize URL paths with IDs or other dynamic values
- **HTTP Requests**: Template request bodies with variable data
- **Configuration**: Store templates with default values that can be overridden
- **Dynamic Content**: Generate strings based on runtime variables while maintaining defaults

## Notes

- Placeholders are case-sensitive: `{id=1}` and `{ID=1}` are different variables
- The original string with placeholders is always preserved
- Default values are extracted during parsing and stored separately
- Injection is additive - you can call `Inject()` multiple times with different variables

