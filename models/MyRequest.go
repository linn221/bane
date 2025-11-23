package models

import (
	"time"
)

// MyRequest represents a request execution with response data
// This model stores all information about HTTP request execution for ethical hacking
type MyRequest struct {
	Id         int      `gorm:"primaryKey"`
	EndpointId int      `gorm:"not null;index"`
	Endpoint   Endpoint `gorm:"foreignKey:EndpointId"`

	// Request information
	RequestMethod  string `gorm:"size:10;not null"`
	RequestUrl     string `gorm:"not null"`
	RequestHeaders string `gorm:"type:text"` // JSON string of headers
	RequestBody    string `gorm:"type:text"`

	// Response information
	ResponseStatus  int    `gorm:"not null"`
	ResponseHeaders string `gorm:"type:text"` // JSON string of response headers
	ResponseBody    string `gorm:"type:longtext"`
	ContentType     string `gorm:"size:255"`
	ContentLength   int64  `gorm:"default:0"`

	// Performance metrics
	Latency int64 `gorm:"not null"`  // Latency in milliseconds
	Size    int64 `gorm:"default:0"` // Response size in bytes

	// Execution metadata
	ExecutedAt  time.Time `gorm:"autoCreateTime"`
	Variables   string    `gorm:"type:text"` // JSON string of variables used
	CurlCommand string    `gorm:"type:text"` // The actual curl command executed

	// Error information
	Error   string `gorm:"type:text"` // Error message if request failed
	Success bool   `gorm:"default:true"`
}

// MyRequestFilter for filtering requests
type MyRequestFilter struct {
	EndpointId int    `json:"endpointId,omitempty"`
	Success    *bool  `json:"success,omitempty"`
	StatusMin  int    `json:"statusMin,omitempty"`
	StatusMax  int    `json:"statusMax,omitempty"`
	DateFrom   string `json:"dateFrom,omitempty"`
	DateTo     string `json:"dateTo,omitempty"`
}
