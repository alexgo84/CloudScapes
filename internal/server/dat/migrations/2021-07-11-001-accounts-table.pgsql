CREATE TABLE IF NOT EXISTS accounts(
   id serial PRIMARY KEY,
   created_at TIMESTAMPTZ DEFAULT NOW()
);