CREATE TABLE IF NOT EXISTS deployments (

    -- generic
    id serial PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by bigint NOT NULL,
    modified_at TIMESTAMPTZ,
    modified_by bigint,
    deleted_at TIMESTAMPTZ,
    deleted_by bigint,

    -- specific
    name text NOT NULL,

    -- deployment
    image_path text,
    image_sha text,
    exlude_from_updates bool,
    planid bigint,
    salesforce_state text,
    replicas int,
    clusterid bigint NOT NULL,

    -- resources
    cpu_limit text,
    mem_limit text,
    cpu_req text,
    mem_req text,

    -- database
    database_service_name text,
    database_service_cloud text,
    database_service_plan text,

    -- high level
    env_vars jsonb NOT NULL default '{}'::jsonb,
    cron_jobs jsonb NOT NULL default '[]'::jsonb,
    config_maps jsonb NOT NULL default '[]'::jsonb,

    CONSTRAINT fk_clients_clusterid
        FOREIGN KEY(clusterid)
	        REFERENCES clusters(id),
    CONSTRAINT fk_clients_planid
        FOREIGN KEY(planid)
	        REFERENCES plans(id)
);