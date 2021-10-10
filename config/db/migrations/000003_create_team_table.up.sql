CREATE TABLE IF NOT EXISTS "team" (
  id uuid DEFAULT uuid_generate_v4(),
  name VARCHAR NOT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
  updated_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (id)
);
