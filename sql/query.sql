-- name: GetEvent :one
SELECT
  *
FROM
  event
WHERE
  url = ?
LIMIT
  1;
