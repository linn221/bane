# KVGroup Documentation

## Overview

`KVGroup` represents a simple collection of key-value pairs where both keys and values are plain strings (not VarStrings). It's useful for storing simple key-value data that doesn't need placeholder functionality.

## Format

KVGroup is stored and serialized as a space-separated string of key-value pairs:
```
key1:value1 key2:value2 key3:value3
```

## How It Works

1. **Structure**: KVGroup contains a slice of `KVPair` structs, each with a `Key` and `Value` (both strings)
2. **Storage**: Serialized as a simple string format when saved to database
3. **Parsing**: String is split by spaces, then each part is split by the first colon
4. **Serialization**: Can be converted to/from string format for GraphQL and database storage

## Usage Examples

### Basic Usage

```go
// Create from string
kv, err := NewKVGroupFromString("name:John age:30 city:NewYork")
if err != nil {
    log.Fatal(err)
}

fmt.Println(kv.String()) // Output: "name:John age:30 city:NewYork"

// Access individual pairs
for _, pair := range kv.KVPairs {
    fmt.Printf("%s: %s\n", pair.Key, pair.Value)
}
```

### Creating Programmatically

```go
kv := KVGroup{
    KVPairs: []KVPair{
        {Key: "name", Value: "John"},
        {Key: "age", Value: "30"},
        {Key: "city", Value: "NewYork"},
    },
}

fmt.Println(kv.String()) // Output: "name:John age:30 city:NewYork"
```

### Database Storage

KVGroup implements GORM's `Valuer` and `Scanner` interfaces:

```go
type Model struct {
    Metadata KVGroup `gorm:"type:text"`
}

model := Model{
    Metadata: KVGroup{
        KVPairs: []KVPair{
            {Key: "version", Value: "1.0"},
            {Key: "author", Value: "Admin"},
        },
    },
}
db.Create(&model)

// When loaded, automatically parsed from string
var loaded Model
db.First(&loaded, model.ID)
```

### GraphQL Integration

KVGroup implements GraphQL's `Marshaler` and `Unmarshaler` interfaces:

```go
// In GraphQL schema:
// scalar KVGroup

// Serialized as: "key1:value1 key2:value2"
// Can be used as input or output type
```

## Methods

### NewKVGroupFromString(input string) (KVGroup, error)

Creates a new KVGroup by parsing the input string. Returns an error if the format is invalid.

### String() string

Returns the string representation of the KVGroup in the format `"key1:value1 key2:value2"`.

### MarshalGQL(w io.Writer)

Implements GraphQL's `Marshaler` interface. Serializes the KVGroup to a quoted string.

### UnmarshalGQL(v interface{}) error

Implements GraphQL's `Unmarshaler` interface. Parses a GraphQL input string into KVGroup.

### Value() (driver.Value, error)

Implements GORM's `Valuer` interface. Returns the string representation for database storage.

### Scan(value interface{}) error

Implements GORM's `Scanner` interface. Parses the stored string back into KVGroup structure.

### ToKVGroup() KVGroup

Returns the KVGroup itself (identity function). Useful for type conversions.

## KVPair Structure

```go
type KVPair struct {
    Key   string
    Value string
}
```

Each key-value pair in a KVGroup is represented as a KVPair struct with simple string fields.

## Related Types

### KVGroupInput

```go
type KVGroupInput struct {
    KVGroup
}
```

A GraphQL input type wrapper around KVGroup. Used when KVGroup needs to be used as a GraphQL input type.

### NewKVPairGroup

```go
type NewKVPairGroup struct {
    KVGroup
}
```

An alias for KVGroup, kept for backward compatibility.

## Use Cases

- **Metadata**: Store simple key-value metadata
- **Tags**: Key-value tags or labels
- **Simple Configuration**: Configuration data that doesn't need templating
- **Attributes**: Object attributes stored as key-value pairs

## Comparison with VarKVGroup

| Feature | KVGroup | VarKVGroup |
|---------|---------|------------|
| Key Type | string | VarString |
| Value Type | string | VarString |
| Placeholders | No | Yes |
| Use Case | Simple key-value data | Dynamic/templated key-value data |

Use `KVGroup` when you have simple, static key-value pairs. Use `VarKVGroup` when you need dynamic values with placeholders and defaults.

## Notes

- Keys and values are separated by a single colon `:`
- Pairs are separated by spaces
- Values cannot contain spaces (they would be interpreted as separate pairs)
- All keys and values are treated as plain strings (no placeholder parsing)
- Empty KVGroup serializes to an empty string

