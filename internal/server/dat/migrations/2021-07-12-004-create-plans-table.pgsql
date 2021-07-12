CREATE TABLE IF NOT EXISTS plans (
    -- generic
    id serial PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    name text NOT NULL,
    replicas int NOT NULL,
    clusterid bigint NOT NULL,

    cpu_limit text NOT NULL,
    mem_limit text NOT NULL,
    cpu_req text NOT NULL,
    mem_req text NOT NULL,

    database_service_name text NOT NULL,
    database_service_cloud text NOT NULL,
    database_service_plan text NOT NULL,

    env_vars jsonb NOT NULL default '{}'::jsonb,
    cron_jobs jsonb NOT NULL default '[]'::jsonb,
    config_maps jsonb NOT NULL default '[]'::jsonb,

    CONSTRAINT fk_plans_clusterid
        FOREIGN KEY(clusterid)
	        REFERENCES clusters(id)
);