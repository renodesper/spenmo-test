CREATE TABLE IF NOT EXISTS "card" (
  id uuid DEFAULT uuid_generate_v4(),
  card_no VARCHAR NOT NULL,
  expiry_month VARCHAR NOT NULL,
  expiry_year VARCHAR NOT NULL,
  cvv VARCHAR NOT NULL,
  daily_limit NUMERIC NOT NULL DEFAULT 0,
  monthly_limit NUMERIC NOT NULL DEFAULT 0,
  wallet_id uuid NOT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
  updated_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (id)
);
