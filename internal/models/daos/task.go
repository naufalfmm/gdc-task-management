package daos

import "time"

const (
	TaskStatusToDo       string = "TODO"
	TaskStatusInProgress string = "IN_PROGRESS"
)

type Task struct {
	ID          int64
	UserID      int64
	TeamID      int64
	AssigneeID  *int64
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type TaskLog struct {
	ID        int64
	TaskID    int64
	ActorID   int64
	Action    string
	Details   map[string]any
	CreatedAt time.Time
}
