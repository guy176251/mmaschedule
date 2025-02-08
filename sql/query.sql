-- name: GetEvent :one
SELECT
  *
FROM
  db_events
WHERE
  url = ?
LIMIT
  1;

-- name: CreateEvent :exec
INSERT INTO
  db_events (url)
VALUES
  (?) ON CONFLICT (url) DO NOTHING;

-- name: UpdateEvent :exec
UPDATE db_events
set
  data = ?,
  history = ?
WHERE
  url = ?;
