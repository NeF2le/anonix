ALTER TABLE mapping.kinds ADD COLUMN IF NOT EXISTS short_name VARCHAR(10) NOT NULL DEFAULT '';

UPDATE mapping.kinds SET short_name = 'fio'  WHERE name = 'name';
UPDATE mapping.kinds SET short_name = 'tel'  WHERE name = 'phone';
UPDATE mapping.kinds SET short_name = 'psp'  WHERE name = 'passport';
UPDATE mapping.kinds SET short_name = 'mail' WHERE name = 'email';
UPDATE mapping.kinds SET short_name = 'adr'  WHERE name = 'address';
UPDATE mapping.kinds SET short_name = 'dob'  WHERE name = 'birth_date';
UPDATE mapping.kinds SET short_name = 'snl'  WHERE name = 'snils';
UPDATE mapping.kinds SET short_name = 'inn'  WHERE name = 'inn';
UPDATE mapping.kinds SET short_name = 'drv'  WHERE name = 'driver_license';
UPDATE mapping.kinds SET short_name = 'crd'  WHERE name = 'bank_card';
UPDATE mapping.kinds SET short_name = 'acc'  WHERE name = 'account_number';
UPDATE mapping.kinds SET short_name = 'ip'   WHERE name = 'ip_address';

ALTER TABLE mapping.mappings ADD COLUMN IF NOT EXISTS token VARCHAR(32);

UPDATE mapping.mappings SET token = encode(gen_random_bytes(4), 'hex') WHERE token IS NULL;

ALTER TABLE mapping.mappings ALTER COLUMN token SET NOT NULL;

ALTER TABLE mapping.mappings ADD CONSTRAINT mappings_token_key UNIQUE (token);
