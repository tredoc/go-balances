-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: GetAllUsers :many
SELECT * FROM users;