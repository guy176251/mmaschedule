-- +goose Up
-- +goose StatementBegin
--PRAGMA journal_mode = WAL;

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

CREATE TABLE tapology (
  name TEXT NOT NULL UNIQUE,
  url TEXT NOT NULL UNIQUE
);

-- +goose StatementEnd
