-- name: ListAllGenres :many
SELECT DISTINCT genre
FROM novel_genres
ORDER BY genre;

-- name: AddGenreToNovel :exec
INSERT INTO novel_genres (novel_id, genre, genre_slug)
VALUES (?, ?, ?)
ON CONFLICT(novel_id, genre) DO NOTHING;

-- name: ListGenresByNovel :many
SELECT novel_id, genre, genre_slug
FROM novel_genres
WHERE novel_id = ?
ORDER BY genre;

-- name: GetGenreBySlug :one
SELECT novel_id, genre, genre_slug
FROM novel_genres
WHERE genre_slug = ?
LIMIT 1;

-- name: DeleteGenreFromNovel :exec
DELETE FROM novel_genres
WHERE novel_id = ? AND genre = ?;

-- name: CountNovelsByGenre :one
SELECT COUNT(*) 
FROM novel_genres
WHERE genre_slug = ?;

-- name: ListNovelsByGenrePaginated :many
SELECT n.*
FROM novels n
JOIN novel_genres ng ON n.id = ng.novel_id
WHERE ng.genre_slug = ?
ORDER BY n.update_time DESC
LIMIT ? OFFSET ?;
