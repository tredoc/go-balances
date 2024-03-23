-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = ?;

-- name: GetAllTransfers :many
SELECT * FROM transfers;