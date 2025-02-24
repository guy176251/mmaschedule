-- name: GetEvent :one
SELECT
  *
FROM
  event
WHERE
  slug = ?
LIMIT
  1;

-- name: ListEvents :many
SELECT
  *
FROM
  event;

-- name: GetUpcomingEvent :one
SELECT
  slug
FROM
  event
WHERE
  date >= ?
ORDER BY
  date ASC
LIMIT
  1;

-- name: ListUpcomingEvents :many
SELECT
  name,
  slug,
  date
FROM
  event
WHERE
  date >= ?
ORDER BY
  date ASC;

-- name: GetTapology :one
SELECT
  *
FROM
  tapology
WHERE
  name = ?
LIMIT
  1;

-- name: CreateTapology :exec
INSERT INTO
  tapology (name, url)
VALUES
  (?, ?);
