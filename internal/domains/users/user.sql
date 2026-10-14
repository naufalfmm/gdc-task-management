-- name: InsertUser :one
INSERT INTO users (email, password_hash, name, team_id)
VALUES ($1, $2, $3, $4)
RETURNING *;