-- @interface UrlRepositoryI
-- name: CreateUrl :one
INSERT INTO urls (original_url, short_code)
VALUES ($1, $2)
RETURNING *;

-- @interface UrlRepositoryI
-- name: GetUrlByCode :one
SELECT * FROM urls
WHERE short_code = $1 LIMIT 1;

-- @interface UrlRepositoryI
-- name: IncrementClicks :exec
UPDATE urls SET clicks = clicks + 1 WHERE id = $1;