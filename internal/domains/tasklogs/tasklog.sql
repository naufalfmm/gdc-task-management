-- name: CreateTaskLog :one
INSERT INTO task_logs (task_id, actor_id, action, details)
VALUES ($1, $2, $3, $4)
RETURNING *;