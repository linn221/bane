# MyString Documentation

## Overview

`MyString` is a flexible string type that can display content from both `VarString` and `VarKVGroup` types in various formats. It provides a unified interface for converting these complex types into readable string representations.

## How It Works

MyString acts as a display adapter that:
1. Takes a `VarString` or `VarKVGroup` as input
2. Applies formatting options (separators, limits, etc.)
3. Produces a formatted string output
4. Supports GraphQL serialization with proper formatting

## Usage Examples

### From VarString

```go
vs, _ := NewVarString("Hello {name=World}!")
ms := NewMyStringFromVarString(*vs)

fmt.Println(ms.String()) // Output: "Hello World!"
// After injection:
vs.Inject(map[string]string{"name": "Go"})
ms = NewMyStringFromVarString(*vs)
fmt.Println(ms.String()) // Output: "Hello Go!"
```

### From VarKVGroup (Key-Value Format)

```go
var vkg VarKVGroup
vkg.UnmarshalGQL("name:John age:30 city:NYC")

// Default format: key-value pairs separated by newlines
ms := NewMyStringFromVarKVGroup(vkg, nil, nil)
fmt.Println(ms.String())
// Output:
// name: John
// age: 30
// city: NYC
```

### From VarKVGroup with Custom Separator

```go
var vkg VarKVGroup
vkg.UnmarshalGQL("name:John age:30 city:NYC")

separator := "\n"
ms := NewMyStringFromVarKVGroup(vkg, &separator, nil)
fmt.Println(ms.String())
// Output (each pair on new line):
// name: John
// age: 30
// city: NYC
```

### From VarKVGroup with Limit

```go
var vkg VarKVGroup
// Assume vkg has 50 items

limit := 10
ms := NewMyStringFromVarKVGroup(vkg, nil, &limit)
fmt.Println(ms.String())
// Output: Only first 10 items displayed
```

### From VarKVGroup as List

```go
var vkg VarKVGroup
vkg.UnmarshalGQL("name:John age:30 city:NYC")

ms := NewMyStringFromVarKVGroupAsList(vkg, ", ")
fmt.Println(ms.String())
// Output: "John, 30, NYC, John, 30, NYC"
// (Note: keys and values are interleaved)
```

### From Single VarKV

```go
kv := VarKV{
    Key:   VarString{OriginalString: "token"},
    Value: VarString{OriginalString: "Bearer{value=abc}"},
}

ms := NewMyStringFromVarKV(kv)
fmt.Println(ms.String())
// Output: "token: Bearerabc"
```

## Methods

### NewMyStringFromVarString(vs VarString) MyString

Creates a MyString from a VarString. The format is set to "string".

```go
vs, _ := NewVarString("Hello {name=World}")
ms := NewMyStringFromVarString(*vs)
```

### NewMyStringFromVarKVGroup(vkg VarKVGroup, sep *string, limit *int) MyString

Creates a MyString from a VarKVGroup in key-value format.

- `sep`: Separator between key-value pairs (default: `"\n"`)
- `limit`: Maximum number of items to display (default: 20, 0 = no limit)

```go
separator := " | "
limit := 5
ms := NewMyStringFromVarKVGroup(vkg, &separator, &limit)
```

### NewMyStringFromVarKVGroupAsList(vkg VarKVGroup, separator string) MyString

Creates a MyString from a VarKVGroup as a flat list. Keys and values are interleaved and separated by the provided separator.

```go
ms := NewMyStringFromVarKVGroupAsList(vkg, ", ")
```

### NewMyStringFromVarKV(kv VarKV) MyString

Creates a MyString from a single VarKV pair.

```go
ms := NewMyStringFromVarKV(kv)
```

### String() string

Returns the string representation. Implements `fmt.Stringer`.

### MarshalGQL(w io.Writer)

Implements GraphQL's `Marshaler` interface. Uses GraphQL block string syntax (`"""`) to preserve newlines and formatting in GraphQL Playground.

```go
// Serializes as:
// """
// name: John
// age: 30
// """
```

### UnmarshalGQL(v interface{}) error

Implements GraphQL's `Unmarshaler` interface. Parses a GraphQL string input into MyString.

## MyString Structure

```go
type MyString struct {
    Content   string // The formatted string content
    Separator string // Separator used (for list format)
    Format    string // "string", "kv", or "list"
    Sep       string // Separator for VarKVGroup (default: "\n")
    Limit     int    // Limit for VarKVGroup items (default: 20)
}
```

## Format Types

### "string"
Used when created from a VarString. Displays the executed VarString directly.

### "kv"
Used when created from a VarKVGroup in key-value format. Displays as:
```
key1: value1
key2: value2
```

### "list"
Used when created from a VarKVGroup as a list. Displays keys and values interleaved:
```
key1, value1, key2, value2
```

## Use Cases

- **Display Formatting**: Convert VarString/VarKVGroup to human-readable strings
- **GraphQL Responses**: Format complex types for GraphQL API responses
- **Logging**: Display variable strings in logs with proper formatting
- **UI Display**: Format data for display in user interfaces

## GraphQL Integration

MyString is particularly useful in GraphQL resolvers:

```go
// In resolver
func (r *endpointResolver) Headers(ctx context.Context, obj *models.Endpoint) (string, error) {
    ms := NewMyStringFromVarKVGroup(obj.HttpHeaders, nil, nil)
    return ms.String(), nil
}
```

The block string format (`"""`) ensures that newlines are preserved when viewing in GraphQL Playground.

## Notes

- MyString is primarily a display/formatting type
- It doesn't modify the original VarString or VarKVGroup
- The `Limit` parameter is useful for displaying large VarKVGroups
- Block string syntax in GraphQL preserves formatting better than regular strings
- The `Separator` field is used differently depending on the format type

