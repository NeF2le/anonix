CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE SCHEMA IF NOT EXISTS mapping;

CREATE TABLE IF NOT EXISTS mapping.mappings
(
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    cipher_text BYTEA,
    dek_wrapped BYTEA,
    deterministic BOOLEAN NOT NULL DEFAULT false,
    reversible BOOLEAN NOT NULL DEFAULT false,
    token_ttl BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_mappings_token_ttl ON mapping.mappings(token_ttl);

CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.roles (
   id SERIAL NOT NULL PRIMARY KEY,
   name VARCHAR(20) NOT NULL UNIQUE
);

INSERT INTO users.roles (name)
SELECT 'admin'
WHERE NOT EXISTS (SELECT 1 FROM users.roles WHERE name = 'admin');
INSERT INTO users.roles (name)
SELECT 'default'
WHERE NOT EXISTS (SELECT 1 FROM users.roles WHERE name = 'default');

CREATE TABLE IF NOT EXISTS users.users (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    login VARCHAR(30) NOT NULL UNIQUE ,
    password_hash VARCHAR(255) NOT NULL ,
    created_at TIMESTAMP DEFAULT now(),
    role_id INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_password_hash ON users.users(password_hash);

ALTER TABLE users.users ADD CONSTRAINT fk_users_role_id
    FOREIGN KEY (role_id) REFERENCES users.roles(id) ON DELETE RESTRICT;