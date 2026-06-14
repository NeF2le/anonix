CREATE TABLE IF NOT EXISTS mapping.audit_log
(
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    action VARCHAR(20) NOT NULL,
    token VARCHAR(100) NOT NULL,
    kind_id INT DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON mapping.audit_log(created_at DESC);
