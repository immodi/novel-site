-- name: AddUserBookmark :exec
INSERT INTO user_bookmarks (user_id, novel_id, created_at)
VALUES (?, ?, ?)
ON CONFLICT(user_id, novel_id) DO NOTHING;

-- name: RemoveUserBookmark :exec
DELETE FROM user_bookmarks
WHERE user_id = ? AND novel_id = ?;

-- name: ListUserBookmarksPaginated :many
SELECT n.*
FROM user_bookmarks ub
JOIN novels n ON ub.novel_id = n.id
WHERE ub.user_id = ?
ORDER BY ub.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountUserBookmarks :one
SELECT COUNT(*)
FROM user_bookmarks
WHERE user_id = ?;

-- name: IsNovelBookmarked :one
SELECT EXISTS (
    SELECT 1
    FROM user_bookmarks
    WHERE user_id = ? AND novel_id = ?
) AS is_bookmarked;

-- name: GetLastReadChapterID :one
SELECT last_read_chapter_id
FROM user_bookmarks
WHERE user_id = ? AND novel_id = ?;

-- name: UpdateLastReadChapter :exec
UPDATE user_bookmarks
SET last_read_chapter_id = ?
WHERE user_id = ? AND novel_id = ?;
