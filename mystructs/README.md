VarString => Dynamic string with variable, can be injected for fuzzing purpose
KVGroup => Static Key Value pair, having slice of KeyPair
VarKVGroup => Dynamic Key Value pair, of VarString type. able to store http headers, injectable
MyString => Adapter type, any of the above type can be converted to MyString, and MyString can do flexible display of text

# mystructs Package

This package provides specialized types for handling variable strings and key-value groups with support for GORM database serialization and GraphQL integration.

## Quick Start

```go
import "github.com/linn221/bane/mystructs"

// Create a VarString with placeholders
vs, _ := mystructs.NewVarString("Hello {name=World}!")
fmt.Println(vs.Exec()) // "Hello World!"

// Inject values
vs.Inject(map[string]string{"name": "Go"})
fmt.Println(vs.Exec()) // "Hello Go!"

// Create a VarKVGroup
var vkg mystructs.VarKVGroup
vkg.UnmarshalGQL("Token:Bearer{token=abc} UA:Mozilla{version=1.0}")
fmt.Println(vkg.Exec()) // "Token:Bearerabc UA:Mozilla1.0"
```

## Types Overview

### VarString
A string with variable placeholders in the format `{name=default}`. Supports injection of runtime values while maintaining defaults.

**See**: [VARSTRING.md](./VARSTRING.md)

### VarKVGroup
A collection of key-value pairs where both keys and values are VarStrings. Useful for HTTP headers, query parameters, etc.

**See**: [VARKVGROUP.md](./VARKVGROUP.md)

### KVGroup
A simple collection of key-value pairs (plain strings). For static key-value data without templating needs.

**See**: [KVGROUP.md](./KVGROUP.md)

### MyString
A flexible display type that can format VarString or VarKVGroup content in various ways for display or GraphQL responses.

**See**: [MYSTRING.md](./MYSTRING.md)

## Common Use Cases

### HTTP Endpoint Configuration

```go
type Endpoint struct {
    HttpPath    mystructs.VarString  `gorm:"not null"`
    HttpHeaders mystructs.VarKVGroup `gorm:"not null"`
    HttpQueries mystructs.VarKVGroup `gorm:"not null"`
}

// Create endpoint with templated values
endpoint := Endpoint{
    HttpPath: mystructs.VarString{OriginalString: "/api/{id=1}/data"},
    HttpHeaders: mystructs.VarKVGroup{
        VarKVs: []mystructs.VarKV{
            {
                Key:   mystructs.VarString{OriginalString: "Authorization"},
                Value: mystructs.VarString{OriginalString: "Bearer{token=default}"},
            },
        },
    },
}

// Use default values
path := endpoint.HttpPath.Exec() // "/api/1/data"

// Inject runtime values
endpoint.HttpPath.Inject(map[string]string{"id": "42"})
path = endpoint.HttpPath.Exec() // "/api/42/data"
```

### GraphQL Integration

All types implement GraphQL's `Marshaler` and `Unmarshaler` interfaces, making them suitable for use as GraphQL scalar types:

```graphql
scalar VarString
scalar VarKVGroup
scalar KVGroup
scalar MyString

type Endpoint {
    httpPath: VarString!
    httpHeaders: VarKVGroup!
}
```

### Database Storage

All types implement GORM's `Valuer` and `Scanner` interfaces, so they can be used directly in GORM models:

```go
type Model struct {
    Template mystructs.VarString `gorm:"type:text"`
    Metadata mystructs.KVGroup   `gorm:"type:text"`
}
```

The original string format is stored in the database, and the types are automatically parsed when loaded.

## Features

- ✅ **Placeholder Support**: Use `{name=default}` syntax for variable strings
- ✅ **Default Values**: Maintain default values that can be overridden
- ✅ **GORM Integration**: Direct use in database models
- ✅ **GraphQL Support**: Native GraphQL scalar type support
- ✅ **Type Safety**: Strongly typed Go structs
- ✅ **Flexible Display**: MyString provides various formatting options

## Documentation

- [VarString Documentation](./VARSTRING.md) - Variable strings with placeholders
- [VarKVGroup Documentation](./VARKVGROUP.md) - Key-value groups with VarString support
- [KVGroup Documentation](./KVGROUP.md) - Simple key-value pairs
- [MyString Documentation](./MYSTRING.md) - Display formatting for complex types

## Examples

See [example.go](./example.go) for complete usage examples.

## Testing

Run tests with:
```bash
go test ./mystructs/...
```

