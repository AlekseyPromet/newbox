-- DCIM queries stub
-- name: ListRegions :many
SELECT * FROM dcim_regions ORDER BY name LIMIT $1 OFFSET $2;
