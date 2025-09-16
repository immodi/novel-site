-- name: GetChapterCommentById :one
SELECT * FROM chapter_comments WHERE id = ?;

-- name: CreateChapterComment :one
INSERT INTO chapter_comments (
    chapter_id, user_id, parent_id, content, last_updated
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateChapterComment :one
UPDATE chapter_comments
SET 
    content = ?, 
    last_updated = ?
WHERE id = ?
RETURNING *;


-- name: GetCommentsByChapter :many
SELECT c.*,
       u.username,
       u.image AS picture_url,
       (SELECT COUNT(*) FROM chapter_comment_reactions r WHERE r.comment_id = c.id AND r.reaction = 'like') AS likes,
       (SELECT COUNT(*) FROM chapter_comment_reactions r WHERE r.comment_id = c.id AND r.reaction = 'dislike') AS dislikes
FROM chapter_comments c
JOIN users u ON u.id = c.user_id
WHERE c.chapter_id = ?
ORDER BY c.last_updated ASC;

-- name: GetRepliesByChapterComment :many
SELECT c.*,
       u.username,
       u.image AS picture_url,
       (SELECT COUNT(*) FROM chapter_comment_reactions r WHERE r.comment_id = c.id AND r.reaction = 'like') AS likes,
       (SELECT COUNT(*) FROM chapter_comment_reactions r WHERE r.comment_id = c.id AND r.reaction = 'dislike') AS dislikes
FROM chapter_comments c
JOIN users u ON u.id = c.user_id
WHERE c.parent_id = ?
ORDER BY c.last_updated ASC;


-- name: GetCommentsByChapterWithUserReactions :many
SELECT
    c.*,
    u.username,
    u.image AS picture_url,
    (SELECT COUNT(*) FROM chapter_comment_reactions
       WHERE comment_id = c.id AND reaction = 'like') AS likes,
    (SELECT COUNT(*) FROM chapter_comment_reactions
       WHERE comment_id = c.id AND reaction = 'dislike') AS dislikes,
    ur.reaction AS user_reaction
FROM chapter_comments c
JOIN users u ON u.id = c.user_id
LEFT JOIN chapter_comment_reactions ur
       ON ur.comment_id = c.id
      AND ur.user_id = ?
WHERE c.chapter_id = ?
ORDER BY c.last_updated ASC;

-- name: GetRepliesByChapterCommentWithUserReactions :many
SELECT
    c.*,
    u.username,
    u.image AS picture_url,
    (SELECT COUNT(*) FROM chapter_comment_reactions
       WHERE comment_id = c.id AND reaction = 'like') AS likes,
    (SELECT COUNT(*) FROM chapter_comment_reactions
       WHERE comment_id = c.id AND reaction = 'dislike') AS dislikes,
    ur.reaction AS user_reaction
FROM chapter_comments c
JOIN users u ON u.id = c.user_id
LEFT JOIN chapter_comment_reactions ur
       ON ur.comment_id = c.id
      AND ur.user_id = ?
WHERE c.parent_id = ?
ORDER BY c.last_updated ASC;

-- name: DeleteChapterComment :exec
DELETE FROM chapter_comments WHERE id = ? AND user_id = ?;

---------------------------------------
-- USER REACTIONS
---------------------------------------

-- name: DeleteChapterReactionIfSame :one
DELETE FROM chapter_comment_reactions
WHERE user_id = ? AND comment_id = ? AND reaction = ?
RETURNING 1;

-- name: UpsertChapterReaction :exec
INSERT INTO chapter_comment_reactions (user_id, comment_id, reaction, last_updated)
VALUES (?, ?, ?, ?)
ON CONFLICT(user_id, comment_id) DO UPDATE
SET reaction = excluded.reaction,
    last_updated = excluded.last_updated;

-- name: RemoveChapterReaction :exec
DELETE FROM chapter_comment_reactions
WHERE user_id = ? AND comment_id = ?;

-- name: GetUserChapterReaction :one
SELECT reaction
FROM chapter_comment_reactions
WHERE user_id = ? AND comment_id = ?;

-- name: GetUserReactionsForChapterComments :many
SELECT 
    comment_id,
    reaction
FROM chapter_comment_reactions
WHERE user_id = ? AND comment_id IN (sqlc.slice('comment_ids'));

-- name: CountChapterReactions :one
SELECT
    SUM(CASE WHEN reaction = 'like' THEN 1 ELSE 0 END) AS likes,
    SUM(CASE WHEN reaction = 'dislike' THEN 1 ELSE 0 END) AS dislikes
FROM chapter_comment_reactions
WHERE comment_id = ?;

