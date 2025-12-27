package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/linn221/bane/mystructs"
)

// ParsedHttpUrl contains all components extracted from an httpUrl VarString
// This struct is designed to make debugging easier by showing all parsed components
type ParsedHttpUrl struct {
	Https       bool                 // true for https, false for http
	HttpDomain  string
	HttpPath    mystructs.VarString
	HttpQueries mystructs.VarKVGroup
	OriginalUrl string // The original URL string for debugging
}

// ParseHttpUrl parses a VarString containing a URL and extracts all components.
// The URL can contain VarString placeholders like {domain=example.com} or {port=8080}.
// Returns a ParsedHttpUrl struct with all extracted components for easy debugging.
//
// URL format: [scheme]://[domain][:port][/path][?query]
// Examples:
//   - "https://example.com/api/users"
//   - "http://{domain=localhost}:{port=8080}/api/{id=1}/data?key=value&token={token=abc}"
//   - "https://api.example.com:443/v1/endpoint?param1=val1&param2={val2=default}"
//
// The function:
// 1. Executes the VarString to get the actual URL (with defaults)
// 2. Parses the URL using Go's net/url package
// 3. Extracts schema, domain, port, path, and query parameters
// 4. Converts query parameters to VarKVGroup
// 5. Preserves the original VarString structure for path and queries
func ParseHttpUrl(httpUrl mystructs.VarString) (*ParsedHttpUrl, error) {
	// Execute the VarString to get the URL with default values
	executedUrl := httpUrl.Exec()
	
	// Parse the URL
	parsedUrl, err := url.Parse(executedUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL '%s': %w", executedUrl, err)
	}

	// Extract schema - convert to bool (true for https, false for http)
	scheme := strings.ToLower(parsedUrl.Scheme)
	var https bool
	switch scheme {
	case "https":
		https = true
	case "http":
		https = false
	case "":
		// Default to https if no scheme provided
		https = true
	default:
		return nil, fmt.Errorf("unsupported URL scheme '%s' (expected http or https)", scheme)
	}

	// Extract domain
	httpDomain := parsedUrl.Hostname()
	if httpDomain == "" {
		return nil, fmt.Errorf("missing domain in URL '%s'", executedUrl)
	}

	// Extract path - preserve as VarString with original placeholders
	httpPath := parsedUrl.Path
	if httpPath == "" {
		httpPath = "/"
	}
	
	// Create VarString for path, preserving the original structure
	// We need to reconstruct the path from the original VarString
	pathVarString, err := extractPathFromOriginalUrl(httpUrl, httpPath)
	if err != nil {
		// Fallback: create a simple VarString from the parsed path
		pathVarString = mystructs.VarString{OriginalString: httpPath}
		// Try to parse it
		parsed, err := mystructs.NewVarString(httpPath)
		if err == nil {
			pathVarString = *parsed
		}
	}

	// Extract query parameters - convert to VarKVGroup
	httpQueries := mystructs.VarKVGroup{VarKVs: []mystructs.VarKV{}}
	queryValues := parsedUrl.Query()
	
	// Build VarKVGroup from query parameters
	// We need to preserve the original VarString structure from the original URL
	for key, values := range queryValues {
		// Use the first value if multiple values exist
		value := ""
		if len(values) > 0 {
			value = values[0]
		}
		
		// Try to extract the original VarString structure for this query param
		keyVarString, valueVarString := extractQueryParamFromOriginalUrl(httpUrl, key, value)
		
		httpQueries.VarKVs = append(httpQueries.VarKVs, mystructs.VarKV{
			Key:   keyVarString,
			Value: valueVarString,
		})
	}

	return &ParsedHttpUrl{
		Https:       https,
		HttpDomain:  httpDomain,
		HttpPath:    pathVarString,
		HttpQueries: httpQueries,
		OriginalUrl: executedUrl,
	}, nil
}

