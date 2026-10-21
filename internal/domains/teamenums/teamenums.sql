-- name: GetTeamEnum :one
SELECT *
FROM team_enums
WHERE team_id = $1
    AND enum = $2;