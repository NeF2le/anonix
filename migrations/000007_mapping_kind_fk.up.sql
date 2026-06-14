ALTER TABLE mapping.mappings
    ADD CONSTRAINT fk_mappings_kind_id
    FOREIGN KEY (kind_id) REFERENCES mapping.kinds(id) ON DELETE RESTRICT;

ALTER TABLE mapping.audit_log
    ADD CONSTRAINT fk_audit_log_kind_id
    FOREIGN KEY (kind_id) REFERENCES mapping.kinds(id) ON DELETE RESTRICT;
