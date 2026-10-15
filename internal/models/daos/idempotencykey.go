package daos

import "time"

const (
	IdempotencyKeyTypeTask string = "TASK"
)

type IdempotencyKey struct {
	Key       string
	Type      string
	TypeID    int64
	ExpiresAt time.Time
	CreatedAt time.Time
}
