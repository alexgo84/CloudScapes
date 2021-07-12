CREATE TABLE IF NOT EXISTS plans (
    id serial PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    modified_at TIMESTAMPTZ,
    created_by bigint NOT NULL,
    modified_by bigint,
    name text NOT NULL,

    replicas int,

    cpu_limit text,
    mem_limit text,
    cpu_req text,
    mem_req text,

    database_service_name text,
    database_service_cloud text,
    database_service_plan text,
    database_connections int,

    env_vars jsonb NOT NULL default '{}'::jsonb,
    cron_jobs jsonb NOT NULL default '[]'::jsonb,
    config_maps jsonb NOT NULL default '[]'::jsonb
);