-- name: CreateIdempotencyKey :one
INSERT INTO idempotency_keys ("key", "type", type_id, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetValidIdempotencyKey :one
SELECT *
FROM idempotency_keys
WHERE key = $1
  AND expires_at > NOW()
FOR UPDATE;
