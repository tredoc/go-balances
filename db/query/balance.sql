-- name: GetBalance :one
SELECT * FROM balances
WHERE id = ?;

-- name: GetAllBalances :many
SELECT * FROM balances;