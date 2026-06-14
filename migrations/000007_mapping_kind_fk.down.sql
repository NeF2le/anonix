ALTER TABLE mapping.audit_log DROP CONSTRAINT IF EXISTS fk_audit_log_kind_id;
ALTER TABLE mapping.mappings DROP CONSTRAINT IF EXISTS fk_mappings_kind_id;
