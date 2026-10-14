package daos

import "time"

type Team struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
