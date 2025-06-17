-- name: InsertTest :one
INSERT INTO test(ID, CONTENT)
VALUES ($1, $2)
RETURNING *;


-- name: InsertEvent :exec
INSERT INTO outbox(ID, EVENT_TYPE, AGGREGATE_TYPE, AGGREGATE_ID, PAYLOAD)
VALUES($1, $2, $3, $4, $5);


-- name: GetUnpublishedEvents :many
SELECT * FROM outbox
WHERE published IS FALSE
ORDER BY created_at
LIMIT 100;


-- name: SetEventAsPublished :exec
UPDATE outbox
SET published = TRUE, published_at = NOW()
WHERE id = $1;

