-- name: CreateChapter :one
INSERT INTO chapters (
    novel_id, title, content
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: GetChapterByID :one
SELECT * FROM chapters
WHERE id = ? LIMIT 1;

-- name: ListChaptersByNovel :many
SELECT * FROM chapters
WHERE novel_id = ?
ORDER BY id ASC;

-- name: DeleteChapter :exec
DELETE FROM chapters WHERE id = ?;

-- name: ListChaptersByNovelPaginated :many
SELECT * FROM chapters
WHERE novel_id = ?
ORDER BY id ASC
LIMIT ? OFFSET ?;

-- name: CountChaptersByNovel :one
SELECT COUNT(*) FROM chapters
WHERE novel_id = ?;
