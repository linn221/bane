package models

import (
	"time"

	"gorm.io/gorm"
)

type TodoStatus string

const (
	ToDoStatusFinished   TodoStatus = "finished"
	ToDoStatusInProgress TodoStatus = "in_progress"
	ToDoStatusCancelled  TodoStatus = "cancelled"
)

type Todo struct {
	Id          int        `gorm:"primaryKey"`
	Title       string     `gorm:"not null"`
	Description string     `gorm:"type:text"`
	Status      TodoStatus `gorm:"not null;default:'in_progress'"`
	Priority    int        `gorm:"default:0"`
	Deadline    time.Time
	Created     time.Time `gorm:"not null"`
	ProjectId   int       `gorm:"index"`
	ParentId    int       `gorm:"index"` // id of todo
}

type TodoInput struct {
	Title        string      `json:"title"`
	Description  string      `json:"description,omitempty"`
	Priority     int         `json:"priority,omitempty"`
	Deadline     *MyDate     `json:"deadline,omitempty"`
	Alias        string      `json:"alias,omitempty"`
	ProjectAlias string      `json:"projectAlias,omitempty"`
	Status       *TodoStatus `json:"status,omitempty"`
}

func (input *TodoInput) Validate(db *gorm.DB, id int) error {
	// Add validation logic here if needed
	return nil
}

type TodoFilter struct {
	Title        string     `json:"title,omitempty"`
	Status       TodoStatus `json:"status,omitempty"`
	ProjectId    int        `json:"projectId,omitempty"`
	ProjectAlias string     `json:"projectAlias,omitempty"`
	ParentId     int        `json:"parentId,omitempty"`
	PriorityMin  int        `json:"priorityMin,omitempty"`
	PriorityMax  int        `json:"priorityMax,omitempty"`
	DeadlineFrom *MyDate    `json:"deadlineFrom,omitempty"`
	DeadlineTo   *MyDate    `json:"deadlineTo,omitempty"`
	Search       string     `json:"search,omitempty"`
}
