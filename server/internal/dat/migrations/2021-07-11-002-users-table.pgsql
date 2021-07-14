CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    name text NOT NULL,
    email text NOT NULL,
    accountid bigint NOT NULL,
    CONSTRAINT fk_users_accountid
        FOREIGN KEY(accountid)
	        REFERENCES accounts(id)
);