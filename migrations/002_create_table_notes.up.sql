CREATE TABLE notes (
  id       SERIAL   PRIMARY KEY,
  user_id  INTEGER  NOT NULL REFERENCES users(id),

  title  TEXT,
  body   TEXT  NOT NULL,

  created_at  TIMESTAMP without time zone  NOT NULL DEFAULT (now() at time zone 'utc'),
  updated_at  TIMESTAMP without time zone  NOT NULL DEFAULT (now() at time zone 'utc'),
  deleted_at  TIMESTAMP without time zone
);

CREATE INDEX notes_user_id_idx ON notes(user_id);
