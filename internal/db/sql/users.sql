-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, created_at, role)
VALUES (?, ?, ?, ?, ?)
RETURNING id;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?;

-- name: GetAllUsers :many
SELECT *
FROM users;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = ?;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE LOWER(email) = LOWER(?);

-- name: UpdateUserImage :exec
UPDATE users
SET image = ?
WHERE id = ?;

-- name: UpdateUserPartial :one
UPDATE users
SET
    username = CASE WHEN :username IS NOT NULL THEN :username ELSE username END,
    email = CASE WHEN :email IS NOT NULL THEN :email ELSE email END,
    password_hash = CASE WHEN :password_hash IS NOT NULL THEN :password_hash ELSE password_hash END,
    role = CASE WHEN :role IS NOT NULL THEN :role ELSE role END,
    image = CASE WHEN :image IS NOT NULL THEN :image ELSE image END,
    created_at = CASE WHEN :created_at IS NOT NULL THEN :created_at ELSE created_at END
WHERE id = :id
RETURNING *;
