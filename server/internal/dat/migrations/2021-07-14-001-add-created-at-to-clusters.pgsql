ALTER TABLE clusters ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();