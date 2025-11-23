package models

import (
	"strings"

	"github.com/linn221/bane/mystructs"
	"github.com/linn221/bane/validate"
	"gorm.io/gorm"
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
	Id                  int                  `gorm:"primaryKey"`
	Name                string               `gorm:"size:255;not null"`
	Description         string               `gorm:"default:null"`
	HttpSchema          HttpSchema           `gorm:"size:10;not null"`
	HttpMethod          HttpMethod           `gorm:"size:10;not null"`
	HttpDomain          string               `gorm:"index;not null"`
	HttpPort            int                  `gorm:"default:0"` // 0 means default port (80/443)
	HttpPath            mystructs.VarString  `gorm:"not null"`
	HttpQueries         mystructs.VarKVGroup `gorm:"not null"`
	HttpHeaders         mystructs.VarKVGroup `gorm:"not null"`
	HttpCookies         mystructs.VarKVGroup `gorm:"not null"`
	HttpBody            mystructs.VarString  `gorm:"not null"`
	HttpTimeout         int                  `gorm:"default:30"`   // Timeout in seconds
	HttpFollowRedirects bool                 `gorm:"default:true"` // Follow redirects
	// Vulns       []Vuln               `gorm:"many2many:endpoint_vulns"`
}

type EndpointInput struct {
	Name                string               `json:"name"`
	Alias               string               `json:"alias,omitempty"`
	Description         string               `json:"description"`
	HttpSchema          HttpSchema           `json:"http_schema"`
	HttpMethod          HttpMethod           `json:"http_method"`
	HttpDomain          string               `json:"http_domain"`
	HttpPort            int                  `json:"http_port,omitempty"`
	HttpPath            mystructs.VarString  `json:"http_path"`
	HttpQueries         mystructs.VarKVGroup `json:"http_queries"`
	HttpHeaders         mystructs.VarKVGroup `json:"http_headers"`
	HttpCookies         mystructs.VarKVGroup `json:"http_cookies"`
	HttpBody            mystructs.VarString  `json:"http_body"`
	HttpTimeout         int                  `json:"http_timeout,omitempty"`
	HttpFollowRedirects bool                 `json:"http_follow_redirects,omitempty"`
	// Vulns       []string             `json:"vulns"` // alias of vulns to attach to the endpoint
}

func (input *EndpointInput) Validate(db *gorm.DB, id int) error {
	var rules []validate.Rule

	// Validate alias uniqueness if provided
	if input.Alias != "" {
		rules = append(rules, validate.NewUniqueRule("endpoints", "alias", input.Alias, nil).Except(id).Say("duplicate alias for endpoint"))
	}

	return validate.Validate(db, rules...)
}

// PatchEndpoint represents partial updates for Endpoint
type PatchEndpoint struct {
	Name                *string               `json:"name,omitempty"`
	Alias               *string               `json:"alias,omitempty"`
	Description         *string               `json:"description,omitempty"`
	HttpSchema          *HttpSchema           `json:"http_schema,omitempty"`
	HttpMethod          *HttpMethod           `json:"http_method,omitempty"`
	HttpDomain          *string               `json:"http_domain,omitempty"`
	HttpPort            *int                  `json:"http_port,omitempty"`
	HttpPath            *mystructs.VarString  `json:"http_path,omitempty"`
	HttpQueries         *mystructs.VarKVGroup `json:"http_queries,omitempty"`
	HttpHeaders         *mystructs.VarKVGroup `json:"http_headers,omitempty"`
	HttpCookies         *mystructs.VarKVGroup `json:"http_cookies,omitempty"`
	HttpBody            *mystructs.VarString  `json:"http_body,omitempty"`
	HttpTimeout         *int                  `json:"http_timeout,omitempty"`
	HttpFollowRedirects *bool                 `json:"http_follow_redirects,omitempty"`
}

func (input *PatchEndpoint) Validate(db *gorm.DB, id int) error {
	var rules []validate.Rule

	// Validate alias uniqueness if provided
	if input.Alias != nil && *input.Alias != "" {
		rules = append(rules, validate.NewUniqueRule("endpoints", "alias", *input.Alias, nil).Except(id).Say("duplicate alias for endpoint"))
	}

	return validate.Validate(db, rules...)
}

type EndpointFilter struct {
	HttpSchema HttpSchema `json:"httpSchema,omitempty"`
	HttpMethod HttpMethod `json:"httpMethod,omitempty"`
	HttpDomain string     `json:"httpDomain,omitempty"`
	Search     string     `json:"search,omitempty"`
}

type CurlImportResult struct {
	Name                string               `json:"name"`
	Description         string               `json:"description"`
	HttpSchema          HttpSchema           `json:"httpSchema"`
	HttpMethod          HttpMethod           `json:"httpMethod"`
	HttpDomain          string               `json:"httpDomain"`
	HttpPort            int                  `json:"httpPort"`
	HttpPath            mystructs.VarString  `json:"httpPath"`
	HttpQueries         mystructs.VarKVGroup `json:"httpQueries"`
	HttpHeaders         mystructs.VarKVGroup `json:"httpHeaders"`
	HttpCookies         mystructs.VarKVGroup `json:"httpCookies"`
	HttpBody            mystructs.VarString  `json:"httpBody"`
	HttpTimeout         int                  `json:"httpTimeout"`
	HttpFollowRedirects bool                 `json:"httpFollowRedirects"`
}

func (e *Endpoint) Text() string {
	return strings.Join([]string{
		e.Name,
		e.Description,
		e.HttpDomain,
		e.HttpPath.Exec(),
		e.HttpQueries.Exec(),
		e.HttpHeaders.Exec(),
		e.HttpCookies.Exec(),
		e.HttpBody.Exec(),
	}, " ")
}
