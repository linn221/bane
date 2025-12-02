package models

import (
	"errors"
	"io"
	"strconv"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusFinished   TaskStatus = "finished"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCancelled  TaskStatus = "cancelled"
	TaskStatusDeadlined  TaskStatus = "deadlined"
)

// MarshalGQL implements the graphql.Marshaler interface.
func (t TaskStatus) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(string(t))))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (t *TaskStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("task status must be a string")
	}
	switch str {
	case "finished":
		*t = TaskStatusFinished
	case "in_progress":
		*t = TaskStatusInProgress
	case "cancelled":
		*t = TaskStatusCancelled
	case "deadlined":
		*t = TaskStatusDeadlined
	default:
		return errors.New("invalid task status")
	}
	return nil
}

type Task struct {
	Id            int    `gorm:"primaryKey"`
	Title         string `gorm:"not null"`
	Description   string `gorm:"type:text"`
	Status        TaskStatus
	Priority      int `gorm:"default:0"`
	Deadline      MyDate
	RemindDate    MyDate
	FinishedDate  MyDate
	CancelledDate MyDate
	Created       MyDate `gorm:"not null"`
	ProjectId     int    `gorm:"index"`
}

type TaskInput struct {
	Title        string  `json:"title"`
	Description  string  `json:"description,omitempty"`
	Priority     int     `json:"priority,omitempty"`
	Deadline     *MyDate `json:"deadline,omitempty"`
	RemindDate   *MyDate `json:"remindDate,omitempty"`
	Alias        string  `json:"alias,omitempty"`
	ProjectAlias string  `json:"projectAlias,omitempty"`
}

func (input *TaskInput) Validate(db *gorm.DB, id int) error {
	// Add validation logic here if needed
	return nil
}

type TaskFilter struct {
	Today   bool   `json:"today,omitempty"`
	Search  string `json:"search,omitempty"`
	Project *string
	Status  *TaskStatus `json:"status,omitempty"`
}
