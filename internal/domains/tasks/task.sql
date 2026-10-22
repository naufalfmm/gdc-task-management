-- name: CreateTask :one
INSERT INTO tasks (title, description, status, user_id, assignee_id)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTaskByID :one
SELECT *
FROM tasks
WHERE id = $1
    AND deleted_at IS NULL
FOR UPDATE;

-- name: GetTaskByCode :one
SELECT *
FROM tasks
WHERE code = $1
    AND deleted_at IS NULL
FOR UPDATE;

-- name: AssignTask :one
UPDATE tasks
SET assignee_id = $1,
    updated_at = NOW()
WHERE id = $2
    AND deleted_at IS NULL
RETURNING *;

-- name: GetLatestActiveTeamTask :one
SELECT *
FROM tasks
WHERE team_id = $1
    AND deleted_at IS NULL
ORDER BY id DESC
LIMIT 1
FOR UPDATE;

-- name: UpdateTask :one
UPDATE tasks
SET
    title = COALESCE(NULLIF(sqlc.arg(title)::text, ''), title),
    description = COALESCE(NULLIF(sqlc.arg(description)::text, ''), description),
    status = COALESCE(NULLIF(sqlc.arg(status)::text, ''), status),
    updated_at = NOW()
WHERE 
    CASE WHEN sqlc.arg(id)::bigint <> 0 THEN id = sqlc.arg(id)::bigint ELSE TRUE END
    AND CASE WHEN sqlc.arg(code)::text <> '' THEN code = sqlc.arg(code)::text ELSE TRUE END
    AND deleted_at IS NOT NULL
RETURNING *;

-- name: DeleteTask :one
UPDATE tasks
SET
    deleted_at = NOW()
WHERE 
    CASE WHEN sqlc.arg(id)::bigint <> 0 THEN id = sqlc.arg(id)::bigint ELSE TRUE END
    AND CASE WHEN sqlc.arg(code)::text <> '' THEN code = sqlc.arg(code)::text ELSE TRUE END
    AND deleted_at IS NOT NULL
RETURNING *;  