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
