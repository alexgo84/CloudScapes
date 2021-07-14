CREATE TABLE IF NOT EXISTS sessions(
    id serial PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    token text NOT NULL,
    userid bigint NOT NULL,
    CONSTRAINT fk_sessions_userid
        FOREIGN KEY(userid)
	        REFERENCES users(id)
);