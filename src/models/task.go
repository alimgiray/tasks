package models

import (
	"time"
)

type Task struct {
	TaskID    int
	ProjectID int
	Key       string

	Deadline  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	Assignee    string
	Title       string
	Status      string
	Description string

	Expanded bool
}
