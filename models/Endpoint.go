package models

import "github.com/linn221/bane/mystructs"

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
	ProgramId   int                  `gorm:"not null;index"`
	Program     Program              `gorm:"foreignKey:ProgramId"`
	Name        string               `gorm:"size:255;not null"`
	Description string               `gorm:"default:null"`
	HttpSchema  HttpSchema           `gorm:"size:10;not null"`
	HttpMethod  HttpMethod           `gorm:"size:10;not null"`
	HttpDomain  string               `gorm:"index;not null"`
	HttpPath    mystructs.VarString  `gorm:"not null"`
	HttpQueries mystructs.VarKVGroup `gorm:"not null"`
	HttpHeaders mystructs.VarKVGroup `gorm:"not null"`
	HttpCookies mystructs.VarKVGroup `gorm:"not null"`
	HttpBody    mystructs.VarString  `gorm:"not null"`
	Vulns       []Vuln               `gorm:"many2many:endpoint_vulns"`
}

type NewEndpoint struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	HttpSchema  HttpSchema           `json:"http_schema"`
	HttpMethod  HttpMethod           `json:"http_method"`
	HttpDomain  string               `json:"http_domain"`
	HttpPath    string               `json:"http_path"`
	HttpQueries mystructs.VarKVGroup `json:"http_queries"`
	HttpHeaders mystructs.VarKVGroup `json:"http_headers"`
	HttpCookies mystructs.VarKVGroup `json:"http_cookies"`
	HttpBody    string               `json:"http_body"`
	Vulns       []string             `json:"vulns"` // alias of vulns to attach to the endpoint
}
