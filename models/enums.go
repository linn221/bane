package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

type WordType string

const (
	WordTypeFuzz   WordType = "fuzz"
	WordTypeAttack WordType = "attack"
	WordTypeRegex  WordType = "regex"
)

func (t WordType) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(string(t))))
}

func (t *WordType) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("word type must be string")
	}
	switch str {
	case "fuzz":
		*t = WordTypeFuzz
	case "attack":
		*t = WordTypeAttack
	case "regex":
		*t = WordTypeRegex
	default:
		return errors.New("invalid word type")
	}
	return nil
}

type MyTime struct {
	time.Time
}

// MarshalGQL implements the graphql.Marshaler interface.
func (u MyTime) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(u.Format(time.DateOnly)))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (u *MyTime) UnmarshalGQL(v interface{}) error {
	_, ok := v.(string)
	if !ok {
		return fmt.Errorf("uint must be a string")
	}

	return nil
}

type MyDate struct {
	time.Time
}

// MarshalGQL implements the graphql.Marshaler interface.
func (u MyDate) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(u.Format("02-01-2006")))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (u *MyDate) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("MyDate must be a string")
	}

	// Parse the date in format "15-08-2025"
	parsedTime, err := time.Parse("02-01-2006", str)
	if err != nil {
		return fmt.Errorf("invalid date format, expected DD-MM-YYYY: %v", err)
	}

	u.Time = parsedTime
	return nil
}

// Value implements the driver.Valuer interface for GORM
func (u MyDate) Value() (driver.Value, error) {
	if u.Time.IsZero() {
		return nil, nil
	}
	return u.Time, nil
}

// Scan implements the sql.Scanner interface for GORM
func (u *MyDate) Scan(value interface{}) error {
	if value == nil {
		u.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		u.Time = v
		return nil
	case []byte:
		// Try to parse as time.Time from bytes
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(v))
		if err != nil {
			// Try alternative format
			parsedTime, err = time.Parse("2006-01-02T15:04:05Z07:00", string(v))
			if err != nil {
				return fmt.Errorf("cannot scan %T into MyDate: %v", value, err)
			}
		}
		u.Time = parsedTime
		return nil
	case string:
		// Try to parse as time.Time from string
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", v)
		if err != nil {
			// Try alternative format
			parsedTime, err = time.Parse("2006-01-02T15:04:05Z07:00", v)
			if err != nil {
				return fmt.Errorf("cannot scan %T into MyDate: %v", value, err)
			}
		}
		u.Time = parsedTime
		return nil
	default:
		return fmt.Errorf("cannot scan %T into MyDate", value)
	}
}

// HttpSchema GraphQL methods
func (h HttpSchema) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(string(h))))
}

func (h *HttpSchema) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("http schema must be string")
	}
	switch str {
	case "http":
		*h = HttpSchemaHttp
	case "https":
		*h = HttpSchemaHttps
	default:
		return errors.New("invalid http schema")
	}
	return nil
}

// HttpMethod GraphQL methods
func (h HttpMethod) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(string(h))))
}

func (h *HttpMethod) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("http method must be string")
	}
	switch str {
	case "GET":
		*h = HttpMethodGet
	case "POST":
		*h = HttpMethodPost
	case "PUT":
		*h = HttpMethodPut
	case "DELETE":
		*h = HttpMethodDelete
	case "PATCH":
		*h = HttpMethodPatch
	case "HEAD":
		*h = HttpMethodHead
	case "OPTIONS":
		*h = HttpMethodOptions
	case "TRACE":
		*h = HttpMethodTrace
	default:
		return errors.New("invalid http method")
	}
	return nil
}

// VulnReferenceType GraphQL methods
func (v VulnReferenceType) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(string(v))))
}

func (v *VulnReferenceType) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("vuln reference type must be string")
	}
	switch str {
	case "programs":
		*v = VulnReferenceTypeProgram
	case "endpoints":
		*v = VulnReferenceTypeEndpoint
	case "requests":
		*v = VulnReferenceTypeRequest
	case "notes":
		*v = VulnReferenceTypeNote
	case "attachments":
		*v = VulnReferenceTypeAttachment
	case "vulns":
		*v = VulnReferenceTypeVuln
	default:
		return errors.New("invalid vuln reference type")
	}
	return nil
}

// TaggableType GraphQL methods
func (t TaggableType) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(string(t))))
}

func (t *TaggableType) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("taggable type must be string")
	}
	switch str {
	case "programs":
		*t = TaggableTypePrograms
	case "endpoints":
		*t = TaggableTypeEndpoints
	case "requests":
		*t = TaggableTypeRequests
	case "vulns":
		*t = TaggableTypeVulns
	case "notes":
		*t = TaggableTypeNotes
	default:
		return errors.New("invalid taggable type")
	}
	return nil
}
