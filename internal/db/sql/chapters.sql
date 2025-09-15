-- name: CreateChapter :one
INSERT INTO chapters (
    novel_id, chapter_number, title, content, release_date
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetChapterByID :one
SELECT * FROM chapters
WHERE id = ? LIMIT 1;

-- name: GetChapterByNumber :one
SELECT * FROM chapters
WHERE novel_id = ? AND chapter_number = ?
LIMIT 1;

-- name: ListChaptersByNovel :many
SELECT * FROM chapters
WHERE novel_id = ?
ORDER BY chapter_number ASC;

-- name: DeleteChapter :exec
DELETE FROM chapters WHERE id = ?;

-- name: ListChaptersByNovelPaginated :many
SELECT * FROM chapters
WHERE novel_id = ?
ORDER BY chapter_number ASC
LIMIT ? OFFSET ?;

-- name: CountChaptersByNovel :one
SELECT COUNT(*) FROM chapters
WHERE novel_id = ?;

-- name: GetNextChapterNumber :one
SELECT COALESCE(MAX(chapter_number), 0) + 1 as next_number
FROM chapters
WHERE novel_id = ?;

-- name: UpdateChapterNumber :one
UPDATE chapters
SET chapter_number = ?
WHERE id = ?
RETURNING *;

-- name: GetLatestChapterByNovel :one
SELECT *
FROM chapters
WHERE novel_id = ?
ORDER BY chapter_number DESC
LIMIT 1;
