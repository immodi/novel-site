-- name: ListAllGenres :many
SELECT genre FROM genres
ORDER BY genre;

-- name: AddGenreToNovel :exec
INSERT INTO novel_genres (novel_id, genre)
VALUES (?, ?)
ON CONFLICT DO NOTHING;

-- name: ListGenresByNovel :many
SELECT genre FROM novel_genres
WHERE novel_id = ?;

-- name: DeleteGenreFromNovel :exec
DELETE FROM novel_genres
WHERE novel_id = ? AND genre = ?;

-- name: CountNovelsByGenre :one
SELECT COUNT(*) 
FROM novel_genres
WHERE genre = ?;

-- name: ListNovelsByGenrePaginated :many
SELECT n.*
FROM novels n
JOIN novel_genres ng ON n.id = ng.novel_id
WHERE ng.genre = ?
ORDER BY n.update_time DESC
LIMIT ? OFFSET ?;

