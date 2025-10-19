package services

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/mystructs"
)

// CurlService handles generation of curl commands from Endpoint data
type CurlService struct{}

var CurlServiceInstance = &CurlService{}

// GenerateCurlCommand creates a curl command from an Endpoint with optional variables
func (cs *CurlService) GenerateCurlCommand(endpoint *models.Endpoint, variables string) (string, error) {
	// Parse variables if provided
	varKVs := make(map[string]string)
	if variables != "" {
		// Parse variables in format "key:value key:value"
		pairs := strings.Fields(variables)
		for _, pair := range pairs {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) == 2 {
				varKVs[parts[0]] = parts[1]
			}
		}
	}

	// Build the curl command
	var curlParts []string
	curlParts = append(curlParts, "curl")

	// Add HTTP method
	curlParts = append(curlParts, "-X", string(endpoint.HttpMethod))

	// Build URL
	url := cs.buildURL(endpoint, varKVs)
	curlParts = append(curlParts, url)

	// Add headers
	headers := cs.buildHeaders(endpoint, varKVs)
	for _, header := range headers {
		curlParts = append(curlParts, "-H", header)
	}

	// Add cookies
	cookies := cs.buildCookies(endpoint, varKVs)
	if cookies != "" {
		curlParts = append(curlParts, "-b", cookies)
	}

	// Add body if present
	body := cs.buildBody(endpoint, varKVs)
	if body != "" {
		curlParts = append(curlParts, "--data-raw", body)
	}

	return strings.Join(curlParts, " "), nil
}

// buildURL constructs the full URL with path and query parameters
func (cs *CurlService) buildURL(endpoint *models.Endpoint, varKVs map[string]string) string {
	// Build base URL
	baseURL := fmt.Sprintf("%s://%s", string(endpoint.HttpSchema), endpoint.HttpDomain)

	// Add path with variable substitution
	path := cs.substituteVariables(endpoint.HttpPath, varKVs)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Add query parameters
	queryParams := cs.buildQueryParams(endpoint, varKVs)
	if queryParams != "" {
		path += "?" + queryParams
	}

	return baseURL + path
}

// buildHeaders creates HTTP headers from the endpoint's httpHeaders field
func (cs *CurlService) buildHeaders(endpoint *models.Endpoint, varKVs map[string]string) []string {
	var headers []string

	for _, kv := range endpoint.HttpHeaders.VarKVs {
		key := cs.substituteVariables(kv.Key, varKVs)
		value := cs.substituteVariables(kv.Value, varKVs)
		headers = append(headers, fmt.Sprintf("%s: %s", key, value))
	}

	return headers
}

// buildCookies creates cookie string from the endpoint's httpCookies field
func (cs *CurlService) buildCookies(endpoint *models.Endpoint, varKVs map[string]string) string {
	var cookies []string

	for _, kv := range endpoint.HttpCookies.VarKVs {
		key := cs.substituteVariables(kv.Key, varKVs)
		value := cs.substituteVariables(kv.Value, varKVs)
		cookies = append(cookies, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(cookies, "; ")
}

// buildQueryParams creates query parameters from the endpoint's httpQueries field
func (cs *CurlService) buildQueryParams(endpoint *models.Endpoint, varKVs map[string]string) string {
	var params []string

	for _, kv := range endpoint.HttpQueries.VarKVs {
		key := cs.substituteVariables(kv.Key, varKVs)
		value := cs.substituteVariables(kv.Value, varKVs)
		params = append(params, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
	}

	return strings.Join(params, "&")
}

// buildBody creates the request body with variable substitution
func (cs *CurlService) buildBody(endpoint *models.Endpoint, varKVs map[string]string) string {
	return cs.substituteVariables(endpoint.HttpBody, varKVs)
}

// substituteVariables replaces variables in a VarString with provided values
func (cs *CurlService) substituteVariables(vs mystructs.VarString, varKVs map[string]string) string {
	// Create a copy of the VarString to avoid modifying the original
	substituted := vs

	// Inject the provided variables
	substituted.Inject(varKVs)

	// Return the executed result
	return substituted.Exec()
}
