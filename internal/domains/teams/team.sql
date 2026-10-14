-- name: GetTeamByName :one
SELECT * FROM teams WHERE name = $1;

-- name: InsertTeam :one
INSERT INTO teams (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: GetTeamByID :one
SELECT * FROM teams WHERE id = $1;