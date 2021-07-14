ALTER TABLE deployments DROP CONSTRAINT IF EXISTS deployments_accountid_name_key;
DROP INDEX IF EXISTS deployments_accountid_name_key;

CREATE UNIQUE INDEX deployments_accountid_name_deleted_is_null_key ON deployments (name, accountid)
    WHERE deleted_at IS NULL;
