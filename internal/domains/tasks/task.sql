-- name: CreateTask :one
INSERT INTO tasks (title, description, status, user_id, assignee_id)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTaskByID :one
SELECT *
FROM tasks
WHERE id = $1
FOR UPDATE;

-- name: AssignTask :one
UPDATE tasks
SET assignee_id = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;