-- name: CreateNovel :one
INSERT INTO novels (
    title, description, cover_image, author, publisher, release_year, is_completed, update_time, view_count
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, 0
)
RETURNING *;

-- name: GetNovelByID :one
SELECT * FROM novels
WHERE id = ? LIMIT 1;

-- name: ListNovels :many
SELECT * FROM novels
ORDER BY update_time DESC;

-- name: CountNovels :one
SELECT COUNT(*) FROM novels;

-- name: CountNovelsByAuthor :one
SELECT COUNT(*)
FROM novels
WHERE author = ?;

-- name: ListNovelsByAuthorPaginated :many
SELECT *
FROM novels
WHERE author = ?
ORDER BY update_time DESC
LIMIT ? OFFSET ?;

-- name: UpdateNovel :one
UPDATE novels
SET
    title = ?,
    description = ?,
    cover_image = ?,
    author = ?,
    publisher = ?,
    release_year = ?,
    is_completed = ?,
    update_time = ?,
    view_count = ?
WHERE id = ?
RETURNING *;

-- name: UpdateNovelPartial :one
UPDATE novels
SET
    title = CASE WHEN :title IS NOT NULL THEN :title ELSE title END,
    description = CASE WHEN :description IS NOT NULL THEN :description ELSE description END,
    cover_image = CASE WHEN :cover_image IS NOT NULL THEN :cover_image ELSE cover_image END,
    author = CASE WHEN :author IS NOT NULL THEN :author ELSE author END,
    publisher = CASE WHEN :publisher IS NOT NULL THEN :publisher ELSE publisher END,
    release_year = CASE WHEN :release_year IS NOT NULL THEN :release_year ELSE release_year END,
    is_completed = CASE WHEN :is_completed IS NOT NULL THEN :is_completed ELSE is_completed END,
    update_time = CASE WHEN :update_time IS NOT NULL THEN :update_time ELSE update_time END,
    view_count = CASE WHEN :view_count IS NOT NULL THEN :view_count ELSE view_count END
WHERE id = :id
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

-- name: ListCompletedNovels :many
SELECT *
FROM novels
WHERE is_completed = 1
ORDER BY update_time DESC;

-- name: ListNewestHomeNovelsPaginated :many
SELECT *
FROM novels
ORDER BY update_time DESC
LIMIT ? OFFSET ?;

-- name: ListHotNovelsPaginated :many
SELECT *
FROM novels
ORDER BY view_count DESC, update_time DESC
LIMIT ? OFFSET ?;

-- name: ListCompletedNovelsPaginated :many
SELECT *
FROM novels
WHERE is_completed = 1
ORDER BY update_time DESC
LIMIT ? OFFSET ?;

-- name: ListOnGoingNovelsPaginated :many
SELECT *
FROM novels
WHERE is_completed = 0
ORDER BY update_time DESC
LIMIT ? OFFSET ?;

-- name: CountCompletedNovels :one
SELECT COUNT(*) 
FROM novels
WHERE is_completed = 1;

-- name: CountOnGoingNovels :one
SELECT COUNT(*) 
FROM novels
WHERE is_completed = 0;

-- name: GetNovelTags :many
SELECT tag
FROM novel_tags
WHERE novel_id = ?;

-- name: IncrementNovelViewCount :exec
UPDATE novels
SET view_count = view_count + 1,
    update_time = update_time
WHERE id = ?;

-- name: SearchNovels :many
SELECT *
FROM novels
WHERE LOWER(title) LIKE '%' || LOWER(sqlc.arg(search)) || '%'
   OR LOWER(author) LIKE '%' || LOWER(sqlc.arg(search)) || '%'
   OR LOWER(description) LIKE '%' || LOWER(sqlc.arg(search)) || '%'
ORDER BY update_time DESC
LIMIT sqlc.arg(limit) OFFSET sqlc.arg(offset);

-- name: CountSearchNovels :one
SELECT COUNT(*) AS total
FROM novels
WHERE LOWER(title) LIKE '%' || LOWER(sqlc.arg(search)) || '%'
   OR LOWER(author) LIKE '%' || LOWER(sqlc.arg(search)) || '%'
   OR LOWER(description) LIKE '%' || LOWER(sqlc.arg(search)) || '%';
