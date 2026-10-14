package daos

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Name         string
	TeamID       int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
