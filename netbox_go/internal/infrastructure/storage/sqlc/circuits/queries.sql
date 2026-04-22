-- Circuits queries stub
-- name: ListProviders :many
SELECT * FROM circuits_providers ORDER BY name LIMIT $1 OFFSET $2;
