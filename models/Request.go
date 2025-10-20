package models

import (
	"time"

	"github.com/linn221/bane/mystructs"
)

type Job struct {
	Id          int       `gorm:"primaryKey"`
	ProgramId   int       `gorm:"not null;index"`
	Program     Program   `gorm:"foreignKey:ProgramId"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"default:null"`
	JobDate     time.Time `gorm:"not null"`
}

type Request struct {
	Id                  int               `gorm:"primaryKey"`
	EndpointId          int               `gorm:"not null;index"`
	Endpoint            Endpoint          `gorm:"foreignKey:EndpointId"`
	JobId               int               `gorm:"not null;index"`
	Job                 Job               `gorm:"foreignKey:JobId"`
	SequenceNumber      int               `gorm:"not null"`
	ProgramId           int               `gorm:"not null;index"`
	Program             Program           `gorm:"foreignKey:ProgramId"`
	Description         string            `gorm:"default:null"`
	HttpSchema          HttpSchema        `gorm:"size:10;not null"`
	HttpMethod          HttpMethod        `gorm:"size:10;not null"`
	HttpDomain          string            `gorm:"index;not null"`
	HttpPath            string            `gorm:"not null"`
	HttpQueries         mystructs.KVGroup `gorm:"not null"`
	HttpHeaders         mystructs.KVGroup `gorm:"not null"`
	HttpCookies         mystructs.KVGroup `gorm:"not null"`
	HttpBody            string            `gorm:"not null"`
	ResponseStatusCode  int               `gorm:"not null"`
	ResponseContentType string            `gorm:"not null"`
	ResponseLatency     time.Duration     `gorm:"not null"`
	ResponseSize        int               `gorm:"not null"`
	ResponseBody        string            `gorm:"not null"`
	ResponseHeaders     mystructs.KVGroup `gorm:"not null"`
	ResponseCookies     mystructs.KVGroup `gorm:"not null"`
}

type NewJob struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	JobDate     time.Time `json:"job_date"`
}
