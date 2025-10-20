package services

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	parsecurl "parse-curl"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/mystructs"
)

// ParseCurlCommand parses a curl command string and returns a CurlImportResult
func ParseCurlCommand(curlCommand string) (*models.CurlImportResult, error) {
	// Remove extra whitespace and normalize
	curlCommand = strings.TrimSpace(curlCommand)

	// Basic validation
	if !strings.HasPrefix(curlCommand, "curl") {
		return nil, fmt.Errorf("command must start with 'curl'")
	}

	// Parse the curl command using the library
	req, ok := parsecurl.Parse(curlCommand)
	if !ok {
		return nil, fmt.Errorf("failed to parse curl command")
	}

	// Initialize the endpoint with defaults
	endpoint := &models.CurlImportResult{
		HttpSchema:          models.HttpSchemaHttps,
		HttpMethod:          models.HttpMethodGet,
		HttpPort:            0, // Default port
		HttpTimeout:         30,
		HttpFollowRedirects: true,
		HttpPath:            mystructs.VarString{OriginalString: "/"},
		HttpQueries:         mystructs.VarKVGroup{VarKVs: []mystructs.VarKV{}},
		HttpHeaders:         mystructs.VarKVGroup{VarKVs: []mystructs.VarKV{}},
		HttpCookies:         mystructs.VarKVGroup{VarKVs: []mystructs.VarKV{}},
		HttpBody:            mystructs.VarString{OriginalString: ""},
	}

	// Extract URL information
	parsedURL, err := url.Parse(req.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Extract schema
	if parsedURL.Scheme == "http" {
		endpoint.HttpSchema = models.HttpSchemaHttp
	}

	// Extract domain
	endpoint.HttpDomain = parsedURL.Hostname()

	// Extract port
	if parsedURL.Port() != "" {
		if port, err := strconv.Atoi(parsedURL.Port()); err == nil {
			endpoint.HttpPort = port
		}
	}

	// Extract path
	if parsedURL.Path == "" {
		endpoint.HttpPath = mystructs.VarString{OriginalString: "/"}
	} else {
		endpoint.HttpPath = mystructs.VarString{OriginalString: parsedURL.Path}
	}

	// Extract query parameters
	if parsedURL.RawQuery != "" {
		queryParams, err := url.ParseQuery(parsedURL.RawQuery)
		if err == nil {
			var varKVs []mystructs.VarKV
			for key, values := range queryParams {
				// Join multiple values with comma
				value := strings.Join(values, ",")
				varKVs = append(varKVs, mystructs.VarKV{
					Key:   mystructs.VarString{OriginalString: key},
					Value: mystructs.VarString{OriginalString: value},
				})
			}
			endpoint.HttpQueries = mystructs.VarKVGroup{VarKVs: varKVs}
		}
	}

	// Extract method
	endpoint.HttpMethod = models.HttpMethod(req.Method)

	// Extract headers
	if len(req.Header) > 0 {
		var varKVs []mystructs.VarKV
		for key, value := range req.Header {
			varKVs = append(varKVs, mystructs.VarKV{
				Key:   mystructs.VarString{OriginalString: key},
				Value: mystructs.VarString{OriginalString: value},
			})
		}
		endpoint.HttpHeaders = mystructs.VarKVGroup{VarKVs: varKVs}
	}

	// Extract body if present
	if req.Body != "" {
		endpoint.HttpBody = mystructs.VarString{OriginalString: req.Body}
	}

	// Set default name based on URL
	endpoint.Name = fmt.Sprintf("%s %s", endpoint.HttpMethod, parsedURL.Hostname())
	if parsedURL.Path != "/" {
		endpoint.Name += parsedURL.Path
	}

	return endpoint, nil
}
