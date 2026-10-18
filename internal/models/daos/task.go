package daos

import (
	"strconv"
	"strings"
	"time"
)

const (
	TaskStatusToDo       string = "TODO"
	TaskStatusInProgress string = "IN_PROGRESS"
)

type Task struct {
	ID          int64
	Code        string
	UserID      int64
	TeamID      int64
	AssigneeID  *int64
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

func (t Task) SplitCode() (string, int64) {
	parts := strings.SplitN(t.Code, "-", 2)
	if len(parts) < 2 {
		return parts[0], 0
	}

	latestNumber, _ := strconv.ParseInt(parts[1], 10, 64)

	return parts[0], latestNumber
}

type TaskLog struct {
	ID        int64
	TaskID    int64
	ActorID   int64
	Action    string
	Details   map[string]any
	CreatedAt time.Time
}
