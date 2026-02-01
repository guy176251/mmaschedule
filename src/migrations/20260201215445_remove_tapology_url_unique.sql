-- +goose Up
-- +goose StatementBegin
ALTER TABLE tapology
RENAME TO old_tapology;

CREATE TABLE tapology (name string NOT NULL UNIQUE, url string NOT NULL);

INSERT INTO
  tapology
SELECT
  *
FROM
  old_tapology;
-- +goose StatementEnd
