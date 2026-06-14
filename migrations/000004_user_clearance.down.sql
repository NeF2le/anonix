ALTER TABLE auth.users DROP CONSTRAINT IF EXISTS chk_users_clearance_level;

ALTER TABLE auth.users DROP COLUMN IF EXISTS clearance_level;
