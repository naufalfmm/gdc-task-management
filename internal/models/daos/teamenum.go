package daos

import "time"

type TeamEnum struct {
	ID           int64
	TeamID       int64
	Enum         string
	Data         []string
	DefaultDatum *string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
