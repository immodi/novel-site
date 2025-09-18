-- name: AddTagToNovel :exec
INSERT INTO novel_tags (novel_id, tag, tag_slug)
VALUES (?, ?, ?)
ON CONFLICT(novel_id, tag) DO NOTHING;

-- name: ListAllTags :many
SELECT DISTINCT tag
FROM novel_tags
ORDER BY tag COLLATE NOCASE;

-- name: ListAllTagSlugs :many
SELECT DISTINCT tag_slug
FROM novel_tags
ORDER BY tag COLLATE NOCASE;

-- name: RemoveTagFromNovel :exec
DELETE FROM novel_tags
WHERE novel_id = ? AND tag = ?;

-- name: ListTagsByNovel :many
SELECT novel_id, tag, tag_slug
FROM novel_tags
WHERE novel_id = ?
ORDER BY tag ASC;

-- name: GetTagBySlug :one
SELECT novel_id, tag, tag_slug
FROM novel_tags
WHERE tag_slug = ?
LIMIT 1;

-- name: ListNovelsByTag :many
SELECT n.*
FROM novels n
JOIN novel_tags t ON n.id = t.novel_id
WHERE t.tag_slug = ?
ORDER BY n.update_time DESC;

-- name: DeleteAllTagsByNovel :exec
DELETE FROM novel_tags
WHERE novel_id = ?;

-- name: CountNovelsByTag :one
SELECT COUNT(*)
FROM novel_tags
WHERE tag_slug = ?;

-- name: ListNovelsByTagPaginated :many
SELECT n.*
FROM novels n
JOIN novel_tags t ON n.id = t.novel_id
WHERE t.tag_slug = ?
ORDER BY n.update_time DESC
LIMIT ? OFFSET ?;
