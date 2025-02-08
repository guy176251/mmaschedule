-- +goose Up
-- +goose StatementBegin
CREATE TABLE event (
  url TEXT NOT NULL UNIQUE,
  slug TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  location TEXT NOT NULL,
  organization TEXT NOT NULL,
  image TEXT NOT NULL,
  date INTEGER NOT NULL,
  fights TEXT NOT NULL,
  history TEXT NOT NULL
);
-- +goose StatementEnd
