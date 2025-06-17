-- name: InsertProcessed :exec
INSERT INTO processed(id, content)
VALUES ($1, $2);

-- name: IsProcessed :one
SELECT EXISTS(SELECT * FROM processed WHERE id = $1);
