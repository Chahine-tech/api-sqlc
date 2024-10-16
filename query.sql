-- name: GetUsers :many
SELECT id, name, email FROM users;

-- name: CreateUser :one
INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, password = $4
WHERE id = $1
RETURNING id, name, email;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
