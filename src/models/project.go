package models

import "time"

type Project struct {
	ProjectID int
	Name      string

	CreatedAt time.Time
	UpdatedAt time.Time
}
