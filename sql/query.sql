-- name: GetEvent :one
SELECT
  *
FROM
  event
WHERE
  slug = ?
LIMIT
  1;

-- name: GetTapology :one
SELECT
  *
FROM
  tapology
WHERE
  name = ?
LIMIT
  1;
