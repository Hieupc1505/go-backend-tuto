-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    email,
    hashed_password
) VALUES (
    $1, $2
) RETURNING *;

-- name: ChangePassword :one
UPDATE users
SET hashed_password = $2
WHERE email = $1
RETURNING *;
