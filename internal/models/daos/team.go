package daos

import "time"

type Team struct {
	ID        int64
	Code      string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
