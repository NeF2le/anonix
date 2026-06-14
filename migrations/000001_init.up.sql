CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE SCHEMA IF NOT EXISTS mapping;

CREATE TABLE IF NOT EXISTS mapping.mappings
(
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    cipher_text BYTEA NOT NULL,
    dek_wrapped BYTEA NOT NULL,
    kind_id INT DEFAULT NULL,
    deterministic BOOLEAN NOT NULL,
    token_ttl BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_mappings_token_ttl ON mapping.mappings(token_ttl);

CREATE TABLE IF NOT EXISTS mapping.kinds
(
    id INT NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    russian_name VARCHAR(100) NOT NULL UNIQUE,
    access_level INT NOT NULL
);

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 1, 'name', 'ФИО', 1
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'name');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 2, 'phone', 'Телефон', 2
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'phone');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 3, 'passport', 'Паспорт', 3
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'passport');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 4, 'email', 'Email', 2
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'email');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 5, 'address', 'Адрес', 2
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'address');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 6, 'birth_date', 'Дата рождения', 2
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'birth_date');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 7, 'snils', 'СНИЛС', 3
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'snils');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 8, 'inn', 'ИНН', 3
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'inn');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 9, 'driver_license', 'Водительское удостоверение', 3
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'driver_license');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 10, 'bank_card', 'Банковская карта', 4
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'bank_card');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 11, 'account_number', 'Номер счета', 4
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'account_number');

INSERT INTO mapping.kinds (id, name, russian_name, access_level)
SELECT 12, 'ip_address', 'IP-адрес', 1
WHERE NOT EXISTS (SELECT 1 FROM mapping.kinds WHERE name = 'ip_address');

CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.roles (
   id INT NOT NULL PRIMARY KEY,
   name VARCHAR(20) NOT NULL UNIQUE
);

INSERT INTO auth.roles (id, name)
SELECT 1, 'admin'
WHERE NOT EXISTS (SELECT 1 FROM auth.roles WHERE name = 'admin');
INSERT INTO auth.roles (id, name)
SELECT 2, 'auditor'
WHERE NOT EXISTS (SELECT 1 FROM auth.roles WHERE name = 'auditor');
INSERT INTO auth.roles (id, name)
SELECT 3, 'specialist'
WHERE NOT EXISTS (SELECT 1 FROM auth.roles WHERE name = 'specialist');

CREATE TABLE IF NOT EXISTS auth.users (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    login VARCHAR(30) NOT NULL UNIQUE ,
    password_hash VARCHAR(255) NOT NULL ,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_password_hash ON auth.users(password_hash);

CREATE TABLE IF NOT EXISTS auth.users_roles (
      id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
      user_id uuid,
      role_id INT
);

ALTER TABLE auth.users_roles ADD CONSTRAINT fk_users_roles_role_id
    FOREIGN KEY (role_id) REFERENCES auth.roles(id) ON DELETE RESTRICT;

ALTER TABLE auth.users_roles ADD CONSTRAINT fk_users_roles_user_id
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX IF NOT EXISTS uniq_user_role
    ON auth.users_roles(user_id, role_id);