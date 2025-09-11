-- name: GetCommentById :one
SELECT * FROM comments WHERE id = ?;

-- name: CreateComment :one
INSERT INTO comments (
    novel_id, user_id, parent_id, content, created_at
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateComment :one
UPDATE comments
SET 
    content = ?, 
    created_at = ?
WHERE id = ?
RETURNING *;


-- name: GetCommentsByNovel :many
SELECT c.*,
       u.username,
       u.image AS picture_url,
       (SELECT COUNT(*) FROM comment_reactions r WHERE r.comment_id = c.id AND r.reaction = 'like') AS likes,
       (SELECT COUNT(*) FROM comment_reactions r WHERE r.comment_id = c.id AND r.reaction = 'dislike') AS dislikes
FROM comments c
JOIN users u ON u.id = c.user_id
WHERE c.novel_id = ?
ORDER BY c.created_at ASC;

-- name: GetRepliesByComment :many
SELECT c.*,
       u.username,
       u.image AS picture_url,
       (SELECT COUNT(*) FROM comment_reactions r WHERE r.comment_id = c.id AND r.reaction = 'like') AS likes,
       (SELECT COUNT(*) FROM comment_reactions r WHERE r.comment_id = c.id AND r.reaction = 'dislike') AS dislikes
FROM comments c
JOIN users u ON u.id = c.user_id
WHERE c.parent_id = ?
ORDER BY c.created_at ASC;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = ? AND user_id = ?;

---------------------------------------
-- USER REACTIONS
---------------------------------------

-- name: UpsertReaction :exec
INSERT INTO comment_reactions (user_id, comment_id, reaction, created_at)
VALUES (?, ?, ?, datetime('now'))
ON CONFLICT(user_id, comment_id) DO UPDATE
SET reaction = excluded.reaction,
    created_at = datetime('now');
-- RETURNING *;

-- name: RemoveReaction :exec
DELETE FROM comment_reactions
WHERE user_id = ? AND comment_id = ?;

-- name: GetUserReaction :one
SELECT reaction
FROM comment_reactions
WHERE user_id = ? AND comment_id = ?;

-- name: GetUserReactionsForComments :many
SELECT 
    comment_id,
    reaction
FROM comment_reactions
WHERE user_id = ? AND comment_id IN (sqlc.slice('comment_ids'));

-- name: CountReactions :one
SELECT
    SUM(CASE WHEN reaction = 'like' THEN 1 ELSE 0 END) AS likes,
    SUM(CASE WHEN reaction = 'dislike' THEN 1 ELSE 0 END) AS dislikes
FROM comment_reactions
WHERE comment_id = ?;
