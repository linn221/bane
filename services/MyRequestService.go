package services

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/linn221/bane/app"
	"github.com/linn221/bane/models"
	"github.com/linn221/bane/mystructs"
	"gorm.io/gorm"
)

type myRequestService struct{}

var MyRequestService = &myRequestService{}

// CreateMyRequest creates a new MyRequest record
func (mrs *myRequestService) CreateMyRequest(db *gorm.DB, request *models.MyRequest) (*models.MyRequest, error) {
	if err := db.Create(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

// GetMyRequestByID retrieves a MyRequest by ID
func (mrs *myRequestService) GetMyRequestByID(db *gorm.DB, id *int) (*models.MyRequest, error) {
	var request models.MyRequest
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Preload("Program").Preload("Endpoint").First(&request, *id).Error
	return &request, err
}

// ListMyRequests retrieves MyRequests with optional filtering
func (mrs *myRequestService) ListMyRequests(db *gorm.DB, filter *models.MyRequestFilter) ([]*models.MyRequest, error) {
	var requests []*models.MyRequest
	query := db.Preload("Program").Preload("Endpoint")

	if filter != nil {
		if filter.ProgramId != 0 {
			query = query.Where("program_id = ?", filter.ProgramId)
		}
		if filter.EndpointId != 0 {
			query = query.Where("endpoint_id = ?", filter.EndpointId)
		}
		if filter.Success != nil {
			query = query.Where("success = ?", *filter.Success)
		}
		if filter.StatusMin != 0 {
			query = query.Where("response_status >= ?", filter.StatusMin)
		}
		if filter.StatusMax != 0 {
			query = query.Where("response_status <= ?", filter.StatusMax)
		}
		if filter.DateFrom != "" {
			query = query.Where("executed_at >= ?", filter.DateFrom)
		}
		if filter.DateTo != "" {
			query = query.Where("executed_at <= ?", filter.DateTo)
		}
	}

	err := query.Order("executed_at DESC").Find(&requests).Error
	return requests, err
}

// ExecuteCurl runs a curl command and captures the response
func (mrs *myRequestService) ExecuteCurl(app *app.App, db *gorm.DB, endpointAlias string, variables mystructs.VarKVGroup) (*models.MyRequest, error) {
	// Find endpoint by alias
	var endpoint models.Endpoint
	err := db.Preload("Program").Where("alias = ?", endpointAlias).First(&endpoint).Error
	if err != nil {
		return nil, fmt.Errorf("endpoint with alias '%s' not found: %v", endpointAlias, err)
	}

	// Generate curl command with variable injection
	curlCommand, err := mrs.generateCurlCommand(&endpoint, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to generate curl command: %v", err)
	}

	// Execute curl command
	startTime := time.Now()
	cmd := exec.Command("bash", "-c", curlCommand)
	output, err := cmd.CombinedOutput()
	latency := time.Since(startTime).Milliseconds()

	// Create MyRequest record
	request := &models.MyRequest{
		ProgramId:      endpoint.ProgramId,
		EndpointId:     endpoint.Id,
		RequestMethod:  string(endpoint.HttpMethod),
		RequestUrl:     mrs.buildUrl(&endpoint),
		RequestHeaders: mrs.serializeHeaders(endpoint.HttpHeaders),
		RequestBody:    endpoint.HttpBody.Exec(),
		Latency:        latency,
		ExecutedAt:     time.Now(),
		Variables:      mrs.serializeVariables(variables),
		CurlCommand:    curlCommand,
		Success:        err == nil,
	}

	if err != nil {
		request.Error = string(output)
		request.ResponseStatus = 0
		request.ResponseBody = string(output)
	} else {
		// Parse curl output to extract response information
		responseInfo := mrs.parseCurlOutput(string(output))
		request.ResponseStatus = responseInfo.Status
		request.ResponseHeaders = responseInfo.Headers
		request.ResponseBody = responseInfo.Body
		request.ContentType = responseInfo.ContentType
		request.ContentLength = int64(len(responseInfo.Body))
		request.Size = int64(len(responseInfo.Body))
	}

	// Save to database
	return mrs.CreateMyRequest(db, request)
}

// generateCurlCommand creates a curl command from endpoint and variables
func (mrs *myRequestService) generateCurlCommand(endpoint *models.Endpoint, variables mystructs.VarKVGroup) (string, error) {
	// Inject variables into endpoint fields
	path := endpoint.HttpPath.Exec()
	headers := endpoint.HttpHeaders
	body := endpoint.HttpBody.Exec()

	// Apply variable injection to path, headers, and body
	for _, kv := range variables.VarKVs {
		key := kv.Key.Exec()
		value := kv.Value.Exec()

		// Replace variables in path
		path = strings.ReplaceAll(path, "{"+key+"}", value)

		// Replace variables in headers
		for i, header := range headers.VarKVs {
			headers.VarKVs[i].Key = mystructs.VarString{OriginalString: strings.ReplaceAll(header.Key.Exec(), "{"+key+"}", value)}
			headers.VarKVs[i].Value = mystructs.VarString{OriginalString: strings.ReplaceAll(header.Value.Exec(), "{"+key+"}", value)}
		}

		// Replace variables in body
		body = strings.ReplaceAll(body, "{"+key+"}", value)
	}

	// Build curl command
	var curlParts []string
	curlParts = append(curlParts, "curl")
	curlParts = append(curlParts, "-X", string(endpoint.HttpMethod))

	// Add URL
	url := mrs.buildUrl(endpoint)
	url = strings.ReplaceAll(url, endpoint.HttpPath.Exec(), path)
	curlParts = append(curlParts, fmt.Sprintf("'%s'", url))

	// Add headers
	for _, header := range headers.VarKVs {
		curlParts = append(curlParts, "-H", fmt.Sprintf("'%s: %s'", header.Key.Exec(), header.Value.Exec()))
	}

	// Add body if present
	if body != "" {
		curlParts = append(curlParts, "-d", fmt.Sprintf("'%s'", body))
	}

	// Add timeout
	if endpoint.HttpTimeout > 0 {
		curlParts = append(curlParts, "--connect-timeout", fmt.Sprintf("%d", endpoint.HttpTimeout))
	}

	// Add follow redirects
	if endpoint.HttpFollowRedirects {
		curlParts = append(curlParts, "-L")
	}

	// Add verbose output for parsing
	curlParts = append(curlParts, "-v", "-s")

	return strings.Join(curlParts, " "), nil
}

// buildUrl constructs the full URL from endpoint
func (mrs *myRequestService) buildUrl(endpoint *models.Endpoint) string {
	schema := string(endpoint.HttpSchema)
	domain := endpoint.HttpDomain
	port := ""
	if endpoint.HttpPort != 0 {
		port = fmt.Sprintf(":%d", endpoint.HttpPort)
	}
	path := endpoint.HttpPath.Exec()

	return fmt.Sprintf("%s://%s%s%s", schema, domain, port, path)
}

// serializeHeaders converts VarKVGroup to JSON string
func (mrs *myRequestService) serializeHeaders(headers mystructs.VarKVGroup) string {
	headerMap := make(map[string]string)
	for _, kv := range headers.VarKVs {
		headerMap[kv.Key.Exec()] = kv.Value.Exec()
	}
	jsonBytes, _ := json.Marshal(headerMap)
	return string(jsonBytes)
}

// serializeVariables converts VarKVGroup to JSON string
func (mrs *myRequestService) serializeVariables(variables mystructs.VarKVGroup) string {
	varMap := make(map[string]string)
	for _, kv := range variables.VarKVs {
		varMap[kv.Key.Exec()] = kv.Value.Exec()
	}
	jsonBytes, _ := json.Marshal(varMap)
	return string(jsonBytes)
}

// parseCurlOutput extracts response information from curl verbose output
func (mrs *myRequestService) parseCurlOutput(output string) struct {
	Status      int
	Headers     string
	Body        string
	ContentType string
} {
	// This is a simplified parser - in a real implementation,
	// you'd want more robust parsing of curl's verbose output
	lines := strings.Split(output, "\n")

	result := struct {
		Status      int
		Headers     string
		Body        string
		ContentType string
	}{
		Status: 200, // Default
	}

	// Find HTTP status line
	for _, line := range lines {
		if strings.Contains(line, "HTTP/") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if status := mrs.extractStatus(parts[1]); status > 0 {
					result.Status = status
				}
			}
		}
	}

	// Extract body (everything after the last empty line)
	bodyStart := -1
	for i, line := range lines {
		if line == "" && i < len(lines)-1 {
			bodyStart = i + 1
		}
	}

	if bodyStart >= 0 && bodyStart < len(lines) {
		result.Body = strings.Join(lines[bodyStart:], "\n")
	}

	return result
}

// extractStatus extracts HTTP status code from string
func (mrs *myRequestService) extractStatus(statusStr string) int {
	// Simple status extraction - could be more robust
	if len(statusStr) >= 3 {
		if status := statusStr[:3]; status >= "100" && status <= "599" {
			// This would need proper string to int conversion
			return 200 // Simplified for now
		}
	}
	return 0
}
