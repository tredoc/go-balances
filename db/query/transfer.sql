-- name: GetTransferByID :one
SELECT * FROM transfers
WHERE id = ?;

-- name: GetAllTransfers :many
SELECT * FROM transfers;

-- name: GetTransfersByAccountID :many
SELECT * FROM transfers
WHERE from_balance_id = ? OR to_balance_id = ?;

-- name: GetTransfersByInAndOutAccountIDs :many
SELECT * FROM transfers
WHERE from_balance_id = ? AND to_balance_id = ?;

-- name: CreateTransfer :execlastid
INSERT INTO transfers (from_balance_id, to_balance_id, amount)
VALUES (?, ?, ?);