-- Tasks queries stub
-- name: ListWorkTypes :many
SELECT * FROM tasks_worktype ORDER BY name LIMIT $1 OFFSET $2;
