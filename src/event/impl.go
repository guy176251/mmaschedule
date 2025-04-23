package event

import (
	"context"
	"encoding/json"
)

type Fight struct {
	Weight   string   `json:"weight"`
	FighterA *Fighter `json:"fighter_a"`
	FighterB *Fighter `json:"fighter_b"`
}

type Fighter struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Country string `json:"country"`
	Link    string `json:"link"`
}

func (e *Event) UnmarshalFights() []*Fight {
	var fights []*Fight
	_ = json.Unmarshal([]byte(e.Fights), &fights)
	return fights
}

func (e *Event) MarshalFights(fights []*Fight) {
	buf, _ := json.Marshal(&fights)
	e.Fights = string(buf)
}

func (e *Event) HasEmptyFights() bool {
	fights := e.UnmarshalFights()

	if len(fights) == 0 {
		return true
	}

	for _, f := range fights {
		if f.IsEmpty() {
			return true
		}
	}

	return false
}

func (f *Fight) IsEmpty() bool {
	return f.FighterA.IsEmpty() && f.FighterB.IsEmpty()
}

func (f *Fighter) IsEmpty() bool {
	return isEmpty(f.Country) && isEmpty(f.Image) && isEmpty(f.Link) && isEmpty(f.Name)
}

func isEmpty(s string) bool {
	return len(s) == 0
}

const upsertEvents = `-- name: UpsertEvents :exec
INSERT INTO
  event (
    url,
    slug,
    name,
    location,
    organization,
    image,
    date,
    fights,
    history
  )
SELECT
  json_extract (e.value, '$.url'),
  json_extract (e.value, '$.slug'),
  json_extract (e.value, '$.name'),
  json_extract (e.value, '$.location'),
  json_extract (e.value, '$.organization'),
  json_extract (e.value, '$.image'),
  json_extract (e.value, '$.date'),
  json_extract (e.value, '$.fights'),
  json_extract (e.value, '$.history')
FROM
  json_each (?) AS e
WHERE true
ON CONFLICT (url) DO UPDATE
SET
  slug = excluded.slug,
  name = excluded.name,
  location = excluded.location,
  organization = excluded.organization,
  image = excluded.image,
  date = excluded.date,
  fights = excluded.fights,
  history = excluded.history
`

func (q *Queries) UpsertEvents(ctx context.Context, e []*Event) error {
	b, err := json.Marshal(&e)
	if err != nil {
		return err
	}
	_, err = q.db.ExecContext(ctx, upsertEvents, string(b))
	return err
}
