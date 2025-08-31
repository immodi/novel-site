-- name: CreateNovel :one
INSERT INTO novels (
    title, description, cover_image, author, status, update_time
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetNovelByID :one
SELECT * FROM novels
WHERE id = ? LIMIT 1;

-- name: ListNovels :many
SELECT * FROM novels
ORDER BY update_time DESC;

-- name: UpdateNovel :one
UPDATE novels
SET
    title = ?,
    description = ?,
    cover_image = ?,
    author = ?,
    status = ?,
    update_time = ?
WHERE id = ?
RETURNING *;

-- name: DeleteNovel :exec
DELETE FROM novels WHERE id = ?;

-- name: GetNovelByNameLike :one
SELECT * FROM novels
WHERE LOWER(title) LIKE LOWER(?)
LIMIT 1;

-- name: ListNewestHomeNovels :many
SELECT * FROM novels
ORDER BY update_time DESC
LIMIT 8;

