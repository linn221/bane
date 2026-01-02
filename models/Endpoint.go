package models

import (
	"strings"

	"github.com/linn221/bane/mystructs"
)

type HttpSchema string

const (
	HttpSchemaHttp  HttpSchema = "http"
	HttpSchemaHttps HttpSchema = "https"
)

type HttpMethod string

const (
	HttpMethodGet     HttpMethod = "GET"
	HttpMethodPost    HttpMethod = "POST"
	HttpMethodPut     HttpMethod = "PUT"
	HttpMethodDelete  HttpMethod = "DELETE"
	HttpMethodPatch   HttpMethod = "PATCH"
	HttpMethodHead    HttpMethod = "HEAD"
	HttpMethodOptions HttpMethod = "OPTIONS"
	HttpMethodTrace   HttpMethod = "TRACE"
)

type Endpoint struct {
	Id          int                  `gorm:"primaryKey"`
	Name        string               `gorm:"size:255;default:null"`
	Description string               `gorm:"default:null"`
	ProjectId   *int                 `gorm:"default:null;index"`          // Optional project reference
	Https       bool                 `gorm:"not null;column:http_schema"` // true for https, false for http
	Method      HttpMethod           `gorm:"size:10;not null;column:http_method"`
	Domain      string               `gorm:"index;not null;column:http_domain"`
	Path        mystructs.VarString  `gorm:"not null;column:http_path"`
	Queries     mystructs.VarKVGroup `gorm:"not null;column:http_queries"`
	Headers     mystructs.VarKVGroup `gorm:"not null;column:http_headers"`
	Body        mystructs.VarString  `gorm:"not null;column:http_body"`
	Input       string               `gorm:"type:text;default:null"` // JSON-encoded EndpointInput for review
	// Vulns       []Vuln               `gorm:"many2many:endpoint_vulns"`
}

type EndpointInput struct {
	Name        string               `json:"name,omitempty"` // Optional
	Description string               `json:"description"`
	ProjectId   *int                 `json:"projectId,omitempty"` // Optional project reference
	Method      *HttpMethod          `json:"method"`              // Required HTTP method
	Url         mystructs.VarString  `json:"url"`                 // Full URL with optional VarString placeholders
	Headers     mystructs.VarKVGroup `json:"headers"`             // HTTP headers
	Body        *mystructs.VarString `json:"body,omitempty"`      // Optional HTTP body
}

type PatchEndpoint struct {
	Name        *string               `json:"name,omitempty"`
	Alias       *string               `json:"alias,omitempty"`
	Description *string               `json:"description,omitempty"`
	Https       *bool                 `json:"https,omitempty"`
	Method      *HttpMethod           `json:"method,omitempty"`
	Domain      *string               `json:"domain,omitempty"`
	Path        *mystructs.VarString  `json:"path,omitempty"`
	Queries     *mystructs.VarKVGroup `json:"queries,omitempty"`
	Headers     *mystructs.VarKVGroup `json:"headers,omitempty"`
	Body        *mystructs.VarString  `json:"body,omitempty"`
}
type EndpointFilter struct {
	Https  *bool      `json:"https,omitempty"` // true for https, false for http, nil for both
	Method HttpMethod `json:"method,omitempty"`
	Domain string     `json:"domain,omitempty"`
	Search string     `json:"search,omitempty"`
}

func (e *Endpoint) Text() string {
	return strings.Join([]string{
		e.Name,
		e.Description,
		e.Domain,
		e.Path.Exec(),
		e.Queries.Exec(),
		e.Headers.Exec(),
		e.Body.Exec(),
	}, " ")
}
