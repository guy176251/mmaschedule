CREATE TABLE db_events (
  url TEXT NOT NULL UNIQUE,
  data BLOB,
  history BLOB
);
