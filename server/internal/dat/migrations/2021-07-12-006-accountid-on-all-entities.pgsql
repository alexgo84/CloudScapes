ALTER TABLE plans 
    ADD COLUMN IF NOT EXISTS accountid bigint NOT NULL;
ALTER TABLE plans 
    ADD CONSTRAINT fk_plans_accountid
        FOREIGN KEY(accountid)
	        REFERENCES accounts(id);
ALTER TABLE plans
ADD UNIQUE (accountid, name);

ALTER TABLE deployments 
    ADD COLUMN IF NOT EXISTS accountid bigint NOT NULL;
ALTER TABLE deployments 
    ADD CONSTRAINT fk_deployments_accountid
        FOREIGN KEY(accountid)
	        REFERENCES accounts(id)
               ON DELETE CASCADE;
ALTER TABLE deployments
ADD UNIQUE (accountid, name);

ALTER TABLE clusters 
    ADD COLUMN IF NOT EXISTS accountid bigint NOT NULL;
ALTER TABLE clusters 
    ADD CONSTRAINT fk_clusters_accountid
        FOREIGN KEY(accountid)
	        REFERENCES accounts(id)
               ON DELETE CASCADE;
ALTER TABLE clusters
ADD UNIQUE (accountid, name);