// extractPathFromOriginalUrl extracts the path portion from the original VarString
// while preserving VarString placeholders
func extractPathFromOriginalUrl(httpUrl mystructs.VarString, parsedPath string) (mystructs.VarString, error) {
	original := httpUrl.OriginalString
	
	// Find the scheme end (://)
	schemeEnd := strings.Index(original, "://")
	if schemeEnd == -1 {
		// No scheme, try to find path starting from beginning
		pathStart := strings.Index(original, "/")
		if pathStart == -1 {
			// No path found, use root
			pathVarString, err := mystructs.NewVarString("/")
			if err != nil {
				return mystructs.VarString{}, err
			}
			return *pathVarString, nil
		}
		// Extract from pathStart to query or end
		queryStart := strings.Index(original[pathStart:], "?")
		pathEnd := len(original)
		if queryStart != -1 {
			pathEnd = pathStart + queryStart
		}
		pathStr := original[pathStart:pathEnd]
		pathVarString, err := mystructs.NewVarString(pathStr)
		if err != nil {
			return mystructs.VarString{}, err
		}
		return *pathVarString, nil
	}
	
	// Find where the domain/port ends (first / or ? after scheme)
	afterScheme := original[schemeEnd+3:]
	pathStart := strings.Index(afterScheme, "/")
	queryStart := strings.Index(afterScheme, "?")
	
	// Determine where path starts
	var actualPathStart int
	if pathStart != -1 {
		actualPathStart = schemeEnd + 3 + pathStart
	} else if queryStart != -1 {
		// No / found, but ? exists - path is empty, use root
		pathVarString, err := mystructs.NewVarString("/")
		if err != nil {
			return mystructs.VarString{}, err
		}
		return *pathVarString, nil
	} else {
		// No / and no ? - path is empty, use root
		pathVarString, err := mystructs.NewVarString("/")
		if err != nil {
			return mystructs.VarString{}, err
		}
		return *pathVarString, nil
	}
	
	// Extract path portion (up to ? if query exists)
	pathEnd := len(original)
	queryIdx := strings.Index(original[actualPathStart:], "?")
	if queryIdx != -1 {
		pathEnd = actualPathStart + queryIdx
	}
	
	pathStr := original[actualPathStart:pathEnd]
	if pathStr == "" {
		pathStr = "/"
	}
	
	pathVarString, err := mystructs.NewVarString(pathStr)
	if err != nil {
		return mystructs.VarString{}, err
	}
	return *pathVarString, nil
}

// extractQueryParamFromOriginalUrl extracts a query parameter from the original VarString
// while preserving VarString placeholders
func extractQueryParamFromOriginalUrl(httpUrl mystructs.VarString, key, value string) (mystructs.VarString, mystructs.VarString) {
	original := httpUrl.OriginalString
	
	// Find the query string portion
	queryStart := strings.Index(original, "?")
	if queryStart == -1 {
		// No query string, create simple VarStrings
		keyVar, _ := mystructs.NewVarString(key)
		valueVar, _ := mystructs.NewVarString(value)
		return *keyVar, *valueVar
	}
	
	queryStr := original[queryStart+1:]
	
	// Try to find this key=value pair in the original query string
	// Look for patterns like "key=value" or "key={value=default}"
	keyPattern := key + "="
	keyIndex := strings.Index(queryStr, keyPattern)
	if keyIndex != -1 {
		// Found the key, extract the value portion
		valueStart := keyIndex + len(keyPattern)
		// Find where this value ends (next & or end of string)
		valueEnd := strings.Index(queryStr[valueStart:], "&")
		if valueEnd == -1 {
			valueEnd = len(queryStr)
		} else {
			valueEnd = valueStart + valueEnd
		}
		
		valueStr := queryStr[valueStart:valueEnd]
		
		// Create VarStrings preserving structure
		keyVar, _ := mystructs.NewVarString(key)
		valueVar, err := mystructs.NewVarString(valueStr)
		if err != nil {
			// Fallback to simple string
			valueVar, _ = mystructs.NewVarString(value)
		}
		
		return *keyVar, *valueVar
	}
	
	// Not found in original, create simple VarStrings
	keyVar, _ := mystructs.NewVarString(key)
	valueVar, _ := mystructs.NewVarString(value)
	return *keyVar, *valueVar
}

