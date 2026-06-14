ALTER TABLE auth.users ADD COLUMN IF NOT EXISTS clearance_level INT NOT NULL DEFAULT 1;

ALTER TABLE auth.users ADD CONSTRAINT chk_users_clearance_level CHECK (clearance_level BETWEEN 1 AND 4);
