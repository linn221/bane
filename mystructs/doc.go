// Package mystructs provides specialized types for handling variable strings and key-value groups
// with support for GORM database serialization and GraphQL integration.
//
// Overview
//
// This package contains types designed to work with dynamic, template-like strings that can
// contain placeholders with default values. These types are particularly useful for HTTP
// request configurations, API endpoints, and other scenarios where you need to
// parameterize strings while maintaining a default configuration.
//
// Main Types
//
//   - VarString: A string with variable placeholders in the format {name=default}
//   - VarKVGroup: A collection of key-value pairs where both keys and values are VarStrings
//   - KVGroup: A collection of simple key-value pairs (strings only)
//   - MyString: A flexible string type that can display VarString or VarKVGroup content
//
// All types support:
//   - GORM database serialization (Value/Scan methods)
//   - GraphQL serialization (MarshalGQL/UnmarshalGQL methods)
//   - String representation via String() method
package mystructs

