-- name: GetEntryByID :one
SELECT * FROM entries
WHERE id = ?;

-- name: GetAllEntries :many
SELECT * FROM entries;

-- name: GetEntriesByBalanceID :many
SELECT * FROM entries
WHERE balance_id = ?;

-- name: CreateEntry :execlastid
INSERT INTO entries (balance_id, amount)
VALUES (?, ?);

-- name: GetLastEntryID :one
SELECT id FROM entries
ORDER BY id DESC
LIMIT 1;