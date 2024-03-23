-- name: GetEntry :one
SELECT * FROM entries
WHERE id = ?;

-- name: GetAllEntries :many
SELECT * FROM entries;