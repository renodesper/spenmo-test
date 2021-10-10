CREATE TABLE IF NOT EXISTS "wallet" (
  id uuid DEFAULT uuid_generate_v4(),
  balance NUMERIC NOT NULL DEFAULT 0,
  daily_limit NUMERIC NOT NULL DEFAULT 0,
  monthly_limit NUMERIC NOT NULL DEFAULT 0,
  team_id uuid,
  user_id uuid,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
  updated_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (id)
);
