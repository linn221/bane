# VarKVGroup Documentation

## Overview

`VarKVGroup` represents a collection of key-value pairs where both the keys and values are `VarString` types. This allows you to have dynamic key-value pairs where both keys and values can contain placeholders with default values.

## Format

VarKVGroup is stored and serialized as a space-separated string of key-value pairs in the format:
```
key1:value1 key2:value2 key3:value3
```

Both keys and values can be VarStrings with placeholders like `{name=default}`.

## How It Works

1. **Structure**: VarKVGroup contains a slice of `VarKV` structs, each with a `Key` and `Value` (both VarStrings)
2. **Storage**: When saved to database, it's serialized as `"key1:value1 key2:value2"`
3. **Parsing**: When loaded, the string is split by spaces, then each part is split by the first colon
4. **Execution**: The `Exec()` method evaluates all VarStrings in keys and values, returning the final string

## Usage Examples

### Basic Usage

```go
// Create from string
var vkg VarKVGroup
err := vkg.UnmarshalGQL("header1:value1 header2:value2")
if err != nil {
    log.Fatal(err)
}

fmt.Println(vkg.Exec()) // Output: "header1:value1 header2:value2"
```

### With VarString Placeholders

```go
// Keys and values can contain VarString placeholders
var vkg VarKVGroup
vkg.UnmarshalGQL("Token:Bearer{token=abc123} UA:Mozilla{version=1.0}")

// Use defaults
fmt.Println(vkg.Exec())
// Output: "Token:Bearerabc123 UA:Mozilla1.0"

// Inject values
vkg.VarKVs[0].Value.Inject(map[string]string{"token": "xyz789"})
vkg.VarKVs[1].Value.Inject(map[string]string{"version": "2.0"})
fmt.Println(vkg.Exec())
// Output: "Token:Bearerxyz789 UA:Mozilla2.0"
```

### Database Storage

VarKVGroup implements GORM's `Valuer` and `Scanner` interfaces:

```go
type Endpoint struct {
    HttpHeaders mystructs.VarKVGroup `gorm:"not null"`
}

endpoint := Endpoint{
    HttpHeaders: mystructs.VarKVGroup{
        VarKVs: []mystructs.VarKV{
            {
                Key:   mystructs.VarString{OriginalString: "Authorization"},
                Value: mystructs.VarString{OriginalString: "Bearer{token=default}"},
            },
        },
    },
}
db.Create(&endpoint)

// When loaded, the string is parsed back into VarKVGroup
var loaded Endpoint
db.First(&loaded, endpoint.ID)
```

### GraphQL Integration

```go
// In GraphQL schema:
// scalar VarKVGroup

// Serialized as: "key1:value1 key2:value2"
// Both keys and values preserve their original VarString format
```

## Methods

### Value() (driver.Value, error)

Implements GORM's `Valuer` interface. Serializes the VarKVGroup to a string format:
```
"key1:value1 key2:value2"
```

Where keys and values are the `OriginalString` from each VarString.

### Scan(value interface{}) error

Implements GORM's `Scanner` interface. Parses a string in the format `"key1:value1 key2:value2"` and:
1. Splits by spaces to get individual pairs
2. Splits each pair by the first colon to separate key and value
3. Parses both key and value as VarStrings
4. Reconstructs the VarKVGroup

### MarshalGQL(w io.Writer)

Implements GraphQL's `Marshaler` interface. Serializes to the string format with original VarString values.

### UnmarshalGQL(v interface{}) error

Implements GraphQL's `Unmarshaler` interface. Parses a GraphQL input string into VarKVGroup.

### Exec() string

Evaluates all VarStrings in keys and values using their `Exec()` method, then returns the final string:
```
"executedKey1:executedValue1 executedKey2:executedValue2"
```

### IsZero() bool

Returns true if the VarKVGroup has no key-value pairs.

## Use Cases

- **HTTP Headers**: Store headers with dynamic values (e.g., tokens, user agents)
- **Query Parameters**: Parameterize URL query strings
- **Cookies**: Dynamic cookie values with placeholders
- **Configuration**: Key-value configurations where values can be templated

## VarKV Structure

```go
type VarKV struct {
    Key   VarString
    Value VarString
}
```

Each key-value pair in a VarKVGroup is represented as a VarKV struct, allowing both the key and value to be VarStrings with placeholders.

## Notes

- Keys and values are separated by a single colon `:`
- Pairs are separated by spaces
- If a value contains spaces, it must be part of a single pair (no spaces between pairs)
- Both keys and values are parsed as VarStrings, so they can contain placeholders
- The `Exec()` method evaluates all placeholders in both keys and values

