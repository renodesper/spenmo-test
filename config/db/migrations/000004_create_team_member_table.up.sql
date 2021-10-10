CREATE TABLE IF NOT EXISTS "team_member" (
  id uuid DEFAULT uuid_generate_v4(),
  team_id uuid NOT NULL,
  user_id uuid NOT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
  updated_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (id)
);
