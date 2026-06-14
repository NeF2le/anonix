ALTER TABLE auth.users_roles DROP CONSTRAINT IF EXISTS fk_users_roles_user_id;
ALTER TABLE auth.users_roles DROP CONSTRAINT IF EXISTS fk_users_roles_role_id;

DROP TABLE IF EXISTS auth.users_roles;
DROP TABLE IF EXISTS auth.users;
DROP TABLE IF EXISTS auth.roles;

DROP TABLE IF EXISTS mapping.mappings;
DROP TABLE IF EXISTS mapping.kinds;

DROP SCHEMA IF EXISTS auth CASCADE;
DROP SCHEMA IF EXISTS mapping CASCADE;