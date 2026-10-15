package daos

import "time"

const (
	TaskStatusToDo       string = "TODO"
	TaskStatusInProgress string = "IN_PROGRESS"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      string
	UserID      int64
	AssigneeID  *int64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
