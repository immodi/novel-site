-- name: CreateNovel :one
INSERT INTO novels (
    title, description, cover_image, author, status, update_time, view_count
) VALUES (
    ?, ?, ?, ?, ?, ?, 0
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
    update_time = ?,
    view_count = ?
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
LIMIT 6;

-- name: ListHotNovels :many
SELECT * FROM novels
ORDER BY view_count DESC, update_time DESC
LIMIT 6;

-- name: GetNovelTags :many
SELECT tag
FROM novel_tags
WHERE novel_id = ?;

-- name: IncrementNovelViewCount :exec
UPDATE novels
SET view_count = view_count + 1,
    update_time = update_time
WHERE id = ?;

