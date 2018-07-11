CREATE EXTENSION IF NOT EXIST CITEXT;

CREATE TABLE users (
  id        SERIAL  PRIMARY KEY,
  username  CITEXT  NOT NULL UNIQUE,
  email     CITEXT  NOT NULL UNIQUE,
  password  TEXT    NOT NULL,

  created_at  TIMESTAMP without time zone  NOT NULL DEFAULT (now() at time zone 'utc'),
  updated_at  TIMESTAMP without time zone  NOT NULL DEFAULT (now() at time zone 'utc'),
  deleted_at  TIMESTAMP without time zone
);
