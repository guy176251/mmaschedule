// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package event

import (
	"context"
)

const getEvent = `-- name: GetEvent :one
SELECT
  url, slug, name, location, organization, image, date, fights, history
FROM
  event
WHERE
  slug = ?
LIMIT
  1
`

func (q *Queries) GetEvent(ctx context.Context, slug string) (Event, error) {
	row := q.db.QueryRowContext(ctx, getEvent, slug)
	var i Event
	err := row.Scan(
		&i.Url,
		&i.Slug,
		&i.Name,
		&i.Location,
		&i.Organization,
		&i.Image,
		&i.Date,
		&i.Fights,
		&i.History,
	)
	return i, err
}

const getTapology = `-- name: GetTapology :one
SELECT
  name, url
FROM
  tapology
WHERE
  name = ?
LIMIT
  1
`

func (q *Queries) GetTapology(ctx context.Context, name string) (Tapology, error) {
	row := q.db.QueryRowContext(ctx, getTapology, name)
	var i Tapology
	err := row.Scan(&i.Name, &i.Url)
	return i, err
}
