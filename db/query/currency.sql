-- name: GetCurrency :one
SELECT * FROM currencies
WHERE id = ?;

-- name: GetAllCurrencies :many
SELECT * FROM currencies;