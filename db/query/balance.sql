-- name: GetBalanceByID :one
SELECT * FROM balances
WHERE id = ?;

-- name: GetAllBalances :many
SELECT * FROM balances;

-- name: GetBalancesByUserID :many
SELECT * FROM balances
WHERE user_id = ?;

-- name: UpdateBalance :exec
UPDATE balances
SET amount = ?
WHERE id = ?;

-- name: GetBalanceByIDForUpdate :one
SELECT * FROM balances
WHERE id = ?
FOR UPDATE;