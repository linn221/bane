package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
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

// monthMap maps lowercase 3-4 letter month abbreviations to month numbers
var monthMap = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "sept": 9, "oct": 10, "nov": 11, "dec": 12,
}

// parseMonthName parses a lowercase 3-4 letter month name and returns the month number
func parseMonthName(monthStr string) (int, error) {
	monthStr = strings.ToLower(strings.TrimSpace(monthStr))
	if len(monthStr) < 3 {
		return 0, fmt.Errorf("month name too short")
	}

	// Try exact match first (for 3-4 letter inputs)
	if month, ok := monthMap[monthStr]; ok {
		return month, nil
	}

	// If longer than 4, try first 4 characters, then first 3
	if len(monthStr) > 4 {
		if month, ok := monthMap[monthStr[:4]]; ok {
			return month, nil
		}
		if month, ok := monthMap[monthStr[:3]]; ok {
			return month, nil
		}
	} else if len(monthStr) == 4 {
		// For 4-letter input, also try first 3 characters
		if month, ok := monthMap[monthStr[:3]]; ok {
			return month, nil
		}
	}

	return 0, fmt.Errorf("invalid month name: %s", monthStr)
}

// parseDateString parses date strings like "jan 8" or "jan 8 2024"
func parseDateString(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)

	// Pattern: month (3-4 letters) day [year]
	// Examples: "jan 8", "jan 8 2024", "december 25 2023"
	re := regexp.MustCompile(`^([a-z]{3,4})\s+(\d{1,2})(?:\s+(\d{4}))?$`)
	matches := re.FindStringSubmatch(strings.ToLower(dateStr))

	if len(matches) == 0 {
		return time.Time{}, fmt.Errorf("invalid date format, expected format like 'jan 8' or 'jan 8 2024'")
	}

	monthStr := matches[1]
	dayStr := matches[2]
	yearStr := matches[3]

	month, err := parseMonthName(monthStr)
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %s", dayStr)
	}

	year := time.Now().Year()
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid year: %s", yearStr)
		}
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}

// parseDateTimeString parses date-time strings like "jan 8 11:01 PM" or "jan 8 2024 11:01 PM"
func parseDateTimeString(dateTimeStr string) (time.Time, error) {
	dateTimeStr = strings.TrimSpace(dateTimeStr)

	// Pattern: month (3-4 letters) day [year] time (HH:MM AM/PM)
	// Examples: "jan 8 11:01 PM", "jan 8 2024 11:01 PM"
	re := regexp.MustCompile(`(?i)^([a-z]{3,4})\s+(\d{1,2})(?:\s+(\d{4}))?\s+(\d{1,2}):(\d{2})\s+(AM|PM)$`)
	matches := re.FindStringSubmatch(dateTimeStr)

	if len(matches) == 0 {
		return time.Time{}, fmt.Errorf("invalid date-time format, expected format like 'jan 8 11:01 PM' or 'jan 8 2024 11:01 PM'")
	}

	monthStr := strings.ToLower(matches[1])
	dayStr := matches[2]
	yearStr := matches[3]
	hourStr := matches[4]
	minuteStr := matches[5]
	ampm := strings.ToLower(matches[6])

	month, err := parseMonthName(monthStr)
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %s", dayStr)
	}

	year := time.Now().Year()
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid year: %s", yearStr)
		}
	}

	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour: %s", hourStr)
	}

	minute, err := strconv.Atoi(minuteStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute: %s", minuteStr)
	}

	// Convert to 24-hour format
	if ampm == "pm" && hour != 12 {
		hour += 12
	} else if ampm == "am" && hour == 12 {
		hour = 0
	}

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC), nil
}

// MarshalGQL implements the graphql.Marshaler interface.
func (u MyTime) MarshalGQL(w io.Writer) {
	// Format: "January 8 11:01 PM" (no year)
	monthName := u.Time.Format("January")
	day := u.Time.Day()
	timeStr := u.Time.Format("3:04 PM")
	formatted := fmt.Sprintf("%s %d %s", monthName, day, timeStr)
	fmt.Fprint(w, strconv.Quote(formatted))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (u *MyTime) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("MyTime must be a string")
	}

	parsedTime, err := parseDateTimeString(str)
	if err != nil {
		return err
	}

	u.Time = parsedTime
	return nil
}

// Value implements the driver.Valuer interface for GORM
func (u MyTime) Value() (driver.Value, error) {
	if u.Time.IsZero() {
		return nil, nil
	}
	return u.Time, nil
}

// Scan implements the sql.Scanner interface for GORM
func (u *MyTime) Scan(value interface{}) error {
	if value == nil {
		u.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		u.Time = v
		return nil
	case []byte:
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(v))
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02T15:04:05Z07:00", string(v))
			if err != nil {
				return fmt.Errorf("cannot scan %T into MyTime: %v", value, err)
			}
		}
		u.Time = parsedTime
		return nil
	case string:
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", v)
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02T15:04:05Z07:00", v)
			if err != nil {
				return fmt.Errorf("cannot scan %T into MyTime: %v", value, err)
			}
		}
		u.Time = parsedTime
		return nil
	default:
		return fmt.Errorf("cannot scan %T into MyTime", value)
	}
}

type MyDate struct {
	time.Time
}

// MarshalGQL implements the graphql.Marshaler interface.
func (u MyDate) MarshalGQL(w io.Writer) {
	// Format: "January 8 2025"
	monthName := u.Time.Format("January")
	day := u.Time.Day()
	year := u.Time.Year()
	formatted := fmt.Sprintf("%s %d %d", monthName, day, year)
	fmt.Fprint(w, strconv.Quote(formatted))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (u *MyDate) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("MyDate must be a string")
	}

	// Parse the date in format "jan 8" or "jan 8 2024"
	parsedTime, err := parseDateString(str)
	if err != nil {
		return err
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

type KVString struct {
	Key   string
	Value string
}

// KVString GraphQL methods
func (t KVString) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(fmt.Sprintf("%s:%v", t.Key, t.Value))))
}

func (t *KVString) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("kvstring type must be string")
	}
	sepIndex := strings.Index(str, ":")
	if sepIndex == -1 {
		return errors.New("invalid string")
	}
	key := str[:sepIndex]
	value := str[sepIndex+1:]
	t.Key = key
	t.Value = value
	return nil
}

type KVInt struct {
	Key   string
	Value int
}

// KVInt GraphQL methods
func (t KVInt) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(fmt.Sprintf("%s:%v", t.Key, t.Value))))
}

func (t *KVInt) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("kvstring type must be string")
	}
	sepIndex := strings.Index(str, ":")
	if sepIndex == -1 {
		return errors.New("invalid string")
	}
	key := str[:sepIndex]
	value := str[sepIndex+1:]
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("error converting the int")
	}
	t.Key = key
	t.Value = valueInt
	return nil
}

type PatchInput struct {
	Values    []KVString `json:"values,omitempty"`
	ValuesInt []KVInt    `json:"valuesInt,omitempty"`
}
